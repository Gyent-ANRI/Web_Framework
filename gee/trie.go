package gee

//前缀树，用于实现动态路由(将动态路由翻译成router内部储存的静态路由)

//树的节点
type node struct {
	pattern  string  //该节点匹配的准确路径
	part     string  //这个节点本身所处的部分
	children []*node //子节点
	isWild   bool    //是否模糊查询
}

//根据part查找子节点并返回第一个满足要求的节点(不允许模糊)
func (n *node) matchfirst(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

//根据part查找子节点并返回所有满足要求的节点
func (n *node) matchall(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		//所有Part满足或者允许模糊查询的子节点都可用
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//通过迭代遍历子节点进行节点插入; height是指当前节点所处的层数
func (n *node) insert(pattern string, parts []string, height int) {
	//递归退出条件：已经在part最高层
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	//下一级的part
	part := parts[height]
	//查找得到的目标节点
	child := n.matchfirst(part)

	//如果查找不到则创建
	if child == nil {
		//如果part是*或者由 : 开头则开启模糊查询
		child = &node{part: part, isWild: part[0] == '*' || part[0] == ':'}
		n.children = append(n.children, child)
	}

	//递归调用
	child.insert(pattern, parts, height+1)
}

//通过递归查询对应节点
func (n *node) search(pattern string, parts []string, height int) *node {

	//退出条件: 触底或者碰到*
	if len(parts) == height || n.part[0] == '*' {
		//pattern为空则是中继节点，查询失败
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	//查询所有满足part的节点
	children := n.matchall(part)

	for _, child := range children {
		result := child.search(pattern, parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
