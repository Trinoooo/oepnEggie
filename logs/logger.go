package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	file   *log.Logger
	strout *log.Logger
}

func (l *Logger) Println(args ...any) {
	if l.file != nil {
		l.file.Println(args)
	}
	l.strout.Println(args)
}

var logV1 *Logger

// 初始化Logger
// 请确保默认路径有读写权限
func init() {
	output := openOrCreateLogOutput()
	logV1 = &Logger{
		strout: log.Default(),
	}
	if output != nil {
		logV1.file = log.New(output, Prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	}
}

func V1() *Logger {
	return logV1
}

// openOrCreateLogOutput
// 检查日志文件是否存在，不存在则创建
func openOrCreateLogOutput() *os.File {
	absPath := buildLogOutputFileAbsPath()
	dirPath := filepath.Dir(absPath)
	if err := os.MkdirAll(dirPath, Perm); err != nil {
		log.Println(err)
		return nil
	}

	file, err := os.OpenFile(absPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, Perm)
	if err != nil {
		log.Println(err)
		return nil
	}

	return file
}

func buildLogOutputFileAbsPath() string {
	now := time.Now()
	return fmt.Sprintf(LogOutputFile, now.Format(DateLayout), now.Hour())
}
