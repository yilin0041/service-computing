# 程序包开发，读简单配置文件 v1
## 概述
配置文件（Configuration File，CF）是一种文本文档，为计算机系统或程序配置参数和初始设置。传统的配置文件就是文本行，在 Unix 系统中随处可见，通常使用 `.conf`,`.config`,`.cfg` 作为后缀，并逐步形成了 `key = value` 的配置习惯。在 Windows 系统中添加了对 section 支持，通常用 .ini 作为后缀。面向对象语言的兴起，程序员需要直接将文本反序列化成内存对象作为配置，逐步提出了一些新的配置文件格式，包括 `JSON`，`YAML`，`TOML` 等。
## 任务目标

- 熟悉程序包的编写习惯（idioms）和风格（convetions）
- 熟悉 io 库操作
- 使用测试驱动的方法
- 简单 Go 程使用
- 事件通知
## 任务内容

在 Gitee 或 GitHub 上发布一个读配置文件程序包，第一版仅需要读 ini 配置，配置文件格式案例：

```bash
# possible values : production, development
app_mode = development

[paths]
# Path to where grafana can store temp files, sessions, and the sqlite3 db (if that is used)
data = /home/git/grafana

[server]
# Protocol (http or https)
protocol = http

# The http port  to use
http_port = 9999

# Redirect to correct domain if host header does not match domain
# Prevents DNS rebinding attacks
enforce_domain = true
```
## 注意事项
- 本次开发使用环境为linux,但已经适配Windows
- api文档使用godoc生成
## api文档
api文档见同目录下：api文档.pdf
## 代码分析
这一部分详细分析了api文档中没有提到的内部函数以及api的实现
### 结构体
首先我们需要一个结构体来保存文件名及上一次的读取信息，定义如下：

```go
type Config struct {
	filepath string
	confList []map[string]map[string]string
}
```
### 初始化函数
初始化函数用于判断当前的系统，这里主要用于不同的注释格式。
```go
func init() {
	sysType := runtime.GOOS
	if sysType == "linux" {
		system = true
	}

	if sysType == "windows" {
		system = false
	}
}
```
### 判断section唯一性
我们需要保证每个section都是唯一的，所以需要一个函数进行判断。这个函数主要适用于不同section的切换。这里实现的方法是将现在的list变量即可。

```go
func (c *Config) unique(conf string) bool {
	for _, v := range c.confList {
		for k := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}
```
### 读配置
这个函数是整个程序中最为关键的一个函数，它用来读取配置文件，并将其写入`Config.confList`中供后续使用。
首先我们打开文件，然后读时间戳用于后续判断文件是否更改

```go
	file, err := os.Open(c.filepath)
	if err != nil {
		return nil, err
	}
	lasttime, err = getFileModTime(c.filepath)
	defer file.Close()
```
生成一些变量供我们使用，这里使用的主要有一个用于记录当前`section`的data，`section`记录当前section，`first`是用于记录是否为全局参数（即前面没有section的参数），最后buf为一个缓冲区。

```go
	var data map[string]map[string]string
	var section string
	first := true
	buf := bufio.NewReader(file)
```
接下来是一个for循环，在循环中有：
首先需要读取一行数据。这在上次实验中已经运用过，这里不再赘述。

```go
	l, err := buf.ReadString('\n')
	l = strings.TrimSpace(l)
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		if len(l) == 0 {
			break
		}
	}
```
首先判断系统信息
```go
		if !system && string(l[0]) == "#" || system && string(l[0]) == ";" {
			return nil, errors.New("System Error")
		}
```
如果读到空行或者注释则直接跳过

```go
if (len(l) == 0) || (string(l[0]) == "#") {
		continue
	}
```
如果读到`[`或`]`则说明读到了section，需要进行初始化设置，并将first置`false`。

```go
else if l[0] == '[' && l[len(l)-1] == ']' {
	section = strings.TrimSpace(l[1 : len(l)-1])
	data = make(map[string]map[string]string)
	data[section] = make(map[string]string)
	first = false
		} 
```
如果`first`为true且经过上两次判断可知不是空行、注释也不是section，那么就说明是全局设置，所以我们要进行单独读取。开辟一个新的section同时把它置为空，然后进行读取即可。

```go
else if first {
	section = ""
	data = make(map[string]map[string]string)
	data[section] = make(map[string]string)
	first = false
	i := strings.IndexAny(l, "=")
	if i == -1 {
		continue
	}
	value := strings.TrimSpace(l[i+1 : len(l)])
	data[section][strings.TrimSpace(l[0:i])] = value
	if c.unique(section) == true {
		c.confList = append(c.confList, data)
	}
		} 
```
最后是为key=value字段，我们直接读取即可。最后需要判断一下该section是否为新的section，如果是则需要将其加入map组中。

```go
else {
	i := strings.IndexAny(l, "=")
	if i == -1 {
		continue
	}
	value := strings.TrimSpace(l[i+1 : len(l)])
	data[section][strings.TrimSpace(l[0:i])] = value
	if c.unique(section) == true {
		c.confList = append(c.confList, data)
	}
		}
```
一切完成后，返回得到的map即可。（错误返回空）
### Get和Set函数
get和set函数的结构是十分相似的，都是需要对map进行遍历。不同的是，get只需要找到即可，而set则需要判断该是否已经有该字段，如果有的话，我们需要对其值进行更改，如果没有则需要新建一个map来存储。

#### Get函数
核心如下，对map进行遍历，如果找到则返回该值，如果循环结束则返回错误。
```go
	for _, v := range conf {
		for k, val := range v {
			if k == sec {
				return val[key], nil
			}
		}
```
#### Set函数
首先利用一个内置函数来进行判断是否存在。
```go
for i, v := range data {
		_, ok = v[section]
		index[i] = ok
	}
	i, ok := func(m map[int]bool) (i int, v bool) {
		for i, v := range m {
			if v == true {
				return i, true
			}
		}
		return 0, false
	}(index)
```
如果存在改变其值

```go
	if ok {
		c.confList[i][section][key] = value
		return true
	}
```
如果不存在则新建

```go
	conf[section] = make(map[string]string)
	conf[section][key] = value
	c.confList = append(c.confList, conf)
	return true
```
### 判断两配置相同
该函数利用get函数进行判断两配置是否相同

```go
func equal(conf1 *Config, conf2 *Config) bool {
	for _, v := range conf1.confList {
		for k, val := range v {
			for key, value1 := range val {
				value2, _ := conf2.GetValue(k, key)
				if value1 != value2 {
					return false
				}
			}
		}
	}
	return true
}
```
### 判断时间戳
利用该函数进行时间戳的判断从而判断文件是否被修改。这里主要利用了函数`fi.ModTime().Unix()`进行返回。

```go
func getFileModTime(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Now().Unix(), err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now().Unix(), err
	}
	return fi.ModTime().Unix(), nil
}
```
### 监听
监听函数使用for循环监听时间戳的修改，如果修改则直接break

```go
	for {
		newtime, err := getFileModTime(infile)
		if err != nil {
		}
		if lasttime != newtime {
			lasttime = newtime
			break
		}
	}
```
### Watch函数
利用listen监听，并调用用户方法，最后返回最新的Configuration。
```go
func Watch(filename string, listener ListenFunc) (Configuration *Config, err error) {
	listener.listen(filename)
	listener(filename)
	Configuration = SetConfig(filename)
	return Configuration, err
}
```
## 测试
本部分分为单元测试和功能测试两部分进行。
### 单元测试与模块测试
本部分主要测试各个函数能够正常运行。首先先来介绍用于测试的代码。
#### SetConfig
这里就直接调用即可
```go
func TestSetConfig(t *testing.T) {
	SetConfig("init.ini")
}

func BenchmarkSetConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetConfig("init.ini")
	}
}
```
#### Unique
这里使用任务内容中的实例ini文件进行测试，分别测试了四个方面。
```go
func TestUnique(t *testing.T) {
	conf := SetConfig("init.ini")
	if !conf.unique("name") {
		t.Errorf("[Error]TestUnique 1")
	}
	if conf.unique("paths") {
		t.Errorf("[Error]TestUnique 2")
	}
	if conf.unique("server") {
		t.Errorf("[Error]TestUnique 3")
	}
	if conf.unique("") {
		t.Errorf("[Error]TestUnique 4")
	}
}

func BenchmarkUnique(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conf := SetConfig("init.ini")
		conf.unique("name")
	}
}
```
#### ReadList
该测试主要保证readList函数正常运行，没有访问空指针，溢出等错误。
```go
func TestReadList(t *testing.T) {
	conf := SetConfig("init.ini")
	conf.readList()
}

func BenchmarkReadList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conf := SetConfig("init.ini")
		conf.readList()
	}
}
```
#### GetValue
通过读取ini文件来测试函数是否正常
```go
func TestGetValue(t *testing.T) {
	conf := SetConfig("init.ini")
	s, _ := conf.GetValue("", "app_mode")
	if s != "development" {
		t.Errorf("[Error]TestGetValue")
	}
}

func BenchmarkGetValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conf := SetConfig("init.ini")
		conf.GetValue("", "app_mode")
	}
}
```
#### SetValue
使用getValue来测试setValue是否正常
```go
func TestSetValue(t *testing.T) {
	conf := SetConfig("init.ini")
	conf.SetValue("", "test", "12345")
	s, _ := conf.GetValue("", "test")
	if s != "12345" {
		t.Errorf("[Error]TestSetValue")
	}
}

func BenchmarkSetValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conf := SetConfig("init.ini")
		conf.SetValue("", "test", "12345")
	}
}
```
#### Equal
通过读取同一个文件来测试equal函数的正确性

```go
func TestEqual(t *testing.T) {
	conf1 := SetConfig("init.ini")
	conf2 := SetConfig("init.ini")
	if !equal(conf1, conf2) {
		t.Errorf("Not Equal")
	}
}

func BenchmarkEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conf1 := SetConfig("init.ini")
		conf2 := SetConfig("init.ini")
		equal(conf1, conf2)
	}
}
```
#### GetFileModTime
这里由于基本都是调用了其他函数，所以我们只需要保证它的可执行性，而具体的功能测试在下一部分实现。
```go
func TestGetFileModTime(t *testing.T) {
	getFileModTime("init.ini")
}

func BenchmarkGetFileModTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFileModTime("init.ini")
	}
}
```
#### Listen
Listen函数同样是注重功能的函数，这里我们只需要保证它可以正常运行即可。由于该测试需要人工输入执行，所以我们不进行`Benchmark`测试

```go
func TestListen(t *testing.T) {
	MyListen := func(string) {
	}
	ListenFunc.listen(MyListen, "init.ini")
}
```
#### Watch
watch测试比较复杂，我们运行测试后，在ini文件中找一条注释然后加个空格保存，此时才会完成测试。由于该测试需要人工输入执行，所以我们不进行`Benchmark`测试

```go
func TestWatch(t *testing.T) {
	MyListen := func(string) {
	}
	conf1 := SetConfig("init.ini")
	Watch("init.ini", conf1, MyListen)
	conf2 := SetConfig("init.ini")
	conf2.readList()
	if !equal(conf1, conf2) {
		t.Errorf("Not Equal")
	}
}
```
#### 单元测试结果
##### Test
![test结果](https://img-blog.csdnimg.cn/2020101915075648.png)
##### Benchmark
![Benchmark](https://img-blog.csdnimg.cn/20201019151017618.png)

#### 额外系统测试
利用随便一个测试(这里使用的TestSetConfig)liunx系统
```go
func TestSetConfig(t *testing.T) {
	SetConfig("init.ini")
	if !system {
		t.Errorf("[Error]System")
	}
}
```
测试结果如下：

![Linux](https://img-blog.csdnimg.cn/20201019212739176.png#pic_center)

发现system为true即Linux系统，Windows系统同理
### 功能测试（简单的使用案例）
功能测试(Linux下)的main函数如下所示：

```go
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
```
测试结果如下：
![测试结果](https://img-blog.csdnimg.cn/20201019155222381.png#pic_center)
经检测，符合功能。

## 自定义错误的使用
使用`errors.New()`，在开发中自定义了一些错误，同时，用户也可以在自定义函数中使用该函数来实现自定义的错误。
## 感悟与总结
本次实验为一个简单的ini读取文件，其难点在于如何进行监听文件的改变，这里我是用了时间戳来进行判断，一旦文件更改，时间戳就会改变，listen函数就会中断来通知用户文件已更改。剩下的文件读取等与上次作业基本相同。
总的来说，通过实验，我学会了map的使用，时间戳的使用以及合理的测试方法。相信这对以后的实验会有很大的帮助。