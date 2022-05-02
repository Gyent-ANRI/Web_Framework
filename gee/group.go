package gee

import (
	"net/http"
	"path"
)

type Group struct {
	prefix     string       //前缀
	midlleware []HandleFunc //中间件
	engine     *Engine      //共有一个指向引擎的指针
}

func (g *Group) AddGroup(pref string) *Group {
	newgroup := &Group{
		prefix: g.prefix + pref,
		engine: g.engine,
	}
	g.engine.groups = append(g.engine.groups, newgroup)
	weblog.Println("Add Group: ", pref)
	return newgroup
}

func (g *Group) addRoute(method string, path string, f HandleFunc) {
	g.engine.router.addRoute(method, g.prefix+path, f)
}

func (g *Group) Get(path string, f HandleFunc) {
	g.addRoute("GET", path, f)
}

func (g *Group) Post(path string, f HandleFunc) {
	g.addRoute("POST", path, f)
}

//使用中间件
func (g *Group) Use(funcs ...HandleFunc) {
	g.midlleware = append(g.midlleware, funcs...)
}

//返回一个方法，用http.FileServer处理请求 relativepath是相对路径，fs是根目录
func (g *Group) GetFileHandle(relativePath string, fs http.FileSystem) HandleFunc {

	//http.StripPrefix把请求Path转换成对应fs的路径
	fileServer := http.StripPrefix(path.Join(g.prefix, relativePath), http.FileServer(fs))
	return func(c *Context) {
		//检查文件是否存在
		file := c.Param["filepath"]
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			weberr.Println("Fail to find the file ", file)
			return
		}
		weblog.Println("Return file: ", file)
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

//增加静态的路径对文件的映射
func (g *Group) Static(relativePath string, root string) {

	fileHandler := g.GetFileHandle(relativePath, http.Dir(root))

	g.Get(path.Join(relativePath, "/*filepath"), fileHandler)
}
