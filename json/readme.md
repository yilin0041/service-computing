# 对象序列化支持包开发
将一个对象写成特定文本格式的字符流，称为序列化。
## 课程任务
- 参考官方 encoding/json 包 Marshal 函数，将结构数据格式化为 json 字符流
	* 必须导出 func JsonMarshal(v interface{}) ([]byte, error)
	* 可以参考、甚至复制原来的代码
	* 支持字段的标签（Tag），标签满足 mytag:"你自己的定义"
	* 不允许使用第三方包
- 包必须包括以下内容：
	* 生成的中文 api 文档
	* 有较好的 Readme 文件，包括一个简单的使用案例
	* 每个go文件必须有对应的测试文件

## 注意事项
- 本次开发使用环境为linux
- api文档使用godoc生成
## api文档
api文档见同目录下：api文档.pdf
## 代码分析
这一部分详细分析了api文档中没有提到的内部函数以及api的实现
### 数据结构
只需要定义一个bytes.buffer即可
```go
type marshalData struct {
	bytes.Buffer
}
```

###  顶层函数Marshal
传入的是一个接口，然后调用`marshal`然后返回字节流即可
```go
func Marshal(v interface{}) ([]byte, error) {
	var data marshalData
	val := reflect.ValueOf(v)
	if err := data.marshal(val); err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}
```
###  选择函数marshal
这里需要根据类型进行选择不同的调用函数
```go
func (data *marshalData) marshal(val reflect.Value) error {
	if val.Kind() == reflect.Int {
		return data.marshalInt(val)
	} else if val.Kind() == reflect.String {
		return data.marshalString(val)
	} else if val.Kind() == reflect.Slice {
		return data.marshalSlice(val)
	} else if val.Kind() == reflect.Array {
		return data.marshalArray(val)
	} else if val.Kind() == reflect.Map {
		return data.marshalMap(val)
	} else if val.Kind() == reflect.Struct {
		return data.marshalStruct(val)
	} else if val.Kind() == reflect.Ptr {
		return data.marshalPtr(val)
	} else if val.Kind() == reflect.Interface {
		return data.marshalInterface(val)
	} else {
		return errors.New("Unknown type  " + val.Kind().String())
	}
}
```
###  处理int
这里只需要写入字节流即可，不需要进行额外的操作
```go
func (data *marshalData) marshalInt(val reflect.Value) error {
	if _, err := data.Write(strconv.AppendInt([]byte{}, val.Int(), 10)); err != nil {
		return err
	}
	return nil
}
```
###  处理string
这里需要先传入双引号，然后再将string值传入字节流
```go
func (data *marshalData) marshalString(val reflect.Value) error {
	if err := data.WriteByte('"'); err != nil {
		return err
	}
	if _, err := data.Write([]byte(val.String())); err != nil {
		return err
	}
	if err := data.WriteByte('"'); err != nil {
		return err
	}
	return nil
}
```
### 处理slice
切片需要使用中括号括起来然后对需要处理的切片中的值分别调用`marshal`函数，剩下的部分与`string`的处理方式类似

```go
func (data *marshalData) marshalSlice(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		if err := data.WriteByte('['); err != nil {
			return err
		}
		for i := 0; i < val.Len(); i++ {
			element := reflect.ValueOf(val.Index(i).Interface())
			if err := data.marshal(element); err != nil {
				return err
			}
			if i != val.Len()-1 {
				if err := data.WriteByte(','); err != nil {
					return err
				}
			}
		}
		return data.WriteByte(']')
	}
	valBytes := val.Bytes()
	if err := data.WriteByte('"'); err != nil {
		return err
	}
	if _, err := data.Write(valBytes); err != nil {
		return err
	}
	if err := data.WriteByte('"'); err != nil {
		return err
	}
	if err := data.WriteByte(':'); err != nil {
		return err
	}
	return nil
}
```
###  处理Array
处理数组的时候与处理切片是类似的，只是切片处理的时候是按照`reflect.Uint8`处理的，因此数组这里我们就不对该类型进行处理。
```go
func (data *marshalData) marshalArray(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		if err := data.WriteByte('['); err != nil {
			return err
		}
		for i := 0; i < val.Len(); i++ {
			element := reflect.ValueOf(val.Index(i).Interface())
			if err := data.marshal(element); err != nil {
				return err
			}
			if i != val.Len()-1 {
				if err := data.WriteByte(','); err != nil {
					return err
				}
			}
		}
		return data.WriteByte(']')
	}
	return errors.New("Unknown type  " + elemKind.String())
}
```
###  处理Map
map的处理是最为复杂的，这里借鉴了网上的写法，同时为了简便，只能够处理string类型的key的map，这里需要注意的是，我们在传入的时候会将原来的顺序打乱，即key与value不对应。因此我们需要使用sort来将其重新对应。剩下的就是写入字符流等与之前类似的操作。

```go
func (data *marshalData) marshalMap(val reflect.Value) error {
	keys := val.MapKeys()
	if err := data.WriteByte('{'); err != nil {
		return err
	}
	raw := make(sortableByteSliceSlice, len(keys))
	for i, key := range keys {
		if key.Kind() != reflect.String {
			e := "Map keys must be 'string' type,your keys is '" + (key.Kind().String()) + "' type"
			return errors.New(e)
		}
		raw[i] = []byte(key.String())
	}
	sort.Sort(raw)
	i := 0
	for _, rawKey := range raw {
		key := string(rawKey)
		vKey := reflect.ValueOf(key)
		if err := data.marshal(vKey); err != nil {
			return err
		}
		if err := data.WriteByte(':'); err != nil {
			return err
		}
		value := val.MapIndex(vKey)
		if err := data.marshal(value); err != nil {
			return err
		}
		if i != raw.Len()-1 {
			if err := data.WriteByte(','); err != nil {
				return err
			}
		}
		i++
	}
	return data.WriteByte('}')
}
```
###  处理Tag
标签的处理使用`regexp.MustCompile`对其标签进行正则匹配，识别标签`json`（方便后续使用源库检查），同时如果为“-”则忽略，下面就是简单的字符串读取匹配与返回。
```go
func myTag(value reflect.Value, name string) string {
	var tag string
	field, hasField := value.Type().FieldByName(name)
	if !hasField {
		tag = ""
	} else {
		tag = string(field.Tag)
	}
	var readTag *string
	const fieldRegexp = `json:"([\w- ]*)"`
	reg := regexp.MustCompile(fieldRegexp)
	if matches := reg.FindStringSubmatch(tag); len(matches) > 2 {
		panic("regexp returns more then two groups!")
	} else if len(matches) == 2 {
		readTag = &matches[1]
	} else {
		readTag = nil
	}
	if readTag == nil {
		return name
	} else if *readTag == "" || *readTag == "-" {
		return ""
	} else {
		return *readTag
	}
}
```
###  处理struct
`struct`的处理过程也比较的容易理解（相对于`map`），就是先进行标签处理，再进行判断，如果首字母是小写则直接忽略，如果是大写则对他们的值调用`marshal`函数，动态循环处理。
```go
func (data *marshalData) marshalStruct(val reflect.Value) error {
	if err := data.WriteByte('{'); err != nil {
		return err
	}
	valType := val.Type()

	fields := positionedFieldsByName{}
	count := 0
	for i := 0; i < val.NumField(); i++ {
		fieldOpt := myTag(val, valType.Field(i).Name)
		if len(fieldOpt) == 0 {
			count++
			continue
		}
		temp := (string)(valType.Field(i).Name)
		if temp[0] < 'A' || temp[0] > 'Z' {
			count++
			continue
		}
		fields = append(fields, positionedField{[]byte(fieldOpt), i})
	}
	for _, f := range fields {
		count++
		if err := data.marshal(reflect.ValueOf(f.name)); err != nil {
			return err
		}
		if err := data.marshal(val.Field(f.pos)); err != nil {
			return err
		}
		if count != val.NumField() {
			if err := data.WriteByte(','); err != nil {
				return err
			}
		}
	}
	return data.WriteByte('}')
}
```
###  处理指针和接口
指针和接口的处理方法完全相同，即直接对其`Elem()`调用`marshal`即可

```go
func (data *marshalData) marshalPtr(val reflect.Value) error {
	return data.marshal(val.Elem())
}

func (data *marshalData) marshalInterface(val reflect.Value) error {
	return data.marshal(val.Elem())
}
```
至此，基本代码已经完成。详细的细节还需参照源码。
## 测试
### 单元函数测试

首先讲一下单元测试的思路，因为函数的结构大部分很相似，故没有特别说明的都是按照该思路。

- 先使用一个结构体定义一个interface，这样可以方便后续的使用
- 再给该interface定义想要测试的数据类型
- 通过调用需要测试的函数得到一个字节流
- 判断调用过程中有没有发生错误
- 判断读取的字节流是否跟期望相同
- 注意
	* 前面函数的测试中，有些是可以传入开头非大写的参数的，但是slice和struct测试中不可以，这里先行说明
	* 在一些测试中，可能定义了一些多余的变量，这些变量是为了方便调试使用的，另外也保持了测试函数的基本统一
	* 每个函数均有测试结果，在最后也会给出最终的全部测试结果
#### MarshalInt测试

```go
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
```
测试结果：

![int](https://img-blog.csdnimg.cn/20201026184125695.png#pic_center)
#### MarshalString测试

```go
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
```
测试结果：

![string](https://img-blog.csdnimg.cn/20201026184538775.png#pic_center)
#### MarshalSlice测试
这里需要注意传入的是一个切片
```go
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
```
测试结果：

![slice](https://img-blog.csdnimg.cn/20201026184752149.png#pic_center)

#### TestMarshalSlideWithLower测试
测试用于输入小写字母的时候的判断
```go
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
```
测试结果：

![slice lower](https://img-blog.csdnimg.cn/20201026185157789.png#pic_center)

#### MarshalArray测试

```go
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
```
测试结果：

![array](https://img-blog.csdnimg.cn/20201026185334960.png#pic_center)

#### MarshalMap测试

```go
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
```
测试结果：

![map](https://img-blog.csdnimg.cn/20201026185410372.png#pic_center)

#### MarshalStruct测试
##### 普通测试
```go
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
```
测试结果：

![struct](https://img-blog.csdnimg.cn/20201026185458901.png#pic_center)
##### 带Tag的测试

```go
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
```
测试结果：

![structTag](https://img-blog.csdnimg.cn/20201026185548906.png#pic_center)
##### 开头小写

```go
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
```
测试结果：

![structLower](https://img-blog.csdnimg.cn/20201026185742139.png#pic_center)
##### Tag中的忽略符

```go
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
```
测试结果：

![struct](https://img-blog.csdnimg.cn/20201026185839733.png#pic_center)

#### MarshalPtr测试

```go
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
```
测试结果：

![ptr](https://img-blog.csdnimg.cn/20201026190004355.png#pic_center)
interface与指针是同理的，故不再进行测试
#### Marshal测试（集成功能）
这里对基本以上所有的类型都进行了测试，同时使用函数`json.Marshal`进行结果对比
```go
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
```
测试结果：

![marshal](https://img-blog.csdnimg.cn/20201026190237377.png#pic_center)

全部测试结果：

![all](https://img-blog.csdnimg.cn/20201026190326524.png#pic_center)

### 功能测试（简单的实例）
- 实例测试中同样使用了`json.Marchal`函数做对比
- 使用接口简化定义
```go
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

```
测试结果：

![test](https://img-blog.csdnimg.cn/20201026191047978.png#pic_center)
## 感悟与总结
本次实验为一个仿写`json.Marchal`函数。其实，如果想完全的写出该函数并不容易，但是写简化版的还是比较简单的。在本次实验中，最重要的其实就是错误的及时捕捉和处理，因为在调用函数时一般都是会返回一个错误，但这个错误有可能不被捕捉，而且没有设置中断，那么我们可能就不知道错误在哪里发生。
另外，我们还要注意格式，这次实验中只实现了简单的格式，更为复杂的格式有可能会出错，但还是完成了老师布置的读取Tag的基本要求。其实，让我印象最深刻的应该是map的处理，查询了大量资料和代码才搞定。
最后，本次实验中积累的经验，相信会对以后的实验会有很大的帮助。