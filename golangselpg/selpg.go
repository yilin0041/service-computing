package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
)

type selpgArgs struct {
	start    int
	end      int
	inFile   string
	outDest  string
	pageLen  int
	pageType bool
}

//IntMax :max int
const IntMax = int(^uint(0) >> 1)

//IntMin :min int
const IntMin = ^IntMax

func getArgs(args *selpgArgs) {
	pflag.IntVarP(&(args.start), "start", "s", IntMin, "Define startPage")
	pflag.IntVarP(&(args.end), "end", "e", IntMin, "Define endPage")
	pflag.IntVarP(&(args.pageLen), "pageLen", "l", 72, "Define pageLength")
	pflag.StringVarP(&(args.outDest), "outDest", "d", "", "Define printDest")
	pflag.BoolVarP(&(args.pageType), "pageType", "f", false, "Define pageType")
	pflag.Parse()
	argLeft := pflag.Args()
	if len(argLeft) > 0 {
		args.inFile = string(argLeft[0])
		args.inFile = "src/github.com/yilin0041/service-computing/golangselpg/" + args.inFile
	} else {
		args.inFile = ""
	}
}

//Usage :the help for user
func Usage() {
	fmt.Fprintf(os.Stderr, "USAGE: golangselpg -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n")
}

func checkArgs(args selpgArgs) {
	if args.start == IntMin {
		fmt.Fprintf(os.Stderr, "[Error] The start page can't be empty! \n")
		os.Exit(2)
		Usage()
	} else if args.end == IntMin {
		fmt.Fprintf(os.Stderr, "[Error] The end page can't be empty! \n")
		os.Exit(3)
		Usage()
	} else if (args.start < 0) || (args.end < 0) {
		fmt.Fprintf(os.Stderr, "[Error] The page number can't be negative!\n")
		Usage()
		os.Exit(4)
	} else if args.start > args.end {
		fmt.Fprintf(os.Stderr, "[Error] The start page  can't be bigger than end page!\n")
		Usage()
		os.Exit(5)
	} else if args.pageLen <= 0 {
		fmt.Fprintf(os.Stderr, "[Error] The page  length can't be less than 1! \n")
		Usage()
		os.Exit(6)
	}
}

func output(infile *os.File, args *selpgArgs) {
	var out io.WriteCloser
	var cmd *exec.Cmd
	var err error
	if len(args.outDest) > 0 {
		cmd = exec.Command("lp", "-d"+args.outDest)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		out, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR]StdinPipe  Error!\n")
		}
	}
	lineCount := 0
	pageCount := 1
	buf := bufio.NewReader(infile)
	for {
		var lineData string
		err = nil
		if args.pageType {
			lineData, err = buf.ReadString('\f')
			pageCount++
		} else {
			lineData, err = buf.ReadString('\n')
			lineCount++
			if lineCount > args.pageLen {
				pageCount++
				lineCount = 1
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR]Read Input File Error!\n")
		}
		err = nil
		if (pageCount >= args.start) && (pageCount <= args.end) {
			if len(args.outDest) == 0 {
				_, err = fmt.Fprintf(os.Stdout, "%s", lineData)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[ERROR]Std Output Error!\n")
				}
			} else {
				_, err = out.Write([]byte(lineData))
				if err != nil {
					fmt.Fprintf(os.Stderr, "[ERROR]Pipe Output Error!\n")
				}
			}
		}
	}
	if len(args.outDest) > 0 {
		out.Close()
		err = nil
		err = cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR]Cmd Run Error!\n")
		}
	}
}

func excute(args *selpgArgs) {
	var infile *os.File
	if args.inFile == "" {
		infile = os.Stdin
	} else {
		var err error
		infile, err = os.Open(args.inFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR]Open Input File!\n")
		}
	}
	output(infile, args)
}

func main() {
	var args selpgArgs
	getArgs(&args)
	checkArgs(args)
	excute(&args)
}
