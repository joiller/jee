package jee

import "strings"

type node struct {
	part     string
	pattern  string
	children []*node
	isWild   bool
}

func splitPattern(pattern string) []string {
	splits := strings.Split(pattern, "/")
	res := make([]string, 0)
	for i, split := range splits {
		if split == "" {
			continue
		}
		if split[0] == '*' {
			if len(splits) > i+1 {
				split += strings.Join(splits[i+1:], "/")
			}
			i = len(splits)
		}
		res = append(res, split)
	}
	return res
}

func (n *node) searchNode(patterns []string, height int) *node {
	if height == len(patterns) {
		return n
	}
	part := patterns[height]
	for _, child := range n.children {
		if child.part == part || child.isWild {
			search := child.searchNode(patterns, height+1)
			if search != nil {
				return search
			}
		}
	}
	return nil
}

// pattern是一个完整路径
func (n *node) addNode(pattern string, patterns []string, height int) {
	if len(patterns) == height {
		n.pattern = pattern
		return
	}
	var next *node
	for _, child := range n.children {
		if child.part == patterns[height] || child.isWild {
			next = child
			break
		}
	}
	if next == nil {
		next = &node{
			part:   patterns[height],
			isWild: patterns[height][0] == '*' || patterns[height][0] == ':',
		}
		n.children = append(n.children, next)
	}
	next.addNode(pattern, patterns, height+1)
}
