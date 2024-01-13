package soup

import (
	"io"

	"golang.org/x/net/html"
)

type Node struct {
	node *html.Node
}

func NewNode(node *html.Node) *Node {
	return &Node{node: node}
}

func Parse(r io.Reader) (*Node, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return NewNode(node), nil
}

func (n *Node) Render(w io.Writer) error {
	return html.Render(w, n.node)
}

func (n *Node) Remove() {
	n.node.Parent.RemoveChild(n.node)
}

func (n *Node) AllElements() []*Node {
	var ret []*Node

	if n.node.Type == html.ElementNode {
		ret = append(ret, n)
	}

	for c := n.node.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, NewNode(c).AllElements()...)
	}

	return ret
}

func (n *Node) FindAll(tag string) []*Node {
	var ret []*Node

	if n.node.Type == html.ElementNode && n.node.Data == tag {
		ret = append(ret, n)
	}

	for c := n.node.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, NewNode(c).FindAll(tag)...)
	}

	return ret
}

func (n *Node) GetAttribute(key string) string {
	for _, attr := range n.node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}

func (n *Node) SetAttribute(key, val string) {
	for _, a := range n.node.Attr {
		if a.Key == key {
			a.Val = val
			return
		}
	}

	attr := html.Attribute{Key: key, Val: val}
	n.node.Attr = append(n.node.Attr, attr)
}

func (n *Node) RemoveAttribute(key string) {
	attrs := make([]html.Attribute, 0, len(n.node.Attr))

	for _, a := range n.node.Attr {
		if a.Key != key {
			attrs = append(attrs, a)
		}
	}

	n.node.Attr = attrs
}
