package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "luvx",
	Short: "luvx CLI tool",
	Long:  `一个例子`,
}

var foo string
var cmd1 = &cobra.Command{
	Use:   "cmd1",
	Short: "cmd1子命令简短描述",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("子命令: cmd1, 参数search: %s, foo: %s, 内容: %s\n", search, foo, strings.Join(args, " "))
	},
}

var bar string
var cmd2 = &cobra.Command{
	Use:   "cmd2",
	Short: "cmd2子命令简短描述",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("子命令: cmd2, 参数search: %s, foo: %s, 内容: %s\n", search, bar, strings.Join(args, " "))
	},
}

var search string

func init0() {
	// 全局持续参数
	//rootCmd.PersistentFlags().StringVarP(&search, "search", "s", "search默认值", "search参数说明")

	// 子命令共有参数
	for _, cmd := range []*cobra.Command{cmd1, cmd2} {
		cmd.Flags().StringVarP(&search, "search", "s", "search默认值", "search参数说明")
	}

	cmd1.Flags().StringVarP(&foo, "foo", "f", "foo默认值", "foo参数说明")

	cmd2.Flags().StringVarP(&bar, "bar", "b", "bar默认值", "bar参数说明")

	// 将子命令添加到主命令中
	rootCmd.AddCommand(cmd1, cmd2)
}

func main_cobra() {
	init0()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
