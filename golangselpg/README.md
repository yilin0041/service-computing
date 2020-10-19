# CLI 命令行实用程序开发基础
## 概述
CLI（Command Line Interface）实用程序是Linux下应用开发的基础。正确的编写命令行程序让应用与操作系统融为一体，通过shell或script使得应用获得最大的灵活性与开发效率。例如：

- Linux提供了cat、ls、copy等命令与操作系统交互；
- go语言提供一组实用程序完成从编码、编译、库管理、产品发布全过程支持；
- 容器服务如docker、k8s提供了大量实用程序支撑云服务的开发、部署、监控、访问等管理任务；
- git、npm等也是大家比较熟悉的工具。

尽管操作系统与应用系统服务可视化、图形化，但在开发领域，CLI在编程、调试、运维、管理中提供了图形化程序不可替代的灵活性与效率。

## 开发题目
使用 golang 开发 [开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html) 中的 selpg

## 开发要求
- 请按文档 使用 selpg 章节要求测试你的程序
- 请使用 pflag 替代 goflag 以满足 Unix 命令行规范， 参考：[Golang之使用Flag和Pflag](https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/)
- golang 文件读写、读环境变量，请自己查 os 包
- “-dXXX” 实现，请自己查 `os/exec` 库，例如案例 [Command](https://godoc.org/os/exec#example-Command)，管理子进程的标准输入和输出通常使用 io.Pipe，具体案例见 [Pipe](https://godoc.org/io#Pipe)
- 请自带测试程序，确保函数等功能正确

## 代码分析
### 参数设置
通过设置一个struct结构来进行参数的设置，包括了开始页，结束页，输入文件，输出目的地，页的行数以及类型（-f或-lNumber）

```go
type selpgArgs struct {
	start    int
	end      int
	inFile   string
	outDest  string
	pageLen  int
	pageType bool
}
```
### 获得参数
我们定义了`selpgArgs`的参数结构体，现在我们需要使用`pflag`来获取参数。pflag的使用可以参考：[Golang之使用Flag和Pflag](https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/)。
这里选择一个解释，如`pflag.IntVarP(&(args.start), "start", "s", IntMin, "Define startPage")`这里表示将参数写入`args.start`，注意使用取地址符，而因为获取的是一个整型参数，所以我们使用的函数是`IntVarP`；`start`为名字；`s`表示在`-s`之后读取；`IntMin`表示如果没有`-s`初始化的值，这个数字的定义后续我们会提到；最后`Define startPage`是`usage`。
更加值得注意的是，我们在定义完所有的`pflag`后，需要使用	`pflag.Parse()`将所有的`pflag`启动。
最后，我们在来说明一个定义：

```go
//IntMax :max int
const IntMax = int(^uint(0) >> 1)

//IntMin :min int
const IntMin = ^IntMax
```
这里定义了一个最大数和一个最小数用于初始化。
下面，给出获得参数的详细代码

```go
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
```
### 帮助
这里使用`Usage`函数来给出用户帮助
```go
//Usage :the help for user
func Usage() {
	fmt.Fprintf(os.Stderr, "USAGE: golangselpg -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n")
}
```
### 检查参数
我们输入完参数后，需要检查参数的合法性。首先检查页是否被赋值，如果没赋值，则返回错误信息；然后检查页是否为正数，若为负数则返回错误信息；最后检查每页行数，若小于零，则返回错误信息。返回错误信息后，自动返回帮助来辅助用户检查自己的命令。

```go
func checkArgs(args selpgArgs) {
	if args.start == IntMin {
		fmt.Fprintf(os.Stderr, "[Error] The start page can't be empty! \n")
		os.Exit(3)
		Usage()
	} else if args.end == IntMin {
		fmt.Fprintf(os.Stderr, "[Error] The end page can't be empty! \n")
		os.Exit(4)
		Usage()
	} else if (args.start < 0) || (args.end < 0) {
		fmt.Fprintf(os.Stderr, "[Error] The page number can't be negative!\n")
		Usage()
		os.Exit(5)
	} else if args.start > args.end {
		fmt.Fprintf(os.Stderr, "[Error] The start page  can't be bigger than end page!\n")
		Usage()
		os.Exit(6)
	} else if args.pageLen <= 0 {
		fmt.Fprintf(os.Stderr, "[Error] The page  length can't be less than 1! \n")
		Usage()
		os.Exit(7)
	}
}
```
### 输出
输出函数是用于输出到对应管程或文件或命令行，这里传入的参数是我们获得的参数结构体和读取的文件。
然后使用`exec.Command("lp", "-d"+args.outDest)`执行管道，将命令行的输入管道`cmd.StdinPipe()`获取的指针赋值给`out`

```go
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
```
利用for循环一直读取到结束为止，这里进行判断：

```go
		if err == io.EOF {
			break
		}
```
由此来退出循环。
新建一个缓冲区`buf := bufio.NewReader(infile)`来读取文件
下面来具体讲解循环。
首先进行判断是以"\f"结尾还是"\n"，不同的类型按照不同的方式读取
```go
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
```
如果读取的页数` (pageCount >= args.start) && (pageCount <= args.end)`，那么需要进行输出。输出的时候需要注意方式，如果没有`outDest`，则` fmt.Fprintf(os.Stdout, "%s", lineData)`直接进行输出，但如果有指定的管道，则`out.Write([]byte(lineData))`进行输出到管道。

在循环结束后，进行判断，如果有`outDest`，则需要`out.Close()`（**重要！如果不关闭则无法输出**），最后`cmd.Run()`即可

### 执行
执行只需要对输入文件进行处理，如果有输入文件则打开该文件，否则从命令行获取。获取完成后直接执行`output`即可

```go
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
```
## 测试
### 单元或集成测试
我们首先应该编写测试文件
#### 测试Usage
对Usage进行函数测试和性能测试
```go
func TestUsage(t *testing.T) {
	Usage()
}
func BenchmarkUsage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Usage()
	}
}
```
运行结果如下：

![TestUsage](https://img-blog.csdnimg.cn/20201011234048398.png)

![BenchmarkUsage](https://img-blog.csdnimg.cn/20201011233855131.png)

#### 测试GetArgs
对GetArgs进行函数测试，这里不能性能测试，因为pflags需要使用命令行传参，所以不进行性能测试。

```go
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
```
运行结果如下：

![TestGetArgs](https://img-blog.csdnimg.cn/20201011234238986.png)
#### 测试CheckArgs
对CheckArgs进行函数测试和性能测试

```go
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
```
测试结果如下：

![TestCheckArgs](https://img-blog.csdnimg.cn/20201011235028273.png)

![BenchmarkCheckArgs](https://img-blog.csdnimg.cn/20201011235106373.png)
#### 测试Output
对Output进行函数测试和性能测试
##### 正常输出
```go
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
```
测试结果如下：

![TestOutput](https://img-blog.csdnimg.cn/2020101123530823.png)

![BenchmarkOutput](https://img-blog.csdnimg.cn/20201011235354516.png)
##### 打印机输出
这里使用了虚拟打印机Cups-PDF进行打印，在我的虚拟机中，"PDF"为一个虚拟打印机
```go
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
```
测试结果如下：

![TestOutputWithDestination](https://img-blog.csdnimg.cn/20201011235703695.png)

![BenchmarkOutputWithDestination](https://img-blog.csdnimg.cn/2020101123580635.png)
#### 测试Excute
对Excute进行函数测试和性能测试

```go
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
```
测试结果如下：

![TestExcute](https://img-blog.csdnimg.cn/20201011235952983.png)

![BenchmarkExcute](https://img-blog.csdnimg.cn/20201012000033553.png)

### 功能测试
测试文件如下所示

![testfile](https://img-blog.csdnimg.cn/20201012001348673.png)
#### 1. `$ selpg -s1 -e1 input_file`
该命令将把“input_file”的第 1 页写至标准输出（也就是屏幕），因为这里没有重定向或管道。
测试结果如下：

![test1](https://img-blog.csdnimg.cn/20201012001428953.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
#### 2. `$ selpg -s1 -e1 < input_file`
该命令与示例 1 所做的工作相同，但在本例中，selpg 读取标准输入，而标准输入已被 shell／内核重定向为来自“input_file”而不是显式命名的文件名参数。输入的第 1 页被写至屏幕。
结果如下所示：

![test2](https://img-blog.csdnimg.cn/2020101200161950.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
#### 3. `$ other_command | selpg -s1 -e2`
“other_command”的标准输出被 shell／内核重定向至 selpg 的标准输入。将第 1页写至 selpg 的标准输出（屏幕）。（为了测试方便进行了修改）
结果如下所示：

![test3](https://img-blog.csdnimg.cn/20201012002245775.png)
#### 4. `$ selpg -s10 -e20 -l1 input_file >output_file`
selpg 将第 10 页到第 20 页写至标准输出；标准输出被 shell／内核重定向至“output_file”。（为了测试方便增加了-l1属性）
结果如下所示：

![test4](https://img-blog.csdnimg.cn/20201012002637914.png)
![test4](https://img-blog.csdnimg.cn/20201012002654527.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
#### 5. `selpg -s10 -e20 input_file 2>error_file`
selpg 将第 10 页到第 20 页写至标准输出（屏幕）；所有的错误消息被 shell／内核重定向至“error_file”。请注意：在“2”和“>”之间不能有空格；这是 shell 语法的一部分（请参阅“man bash”或“man sh”）。
结果如下所示：

![test5](https://img-blog.csdnimg.cn/20201012003035335.png)
![test5](https://img-blog.csdnimg.cn/20201012003051963.png)
#### 6. `$ selpg -s10 -e20 -l1 input_file >output_file 2>error_file`
selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至“error_file”。当“input_file”很大时可使用这种调用；您不会想坐在那里等着 selpg 完成工作，并且您希望对输出和错误都进行保存。（为了测试方便增加了-l1属性）
结果如下所示：

![test6](https://img-blog.csdnimg.cn/20201012003358177.png)
![test6](https://img-blog.csdnimg.cn/20201012003411593.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
#### 7. `$ selpg -s10 -e20 -l1 input_file >output_file 2>/dev/null`
selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至 /dev/null（空设备），这意味着错误消息被丢弃了。设备文件 /dev/null 废弃所有写至它的输出，当从该设备文件读取时，会立即返回 EOF。（为了测试方便增加了-l1属性）
测试了两种情况，分别为有错和无错
结果如下所示：

![test7](https://img-blog.csdnimg.cn/20201012003735887.png)

![test7](https://img-blog.csdnimg.cn/20201012003656744.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
![test7](https://img-blog.csdnimg.cn/20201012003720631.png)

#### 8. `$ selpg -s10 -e20 -l1 input_file >/dev/null`
selpg 将第 10 页到第 20 页写至标准输出，标准输出被丢弃；错误消息在屏幕出现。这可作为测试 selpg 的用途，此时您也许只想（对一些测试情况）检查错误消息，而不想看到正常输出。（为了测试方便增加了-l1属性）
结果如下所示：

![test8](https://img-blog.csdnimg.cn/20201012003857471.png)
![test8](https://img-blog.csdnimg.cn/20201012003914375.png)


#### 9. `$ selpg -s10 -e20 -l1 input_file | other_command`
selpg 的标准输出透明地被 shell／内核重定向，成为“other_command”的标准输入，第 10 页到第 20 页被写至该标准输入。“other_command”的示例可以是 lp，它使输出在系统缺省打印机上打印。“other_command”的示例也可以 wc，它会显示选定范围的页中包含的行数、字数和字符数。“other_command”可以是任何其它能从其标准输入读取的命令。错误消息仍在屏幕显示。（为了测试方便增加了-l1属性）
结果如下所示：

![test9](https://img-blog.csdnimg.cn/20201012004207680.png)
![test9](https://img-blog.csdnimg.cn/20201012004239868.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

#### 10. `$ selpg -s10 -e20 -l1 input_file 2>error_file | other_command`
与上面的示例 9 相似，只有一点不同：错误消息被写至“error_file”。（为了测试方便增加了-l1属性）
结果如下所示：

这里没有输入`-s10`以测试错误信息

![test10](https://img-blog.csdnimg.cn/20201012004415289.png)
![test10](https://img-blog.csdnimg.cn/20201012004515132.png)
#### 11. `$ selpg -s10 -e20 -dlp1 input_file`
第 10 页到第 20 页由管道输送至命令“lp -dlp1”，该命令将使输出在打印机 lp1 上打印。（为了测试方便增加了-l1属性）,同时，使用虚拟打印机PDF
测试结果如下：

![test11](https://img-blog.csdnimg.cn/20201012004841984.png)

![test11](https://img-blog.csdnimg.cn/20201012004909640.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
#### 12. `$ selpg -s10 -e20 input_file > output_file 2>error_file &`
该命令利用了 Linux 的一个强大特性，即：在“后台”运行进程的能力。在这个例子中发生的情况是：“进程标识”（pid）如 1234 将被显示，然后 shell 提示符几乎立刻会出现，使得您能向 shell 输入更多命令。同时，selpg 进程在后台运行，并且标准输出和标准错误都被重定向至文件。这样做的好处是您可以在 selpg 运行时继续做其它工作。
测试结果如下：

![test12](https://img-blog.csdnimg.cn/20201012005037639.png)

### 总结
本次实验是go语言的Linux命令行开发任务，即实现一个selpg程序用来读取文件页到指定位置。实验的难点在于将内容通过管道输出到打印机，需要阅读大量的代码和源文件才能解决。这里，我使用了selpg的c语言代码进行修改和优化来完成go语言版本，虽然在编程过程中可能会遇到一些困难，但是只要仔细阅读老师给出的参考资料，问题还是很容易解决的。
