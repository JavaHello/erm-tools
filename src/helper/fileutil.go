package helper

import (
	"io/ioutil"
	"log"
)

func ReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("util.ReadFileToString Open error ", err)
		return nil
	}

	return data
}
