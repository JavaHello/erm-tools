package helper

import (
	"fmt"
	"io/ioutil"
	"os"
)

var confFile = "../conf/erm-tools.conf"

func init() {
	fmt.Println("env init")
	fp, err := os.Open(confFile)
	defer func() {
		if fp != nil {
			fp.Close()
		}
	}()
	if err != nil {
		fmt.Println("读取配置文件失败", err)
		return
	}

	data, _ := ioutil.ReadAll(fp)
	fmt.Println(string(data))
}

func readEnv() {

}
