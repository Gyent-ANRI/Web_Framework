package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	roots  map[string]*node      //用于储存不同方法对应的根节点
	routes map[string]HandleFunc //静态路由以及对应处理函数
}

func NewRouter() *Router {
	return &Router{
		make(map[string]*node),
		make(map[string]HandleFunc),
	}
}

func parsePattern(pattern string) []string {
	split := strings.Split(pattern, "/")
	result := make([]string, 0)

	//去除为空的部分
	for _, s := range split {
		if s != "" {
			result = append(result, s)
		}
	}
	return result

}

func (r *Router) addRoute(method string, path string, handlefunc HandleFunc) {
	key := method + "-" + path
	r.routes[key] = handlefunc

	weblog.Println("Add Route: ", key)
	//添加对应节点

	root, exist := r.roots[method]

	if !exist {
		root = &node{part: "root"}
		r.roots[method] = root
	}

	root.insert(path, parsePattern(path), 0)

}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	parts := parsePattern(path)
	root, exist := r.roots[method]
	param := make(map[string]string)

	if !exist {
		return nil, nil
	}

	n := root.search(path, parts, 0)
	if n == nil {
		return nil, nil
	}
	nodeparts := parsePattern(n.pattern)
	for i, nodepart := range nodeparts {
		if nodepart[0] == ':' {
			param[nodepart[1:]] = parts[i]
		}
		if nodepart[0] == '*' && len(nodepart) > 1 {
			//碰到*把后续所有part合并赋值
			param[nodepart[1:]] = strings.Join(parts[i:], "/")
		}
	}
	return n, param

}

func (r *Router) handle(c *Context) {

	n, param := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Param = param
		key := c.Request.Method + "-" + n.pattern
		if handle, OK := r.routes[key]; OK {
			c.middleware = append(c.middleware, handle)
		} else {
			c.HTML(500, "", "Internal problem")
			weberr.Println("Fail to find handlefunc")
			return
		}
	} else {
		c.middleware = append(c.middleware, func(ctx *Context) {
			c.HTML(http.StatusNotFound, "", "404 NOT FOUND")
			weberr.Println("Fail to get route ", c.Path)
		})
	}
	c.Next()
}
