package cookie

import (
    "bytes"
    "context"
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha1"
    "database/sql"
    "errors"
    "fmt"
    "github.com/allegro/bigcache/v3"
    "github.com/bytedance/sonic"
    lcommon "github.com/luvx21/coding-go/coding-common/common_x"
    "github.com/luvx21/coding-go/coding-common/dbs"
    "github.com/luvx21/coding-go/coding-common/maps_x"
    "golang.org/x/crypto/pbkdf2"
    "log"
    "luvx/gin/db"
    "os/exec"
    "time"
)

var (
    passwordByte = []byte("")
    client       *sql.DB
    cache, _     = bigcache.New(context.Background(), bigcache.DefaultConfig(60*time.Minute))
)

func ClearCache() error {
    return cache.Reset()
}

func GetCookieStrByHost(hosts ...string) string {
    resultMap := make(map[string]string)
    for _, host := range hosts {
        tt := GetCookieByHostCache(host)
        for k, v := range tt {
            resultMap[k] = v
        }
    }
    return maps_x.Join(resultMap, "=", "; ")
}

func GetCookieByHostCache(host string) map[string]string {
    b, err := cache.Get(host)
    if b == nil || err != nil {
        m := GetCookieByHost(host)
        b, _ = sonic.Marshal(m)
        _ = cache.Set(host, b)
    }
    tt := make(map[string]string)
    _ = sonic.Unmarshal(b, &tt)
    return tt
}

func GetCookieByHost(hosts ...string) map[string]string {
    rows, _ := readDb(hosts...)
    result := make(map[string]string)
    for _, row := range rows {
        result[row["name"].(string)] = row["value"].(string)
    }
    return result
}

func Sync2Turso(hosts ...string) {
    key := masterKey()
    rowsMap, err := readDb(hosts...)
    if err != nil {
        log.Fatalln("读取cookie异常", err)
    }

    rows := make([][]any, 0, len(rowsMap))
    for _, row := range rowsMap {
        value, _ := DecryptWithChromium(key, row["encrypted_value"].([]byte))
        values := []any{row["host_key"], row["name"], string(value)}
        rows = append(rows, values)
    }
    _, _ = db.Turso.Exec("delete from cookies;")
    for i, row := range rows {
        _, _ = db.Turso.Exec("insert into cookies(host_key, name, value) values(?, ?, ?)", row...)
        fmt.Println("行:", i)
    }
}

func readDb(hosts ...string) ([]map[string]any, error) {
    _sql := `
select *
from cookies
where true
and host_key in (%s)
order by host_key, name
-- limit 1
;
`
    var args string
    for i := 0; i < len(hosts); i++ {
        args += lcommon.IfThen(i == 0, "", ", ") + fmt.Sprintf("'%s'", hosts[i])
    }
    _sql = fmt.Sprintf(_sql, args)

    getClient()
    rowsMap, err := dbs.RowsMap(context.TODO(), client, _sql)
    if err == nil {
        key := masterKey()
        for _, row := range rowsMap {
            value, _ := DecryptWithChromium(key, row["encrypted_value"].([]byte))
            row["value"] = string(value)
        }
    }
    return rowsMap, err
}

func getClient() *sql.DB {
    if client == nil {
        home, _ := lcommon.Dir()
        client, _ = db.GetDataSource(home + "/data/sqlite/Cookies")
    }
    return client
}

func masterKey() []byte {
    if len(passwordByte) == 0 {
        var (
            stdout, stderr bytes.Buffer
        )
        cmd := exec.Command("security", "find-generic-password", "-wa", "Microsoft Edge")
        cmd.Stdout = &stdout
        cmd.Stderr = &stderr
        if err := cmd.Run(); err != nil {
            fmt.Println(cmd.String())
            //fmt.Printf("run security command failed: %w, message %s", err, stderr.String())
        }
        passwordByte = stdout.Bytes()
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
