//封装 http.ResponseWriter 和 http.Request，并提供http和json的直接生成方法

package gee

import (
	"encoding/json"
	"net/http"
)

//json键值对
type H map[string]interface{}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	engine  *Engine

	//提供Path和Method的直接访问途径
	Path, Method string

	//Request的状态码
	StatusCode int

	//动态路由的参数
	Param map[string]string

	//需要实行的中间件应该保存在context中
	middleware []HandleFunc
	index      int
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:     w,
		Request:    req,
		Path:       req.URL.Path,
		Method:     req.Method,
		middleware: make([]HandleFunc, 0),
		index:      -1,
	}
}

//执行下一个handlefunc,用于将next之后的代码块置后
func (c *Context) Next() {
	//取出下一个handlefunc的index
	c.index++
	max := len(c.middleware) - 1
	for ; c.index <= max; c.index++ {
		c.middleware[c.index](c)
	}
}

//根据键值获取POST的Form值
func (c *Context) PostForm(key string) interface{} {
	return c.Request.FormValue(key)
}

//设置头部的键值对
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//设置头部的状态码
func (c *Context) SetStatus(status int) {
	c.Writer.WriteHeader(status)
	c.StatusCode = status
}

//生成HTML响应
func (c *Context) HTML(status int, filename string, html interface{}) {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.SetStatus(status)
	c.engine.httptemplate.ExecuteTemplate(c.Writer, filename, html)
}

func (c *Context) JSON(status int, j interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(status)
	coder := json.NewEncoder(c.Writer)
	err := coder.Encode(j)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
