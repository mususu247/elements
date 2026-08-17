package elements

import (
	"fmt"
	"log"
	"path/filepath"
)

type friend struct {
	parent *Node
	self   *Node
}

func (nd *Node) Friend() *friend {
	err := nd.check()
	if err != nil {
		return nil
	}

	var fd friend
	fd.parent = nd
	fd.self = nd.firstByTagName("$friend")
	return &fd
}

func (fd *friend) Set() *Node {
	if fd == nil {
		log.Println("(Error) Friend is null")
		return nil
	}

	if fd.self == nil {
		nd := fd.parent
		doc := nd.earth
		fd.self = doc.CreateElement("friend")
		fd.self.name = "$friend"
		nd.AppendChild(fd.self)
	}
	return fd.self
}

func (fd *friend) append(ref int64, pendant string) (*Node, error) {
	if fd == nil {
		return nil, fmt.Errorf("(Error) Friend is null")
	}

	nd := fd.parent
	doc := nd.earth
	link := doc.CreateElement("link")
	link.name = "$link"
	link.attri["ref"] = ref
	link.attri["pendant"] = pendant
	fd.self.AppendChild(link)
	return link, nil
}

func (fd *friend) Add(node *Node, from string, to string) error {
	if fd == nil {
		return fmt.Errorf("(Error) Friend is null")
	}

	if fd.parent.num == node.num {
		return fmt.Errorf("(Error) not self link")
	}

	nf := node.Friend().Set()
	nfLinks := nf.getNodesByTagName("$link")
	for i := range nfLinks {
		if v, ok := nfLinks[i].attri["ref"]; ok {
			var ref int64

			switch x := v.(type) {
			case int:
				ref = int64(x)
			case int64:
				ref = x
			case float64:
				ref = int64(x)
			}
			if ref == fd.self.num {
				return fmt.Errorf("(Error) is linked")
			}
		}
	}

	fd.Set()
	sfLinks := fd.self.getNodesByTagName("$link")
	for i := range sfLinks {
		if v, ok := sfLinks[i].attri["ref"]; ok {
			switch ref := v.(type) {
			case int64:
				if ref == nf.num {
					return fmt.Errorf("(Error) is linked")
				}
			}
		}
	}

	fd.append(nf.num, from)
	node.Friend().append(fd.self.num, to)
	fd.self.attri["count"] = len(fd.self.children)
	nf.attri["count"] = len(nf.children)
	return nil
}

func (fd *friend) Remove(node *Node) error {
	if fd == nil {
		return fmt.Errorf("(Error) Friend is null")
	}

	if fd.parent.num == node.num {
		return fmt.Errorf("(Error) not self link")
	}

	//nd := fd.self
	//doc := nd.earth

	nf := node.Friend().Set()
	nfLinks := nf.getElementsByTagName("$link")
	for i := range nfLinks {
		if v, ok := nfLinks[i].attri["ref"]; ok {
			var ref int64

			switch x := v.(type) {
			case int:
				ref = int64(x)
			case int64:
				ref = x
			case float64:
				ref = int64(x)
			}

			if ref == fd.self.num {
				nf.RemoveChild(nfLinks[i])
			}
		}
	}

	sfLinks := fd.self.getElementsByTagName("$link")
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

			if ref == nf.num {
				fd.self.RemoveChild(sfLinks[i])
			}
		}
	}

	fd.self.attri["count"] = len(fd.self.children)
	nf.attri["count"] = len(nf.children)
	return nil
}

func (fd *friend) Links(find string) []*Node {
	var list []*Node

	nd := fd.parent
	doc := nd.earth

	links := fd.self.getElementsByTagName("$link")
	for i := range links {
		var ref int64
		var pendant string

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

		if v, ok := links[i].attri["pendant"]; ok {
			switch x := v.(type) {
			case string:
				pendant = x
			}
		}

		_, err := filepath.Match(find, pendant)
		if err == nil {
			if friend, ok := doc.flat[ref]; ok {
				if node, ok := doc.flat[friend.parent]; ok {
					list = append(list, node)
				}
			}
		}
	}

	return list
}

func (fd *friend) Noting() {
	fd.parent = nil
	fd.self = nil
	fd = nil
}
