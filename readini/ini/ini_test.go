package ini

import (
	"testing"
)

func TestSetConfig(t *testing.T) {
	SetConfig("init.ini")
}

func BenchmarkSetConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SetConfig("init.ini")
	}
}

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
func TestGetFileModTime(t *testing.T) {
	getFileModTime("init.ini")
}

func BenchmarkGetFileModTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFileModTime("init.ini")
	}
}
func TestListen(t *testing.T) {
	MyListen := func(string) {
	}
	ListenFunc.listen(MyListen, "init.ini")
}

// func BenchmarkListen(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		MyListen := func(string) {
// 		}
// 		ListenFunc.listen(MyListen, "init.ini")
// 	}
// }
func TestWatch(t *testing.T) {
	MyListen := func(string) {
	}
	conf1 := SetConfig("init.ini")
	conf2, _ := Watch("init.ini", MyListen)
	if !equal(conf1, conf2) {
		t.Errorf("Not Equal")
	}
}
