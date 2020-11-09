package framework

import "strings"

type node struct {
	path             string
	name             string
	paramKey         []string
	paramValue       []string
	isLeaf           bool
	handler          HandlerFunc
	beforeMiddleware []HandlerFunc
	afterMiddleware  []HandlerFunc
	children         []*node
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (n *node) matchChild(path string) *node {
	for _, child := range n.children {
		if child.path == path {
			return child
		}
	}
	return nil
}

func (n *node) isLeafNode(parts []string) (bool, []string) {
	if len(parts) == 1 {
		return true, []string{}
	}
	if len(parts) > 1 && parts[1][0] == ':' {
		paramKey := []string{}
		for i := 1; i < len(parts); i++ {
			paramKey = append(paramKey, strings.Replace(parts[i], ":", "", 1))
		}
		return true, paramKey
	}
	return false, []string{}
}

func (n *node) insert(parts []string, name string, handler HandlerFunc) {
	part := parts[0]
	child := n.matchChild(part)
	if child == nil {
		child = &node{path: part}
		n.children = append(n.children, child)
	}
	if child.isLeaf {
		panic("repeated route: " + name)
	}
	if ok, paramKey := n.isLeafNode(parts); ok {
		child.name = name
		child.paramKey = paramKey
		child.handler = handler
		child.isLeaf = true
	} else {
		child.insert(parts[1:], name, handler)
	}
}

func (n *node) search(parts []string) *node {
	if len(parts) == 0 {
		return nil
	}
	child := n.matchChild(parts[0])
	if child == nil {
		return nil
	}
	if child.isLeaf {
		if len(parts) > 1 {
			child.paramValue = parts[1:]
		}
		return child
	}
	return child.search(parts[1:])
}

func (n *node) insertMiddleware(parts []string, beforeMiddleware HandlerFunc, afterMiddleware HandlerFunc) {
	if len(parts) == 0 {
		return
	}
	part := parts[0]
	child := n.matchChild(part)
	if child == nil {
		child = &node{path: part}
		n.children = append(n.children, child)
	}
	if ok, _ := n.isLeafNode(parts); ok {
		if beforeMiddleware != nil {
			child.beforeMiddleware = append(child.beforeMiddleware, beforeMiddleware)
		}
		if afterMiddleware != nil {
			child.afterMiddleware = append(child.afterMiddleware, afterMiddleware)
		}
	} else {
		child.insertMiddleware(parts[1:], beforeMiddleware, afterMiddleware)
	}
}

func (n *node) searchMiddleware(parts []string, beforeMiddleware []HandlerFunc, afterMiddleware []HandlerFunc) (*node, []HandlerFunc, []HandlerFunc) {
	if len(parts) == 0 {
		return nil, beforeMiddleware, afterMiddleware
	}
	child := n.matchChild(parts[0])
	if child == nil {
		return nil, beforeMiddleware, afterMiddleware
	}
	beforeMiddleware = append(beforeMiddleware, child.beforeMiddleware...)
	afterMiddleware = append(afterMiddleware, child.afterMiddleware...)
	if child.isLeaf {
		if len(parts) > 1 {
			child.paramValue = parts[1:]
		}
		return child, beforeMiddleware, afterMiddleware
	}
	return child.searchMiddleware(parts[1:], beforeMiddleware, afterMiddleware)
}
