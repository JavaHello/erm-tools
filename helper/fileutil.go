package helper

import (
	"erm-tools/logger"
	"io/ioutil"
)

func ReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error.Println("util.ReadFileToString Open error ", err)
		return nil
	}

	return data
}
