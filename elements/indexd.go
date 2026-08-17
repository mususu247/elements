package elements

import (
	"fmt"
	"log"
)

type indexd struct {
	parent *Node
	self   *Node
}

func (nd *Node) Indexd() *indexd {
	err := nd.check()
	if err != nil {
		return nil
	}

	var ix indexd
	ix.parent = nd
	ix.self = nd.firstByTagName("$index")
	return &ix
}

func (ix *indexd) Set() *Node {
	if ix == nil {
		log.Println("(Error) Indexd is null")
		return nil
	}

	if ix.self == nil {
		nd := ix.parent
		doc := nd.earth
		ix.self = doc.CreateElement("index")
		ix.self.name = "$index"
		nd.AppendChild(ix.self)
	}
	return ix.self
}

func (ix *indexd) append(ref int64) (*Node, error) {
	if ix == nil {
		return nil, fmt.Errorf("(Error) Friend is null")
	}

	nd := ix.parent
	doc := nd.earth
	link := doc.CreateElement("link")
	link.name = "$link"
	link.attri["ref"] = ref
	link.attri["pendant"] = ""
	ix.self.AppendChild(link)
	return link, nil
}

func (ix *indexd) Add(node *Node) error {
	if ix == nil {
		return fmt.Errorf("(Error) Friend is null")
	}

	if ix.parent.num == node.num {
		return fmt.Errorf("(Error) not self link")
	}

	ix.Set()
	sfLinks := ix.self.getNodesByTagName("$link")
	for i := range sfLinks {
		if v, ok := sfLinks[i].attri["ref"]; ok {
			var ref int64

			switch x := v.(type) {
			case int:
				ref = int64(x)
			case int64:
				ref = x
			case float64:
				ref = int64(x)
			}
			if ref == ix.self.num {
				return fmt.Errorf("(Error) is linked")
			}
		}
	}

	ix.append(node.num)
	ix.self.attri["count"] = len(ix.self.children)
	return nil
}

func (ix *indexd) Remove(node *Node) error {
	if ix == nil {
		return fmt.Errorf("(Error) Indexd is null")
	}

	if ix.parent.num == node.num {
		return fmt.Errorf("(Error) not self link")
	}

	sfLinks := ix.self.getElementsByTagName("$link")
	for i := range sfLinks {
		if v, ok := sfLinks[i].attri["ref"]; ok {
			var ref int64

			switch x := v.(type) {
			case int:
				ref = int64(x)
			case int64:
				ref = x
			case float64:
				ref = int64(x)
			}

			if ref == node.num {
				ix.self.RemoveChild(sfLinks[i])
			}
		}
	}

	ix.self.attri["count"] = len(ix.self.children)
	return nil
}

func (ix *indexd) Links() []*Node {
	var list []*Node

	nd := ix.parent
	doc := nd.earth

	links := ix.self.getElementsByTagName("$link")
	for i := range links {
		var ref int64

		if v, ok := links[i].attri["ref"]; ok {
			switch x := v.(type) {
			case int:
				ref = int64(x)
			case int64:
				ref = x
			case float64:
				ref = int64(x)
			}
		}

		if indexd, ok := doc.flat[ref]; ok {
			if node, ok := doc.flat[indexd.parent]; ok {
				list = append(list, node)
			}
		}
	}

	return list
}

func (ix *indexd) Noting() {
	ix.parent = nil
	ix.self = nil
	ix = nil
}
