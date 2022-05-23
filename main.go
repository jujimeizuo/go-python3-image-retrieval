package main

import (
	"encoding/json"
	"fmt"
	"github.com/DataDog/go-python3"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
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

// getImportModule
// @Description: 返回py导入包
// return []*python3.PyObject
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

// getMainFunc
// @Description: 返回包函数
// return []*python3.PyObject
func getMainFunc(mainPy *python3.PyObject) (*python3.PyObject, *python3.PyObject, *python3.PyObject, *python3.PyObject, *python3.PyObject) {
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

	return getCntFunc, getTotCntFunc, uploadFunc, reduceFunc, findFunc
}

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type InfoAll struct {
	P50_ac  string `json:"p50_ac"`
	P50_rc  string `json:"p50_rc"`
	Tp10_ac string `json:"tp10_ac"`
	Tp10_rc string `json:"tp10_rc"`
	Tp5_ac  string `json:"tp5_ac"`
	Tp5_rc  string `json:"tp5_rc"`
}

type DataAll struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type ResInfoData struct {
	Info InfoAll   `json:"info"`
	Data []DataAll `json:"data"`
}

func getRes(s string) (Res Response) {
	str := make([]string, 0)
	flag := false
	var j int
	for i := 0; i < len(s)-1; i++ {
		if s[i] == ':' {
			flag = true
			j = i + 2
		} else if flag && s[i] == ',' {
			str = append(str, s[j:i])
			flag = false
		} else if flag && len(str) == 2 {
			str = append(str, s[j:len(s)-1])
			break
		}
	}
	if len(str) < 3 {
		return
	}
	Res.Code = str[0]
	Res.Msg = str[1]
	Res.Data = str[2]
	//fmt.Println(Res)
	return
}

func getResInfoAndData(s string) (res ResInfoData) {
	s = strings.Replace(s, "'", "\"", -1)
	_ = json.Unmarshal([]byte(s), &res)
	//fmt.Println(res)
	return res
}

func main() {

	p := "/Users/fengzetao/GolandProjects/go-python/venv/lib/python3.7/site-packages"
	InsertBeforeSysPath(p)

	mainPy, _ := getImportModule()

	getCntFunc, getTotCntFunc, uploadFunc, reduceFunc, findFunc := getMainFunc(mainPy)

	fmt.Println(getCntFunc, getTotCntFunc, uploadFunc, reduceFunc, findFunc)

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

	// 获取Cache中索引图片的数量
	go func() {
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

	// 获取已上传的图片数量
	go func() {
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

	//获取选择的图片
	r.GET("/api/img/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		// fmt.Printf("[VARS] filename = %s\n", filename)
		c.File("./img/" + filename)
	})

	//  图片上传
	go func() {
		r.POST("/api/upload", func(c *gin.Context) {
			form, _ := c.MultipartForm()
			files := form.File["file"]
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
				}()
				dst := fmt.Sprintf("./img/%s", file.Filename)
				_ = c.SaveUploadedFile(file, dst)
			}
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("%d files uploaded!", len(files)),
			})
		})
	}()

	// 检索图片（添加索引）
	go func() {
		r.POST("/api/reduce", func(c *gin.Context) {
			defer wg.Done()
			_state := python3.PyGILState_Ensure()
			defer python3.PyGILState_Release(_state)
			res := reduceFunc.Call(python3.Py_None, python3.Py_None)
			resJson, _ := pythonRepr(res)
			// fmt.Printf("[VARS] reduceJson = %s\n", resJson)
			Res := getRes(resJson)
			c.JSON(http.StatusOK, gin.H{
				"code": Res.Code,
				"msg":  Res.Msg,
				"data": Res.Data,
			})
		})
	}()

	// 图像查找
	go func() {
		r.POST("/api/find", func(c *gin.Context) {
			defer wg.Done()
			_state := python3.PyGILState_Ensure()
			defer python3.PyGILState_Release(_state)
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
			dst := fmt.Sprintf("./img/%s", file.Filename)
			_ = c.SaveUploadedFile(file, dst)
			var resJson string

			args := python3.PyTuple_New(1)
			python3.PyTuple_SetItem(args, 0, python3.PyUnicode_FromString(file.Filename))
			res := findFunc.Call(args, python3.Py_None)
			resJson, _ = pythonRepr(res)
			//fmt.Printf("[VARS] findJson = %s\n", resJson)
			Res := getRes(resJson)
			InfoAndData := getResInfoAndData(Res.Data)
			c.JSON(http.StatusOK, gin.H{
				"code": Res.Code,
				"msg":  Res.Msg,
				"data": InfoAndData,
			})
		})
	}()

	err := r.Run(":5000")
	if err != nil {
		fmt.Printf("gin run failed, err:%v\n", err)
		return
	}

	wg.Wait()
	python3.PyEval_RestoreThread(state)

	python3.Py_Finalize()
}
