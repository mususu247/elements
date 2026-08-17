package elements

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Dom struct {
	count    int64
	flat     map[int64]*Node
	fileName string
}

func (doc *Dom) Init(indent string) {
	doc.count = 0
	doc.flat = make(map[int64]*Node)
	root := doc.CreateElement("root")
	root.name = "$root"
	root.SetAttribute("indent", indent)
}

func (doc *Dom) check() error {
	if doc == nil {
		return fmt.Errorf("(Error) Dom is null")
	}
	return nil
}

func (doc *Dom) Root() *Node {
	err := doc.check()
	if err != nil {
		return nil
	}

	if root, ok := doc.flat[1]; ok {
		return root
	}
	return nil
}

func (doc *Dom) Head() *Node {
	err := doc.check()
	if err != nil {
		return nil
	}

	root := doc.Root()
	nodes := root.getElementsByTagName("head")
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

func (doc *Dom) Body() *Node {
	err := doc.check()
	if err != nil {
		return nil
	}

	root := doc.Root()
	nodes := root.getElementsByTagName("body")
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

func (doc *Dom) Open(fileName string) error {
	err := doc.check()
	if err != nil {
		return err
	}

	fullPath, err := GetAbsolutePathName(fileName)
	if err != nil {
		return err
	}

	if !FileExists(fullPath) {
		return fmt.Errorf("(Error) not found file:'%v'", fullPath)
	}

	var jsonData []byte
	jsonData, err = os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	var tw Twig
	err = json.Unmarshal(jsonData, &tw)
	if err != nil {
		return err
	}

	numList := make(map[int64]int64)
	root := doc.Root()
	root.fromTwig(&tw, numList)

	links := root.getNodesByTagName("$link")
	for i := range links {
		if v, ok := links[i].attri["ref"]; ok {
			var Num int64
			switch x := v.(type) {
			case int:
				Num = int64(x)
			case int64:
				Num = x
			case float64:
				Num = int64(x)
			}
			ref := numList[Num]
			links[i].attri["ref"] = ref
		}
	}
	return nil
}

func (doc *Dom) SaveAs(fileName string) error {
	err := doc.check()
	if err != nil {
		return err
	}

	fullPath, err := GetAbsolutePathName(fileName)
	if err != nil {
		return err
	}

	root := doc.Root()
	indent := root.GetAttribute("indent")
	tw := root.toTwig()

	var jsonData []byte
	if indent != "" {
		jsonData, err = json.MarshalIndent(tw, "", indent)
		if err != nil {
			return err
		}
	} else {
		jsonData, err = json.Marshal(tw)
		if err != nil {
			return err
		}
	}

	if FileExists(fullPath) {
		err = os.Remove(fullPath)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return err
	}

	tw.Nothing()
	doc.fileName = fullPath
	return nil
}

func (doc *Dom) Close() {
	// remove All
	for i := range doc.flat {
		if _, ok := doc.flat[i]; ok {
			doc.flat[i].earth = nil
			doc.flat[i].attri = nil
			doc.flat[i].children = nil
			delete(doc.flat, i)
		}
	}
	doc.flat = nil
	doc.count = 0
}

// Create
func (doc *Dom) CreateElement(tag string) *Node {
	err := doc.check()
	if err != nil {
		return nil
	}
	var node Node

	if tag, ok := checkTag(tag); !ok {
		log.Printf("(Error) tag is invalid: %s", tag)
		return nil
	}

	node.name = tag
	node.attri = make(map[string]any)
	node.earth = doc

	doc.count++
	node.num = doc.count
	doc.flat[node.num] = &node
	return &node
}

func (doc *Dom) CreateElementNS(space string, tag string) *Node {
	err := doc.check()
	if err != nil {
		return nil
	}
	var node Node

	if tag, ok := checkTag(tag); !ok {
		log.Printf("(Error) tag is invalid: %s", tag)
		return nil
	}

	node.space = space
	node.name = tag
	node.attri = make(map[string]any)
	node.earth = doc

	doc.count++
	node.num = doc.count
	doc.flat[node.num] = &node
	return &node
}

func (doc *Dom) CreateTextNode(text string) *Node {
	err := doc.check()
	if err != nil {
		return nil
	}
	var node Node

	node.name = "$text"
	node.attri = make(map[string]any)
	node.attri["text"] = text
	node.earth = doc

	doc.count++
	node.num = doc.count
	doc.flat[node.num] = &node
	return &node
}

func (doc *Dom) CreateComment(text string) *Node {
	err := doc.check()
	if err != nil {
		return nil
	}
	var node Node

	var kind string
	if strings.Contains(text, "<!") {
		kind = "!"
	} else {
		kind = "?"
	}

	node.name = "$comment"
	node.attri = make(map[string]any)
	node.attri["kind"] = kind
	node.attri["text"] = text
	node.earth = doc

	doc.count++
	node.num = doc.count
	doc.flat[node.num] = &node
	return &node
}

func (doc *Dom) CreateDocumentFragment() *Node {
	err := doc.check()
	if err != nil {
		return nil
	}
	var node Node

	node.name = "$frag"
	node.attri = make(map[string]any)
	node.earth = doc

	doc.count++
	node.num = doc.count
	doc.flat[node.num] = &node
	return &node
}

// Alias
func (doc *Dom) SetAlias(name string, alias ...string) error {
	err := doc.check()
	if err != nil {
		return err
	}

	root := doc.Root()
	node := root.firstByTagName("$alias")
	if node == nil {
		node = doc.CreateElement("alias")
		node.name = "$alias"
		root.AppendChild(node)
	}

	if attr, ok := node.attri[name]; ok {
		switch x := attr.(type) {
		case string:
			list := strings.Split(x, " ")
			list = append(list, alias...)
			list = removeDupStr(list)
			line := strings.Join(list, " ")
			node.attri[name] = line
		}
	} else {
		slices.Sort(alias)
		line := strings.Join(alias, " ")
		node.attri[name] = line
	}
	return nil
}

func (doc *Dom) GetAlias(name string, inc bool) []string {
	err := doc.check()
	if err != nil {
		return nil
	}
	var list []string

	root := doc.Root()
	node := root.firstByTagName("$alias")
	if node == nil {
		if inc {
			list = append(list, name)
		}
		return list
	}

	if attr, ok := node.attri[name]; ok {
		switch x := attr.(type) {
		case string:
			list := strings.Split(x, " ")
			if inc {
				list = append(list, name)
				slices.Sort(list)
			}
			return list
		}
	}

	if inc {
		list = append(list, name)
	}
	return list
}

// ReIndex
func (doc *Dom) ReIndex() error {
	err := doc.check()
	if err != nil {
		return err
	}

	root := doc.Root()
	root.Reindex()
	return nil
}

// Debug.Print
func (nd *Node) Debug(level int) {
	fmt.Printf("%v num:%v space:'%v' name:'%v' index:%v attr:%v parent:%v \n", strings.Repeat(".", level), nd.num, nd.space, nd.name, nd.index, nd.GetAttributeNames(), nd.parent)
	for i := range nd.children {
		nd.earth.flat[nd.children[i]].Debug(level + 1)
	}
}
