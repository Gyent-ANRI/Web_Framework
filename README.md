## Gee 简易Web框架
**参考极客兔兔的[Gee](https://geektutu.com/post/gee.html)**  
自学golang过程中找的学习用项目，旨在熟悉golang语法和了解golang对http的处理方式
### 代码介绍
#### 框架部分
web框架部分基本和教程没有太大的区别，在engine.go中定义了两个Logger用于处理日志
- context: 封装http的请求和响应，并且提供了生成html和json响应的方法。 (没有系统学过html和json，只是写了个大概)
- group: 分组控制，根据路径前缀分组，可以以组为单位添加中间件
- engine: http.Handler，能接管指定端口的http处理，内嵌了group, 可以视为 root 组
- router: 实现了路由功能，储存着路径和处理方法的映射，可以为路径调用对应的方法，借助前缀树trie.go解析动态路由
- trie: 以前缀为键的树，实现了动态路由
- recovery: 错误恢复，防止某个部分出错后整个程序中止，以中间件的形式应用在engine上
- template: 通过engine的LoadHTMLTemp方法可以加载模板文件夹并且使用其中的html模板
#### demo部分
主要是为了测试框架的各项功能
- '/': 属于engine组，应用default模板
- '/user': engine的子组，应用usertemp模板,组下路径/user/:name/space 是动态路由，:name会匹配用户名并储存
- '/user/vip': user的子组，测试子组会应用父组的所有中间件，也能有自己独特的中间件
- '/static': 配置请求路径和本地路径的映射，然后用FileServer处理请求
- '/json': engine组下的路径，用于测试POST访问的响应