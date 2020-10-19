package ini

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

//Config : ini的结构体
type Config struct {
	filepath string
	confList []map[string]map[string]string
}

var lasttime int64
var system bool

func init() {
	sysType := runtime.GOOS
	if sysType == "linux" {
		system = true
	}

	if sysType == "windows" {
		system = false
	}
}

const (
	notFindValue  = "[Error]No This Value\n"
	keyNotUnique  = "[Error]Key is repetitive\n"
	fileError     = "[Error]The file can't open\n"
	fileReadError = "[Error]The file can't read\n"
)

//SetConfig ：初始化一个设置文件
func SetConfig(filepath string) *Config {
	conf := new(Config)
	conf.filepath = filepath
	conf.readList()
	return conf
}

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

func (c *Config) readList() ([]map[string]map[string]string, error) {
	file, err := os.Open(c.filepath)
	if err != nil {
		fmt.Printf("not open")
		return nil, err
	}
	lasttime, err = getFileModTime(c.filepath)
	defer file.Close()
	var data map[string]map[string]string
	var section string
	first := true
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		l = strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("no this value")
				return nil, err
			}
			if len(l) == 0 {
				break
			}
		}
		if (len(l) != 0) && (!system && string(l[0]) == "#" || system && string(l[0]) == ";") {
			return nil, errors.New("System Error")
		}
		if (len(l) == 0) || (system && string(l[0]) == "#") || (!system && string(l[0]) == ";") {
			continue
		} else if l[0] == '[' && l[len(l)-1] == ']' {
			section = strings.TrimSpace(l[1 : len(l)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
			first = false
		} else if first {
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
		} else {
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
	}
	return c.confList, nil
}

//GetValue : 通过section和key来查找一个value
func (c *Config) GetValue(sec string, key string) (string, error) {
	c.readList()
	conf, err := c.readList()
	if err != nil {
		return "", err
	}
	for _, v := range conf {
		for k, val := range v {
			if k == sec {
				return val[key], nil
			}
		}
	}
	return "", errors.New(notFindValue)
}

//SetValue :通过section和key来设置一个value
func (c *Config) SetValue(section, key, value string) bool {
	c.readList()
	data := c.confList
	var ok bool
	var index = make(map[int]bool)
	var conf = make(map[string]map[string]string)
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
	if ok {
		c.confList[i][section][key] = value
		return true
	}
	conf[section] = make(map[string]string)
	conf[section][key] = value
	c.confList = append(c.confList, conf)
	return true
}

//ListenFunc ：实现接口方法 listen 直接调用函数
type ListenFunc func(string)

//Listener  :一个特殊的接口，用来监听配置文件是否被修改，让开发者自己决定如何处理配置变化
type Listener interface {
	listen(inifile string)
}

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

func (f ListenFunc) listen(infile string) {
	for {
		newtime, err := getFileModTime(infile)
		if err != nil {
		}
		if lasttime != newtime {
			lasttime = newtime
			break
		}
	}
}

//Watch :查看更改，若有更改则调用用户函数然后返回最新的更改
func Watch(filename string, listener ListenFunc) (Configuration *Config, err error) {
	listener.listen(filename)
	listener(filename)
	Configuration = SetConfig(filename)
	return Configuration, err
}
