package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/leafney/seine/internal/generator"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [项目名称]",
	Short: "创建新的Go Web项目",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		projectPath, _ := filepath.Abs(projectName)
		
		g := generator.NewGenerator(projectName, projectPath)
		if err := g.Generate(); err != nil {
			return err
		}
		
		fmt.Printf("项目 %s 创建成功!\n", projectName)
		return nil
	},
}
