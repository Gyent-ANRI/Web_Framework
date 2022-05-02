package gee

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	weblog log.Logger
	weberr log.Logger
)

type HandleFunc func(*Context)

type Engine struct {
	Group

	groups []*Group
	//用router存储路径和方法的对应关系
	router *Router

	//
	httptemplate *template.Template
	funcmap      template.FuncMap
}

func init() {
	logfile, err := os.OpenFile("log/log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	errfile, e := os.OpenFile("log/err.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil || e != nil {
		panic("Fail to open log.txt")
	}
	weblog = *log.New(logfile, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	weberr = *log.New(io.MultiWriter(errfile, os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func New() *Engine {
	e := Engine{
		Group: Group{
			"",
			make([]HandleFunc, 0),
			nil,
		},
		groups: make([]*Group, 0),
		router: NewRouter(),
	}

	weblog.Println("New Engine")

	e.engine = &e
	e.groups = append(e.groups, &e.Group)

	return &e
}

//开始监听Port端口
func (e *Engine) Run(port string) error {
	addr := ":" + port
	weblog.Println("Listening ", addr)
	return http.ListenAndServe(addr, e)
}

//实现Handler接口使得Engine可以接管http请求
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	weblog.Printf("New request. Method: %v, Path: %v, FROM: %v\n", req.Method, req.URL.Path, req.Header["User-Agent"])
	c := NewContext(w, req)
	c.engine = e
	//遍历组别，如果包含前缀则把中间件加入context的中间件列表
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			c.middleware = append(c.middleware, group.midlleware...)
		}
	}
	e.router.handle(c)
}

//设置模板的方法
func (engine *Engine) SetFuncMap(fm template.FuncMap) {
	engine.funcmap = fm
}

//使用模板
func (engine *Engine) LoadHTMLTemp(pattern string) {
	//template.Must处理了生成失败的err, 用New创建一个name为空的模板，采用funcmap,并且解析pattern对应的模板文件
	engine.httptemplate = template.Must(template.New("").Funcs(engine.funcmap).ParseGlob(pattern))
}
