package main

import (
	"fmt"
	"gee/gee"
	"time"
)

type client struct {
	Name     string
	Group    string
	Operator string
}

func OnlyForVip() gee.HandleFunc {
	return func(ctx *gee.Context) {

	}
}

func OnlyForUser() gee.HandleFunc {
	return func(ctx *gee.Context) {
		t := time.Now()
		fmt.Println("Handle started")
		ctx.Next()
		fmt.Printf("Handle ended, running time: %v\n", time.Since(t))
	}
}

func main() {
	engine := gee.New()

	engine.LoadHTMLTemp("./template/*")

	//以中间件的形式添加错误恢复
	engine.Use(gee.Recovery())

	var groups []string

	engine.Get("/", func(ctx *gee.Context) {
		ctx.HTML(200, "default.tmpl", groups)
	})

	static := engine.AddGroup("/static")
	{
		static.Static("", "./static")
		static.Get("/", func(ctx *gee.Context) {
			ctx.HTML(200, "default.tmpl", []string{
				"/photo",
				"/shell",
				"/video",
			})
		})
	}
	groups = append(groups, "/static")

	user := engine.AddGroup("/user")
	{
		user.Get("/:name/space", func(c *gee.Context) {
			if n, exist := c.Param["name"]; exist {
				data := client{
					Name:     n,
					Group:    "user",
					Operator: "Read Only",
				}
				c.HTML(200, "usertemp.tmpl", data)

			} else {
				c.HTML(500, "", "Fail to find name")
			}
		})
		user.Use(OnlyForUser())
	}
	groups = append(groups, "/user")

	vip := user.AddGroup("/vip")
	{
		vip.Get("/:name/space", func(c *gee.Context) {
			if _, exist := c.Param["name"]; exist {

				c.HTML(200, "default.tmpl", fmt.Sprintf("You are in %v's space", c.Param["name"]))

			} else {
				c.HTML(500, "", "Fail to find name")
			}
		})
		vip.Use(OnlyForVip())
	}
	groups = append(groups, "/user/vip")

	engine.Post("/json", func(c *gee.Context) {
		c.JSON(200, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	groups = append(groups, "/json")

	fmt.Println("开始监听localhost:9999...")
	engine.Run("9999")
}
