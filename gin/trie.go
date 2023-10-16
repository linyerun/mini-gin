package gee

import "strings"

//核心方法就两个：search(路由选择)、insert(注册路由)
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang	(只有当注册进来的pattern到头了才设置,不然为"")
	part     string  // 路由中的一部分，例如 :lang
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
	children []*node // 子节点，例如 [doc, tutorial, intro]
}

func newNode(part string, isWild bool) *node {
	return &node{
		part:   part,
		isWild: isWild,
	}
}

//用于插入
func (n *node) matchChild(part string) (res *node) {
	// /:name/hello 和 /go 不一样，go不会在:name那里停留
	for _, child := range n.children {
		if child.part == part {
			res = child
			return
		}
	}
	return
}
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	//查找子那里是否存在一样的
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		//不存在
		child = newNode(part, strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*"))
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//用于查询
func (n *node) matchChildren(part string) (res []*node) {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			res = append(res, child)
		}
	}
	return
}
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern != "" {
			return n
		}
		return nil
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		if node := child.search(parts, height+1); node != nil {
			return node
		}
	}
	return nil
}
