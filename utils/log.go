package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	logger  *log.Logger
	logFile *os.File
	nonce   sync.Once
)

// LogMessage 记录日志消息到指定文件
func WriteLog(message string) {
	nonce.Do(func() {

		// 2. 指定日志文件位置
		logFilePath := "data/logs/error.log"
		// 检查目录是否存在，不存在则创建
		if _, err := os.Stat("data/logs"); os.IsNotExist(err) {
			os.Mkdir("data/logs", os.ModePerm)
		}

		// 3. 创建或打开日志文件
		var err error
		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Failed to open log file: %v, logging to stderr\n", err)
			logFile = os.Stderr // 如果打开文件失败，则输出到标准错误
		}

		// 4. 创建 logger 实例
		logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	})

	// 5. 写入日志消息
	logger.Println(message)
}

// CloseLogFile 关闭日志文件,正常来说应该在程序退出时调用
func CloseLogFile() {
	if logFile != nil && logFile != os.Stderr {
		logFile.Close()
	}
}
