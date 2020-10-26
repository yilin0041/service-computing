package main

import (
	"encoding/json"
	"fmt"

	"github.com/yilin0041/service-computing/json/myjson"
)

func main() {

	type StuRead struct {
		Name  interface{} `json:"name"`
		Age   interface{}
		High  interface{}
		sex   interface{}
		Class interface{} `json:"-"`
		Test  interface{}
	}
	var stus1 []*StuRead
	stu1 := &StuRead{"asd1", [...]int{12, 2, 3, 7, 50}, 1, 1, 1, map[string]string{"1": "a", "2": "b"}}
	stus1 = append(stus1, stu1)
	fmt.Println("json marshal is")
	json1, _ := json.Marshal(stus1)
	fmt.Println(string(json1))
	fmt.Println("my marshal is")
	b, err := myjson.Marshal(stus1)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))
}
