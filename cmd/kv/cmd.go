package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/luvx21/coding-go/coding-common/fmt_x"
	"github.com/urfave/cli/v3"
)

func buildCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "列举所有键值对",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "bucket", Aliases: []string{"b"}, Usage: "所有bucket"},
				&cli.BoolFlag{Name: "key", Aliases: []string{"k"}, Usage: "所有key, 可使用@xxx指定bucket"},
				&cli.BoolFlag{Name: "value", Aliases: []string{"v"}, Usage: "所有value, 可使用@xxx指定bucket"},
			},
			Action: list,
		},
		{
			Name:   "get",
			Usage:  "根据key获取value",
			Action: func(ctx context.Context, c *cli.Command) error { return getCmd(ctx, c, false) },
		},
		{
			Name:   "cp",
			Usage:  "获取一个值并放置到剪贴板",
			Action: func(ctx context.Context, c *cli.Command) error { return getCmd(ctx, c, true) },
		},
		{
			Name:  "set",
			Usage: "设置键值对,支持同时设置多对",
			Action: func(ctx context.Context, c *cli.Command) error {
				if c.Args().Len() < 2 {
					fmt_x.Warningln("You must give a key-value pair")
					return nil
				} else if c.Args().Len()%2 != 0 {
					fmt_x.Warningln("键值对数量不匹配")
					return nil
				}

				for i := 0; i < c.Args().Len()-1; i = i + 2 {
					k, v := c.Args().Get(i), c.Args().Get(i+1)
					k, bucket := getBucket(k)
					if err := set(bucket, k, v); err != nil {
						fmt_x.Errorln(err.Error())
					} else {
						fmt.Printf("%s = %s\n", k, v)
					}
				}
				return nil
			},
		},
		{
			Name:  "del",
			Usage: "移除键值对,支持删除多个",
			Action: func(ctx context.Context, c *cli.Command) error {
				if c.Args().Len() < 1 {
					fmt_x.Warningln("You must give a key to delete")
					return nil
				}

				for _, k := range c.Args().Slice() {
					k, bucket := getBucket(k)
					if err := del(bucket, k); err != nil {
						fmt_x.Errorln(err.Error())
					} else {
						fmt.Printf("key: %s 被移除\n", k)
					}
				}
				return nil
			},
		},
	}
}

func list(ctx context.Context, c *cli.Command) error {
	if c.Bool("bucket") {
		list_b(ctx, c, true)
		return nil
	}

	bucket := ""
	if c.Args().Len() > 0 {
		parts := strings.Split(c.Args().Get(0), "@")
		if len(parts) >= 2 {
			bucket = parts[1]
		}
	}
	pairs, err := getAll(bucket)
	if err != nil {
		fmt_x.Errorln(err.Error())
		return nil
	}

	if c.Bool("key") {
		list_k(pairs)
		return nil
	}
	if c.Bool("value") {
		list_v(pairs)
		return nil
	}

	for k, v := range pairs {
		fmt.Printf("%s %s\n", k, v)
	}

	return nil
}

func list_b(context.Context, *cli.Command, bool) error {
	bs, err := listAllBuckets()
	if err != nil {
		fmt_x.Errorln(err.Error())
		return nil
	}
	for _, b := range bs {
		fmt.Println(b)
	}
	return nil
}
func list_k(pairs map[string]string) {
	for k := range pairs {
		fmt.Println(k)
	}
}
func list_v(pairs map[string]string) {
	for _, v := range pairs {
		fmt.Println(v)
	}
}

func getCmd(_ context.Context, c *cli.Command, cp bool) error {
	if c.Args().Len() < 1 {
		fmt_x.Warningln("You must give a key")
		return nil
	}

	k, bucket := getBucket(c.Args().Get(0))
	if v, err := get(bucket, k); err != nil {
		fmt_x.Errorln(err.Error())
	} else {
		if cp {
			clipboard.WriteAll(string(v))
			fmt_x.Successf("剪贴板: %s\n", v)
		} else {
			fmt.Println(v)
		}
	}
	return nil
}

func getBucket(k string) (string, string) {
	bucket, parts := "", strings.Split(k, "@")
	if len(parts) >= 2 {
		k, bucket = parts[0], parts[1]
	}
	return k, bucket
}
