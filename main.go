package main

import (
	"fmt"
	"github.com/DataDog/go-python3"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"sync"
)

func init() {
	python3.Py_Initialize()
	if !python3.Py_IsInitialized() {
		fmt.Printf("Error initializing the python interpreter")
		os.Exit(1)
	}
}

// InsertBeforeSysPath
// @Description: 添加site-packages路径即包的查找路径
// @param p
func InsertBeforeSysPath(p string) {
	sysModule := python3.PyImport_ImportModule("sys")
	path := sysModule.GetAttrString("path")
	python3.PyList_Append(path, python3.PyUnicode_FromString(p))
}

// ImportModule
// @Description: 导入一个包
// @param dir
// @param name
// @return *python3.PyObject
func ImportModule(dir, name string) *python3.PyObject {
	sysModule := python3.PyImport_ImportModule("sys")
	path := sysModule.GetAttrString("path")
	//pathStr, _ := pythonRepr(path)
	//log.Println("before add path is " + pathStr)
	//python3.PyList_Insert(path, 0, python3.PyUnicode_FromString(""))
	python3.PyList_Insert(path, 0, python3.PyUnicode_FromString(dir))
	//pathStr, _ = pythonRepr(path)
	//log.Println("after add path is " + pathStr)
	return python3.PyImport_ImportModule(name)
}

// pythonRepr
// @Description: PyObject转换为string
// @param o
// @return string, error
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

//
func getImportModule() (*python3.PyObject, *python3.PyObject) {

	mainPy := ImportModule("./py3", "main")
	if mainPy == nil {
		log.Fatalf("mainPy is nil")
		return nil, nil
	}
	mainPyRepr, err := pythonRepr(mainPy)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[MODULE] repr(main.py) = %s\n", mainPyRepr)

	siftPy := ImportModule("./py3", "sift")
	if siftPy == nil {
		log.Fatalf("siftPy is nil")
		return nil, nil
	}
	siftPyRepr, err := pythonRepr(mainPy)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[MODULE] repr(sift.py) = %s\n", siftPyRepr)

	return mainPy, siftPy
}

// getSiftFunc
// @param siftPy
// @return []*python3.PyObject
func getSiftFunc(siftPy *python3.PyObject) (*python3.PyObject, *python3.PyObject, *python3.PyObject, *python3.PyObject, *python3.PyObject, *python3.PyObject) {

	getDesFunc := siftPy.GetAttrString("get_des")
	if getDesFunc == nil {
		log.Fatalf("getDesFunc is nil")
	}
	getDesFuncRepr, err := pythonRepr(getDesFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(get_des) = %s\n", getDesFuncRepr)

	encodeFunc := siftPy.GetAttrString("encode")
	if encodeFunc == nil {
		log.Fatalf("encodeFunc is nil")
	}
	encodeFuncRepr, err := pythonRepr(encodeFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(encode) = %s\n", encodeFuncRepr)

	getWordSummaryFunc := siftPy.GetAttrString("get_word_summary")
	if getWordSummaryFunc == nil {
		log.Fatalf("getWordSummaryFunc is nil")
	}
	getWordSummaryFuncRepr, err := pythonRepr(getWordSummaryFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(get_word_summary) = %s\n", getWordSummaryFuncRepr)

	tfidfFunc := siftPy.GetAttrString("tf_idf")
	if tfidfFunc == nil {
		log.Fatalf("tfidfFunc is nil")
	}
	tfidfFuncRepr, err := pythonRepr(tfidfFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(tf_idf) = %s\n", tfidfFuncRepr)

	idfRenderFunc := siftPy.GetAttrString("idf_render")
	if idfRenderFunc == nil {
		log.Fatalf("idfRenderFunc is nil")
	}
	idfRenderFuncRepr, err := pythonRepr(tfidfFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(idf_render) = %s\n", idfRenderFuncRepr)

	summaryMatchFunc := siftPy.GetAttrString("summary_match")
	if summaryMatchFunc == nil {
		log.Fatalf("summaryMatchFunc is nil")
	}
	summaryMatchFuncRepr, err := pythonRepr(tfidfFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(summary_match) = %s\n", summaryMatchFuncRepr)

	return getDesFunc, encodeFunc, getWordSummaryFunc, tfidfFunc, idfRenderFunc, summaryMatchFunc
}

func getMainFunc(mainPy *python3.PyObject) (*python3.PyObject, *python3.PyObject, *python3.PyObject, *python3.PyObject, *python3.PyObject, *python3.PyObject) {
	getCntFunc := mainPy.GetAttrString("get_cnt")
	if getCntFunc == nil {
		log.Fatalf("get_cnt is nil")
	}
	getCntFuncRepr, err := pythonRepr(getCntFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(get_cnt) = %s\n", getCntFuncRepr)

	getTotCntFunc := mainPy.GetAttrString("get_tot_cnt")
	if getTotCntFunc == nil {
		log.Fatalf("get_tot_cnt is nil")
	}
	getTotCntFuncRepr, err := pythonRepr(getTotCntFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(get_tot_cnt) = %s\n", getTotCntFuncRepr)

	getFileFunc := mainPy.GetAttrString("get_file")
	if getFileFunc == nil {
		log.Fatalf("get_file is nil")
	}
	getFileFuncRepr, err := pythonRepr(getFileFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(getFileFuncRepr) = %s\n", getFileFuncRepr)

	uploadFunc := mainPy.GetAttrString("upload_file")
	if uploadFunc == nil {
		log.Fatalf("upload_file is nil")
	}
	uploadFuncRepr, err := pythonRepr(uploadFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(uploadFuncRepr) = %s\n", uploadFuncRepr)

	reduceFunc := mainPy.GetAttrString("reduce")
	if reduceFunc == nil {
		log.Fatalf("reduce is nil")
	}
	reduceFuncRepr, err := pythonRepr(reduceFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(reduceFuncRepr) = %s\n", reduceFuncRepr)

	findFunc := mainPy.GetAttrString("find_image")
	if findFunc == nil {
		log.Fatalf("find is nil")
	}
	findFuncRepr, err := pythonRepr(reduceFunc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[FUNC] repr(findFuncRepr) = %s\n", findFuncRepr)

	return getCntFunc, getTotCntFunc, getFileFunc, uploadFunc, reduceFunc, findFunc
}

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func getRes(s string) (Res Response) {
	//res = res[1 : len(res)-1]
	str := make([]string, 0)
	j := len(s) - 1
	for i := len(s) - 2; i >= 0; i-- {
		if s[i] == ':' {
			str = append(str, s[i+2:j])
		} else if s[i] == ',' {
			j = i
		}
	}
	if len(str) < 3 {
		return
	}
	Res.Code = str[2]
	Res.Msg = str[1]
	Res.Data = str[0]
	// fmt.Println(Res)
	return
}

func getFindRes(s string) {

}

func main() {

	p := "/Users/fengzetao/GolandProjects/go-python/venv/lib/python3.7/site-packages"
	InsertBeforeSysPath(p)

	mainPy, _ := getImportModule()

	//getDesFunc, encodeFunc, getWordSummaryFunc, tfidfFunc, idfRenderFunc, summaryMatchFunc := getSiftFunc(siftPy)
	//defer getDesFunc.DecRef()
	//defer encodeFunc.DecRef()
	//defer getWordSummaryFunc.DecRef()
	//defer tfidfFunc.DecRef()
	//defer idfRenderFunc.DecRef()
	//defer summaryMatchFunc.DecRef()
	//fmt.Println(getDesFunc, encodeFunc, getWordSummaryFunc, tfidfFunc, idfRenderFunc, summaryMatchFunc)

	getCntFunc, getTotCntFunc, getFileFunc, uploadFunc, reduceFunc, findFunc := getMainFunc(mainPy)
	//defer getCntFunc.DecRef()
	//defer getTotCntFunc.DecRef()
	//defer getFileFunc.DecRef()
	fmt.Println(getCntFunc, getTotCntFunc, getFileFunc, uploadFunc, reduceFunc, findFunc)

	var wg sync.WaitGroup
	wg.Add(500)
	state := python3.PyEval_SaveThread()

	r := gin.Default()

	// 加载静态资源
	r.Static("/assets", "./statics")
	r.LoadHTMLFiles("templates/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	go func() {
		// 获取Cache中索引图片的数量
		r.GET("/api/num", func(c *gin.Context) {
			defer wg.Done()
			_state := python3.PyGILState_Ensure()
			defer python3.PyGILState_Release(_state)
			res := getCntFunc.Call(python3.Py_None, python3.Py_None)
			resJson, _ := pythonRepr(res)
			// fmt.Printf("[VARS] CntJson = %s\n", resJson)
			Res := getRes(resJson)
			c.JSON(http.StatusOK, gin.H{
				"code": Res.Code,
				"msg":  Res.Msg,
				"data": Res.Data,
			})
		})
	}()

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

	// 获取选择的图片
	r.GET("/api/img/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		args := python3.PyTuple_New(1)
		python3.PyTuple_SetItem(args, 0, python3.PyUnicode_FromString(filename))
		c.JSON(http.StatusOK, gin.H{
			"data": getFileFunc.Call(args, python3.Py_None),
		})
	})

	go func() {
		//  图片上传
		r.POST("/api/upload", func(c *gin.Context) {
			form, _ := c.MultipartForm()
			files := form.File["file"]
			// wg.Add(len(files))
			for _, file := range files {
				file := file
				wg.Add(1)
				go func() {
					defer wg.Done()
					_state := python3.PyGILState_Ensure()
					defer python3.PyGILState_Release(_state)
					args := python3.PyTuple_New(1)
					python3.PyTuple_SetItem(args, 0, python3.PyUnicode_FromString(file.Filename))
					uploadFunc.Call(args, python3.Py_None)
					dst := fmt.Sprintf("./train/%s", file.Filename)
					_ = c.SaveUploadedFile(file, dst)
				}()
			}
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("%d files uploaded!", len(files)),
			})
		})
	}()

	// 检索图片（添加索引）
	r.POST("/api/reduce", func(c *gin.Context) {
		wg.Add(1)
		defer wg.Done()
		_state := python3.PyGILState_Ensure()
		defer python3.PyGILState_Release(_state)
		res := reduceFunc.Call(python3.Py_None, python3.Py_None)
		resJson, _ := pythonRepr(res)
		fmt.Printf("[VARS] reduceJson = %s\n", resJson)
		Res := getRes(resJson)
		c.JSON(http.StatusOK, gin.H{
			"code": Res.Code,
			"msg":  Res.Msg,
			"data": Res.Data,
		})
	})

	// 图像查找
	r.POST("/api/find", func(c *gin.Context) {
		wg.Add(1)
		defer wg.Done()
		_state := python3.PyGILState_Ensure()
		defer python3.PyGILState_Release(_state)
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		dst := fmt.Sprintf("./test/%s", file.Filename)
		_ = c.SaveUploadedFile(file, dst)

		args := python3.PyTuple_New(1)
		python3.PyTuple_SetItem(args, 0, python3.PyUnicode_FromString(file.Filename))
		res := findFunc.Call(args, python3.Py_None)
		resJson, _ := pythonRepr(res)
		fmt.Printf("[VARS] findJson = %s\n", resJson)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})

	err := r.Run(":5000")
	if err != nil {
		fmt.Printf("gin run failed, err:%v\n", err)
		return
	}

	wg.Wait()
	python3.PyEval_RestoreThread(state)

	python3.Py_Finalize()
}
