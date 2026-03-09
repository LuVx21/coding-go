package cookie

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/dao/redis_dao"
	"luvx/gin/db"

	"github.com/allegro/bigcache/v3"
	"github.com/bytedance/sonic"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/common_x/alias_x"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/infra/infra_sql"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/pbkdf2"
)

var (
	passwordByte = []byte("")
	cache, _     = bigcache.New(context.Background(), bigcache.DefaultConfig(60*time.Minute))
)

func ClearCache() error {
	return cache.Reset()
}

func GetCookieStrByHost(hosts ...string) map[string]string {
	r := make(map[string]string, len(hosts))
	table := GetCookieFromCache(hosts...)
	for _, host := range hosts {
		r[host] = maps_x.Join(table[host], "=", "; ")
	}
	return r
}

func GetCookieFromCache(hosts ...string) alias_x.Table[string] {
	r := make(alias_x.Table[string], len(hosts))
	readDb := make([]string, 0)
	for _, host := range hosts {
		b, err := cache.Get(host)
		if err != nil || len(b) == 0 {
			readDb = append(readDb, host)
		} else {
			kv := make(map[string]any)
			_ = sonic.Unmarshal(b, &kv)
			r[host] = kv
		}
	}
	for host, cookie := range GetCookieFromDb(readDb...) {
		r[host] = cookie
		b, _ := sonic.Marshal(cookie)
		_ = cache.Set(host, b)
	}
	return r
}

func GetCookieFromDb(hosts ...string) alias_x.Table[string] {
	r := make(alias_x.Table[string], len(hosts))
	if len(hosts) == 0 {
		return r
	}
	rows, _ := readDb(hosts...)
	for _, row := range rows {
		host_key := row["host_key"].(string)
		cookie, ok := r[host_key]
		if !ok {
			cookie = make(alias_x.Row)
			r[host_key] = cookie
		}
		cookie[row["name"].(string)] = row["value"].(string)
	}
	return r
}

func Sync2Yun(hosts ...string) {
	table := GetCookieFromDb(hosts...)
	rows := make(alias_x.Rows, 0, len(table))
	for host_key, cookie := range table {
		mongo_dao.CookieCol.DeleteMany(context.TODO(), bson.M{"host_key": host_key})
		doc := map[string]any{
			"_id":      consts.IdWorker.NextId(),
			"host_key": host_key,
			"cookie":   cookie,
		}
		rows = append(rows, doc)
		slog.Debug("推送cookie", host_key, len(cookie))
	}

	mongo_dao.CookieCol.InsertMany(context.TODO(), rows)
	_ = ClearCache()
}

func readDb(hosts ...string) (alias_x.Rows, error) {
	if len(hosts) == 0 {
		return alias_x.Rows{}, nil
	}
	if redis_dao.GetSwitch("remote_cookie") {
		datas, err := mongodb.RowsMap(context.TODO(), mongo_dao.CookieCol, bson.M{"host_key": bson.M{"$in": hosts}})
		return slices_x.Transfer(func(m bson.M) alias_x.Row { return alias_x.Row(m) }, *datas...), err
	}

	_sql := `
select *
from cookies
where true
and host_key in (%s)
order by host_key, name
-- limit 1
;
`
	var args strings.Builder
	for i := range hosts {
		args.WriteString(common_x.IfThen(i == 0, "", ", "))
		fmt.Fprintf(&args, "'%s'", hosts[i])
	}
	_sql = fmt.Sprintf(_sql, args.String())

	rowsMap, err := infra_sql.RowsMap(context.TODO(), db.CookieDb, _sql)
	if err == nil {
		key := masterKey()
		for _, row := range rowsMap {
			encryptedValue, ok := row["encrypted_value"]
			if !ok {
				continue
			}

			real, ok := encryptedValue.([]byte)
			if !ok {
				_, ok = encryptedValue.(string)
				if !ok {
					continue
				}
				real = []byte(encryptedValue.(string))
			}
			value, _ := DecryptWithChromium(key, real)
			// 正式版
			if len(value) > 32 {
				value = value[32:]
			}
			row["value"] = string(value)
		}
	}
	return rowsMap, err
}

func masterKey() []byte {
	if len(passwordByte) == 0 {
		passwordStr := consts.PasswordStr
		if len(passwordStr) == 0 {
			var (
				stdout, stderr bytes.Buffer
			)
			cmd := exec.Command("security", "find-generic-password", "-wa", "Microsoft Edge")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				log.Warnln("执行命令失败:", cmd.String(), stderr.String(), err)
			}
			passwordByte = stdout.Bytes()
		} else {
			passwordByte = []byte(passwordStr)
		}
	}

	secret := bytes.TrimSpace(passwordByte)
	key := pbkdf2.Key(secret, []byte("saltysalt"), 1003, 16, sha1.New)
	return key
}

func DecryptWithChromium(key, password []byte) ([]byte, error) {
	if len(password) <= 3 {
		return nil, errors.New("ciphertext length is invalid")
	}
	var iv = []byte{32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
	return AES128CBCDecrypt(key, iv, password[3:])
}

func AES128CBCDecrypt(key, iv, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Check ciphertext length
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("AES128CBCDecrypt: ciphertext too short")
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("AES128CBCDecrypt: ciphertext is not a multiple of the block size")
	}

	decryptedData := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decryptedData, ciphertext)

	decryptedData, err = pkcs5UnPadding(decryptedData)
	if err != nil {
		return nil, fmt.Errorf("AES128CBCDecrypt: %w", err)
	}

	return decryptedData, nil
}

func pkcs5UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("pkcs5UnPadding: src should not be empty")
	}
	padding := int(src[length-1])
	if padding < 1 || padding > aes.BlockSize {
		return nil, errors.New("pkcs5UnPadding: invalid padding size")
	}
	return src[:length-padding], nil
}
