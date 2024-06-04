/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/senyu-up/toolbox/example/boot"
	"github.com/senyu-up/toolbox/example/global"
	"github.com/senyu-up/toolbox/example/index"
	middleware2 "github.com/senyu-up/toolbox/tool/http/gin_server/middleware"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "",
	TraverseChildren: true,
	Short:            "default cmd, start http、grpc server",
	Long:             ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("server pre run called")
		fmt.Printf("root cmd init called, %s \n", global.ConfigPath)

		// 初始化服务的资源
		err := boot.Boot(global.ConfigPath)
		if err != nil {
			global.ErrChan <- fmt.Errorf("boot err %v", err)
			return
		}

		// 获取 fiber app
		//var app = global.GetFacade().GetFiber()
		//if app == nil {
		//	global.ErrChan <- fmt.Errorf("fiber app is nil")
		//	return
		//}

		// fiber app 注册中间件
		//app.Fiber().Use(
		//	middleware.Cors(),
		//	middleware.PanicRecover())
		//
		//// 注册路由
		//index.RegisterRouter(app.Fiber())

		//获取gin app
		var ginApp = global.GetFacade().GetGin()
		if ginApp == nil {
			global.ErrChan <- fmt.Errorf("gin app is nil")
			return
		}
		//注册gin 中间件
		ginApp.Gin().Use(middleware2.SetRequestId())
		index.RegisterRouterG(ginApp.Gin())
		// 检查 grpc 初始化与否
		//var grpcSer = global.GetFacade().GetGrpcServer()
		//if grpcSer == nil {
		//	global.ErrChan <- fmt.Errorf("grpc server is nil")
		//	return
		//}
		// 注册 grpc handler
		//grpcSer.Register(index.RegisterHandler)
	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		go func() {
			// 启动服务, 接收错误并投递到全局错误通道
			global.ErrChan <- global.GetFacade().StartFiber()
		}()

		go func() {
			// 启动服务, 接收错误并投递到全局错误通道
			global.ErrChan <- global.GetFacade().StartGin()
		}()
		//go func() {
		//	// 启动 grpc
		//	global.ErrChan <- global.GetFacade().StartGrpc()
		//}()
		go func() {
			global.ErrChan <- global.GetFacade().StartHealthChecker()
		}()
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if cmd.Use == "script" || cmd.Use == "version" { // 如果是脚本命令，执行后直接退出程序
			return
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			select {
			case err := <-global.ErrChan:
				log.Printf("get global error %v\n", err)
				global.GetFacade().Shutdown(global.Ctx)
				global.Canel() // 全局 ctx 取消
				return
			case <-global.Ctx.Done():
				log.Printf("get global ctx canceld\n")
				global.GetFacade().Shutdown(context.Background())
				return
			case s := <-c:
				log.Printf("get a signal %s\n", s.String())
				global.GetFacade().Shutdown(global.Ctx)
				global.Canel() // 全局 ctx 取消
				return
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&global.ConfigPath, "conf", "c", ".", "Config file path")
	return
}
