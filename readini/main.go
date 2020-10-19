package main

import (
	"fmt"

	"github.com/yilin0041/service-computing/readini/ini"
)

func main() {
	var c *ini.Config
	c = ini.SetConfig("init.ini")
	MyListen := func(string) {
		fmt.Printf("My listen test!\n")
	}
	c, err := ini.Watch("init.ini", MyListen)
	//人工在server中加入test=123
	if err != nil {
		fmt.Printf("error in watch!\n")
	}
	v1, _ := c.GetValue("", "app_mode")
	v2, _ := c.GetValue("paths", "data")
	v3, _ := c.GetValue("server", "test")
	fmt.Printf("%s\n", v1)
	fmt.Printf("%s\n", v2)
	fmt.Printf("%s\n", v3)
}
