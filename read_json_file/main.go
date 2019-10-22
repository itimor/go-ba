package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//定义配置文件解析后的结构
type MongoConfig struct {
	MongoAddr      string
	MongoPoolLimit int
	MongoDb        string
	MongoCol       string
}

type Config struct {
	Addr  string
	Mongo MongoConfig
}

func main() {
	v := &Config{}
	//指定对应的json配置文件
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("Sys config read err")
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		panic(err)
	}
	fmt.Println(v.Addr)
	fmt.Println(v.Mongo.MongoDb)
}
