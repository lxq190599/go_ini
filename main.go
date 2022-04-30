package main

import (
	"fmt"
	"go_ini/my_ini"
)

type Appconfig struct {
	MysqlConfig `ini:"mysql"`
	RedisConfig `ini:"redis"`
}

//ini配置文件反射
type MysqlConfig struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Passwd   string `ini:"passwd"`
}
type RedisConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Passwd   string `ini:"passwd"`
	Database int    `ini:"database"`
	test     bool   `ini:"test"`
}

func main() {
	var app = new(Appconfig)
	err := my_ini.LoadIni("./my_ini/conf.ini", app)
	if err != nil {
		fmt.Println("加载文件失败:", err)
		return
	}
	fmt.Printf("%#v\n", app)
}
