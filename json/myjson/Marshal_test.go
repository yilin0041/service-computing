package myjson

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMarshalInt(t *testing.T) {
	type test struct {
		testValue interface{}
	}
	var tests1 []*test
	test1 := &test{1}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalInt(reflect.ValueOf(test1.testValue))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "1" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "1")
	}
}

func TestMarshalString(t *testing.T) {
	type test struct {
		testValue interface{}
	}
	var tests1 []*test
	test1 := &test{"1"}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalString(reflect.ValueOf(test1.testValue))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "\"1\"" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "\"1\"")
	}
}
func TestMarshalSlice(t *testing.T) {
	type test struct {
		TestValue interface{}
	}
	var tests1 []*test
	test1 := &test{"1"}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalSlice(reflect.ValueOf(tests1))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "[{\"TestValue\":\"1\"}]" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "[{\"TestValue\":\"1\"}]")
	}
}

func TestMarshalSlideWithLower(t *testing.T) {
	type test struct {
		testValue interface{}
	}
	var tests1 []*test
	test1 := &test{"1"}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalSlice(reflect.ValueOf(tests1))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "[{}]" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "[{}]")
	}
}
func TestMarshalArray(t *testing.T) {
	type test struct {
		testValue interface{}
	}
	var tests1 []*test
	test1 := &test{[...]int{1, 2, 3}}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalArray(reflect.ValueOf(test1.testValue))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "[1,2,3]" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "[1,2,3]")
	}
}

func TestMarshalMap(t *testing.T) {
	type test struct {
		testValue interface{}
	}
	var tests1 []*test
	test1 := &test{map[string]int{"a": 1, "b": 2, "c": 3}}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalMap(reflect.ValueOf(test1.testValue))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "{\"a\":1,\"b\":2,\"c\":3}" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "{\"a\":1,\"b\":2,\"c\":3}")
	}
}

func TestMarshalStruct(t *testing.T) {
	type test struct {
		TestValue interface{}
	}
	var tests1 []*test
	test1 := &test{map[string]int{"a": 1, "b": 2, "c": 3}}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalStruct(reflect.ValueOf(*test1))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "{\"TestValue\":{\"a\":1,\"b\":2,\"c\":3}}" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "{\"TestValue\":{\"a\":1,\"b\":2,\"c\":3}}")
	}
}

func TestMarshalStructWithTag(t *testing.T) {
	type test struct {
		TestValue interface{} `json:"testTag"`
	}
	var tests1 []*test
	test1 := &test{map[string]int{"a": 1, "b": 2, "c": 3}}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalStruct(reflect.ValueOf(*test1))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "{\"testTag\":{\"a\":1,\"b\":2,\"c\":3}}" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "{\"testTag\":{\"a\":1,\"b\":2,\"c\":3}}")
	}
}
func TestMarshalStructWithLower(t *testing.T) {
	type test struct {
		testValue interface{}
	}
	var tests1 []*test
	test1 := &test{map[string]int{"a": 1, "b": 2, "c": 3}}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalStruct(reflect.ValueOf(*test1))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "{}" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "{}")
	}
}
func TestMarshalStructWith_AndTag(t *testing.T) {
	type test struct {
		TestValue interface{} `json:"-"`
	}
	var tests1 []*test
	test1 := &test{map[string]int{"a": 1, "b": 2, "c": 3}}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalStruct(reflect.ValueOf(*test1))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "{}" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "{}")
	}
}
func TestMarshalPtr(t *testing.T) {
	type test struct {
		TestValue interface{}
	}
	var tests1 []*test
	test1 := &test{map[string]int{"a": 1, "b": 2, "c": 3}}
	tests1 = append(tests1, test1)
	var data marshalData
	err := data.marshalPtr(reflect.ValueOf(test1))
	if err != nil {
		t.Fatal(err)
	}
	testB := data.Bytes()
	if (string)(testB) != "{\"TestValue\":{\"a\":1,\"b\":2,\"c\":3}}" {
		t.Fatalf("[Error]Difference between output and expected,\noutput:%s\nexpected:%s\n", (string)(testB), "{\"TestValue\":{\"a\":1,\"b\":2,\"c\":3}}")
	}
}

func TestMarshal(t *testing.T) {
	type StuRead struct {
		TESTINT    interface{}
		TESTSTRING interface{}
		TESTARRAY  interface{}
		TESTMAP    interface{}
		TESTTAG    interface{} `json:"testTag"`
		TESTPTR    *int
		test       interface{}
	}
	a := 123456
	var tests1 []*StuRead
	test1 := &StuRead{123456, "123456", [...]int{1, 2, 3, 4, 5, 6}, map[string]int{"1": 1, "2": 2}, 123456, &a, 123456}
	tests1 = append(tests1, test1)
	output, err := Marshal(tests1)
	if err != nil {
		t.Fatal(err)
	}
	expect, _ := json.Marshal(tests1)
	if string(output) != string(expect) {
		t.Fatalf("[Error]Difference between output and expected: \n%v\n%v\n", string(output), string(expect))
	}
}
