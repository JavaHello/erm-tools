package helper

import (
	"encoding/json"
	"erm-tools/logger"
	"os"
	"path/filepath"
	"strings"
)

var confFile = "./erm-tools.json"

const (
	ERM_ERM = "ERM-ERM"
	ERM_DB  = "ERM-DB"
)

type DBType string

const (
	MYSQL DBType = "mysql"
)

type EnvModel struct {
	ErmFile    string `json:"ermFile"`
	NewErmPath string `json:"newErmPath"`
	OldErmPath string `json:"oldErmPath"`
	DbName     string `json:"dbName"`
	DbHost     string `json:"dbHost"`
	DbUser     string `json:"dbUser"`
	DbPass     string `json:"dbPass"`
	DbType     string `json:"dbType"`
	Type       string `json:"type"`
	OutPath    string `json:"outPath"`
	GenDdl     bool   `json:"genDdl"`
}

var Env EnvModel

func init() {
	logger.Info.Println("env init")
	err := json.Unmarshal(ReadFile(confFile), &Env)
	if err != nil {
		logger.Warn.Println("解析配置文件失败", err)
	}
	Env.verifyEnv()
}

func (env *EnvModel) verifyEnv() {
	if env.Type == "" {
		env.Type = ERM_ERM
	} else {
		env.Type = strings.ToUpper(env.Type)
	}
	if env.DbType == "" {
		env.DbType = string(MYSQL)
	}
	if env.OutPath == "" {
		env.OutPath = "./"
	}
	assertNotEmpty("ermFile 配置错误", env.ErmFile)
	assertNotEmpty("newErmPath 配置错误", env.NewErmPath)
	if env.Type == ERM_ERM {
		assertNotEmpty("oldErmPath 配置错误", env.OldErmPath)
	} else if env.Type == ERM_DB {
		assertNotEmpty("dbName 配置错误", env.DbName)
		assertNotEmpty("dbHost 配置错误", env.DbHost)
		assertNotEmpty("dbUser 配置错误", env.DbUser)
		assertNotEmpty("dbPass 配置错误", env.DbPass)
	} else {
		assertNotEmpty("type 配置错误，可选范围(ERM-ERM|ERM-DB)", env.Type)
	}
}

func (env *EnvModel) NewErmFiles() []string {
	var fileNames []string
	if env.ErmFile != "*.erm" {
		for _, name := range strings.Split(env.ErmFile, ",") {
			fileNames = append(fileNames, env.NewErmPath+string(os.PathSeparator)+name)

		}
		return fileNames
	}
	return findFiles(env.NewErmPath, env.ErmFile)
}

func (env *EnvModel) OldErmFiles() []string {
	var fileNames []string
	if env.ErmFile != "*.erm" {
		for _, name := range strings.Split(env.ErmFile, ",") {
			fileNames = append(fileNames, env.OldErmPath+string(os.PathSeparator)+name)

		}
		return fileNames
	}
	return findFiles(env.OldErmPath, env.ErmFile)
}

func assertNotEmpty(msg, str string) {
	if str == "" || strings.Trim(str, " ") == "" || strings.Trim(str, "　") == "" {
		panic(msg)
	}
}

func findFiles(path, ext string) []string {
	var filesName []string
	fp, err := os.Open(path)
	if err != nil {
		logger.Error.Println("读取文件错误", err)
	}
	defer fp.Close()

	stat, _ := fp.Stat()
	if stat.IsDir() {
		mfs, _ := filepath.Glob(path + string(os.PathSeparator) + "*." + ext)
		filesName = mfs
	} else {
		filesName = append(filesName, path)
	}
	return filesName
}
