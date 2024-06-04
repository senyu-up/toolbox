/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/senyu-up/toolbox/example/boot"
	"github.com/senyu-up/toolbox/example/global"
	"github.com/senyu-up/toolbox/example/index"
	"github.com/spf13/cobra"
)

// cronjobCmd represents the cronjob command
var cronjobCmd = &cobra.Command{
	Use:              "cron",
	TraverseChildren: true,
	Short:            "cron, can only run as a single instance!",
	Long:             ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		// 初始化资源
		boot.Script(global.ConfigPath)

		// 注册 cron job 表达式与任务
		index.CronRegister(global.GetFacade().GetCronClient())
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cronjob called")

		// 启动 cron job
		global.GetFacade().StartCronAsync()

		go func() {
			// 健康检查
			global.ErrChan <- global.GetFacade().StartHealthChecker()
		}()
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
}
