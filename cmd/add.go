package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/leafney/seine/internal/generator"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [组件类型] [组件名称]",
	Short: "添加项目组件",
	Long:  `添加项目组件，如服务(service)、控制器(controller)或数据库(db)`,
}

var addServiceCmd = &cobra.Command{
	Use:   "service [服务名称]",
	Short: "添加服务组件",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName := args[0]
		currentDir, _ := filepath.Abs(".")

		g := generator.NewComponentGenerator(filepath.Base(currentDir), currentDir)
		if err := g.AddService(serviceName); err != nil {
			return err
		}

		fmt.Printf("服务 %s 添加成功!\n", serviceName)
		return nil
	},
}

var addControllerCmd = &cobra.Command{
	Use:   "controller [控制器名称]",
	Short: "添加控制器组件",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		controllerName := args[0]
		currentDir, _ := filepath.Abs(".")

		g := generator.NewComponentGenerator(filepath.Base(currentDir), currentDir)
		if err := g.AddController(controllerName); err != nil {
			return err
		}

		fmt.Printf("控制器 %s 添加成功!\n", controllerName)
		return nil
	},
}

var addDBCmd = &cobra.Command{
	Use:   "db [数据库类型]",
	Short: "添加数据库支持",
	Long:  `添加数据库支持，支持的类型有: mysql, sqlite, postgresql, mongodb`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dbType := args[0]

		// 验证数据库类型
		validTypes := map[string]bool{
			"mysql":      true,
			"sqlite":     true,
			"postgresql": true,
			"mongodb":    true,
		}

		if !validTypes[dbType] {
			return fmt.Errorf("不支持的数据库类型: %s，支持的类型有: mysql, sqlite, postgresql, mongodb", dbType)
		}

		currentDir, _ := filepath.Abs(".")

		g := generator.NewComponentGenerator(filepath.Base(currentDir), currentDir)
		if err := g.AddDatabase(dbType); err != nil {
			return err
		}

		fmt.Printf("数据库 %s 支持添加成功!\n", dbType)
		return nil
	},
}

var addRouteCmd = &cobra.Command{
	Use:   "route [路由名称]",
	Short: "添加API路由",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		routeName := args[0]
		currentDir, _ := filepath.Abs(".")

		g := generator.NewComponentGenerator(filepath.Base(currentDir), currentDir)
		if err := g.AddRoute(routeName); err != nil {
			return err
		}

		fmt.Printf("路由 %s 添加成功!\n", routeName)
		return nil
	},
}

func init() {
	addCmd.AddCommand(addServiceCmd)
	addCmd.AddCommand(addControllerCmd)
	addCmd.AddCommand(addDBCmd)
	addCmd.AddCommand(addRouteCmd)
}
