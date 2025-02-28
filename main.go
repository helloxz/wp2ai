package main

import (
	"fmt"
	"os"
	"wp2ai/router"
	"wp2ai/utils"
)

func main() {
	//获取命令行参数
	args := os.Args
	//获取切片长度
	args_len := len(args)

	//如果参数是1，则没有额外参数
	if args_len == 1 {
		fmt.Printf("Please enter the parameters!\n")
		os.Exit(0)
	} else if args_len == 2 {
		//启动程序
		if args[1] == "start" {
			utils.InitConfig()
			//启动Gin
			router.Start()
		} else {
			fmt.Printf("Please enter the correct parameters!\n")
			os.Exit(0)
		}
	}
}
