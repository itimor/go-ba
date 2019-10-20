package conf

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

var Sysconfig = &sysconfig{}

func init() {
	//指定对应的json配置文件
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Sys config read err")
	}
	err = jsoniter.Unmarshal(b, Sysconfig)
	if err != nil {
		panic(err)
	}
}

type AppConfigStruct struct {
	ServerPort string   `json:"ServerPort"`
	LogLevel   string   `json:"LogLevel"`
	Secret     string   `json:"Secret"`
	IgnoreURLs []string `json:"IgnoreURLs"`
	JWTTimeout string   `json:"JWTTimeout"`
	Port       string   `json:"Port"`
}

type TestDBConfigStruct struct {
	DriverName string `json:"DriverName"`
	DBFile     string `json:"DBFile"`
	Charset    string `json:"Charset"`
}

type sysconfig struct {
	App AppConfigStruct
	DB  TestDBConfigStruct
}
