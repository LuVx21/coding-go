package os_x

import (
	"bytes"
	"os"
	"os/exec"
)

const (
	ENV_HTTP_PROXY  = "http_proxy"
	ENV_HTTPS_PROXY = "https_proxy"
	ENV_ALL_PROXY   = "all_proxy"
)

func Getenv(key string) string {
	v, b := LookupEnv(key, "")
	if !b {
		return ""
	}
	return v
}

func LookupEnv(key string, defaultValue string) (string, bool) {
	value, exist := os.LookupEnv(key)
	if exist {
		return value, true
	}

	cmd := exec.Command("sh", "-c", "kv get "+key)
	var out bytes.Buffer
	cmd.Stderr, cmd.Stdout = os.Stderr, &out
	err := cmd.Run()
	if err == nil {
		return out.String(), true
	}

	return defaultValue, false
}
