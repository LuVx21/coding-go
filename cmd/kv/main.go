package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:        "kv",
		Usage:       "简单的key-value存储cli工具, 支持读取到剪贴板",
		Description: "set key[@DB] value...;\nget key[@DB];\nlist [-b|-k|-v] [@DB];\ndel key[@DB]...",
		Commands:    buildCommands(),
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
