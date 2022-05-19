**视觉智慧实践课上需要写深度学习-图像检索的大作业。Python是机器/深度学习御用开发语言，Golang是新时代后端开发语言。Python很适合算法写模型，而Golang很适合提供API服务，两位同志都红的发紫。（借鉴）**

**所以算法方面用python写的sift，qt不想写（不喜欢），最近又在学go，所以准备用gin搭建web端，不过go和python之间好像有什么联系，没错，就是go-python，所以开始吧。**

## go-python

- [go-python](github.com/sbinet/go-python)主要用在python2.x
- [go-python3](github.com/DataDog/go-python3)主要用来在python3.x

```bash
go get -u github.com/DataDog/go-python3
```

我的mac是python3.9版本的，所以在`go build`的时候报错：

    # github.com/DataDog/go-python3
    ../../go/pkg/mod/github.com/!data!dog/go-python3@v0.0.0-20211102160307-40adc605f1fe/dict.go:141:13: could not determine kind of name for C.PyDict_ClearFreeList

<p style="color: red"> 意思就是python3.9的这个函数PyDict_ClearFreeList被删除了，找不到 </p>

通过查看这个[issue](https://github.com/DataDog/go-python3/issues/38)可以发现：

- python的版本需要用3.7，所以需要配置python3.7虚拟环境
- 下好之后，需要用到pkg-config

解决方案：

下载[python3.7](python.org/downloads/macos/)，然后用anaconda或者venv创建虚拟环境。

我通过`brew install pkg-config`，然后`export $PKG_CONFIG_PATH=/usr/local/bin/pkg-config`，还是不太行。
原来需要用到python3.7里的pkg-config。

直接执行下面:
```bash
sudo vim ~/.bash_profile
export $PKG_CONFIG_PATH=/Library/Frameworks/Python.framework/Versions/3.7/lib/pkgconfig
source ~/.bash_profile
```

最后`go build`可以顺利完成。

## go调用python的过程

1. 初始化python环境
2. 引入模块py对象
3. 使用该模块的变量与函数
4. 解析结果
5. 销毁python3运行环境

## 调用中必用到到几个函数
`init func()`
```go
func init() {
	python3.Py_Initialize()
	if !python3.Py_IsInitialized() {
		fmt.Printf("Error initializing the python interpreter")
		os.Exit(1)
	}
}
```

`ImportModule func(dir, name strng)  *python3.PyObject`
```go
// ImportModule
// @Description: 导入一个包
// @param dir
// @param name
// @return *python3.PyObject
func ImportModule(dir, name string) *python3.PyObject {
	sysModule := python3.PyImport_ImportModule("sys")
	path := sysModule.GetAttrString("path")
	pathStr, _ := pythonRepr(path)
	log.Println("before add path is " + pathStr)
	python3.PyList_Insert(path, 0, python3.PyUnicode_FromString(""))
	python3.PyList_Insert(path, 0, python3.PyUnicode_FromString(dir))
	pathStr, _ = pythonRepr(path)
	log.Println("after add path is " + pathStr)
	return python3.PyImport_ImportModule(name)
}
```

`pythonRepr func(o *python3.PyObject) (string, error)`
```go
// pythonRepr
// @Description: PyObject转换为string
// @param o
// @return string
// @return error
func pythonRepr(o *python3.PyObject) (string, error) {
	if o == nil {
		return "", fmt.Errorf("object is nil")
	}
	s := o.Repr()
	if s == nil {
		python3.PyErr_Clear()
		return "", fmt.Errorf("failed to call Repr object method")
	}
	defer s.DecRef()

	return python3.PyUnicode_AsUTF8(s), nil
}
```


## 协程并发

协程中多次调用Python代码会panic。

错误大概这样：
```bash
[FUNC] hi = &python3.PyObject{ob_refcnt:2, ob_type:(*python3._Ctype_struct__typeobject)(0x4390fa0)}
fatal error: unexpected signal during runtime execution
fatal error: unexpected signal during runtime execution
[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x421f7e6]

runtime stack:
runtime.throw(0x40e1cd5, 0x2a)
        /usr/local/go/src/runtime/panic.go:1116 +0x72
runtime.sigpanic()
        /usr/local/go/src/runtime/signal_unix.go:704 +0x48c

goroutine 20 [syscall]:
runtime.cgocall(0x40b0170, 0xc000038ee0, 0x0)
        /usr/local/go/src/runtime/cgocall.go:133 +0x5b fp=0xc000038eb0 sp=0xc000038e78 pc=0x40048db
github.com/DataDog/go-python3._Cfunc_PyObject_Call(0xbc14a70, 0xd1772d0, 0x4398138, 0x0)
        _cgo_gotypes.go:4351 +0x4e fp=0xc000038ee0 sp=0xc000038eb0 pc=0x40aaaae
github.com/DataDog/go-python3.(*PyObject).Call.func1(0xbc14a70, 0xd1772d0, 0x4398138, 0x0)
        /Users/xiangxianzhang/go/pkg/mod/github.com/!data!dog/go-python3@v0.0.0-20211102160307-40adc605f1fe/object.go:160 +0xab fp=0xc000038f10 sp=0xc000038ee0 pc=0x40abacb
github.com/DataDog/go-python3.(*PyObject).Call(0xbc14a70, 0xd1772d0, 0x4398138, 0x0)
        /Users/xiangxianzhang/go/pkg/mod/github.com/!data!dog/go-python3@v0.0.0-20211102160307-40adc605f1fe/object.go:160 +0x3f fp=0xc000038f40 sp=0xc000038f10 pc=0x40ab49f
main.main.func1(0xc0000aa030, 0xbc14a70, 0xd1772d0)
        /Users/xiangxianzhang/go/src/awesomeProject/main.go:33 +0x76 fp=0xc000038fc8 sp=0xc000038f40 pc=0x40ad736
runtime.goexit()
        /usr/local/go/src/runtime/asm_amd64.s:1374 +0x1 fp=0xc000038fd0 sp=0xc000038fc8 pc=0x4064921
created by main.main
        /Users/xiangxianzhang/go/src/awesomeProject/main.go:31 +0x205

goroutine 1 [semacquire]:
sync.runtime_Semacquire(0xc0000aa038)
        /usr/local/go/src/runtime/sema.go:56 +0x45
sync.(*WaitGroup).Wait(0xc0000aa030)
        /usr/local/go/src/sync/waitgroup.go:130 +0x65
main.main()
        /Users/xiangxianzhang/go/src/awesomeProject/main.go:37 +0x225
```

可以看这个[issue](https://github.com/go-python/cpy3/issues/3)

这是由于直接使用go来调用python函数，不会产生GIL锁，这可能会导致竞争条件，从而导致致命的运行时错误，并且很可能出现分段错误导致整个 Go 应用程序崩溃。

解决方案：
```go
var wg sync.WaitGroup
state := python3.PyEval_SaveThread()

go func() {
    // 获取已上传的图片数量
    r.GET("/api/tot_num", func(c *gin.Context) {
        defer wg.Done()
        _state := python3.PyGILState_Ensure()
        defer python3.PyGILState_Release(_state)
        res := getTotCntFunc.Call(python3.Py_None, python3.Py_None)
        resJson, _ := pythonRepr(res)
        // fmt.Printf("[VARS] TotCntJson = %s\n", resJson)
        Res := getRes(resJson)
        c.JSON(http.StatusOK, gin.H{
        "code": Res.Code,
        "msg":  Res.Msg,
        "data": Res.Data,
        })
    })
}()

wg.Wait()
python3.PyEval_RestoreThread(state)
```