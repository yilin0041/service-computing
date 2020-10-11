package main

import (
	"os"
	"testing"
)

func TestUsage(t *testing.T) {
	Usage()
}
func BenchmarkUsage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Usage()
	}
}
func TestGetArgs(t *testing.T) {
	var args selpgArgs
	getArgs(&args)
	if args.start != IntMin {
		t.Errorf("start init error")
	}
	if args.end != IntMin {
		t.Errorf("end init error")
	}
	if args.pageLen != 72 {
		t.Errorf("pageLen init error")
	}
	if args.outDest != "" {
		t.Errorf("outDest init error")
	}
	if args.pageType != false {
		t.Errorf("pageType init error")
	}
	if args.inFile != "" {
		t.Errorf("inFile init error")
	}
}

func TestCheckArgs(t *testing.T) {
	var args selpgArgs
	args.start = 1
	args.end = 1
	args.pageLen = 1
	checkArgs(args)
}
func BenchmarkCheckArgs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var args selpgArgs
		args.start = 1
		args.end = 1
		args.pageLen = 1
		checkArgs(args)
	}
}
func TestOutput(t *testing.T) {
	var args selpgArgs
	args.start = 1
	args.end = 1
	args.pageLen = 1
	var infile *os.File
	infile, _ = os.Open("test.txt")
	output(infile, &args)
}
func BenchmarkOutput(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var args selpgArgs
		args.start = 1
		args.end = 1
		args.pageLen = 12
		var infile *os.File
		infile, _ = os.Open("test.txt")
		output(infile, &args)
	}
}
func TestOutputWithDestination(t *testing.T) {
	var args selpgArgs
	args.start = 1
	args.end = 1
	args.pageLen = 1
	args.outDest = "PDF"
	var infile *os.File
	infile, _ = os.Open("test.txt")
	output(infile, &args)
}

func BenchmarkOutputWithDestination(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var args selpgArgs
		args.start = 1
		args.end = 1
		args.pageLen = 12
		args.outDest = "PDF"
		var infile *os.File
		infile, _ = os.Open("test.txt")
		output(infile, &args)
	}
}
func TestExcute(t *testing.T) {
	var args selpgArgs
	args.start = 1
	args.end = 1
	args.pageLen = 1
	args.inFile = "test.txt"
	excute(&args)
}

func BenchmarkExcute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var args selpgArgs
		args.start = 1
		args.end = 1
		args.pageLen = 1
		args.inFile = "test.txt"
		excute(&args)
	}
}
