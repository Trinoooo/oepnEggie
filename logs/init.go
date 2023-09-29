package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logV1 *log.Logger

// 初始化Logger
// 确保默认路径
func init() {
	output := openOrCreateLogOutput()
	logV1 = log.New(output, Prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}

func V1() *log.Logger {
	return logV1
}

// openOrCreateLogOutput
// 检查日志文件是否存在，不存在则创建
func openOrCreateLogOutput() *os.File {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	absPath := buildLogOutputFileAbsPath()
	dirPath := filepath.Dir(absPath)
	if err := os.MkdirAll(dirPath, Perm); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(absPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, Perm)
	if err != nil {
		panic(err)
	}

	return file
}

func buildLogOutputFileAbsPath() string {
	now := time.Now()
	return fmt.Sprintf(LogOutputFile, now.Format(DateLayout), now.Hour())
}
