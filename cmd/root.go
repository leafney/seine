package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seine",
	Short: "Seine - Go Web项目基础框架",
	Long:  `Seine是一个用于快速创建Go Web项目的基础框架生成工具`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(addCmd)
}
