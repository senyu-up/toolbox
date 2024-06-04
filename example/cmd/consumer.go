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

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:              "consumer",
	TraverseChildren: true,
	Short:            "message queue consumer",
	Long:             ``,

	PreRun: func(cmd *cobra.Command, args []string) {
		// 初始化资源与门面
		boot.Consumer(global.ConfigPath)

		// 获取 kafka client
		var kafka = global.GetFacade().GetKafkaClient()

		var awsKafka = global.GetFacade().GetAwsKafkaClient()

		// 注册消费者
		if err := index.RegisterConsumer(kafka); err != nil {
			fmt.Printf("register consumer err %v", err)
			global.ErrChan <- err
			return
		}

		if err := index.RegisterKafkaConsumer(awsKafka); err != nil {
			fmt.Printf("register consumer err %v", err)
			global.ErrChan <- err
			return
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consumer called")

		global.GetFacade().StartKafkaConsumeAsync()

		global.GetFacade().StartAwsKafkaConsumeAsync()

		go func() {
			// 健康检查
			global.ErrChan <- global.GetFacade().StartHealthChecker()
		}()
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)
}
