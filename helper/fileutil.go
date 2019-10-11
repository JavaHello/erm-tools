package helper

import (
	"erm-tools/logger"
	"io/ioutil"
)

func ReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("读取文件失败", err)
		return nil
	}

	return data
}
