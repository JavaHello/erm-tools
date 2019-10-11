package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	info   *log.Logger
	warn   *log.Logger
	errLog *log.Logger
)

func init() {

	info = log.New(os.Stdout, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	warn = log.New(os.Stdout, "Warning:", log.Ldate|log.Ltime|log.Lshortfile)

}

func Info(msg ...interface{}) {
	info.Println(msg)
}

func Warn(msg ...interface{}) {
	warn.Println(msg)
}
func Error(msg ...interface{}) {
	if errLog == nil {
		err := errorInit()
		if err != nil {
			fmt.Println(msg)
			return
		}
	}
	errLog.Println(msg)
}

func errorInit() error {
	errFile, err := os.OpenFile("errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
		return err
	}
	errLog = log.New(io.MultiWriter(os.Stderr, errFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
