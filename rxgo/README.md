# 修改、改进 RxGo 包
## 题目要求
- 阅读 ReactiveX 文档。请在 pmlpml/RxGo 基础上，

1. 修改、改进它的实现
2. 或添加一组新的操作，如 filtering
- 该库的基本组成：

	* rxgo.go 给出了基础类型、抽象定义、框架实现、Debug工具等

	* generators.go 给出了 sourceOperater 的通用实现和具体函数实现

	* transforms.go 给出了 transOperater 的通用实现和具体函数实现

## 注意事项
- 本次开发使用环境为linux
- api文档使用godoc生成
- 说明文档编辑采用csdn编辑，故图片可能会有水印
- 本次项目参考了老师已给函数的实现形式
- 本次项目参考了原有RxGo和RxJava等包的实现以及部分GitHub和Gitee中的自建包
- 由于老师给的文档里已经有filter函数，故在本项目中不再实现
- 本项目实现如下的八个函数
	* Debounce — only emit an item from an Observable if a particular timespan has passed without it emitting another item
	* Distinct — suppress duplicate items emitted by an Observable
	* ElementAt — emit only item n emitted by an Observable
	* First — emit only the first item, or the first item that meets a condition, from an Observable
	* IgnoreElements — do not emit any items from an Observable but mirror its termination notification
	* Last — emit only the last item emitted by an Observable
	* Sample — emit the most recent item emitted by an Observable within periodic time intervals
	* Skip — suppress the first n items emitted by an Observable
	* SkipLast — suppress the last n items emitted by an Observable
	* Take — emit only the first n items emitted by an Observable
	* TakeLast — emit only the last n items emitted by an Observable
## api文档
api文档见同目录下：api文档.pdf
## 代码分析
### 变量定义

定义10个变量用来表示将要实现的功能。

- debounce用于记录延迟接受的时间
- distinct用于是否跳过重复
- elementAt用于指定输出特定的位置
- ignoreElement用于是否跳过全部
- first用于是否输出第一个
- last用于是否输出最后一个
- sample用于记录定时接受最近发出的时间
- skip用于记录跳过的个数
- take用于记录获得的个数
- takeOrSkip用于记录是take还是skip
```go
	debounce          time.Duration
	distinct          bool
	elementAt         int
	ignoreElement     bool
	first             bool
	last              bool
	sample            time.Duration
	skip              int
	take              int
	takeOrSkip       bool
```
### 转换节点实现
模仿老师已给的`transform.go`中的`transOperater`结构体，直接写出我们需要的`filteringOperator`结构体
```go
type filteringOperator struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}
```

### 初始化
初始化函数完全参考老师已给的`transform.go`中的`newTransformObservable`函数，就是一个简单的初始化过程
```go
func (parent *Observable) newFilteringObservable(name string) (o *Observable) {
	//new Observable
	o = newObservable()
	o.Name = name

	//chain Observables
	parent.next = o
	o.pred = parent
	o.root = parent.root

	//set options
	o.buf_len = BufferLen
	return o
}
```
### filteringTotalOperator 
由于将要实现的函数中都是类似的filter操作，故我们定义一个统一的operator进行操作,然后把不同实现在op函数中。
```go
var filteringTotalOperator = filteringOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
	var params = []reflect.Value{x}
	x = params[0]
	if !end {
		end = o.sendToFlow(ctx, x.Interface(), out)
	}
	return
},
}
```

### Debounce
![debounce](https://img-blog.csdnimg.cn/20201108013521325.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
Debounce用于在一个数据发出后的指定延迟后获取它，如果在该延迟时间内有新的数据，则重新计时。我们对`Observable`中的`debounce`变量传入需要延迟的时间，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量
```go
func (parent *Observable) Debounce(_debounce time.Duration) (o *Observable) {
	o = parent.newFilteringObservable("debounce")
	o.first, o.last, o.ignoreElement, o.distinct = false, false, false, false
	o.debounce, o.take, o.skip = _debounce, 0, 0
	o.operator = filteringTotalOperator
	return o
}
```
### Distinct
![distinct](https://img-blog.csdnimg.cn/2020110823282113.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

Distinct用于去重，即如果有重复的数据我们就不再对它输出。我们对`Observable`中的`distinct `变量置位，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量
```go
func (parent *Observable) Distinct() (o *Observable) {
	o = parent.newFilteringObservable("distinct")
	o.ignoreElement, o.first, o.last, o.distinct = false, false, false, true
	o.debounce, o.take, o.skip = 0, 0, 0
	o.operator = filteringTotalOperator
	return o
}
```
### ElementAt
![elementAt](https://img-blog.csdnimg.cn/20201108232922781.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)


ElementAt用于输出第x个数据，x从0计数。我们对`Observable`中的`elementAt`变量赋值指定的数据位置，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量
```go
func (parent *Observable) ElementAt(index int) (o *Observable) {
	o = parent.newFilteringObservable("elementAt")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, false
	o.debounce, o.skip, o.take, o.elementAt = 0, 0, 0, index
	o.operator = filteringTotalOperator
	return
}
```
### First
![first](https://img-blog.csdnimg.cn/20201108233203631.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
First用于输出第1个数据。我们对`Observable`中的`first`变量置位，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量
```go
func (parent *Observable) First() (o *Observable) {
	o = parent.newFilteringObservable("first")
	o.first, o.last, o.ignoreElement, o.distinct = true, false, false, false
	o.debounce, o.take, o.skip = 0, 0, 0
	o.operator = filteringTotalOperator
	return o

}
```
### IgnoreElement
![ignoreElement](https://img-blog.csdnimg.cn/20201108234325746.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

IgnoreElement用于忽略全部数据。我们对`Observable`中的`ignoreElement `变量置位，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量
```go
func (parent *Observable) IgnoreElement() (o *Observable) {
	o = parent.newFilteringObservable("ignoreElement")
	o.first, o.last, o.distinct, o.ignoreElement = false, false, false, true
	o.debounce, o.take, o.skip = 0, 0, 0
	o.operator = filteringTotalOperator
	return o
}
```
### Last
![Last](https://img-blog.csdnimg.cn/20201109001953356.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)


Last用于获取最后一个数据。我们对`Observable`中的`last`变量置位，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量
```go
func (parent *Observable) Last() (o *Observable) {
	o = parent.newFilteringObservable("last")
	o.first, o.last, o.distinct, o.ignoreElement = false, true, false, false
	o.debounce, o.take, o.skip = 0, 0, 0
	o.operator = filteringTotalOperator
	return o
}

```
### Sample
![Sample](https://img-blog.csdnimg.cn/20201109002128886.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
Sample在指定时间段内获取最新数据。我们对`Observable`中的`sample `变量赋值指定的时间，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量
```go
func (parent *Observable) Sample(_sample time.Duration) (o *Observable) {
	o = parent.newFilteringObservable("sample")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, false
	o.debounce, o.skip, o.take, o.elementAt, o.sample = 0, 0, 0, 0, _sample
	o.operator = filteringTotalOperator
	return o
}
```
### Skip And SkipLast
![skip](https://img-blog.csdnimg.cn/20201109002631933.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
![skipLast](https://img-blog.csdnimg.cn/20201109002649281.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)


Skip在前面跳过指定的个数的数据。我们对`Observable`中的`skip  `变量赋值要跳过的个数，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量。
SkipLast同理是指定从最后开始往前跳过的个数。
```go
func (parent *Observable) Skip(num int) (o *Observable) {
	o = parent.newFilteringObservable("skip")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, false
	o.debounce, o.take, o.skip = 0, 0, num
	o.operator = filteringTotalOperator
	return o
}

func (parent *Observable) SkipLast(num int) (o *Observable) {
	o = parent.newFilteringObservable("skipLast")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, false
	o.debounce, o.take, o.skip = 0, 0, -num
	o.operator = filteringTotalOperator
	return o
}
```
### Take And TakeLast
![take](https://img-blog.csdnimg.cn/20201109003100389.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

![takeLast](https://img-blog.csdnimg.cn/20201109003132675.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
Take在前面取指定的个数的数据。我们对`Observable`中的`take`变量赋值要取的个数，然后将其它我们定义的变量都复位即可。`operator `需要赋值`filteringTotalOperator`也就是刚才定义的变量。
TakeLast同理是指定从最后开始往前要取的个数。
```go
func (parent *Observable) Take(num int) (o *Observable) {
	o = parent.newFilteringObservable("Take")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, true
	o.debounce, o.skip, o.take = 0, 0, num
	o.operator = filteringTotalOperator
	return o
}

func (parent *Observable) TakeLast(num int) (o *Observable) {
	o = parent.newFilteringObservable("takeLast")
	o.first, o.last, o.distinct, o.ignoreElement, o.takeOrSkip = false, false, false, false, true
	o.debounce, o.skip, o.take = 0, 0, -num
	o.operator = filteringTotalOperator
	return o
}
```
### 操作接口op
首先模仿老师已给的文件进行变量定义，这些变量都是用于输入输出

```go
	in := o.pred.outflow
	out := o.outflow
	var _out []interface{}
	var wg sync.WaitGroup
```
然后我们需要调用一个函数，首先定义临时变量，变量用于计时和记录结束，flag用于记录是否存在（实现distinct）

```go
		end := false
		flag := make(map[interface{}]bool)
		timeStart := time.Now()
		timeSample := time.Now()
```
对于输入，我们进行一次遍历，在遍历中首先确定如果结束或跳过全部或时间不符（sample和debounce操作）

```go
			if end {
				continue
			}
			if o.ignoreElement {
				continue
			}
			if o.sample > 0 && timeSampleFromStart < o.sample {
				continue
			}
			if o.debounce > time.Duration(0) && timeFromStart < o.debounce {
				continue
			}
```
与其它op类似，我们在这也要进行错误匹配，并将其输出到字节流

```go
			xv := reflect.ValueOf(x)
			// send an error to stream if the flip not accept error
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
```
下面对剩下的操作进行判断
```go
			if o.elementAt > 0 || o.take != 0 || o.skip != 0 || o.last || (o.distinct && flag[xv.Interface()]) {
				continue
			}
```
为了实现`distinct`，需要将当前的元素的`flag`标记为`true`

```go
			flag[xv.Interface()] = true
```
下面是一个线程管理，除了需要在`ThreadingDefault`中处理一下`sample`的时间计数之外，其余可以完全参考老师`transform`功能包中的`op`函数的处理

```go
			switch threading := o.threading; threading {
			case ThreadingDefault:
				if o.sample > 0 {
					timeSample = timeSample.Add(o.sample)
				}
				if fop.opFunc(ctx, o, xv, out) {
					end = true
				}
			case ThreadingIO:
				fallthrough
			case ThreadingComputing:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if fop.opFunc(ctx, o, xv, out) {
						end = true
					}
				}()
```
下面进行`first`和`last`的处理。`first`处理直接`break`即可，因为第一次运行到这里的就是第一个数据，我们需要的数据已经拿到，就不再需要往下进行遍历了。last则是处理当不是第一个的时候直接用`fop.opFunc(ctx, o, xv, out)`处理

```go
			if o.first {
				break
			}
		}
		if o.last && len(_out) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				xv := reflect.ValueOf(_out[len(_out)-1])
				fop.opFunc(ctx, o, xv, out)
			}()
		}
```
下面处理take和skip，首先利用`take`和`skip`不为0和`takeOrSkip`的值共同确定是什么操作。首先获得要处理的步长step

```go
			var step int
			if o.takeOrSkip {
				step = o.take
			} else {
				step = o.skip
			}
```
如果是take或skiplast操作，本质是一样的，都是忽略了后面的几个，故在一起处理

```go
				if (o.takeOrSkip && step > 0) || (!o.takeOrSkip && step < 0) {
					if !o.takeOrSkip {
						step = len(_out) + step
					}
					if step >= len(_out) || step <= 0 {
						newIn, err = nil, errors.New("OutOfBound")
					} else {
						newIn, err = _out[:step], nil
					}
```
如果是takelast或skip操作，本质是一样的，都是忽略了前面的几个，故在一起处理

```go
			if (o.takeOrSkip && step < 0) || (!o.takeOrSkip && step > 0) {
				if o.takeOrSkip {
					step = len(_out) + step
				}
				if step >= len(_out) || step <= 0 {
					newIn, err = nil, errors.New("OutOfBound")
				}
```
最后写入即可

```go
				if err != nil {
					o.sendToFlow(ctx, err, out)
				} else {
					xv := newIn
					for _, val := range xv {
						fop.opFunc(ctx, o, reflect.ValueOf(val), out)
					}
				}
```
最后处理取特定位置的`elementAt`，这里只需要取`_out[o.elementAt-1])`位即可

```go
		if o.elementAt != 0 {
			if o.elementAt < 0 || o.elementAt > len(_out) {
				o.sendToFlow(ctx, errors.New("OutOfBound"), out)
			} else {
				xv := reflect.ValueOf(_out[o.elementAt-1])
				fop.opFunc(ctx, o, xv, out)
			}
		}
```
以上就是op的大概过程（部分不太重要的细节因篇幅问题没有展示）
## 测试部分
### 单元测试
单元测试部分针对每个函数进行了测试，我们仿照老师给出的测试格式，来编写自己的测试文件。
#### TestDebounce
使用100s进行测试，由于时间很长，所以返回的应该是空。
```go
func TestDebounce(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).Debounce(100 * time.Millisecond)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{}, res, "Debounce Test Error!")
}
```
测试结果：

![TestDebounce](https://img-blog.csdnimg.cn/20201109195248225.png)
#### TestDistinct
用一串有重复的数字对其进行测试，看其是否能达到去重的效果
```go
func TestDistinct(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5, 6).Map(func(x int) int {
		return x
	}).Distinct()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6}, res, "Distinct Test Error!")
}
```
测试结果：

![TestDistinct](https://img-blog.csdnimg.cn/20201109195532651.png)
#### TestElementAt
用返回第五个数即4进行测试`ElementAt`是否能争取取值
```go
func TestElementAt(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).ElementAt(5)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{4}, res, "ElementAt Test Error!")
}
```
测试结果：

![TestElementAt](https://img-blog.csdnimg.cn/20201109195733919.png)
#### TestIgnoreElement
看是否返回空来测试该函数的正确性
```go
func TestIgnoreElement(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).IgnoreElement()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{}, res, "IgnoreElement Test Error!")
}
```
测试结果：

![TestIgnoreElement](https://img-blog.csdnimg.cn/20201109195847482.png)
#### TestFirst
输入一组数据看返回的是不是第一个数据
```go
func TestFirst(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).First()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{0}, res, "First Test Error!")
}
```
测试结果：



![TestFirst](https://img-blog.csdnimg.cn/20201109201629872.png)

#### TestLast
与First同理，输入一组数据看返回的是不是最后一个数据

```go
func TestLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).Last()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{5}, res, "Last Test Error!")
}
```
测试结果：

![TestLast](https://img-blog.csdnimg.cn/20201109201758189.png)
#### TestSample
用很大的时间段来对Sample测试，看其是否返回空
```go
func TestSample(t *testing.T) {
	res := []int{}
	rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		time.Sleep(2 * time.Millisecond)
		return x
	}).Sample(20* time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{}, res, "SkipLast Test Error!")

}
```
测试结果：

![TestSample](https://img-blog.csdnimg.cn/20201109202152119.png)
#### TestSkip
用一组数据跳过前两个来测试结果是否正确
```go
func TestSkip(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).Skip(2)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{2, 3, 4, 5}, res, "Skip Test Error!")
}
```
测试结果：

![TestSkip](https://img-blog.csdnimg.cn/20201109202341298.png)
#### TestSkipLast
用一组数据跳过最后三个来测试结果是否正确
```go
func TestSkipLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).SkipLast(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{0, 1, 2}, res, "SkipLast Test Error!")
}
```

测试结果：

![TestSkipLast](https://img-blog.csdnimg.cn/20201109202435399.png)
#### TestTake
用一组数据取出前两个来测试结果是否正确
```go
func TestTake(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).Take(2)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{0, 1}, res, "Take Test Error!")
}
```

测试结果：

![TestTake](https://img-blog.csdnimg.cn/20201109202601286.png)
#### TestTakeLast

用一组数据取出最后三个来测试结果是否正确

```go
func TestTakeLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5).Map(func(x int) int {
		return x
	}).TakeLast(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{3, 4, 5}, res, "TakeLast Test Error!")
}
```
测试结果：

![TestTakeLast](https://img-blog.csdnimg.cn/20201109202700170.png#pic_center)
#### 文件测试结果
执行整个测试文件，结果如下：

![filetest](https://img-blog.csdnimg.cn/2020110920281173.png#pic_center)
#### 包测试结果
执行整个测试文件包，结果如下：
![packagetest](https://img-blog.csdnimg.cn/2020110920291132.png#pic_center)

可以看出达到了85%的代码覆盖率，这是因为老师已给文件中一些异常处理以及自己编写的代码中的异常处理部分无法完全覆盖，其余部分都已经覆盖测试。

### 功能测试（使用案例）
使用一个简单的main函数对所有的功能进行测试，代码如下：

```go
package main

import (
	"fmt"
	"time"

	rxgo "github.com/yilin0041/service-computing/rxgo"
)

func main() {
	fmt.Println("测试数据：0,1,2,3,4,5,3,4,5")
	res := []int{}
	ob := rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Debounce(999999)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Debounce(999999): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Distinct()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Distinct(): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).ElementAt(5)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("ElementAt(5): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).First()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("First(): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Last()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Last(): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		time.Sleep(4 * time.Millisecond)
		return x
	}).Sample(40 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Sample(40): (sleep (4))")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Skip(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Skip(4): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).SkipLast(2)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("SkipLast(2): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")

	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).Take(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Take(4) :")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
	res = []int{}
	ob = rxgo.Just(0, 1, 2, 3, 4, 5, 3, 4, 5).Map(func(x int) int {
		return x
	}).TakeLast(2)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("TakeLast(2): ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
}

```
测试结果如下：

![maintest](https://img-blog.csdnimg.cn/20201109214046998.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

经验证，符合要求，测试成功。

## 总结
本次实验是一个包完善的实验，通过老师已经写好的包，来写一个新的功能文件。其实参考老师已经给出的文件，我们就可以知道写函数的大体格式，然后参考网上已有的代码和要实现函数的逻辑就可以具体进行实现了。
通过本次实验，让我更好的理解了线程的使用以及并发的设计思维与模式，会让以后的编程工作有更多的选择