package jee

import "testing"

func newNode() *node {
	return &node{
		isWild: false,
	}
}

func TestAddNode(t *testing.T) {
	n := newNode()
	n.addNode("/hello", splitPattern("/hello"), 0)
	c := n.children[0]
	n.addNode("/hello/omo", splitPattern("/hello/omo"), 0)
	if c.children[0].pattern != "/hello/omo" || c.pattern != "/hello" {
		t.Fail()
	}
}
