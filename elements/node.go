package elements

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Node struct {
	num      int64
	space    string
	name     string
	index    int64
	attri    map[string]any
	parent   int64
	children []int64
	earth    *Dom
}

var errHead = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".", "-"}
var errName = []string{"!", "\"", "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", "/", ";", "=", "?", "@", "[", "\\", "]", "^", "`", "{", "|", "}", "~"}
var SelfClosing = []string{"area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr"}

// Property
func (nd *Node) Name() string {
	err := nd.check()
	if err != nil {
		return ""
	}
	return nd.name
}

func (nd *Node) Space() string {
	err := nd.check()
	if err != nil {
		return ""
	}
	return nd.space
}

func (nd *Node) Num() int64 {
	err := nd.check()
	if err != nil {
		return 0
	}
	return nd.num
}

func (nd *Node) Index() int64 {
	err := nd.check()
	if err != nil {
		return 0
	}
	return nd.index
}

func (nd *Node) getParent() *Node {
	doc := nd.earth
	num := nd.parent
	if parent, ok := doc.flat[num]; ok {
		return parent
	}
	return nil
}

func (nd *Node) Parent() *Node {
	err := nd.check()
	if err != nil {
		return nil
	}

	return nd.getParent()
}

func (nd *Node) getChildren() []*Node {
	var list []*Node

	doc := nd.earth
	for i := range nd.children {
		num := nd.children[i]
		if child, ok := doc.flat[num]; ok {
			list = append(list, child)
		}
	}

	return list
}

func (nd *Node) Children() []*Node {
	err := nd.check()
	if err != nil {
		return nil
	}

	return nd.getChildren()
}

// Chcek
func (nd *Node) check() error {
	if nd == nil {
		log.Println("(Error) Node is null")
		return fmt.Errorf("(Error) Node is null")
	}

	if nd.earth == nil {
		log.Println("(Error) Earth is null")
		return fmt.Errorf("(Error) Earth is null")
	}
	return nil
}

func (nd *Node) nodeCheck(node *Node) bool {
	if node == nil {
		return false
	}

	if node.name[:1] == "$" {
		return true
	}
	return false
}

func (nd *Node) attrCheck(alias []string, value any) bool {
	if nd == nil {
		return false
	}

	var key string
	sw := false
	names := nd.GetAttributeNames()
	for i := range names {
		for j := range alias {
			if names[i] == alias[j] {
				key = alias[j]
				sw = true
			}
		}
	}

	if sw {
		v := nd.attri[key]
		return anyComp(v, value)
	}
	return false
}

// Append and Insert
func (nd *Node) AppendChild(child *Node) error {
	err := nd.check()
	if err != nil {
		return err
	}

	err = child.check()
	if err != nil {
		return err
	}

	if nd.earth != child.earth {
		log.Printf("(Error) not Earth")
		return fmt.Errorf("(Error) not Earth")
	}

	child.parent = nd.num
	nd.children = append(nd.children, child.num)
	return nil
}

func (nd *Node) InsertAdjacentElemnt(position string, child *Node) error {
	err := nd.check()
	if err != nil {
		return err
	}

	err = child.check()
	if err != nil {
		return err
	}

	if nd.earth != child.earth {
		log.Printf("(Error) not Earth")
		return fmt.Errorf("(Error) not Earth")
	}

	pos := strings.ToLower(position)
	switch pos {
	case "beforebegin":
		var backup []int64
		parent := nd.Parent()
		backup = append(backup, parent.children...)
		parent.children = nil

		for i := range backup {
			if backup[i] == nd.num {
				parent.children = append(parent.children, child.num)
				child.parent = parent.num
			}
			parent.children = append(parent.children, backup[i])
		}
	case "afterbegin":
		var backup []int64
		backup = append(backup, nd.children...)
		nd.children = nil
		nd.children = append(nd.children, child.num)
		nd.children = append(nd.children, backup...)
		child.parent = nd.num
	case "beforeend":
		nd.children = append(nd.children, child.num)
		child.parent = nd.num
	case "afterend":
		var backup []int64
		parent := nd.Parent()
		backup = append(backup, parent.children...)
		parent.children = nil

		for i := range backup {
			parent.children = append(parent.children, backup[i])
			if backup[i] == nd.num {
				parent.children = append(parent.children, child.num)
				child.parent = parent.num
			}
		}
	default:
		return fmt.Errorf("(Error) position:'%v' \n", position)
	}
	return nil
}

func (nd *Node) InsertAdjacentHTML(position string, html string) error {
	err := nd.check()
	if err != nil {
		return err
	}

	doc := nd.earth
	frag := doc.CreateDocumentFragment()
	frag.InnerHTML(html)
	defer frag.remove(false)

	pos := strings.ToLower(position)
	switch pos {
	case "beforebegin":
		var backup []int64
		parent := nd.Parent()
		backup = append(backup, parent.children...)
		parent.children = nil

		for i := range backup {
			if backup[i] == nd.num {

				//frag.Children
				childs := frag.getChildren()
				for j := range childs {
					child := childs[j]
					parent.children = append(parent.children, child.num)
					child.parent = parent.num
				}
			}
			parent.children = append(parent.children, backup[i])
		}
	case "afterbegin":
		var backup []int64
		backup = append(backup, nd.children...)
		nd.children = nil

		//frag.Children
		childs := frag.getChildren()
		for j := range childs {
			child := childs[j]
			nd.children = append(nd.children, child.num)
			child.parent = nd.num
		}

		nd.children = append(nd.children, backup...)
	case "beforeend":
		//frag.Children
		childs := frag.getChildren()
		for j := range childs {
			child := childs[j]
			nd.children = append(nd.children, child.num)
			child.parent = nd.num
		}
	case "afterend":
		var backup []int64
		parent := nd.Parent()
		backup = append(backup, parent.children...)
		parent.children = nil

		for i := range backup {
			parent.children = append(parent.children, backup[i])
			if backup[i] == nd.num {

				//frag.Children
				childs := frag.getChildren()
				for j := range childs {
					child := childs[j]
					parent.children = append(parent.children, child.num)
					child.parent = parent.num
				}
			}
		}
	default:
		return fmt.Errorf("(Error) position:'%v' \n", position)
	}
	return nil
}

// Remove
func (nd *Node) removeChildAll() error {
	doc := nd.earth
	for i := range nd.children {
		num := nd.children[i]
		if child, ok := doc.flat[num]; ok {
			if child.name == "$index" {
				if nd.name == "$friend" {
					//nd.Friend.RemoveAll()
				}
			}
			child.removeChildAll()
			child.earth = nil
			child.attri = nil
			child.children = nil
			doc.flat[num] = nil
			delete(doc.flat, num)
		}
	}
	nd.children = nil
	return nil
}

func (nd *Node) remove(deep bool) error {
	if deep {
		nd.removeChildAll()
	}

	doc := nd.earth
	if _, ok := doc.flat[nd.num]; ok {
		nd.earth = nil
		nd.attri = nil
		nd.children = nil
		doc.flat[nd.num] = nil
		delete(doc.flat, nd.num)
	}
	nd = nil
	return nil
}

func (nd *Node) RemoveChild(child *Node) error {
	err := nd.check()
	if err != nil {
		return err
	}

	err = child.check()
	if err != nil {
		return err
	}

	var backup []int64
	backup = append(backup, nd.children...)
	slices.Sort(backup)

	nd.children = nil
	for i := range backup {
		if backup[i] != child.num {
			nd.children = append(nd.children, backup[i])
		}
	}

	child.remove(true)
	return nil
}

func (nd *Node) Remove() error {
	err := nd.check()
	if err != nil {
		return err
	}

	nd.remove(true)
	return nil
}

// Attribute
func (nd *Node) SetAttribute(key string, value any) error {
	err := nd.check()
	if err != nil {
		return err
	}

	if key, ok := checkAttr(key); !ok {
		log.Printf("(Error) key is invalid: %s", key)
		return fmt.Errorf("(Error) key is invalid: %s", key)
	}

	nd.attri[key] = value
	return nil
}

func (nd *Node) getAttribute(key string) any {
	err := nd.check()
	if err != nil {
		return ""
	}

	if v, ok := nd.attri[key]; ok {
		return any2str(v, false)
	}
	return ""
}

func (nd *Node) GetAttribute(key string) string {
	err := nd.check()
	if err != nil {
		return ""
	}

	if v, ok := nd.attri[key]; ok {
		return any2str(v, false)
	}
	return ""
}

func (nd *Node) GetAttributeNames() []string {
	err := nd.check()
	if err != nil {
		return nil
	}

	var list []string
	for k := range nd.attri {
		list = append(list, k)
	}
	return list
}

func (nd *Node) RemoveAttribute(key string) error {
	err := nd.check()
	if err != nil {
		return err
	}

	delete(nd.attri, key)
	return nil
}

// Reindex
func (nd *Node) Reindex() {
	var tagList map[string]int64
	tagList = make(map[string]int64)

	childs := nd.Children()
	for i := range childs {
		tag := childs[i].name
		if _, ok := tagList[tag]; ok {
			count := tagList[tag]
			count++
			tagList[tag] = count
		} else {
			tagList[tag] = 1
		}
	}

	for k := range tagList {
		if tagList[k] == 1 {
			tagList[k] = 0
		}
	}

	count := len(childs)
	for i := range childs {
		ii := count - i - 1
		tag := childs[ii].name
		xx := tagList[tag]
		childs[ii].index = xx
		xx--
		tagList[tag] = xx

		childs[i].Reindex()
	}
}

// GetElementBy
func (nd *Node) firstChild() *Node {
	doc := nd.earth
	count := len(nd.children)
	if count > 0 {
		num := nd.children[0]
		if node, ok := doc.flat[num]; ok {
			return node
		}
	}
	return nil
}

func (nd *Node) lastChild() *Node {
	doc := nd.earth
	count := len(nd.children)
	if count > 0 {
		i := count - 1
		num := nd.children[i]
		if node, ok := doc.flat[num]; ok {
			return node
		}
	}
	return nil
}

func (nd *Node) firstByTagName(tag string) *Node {
	doc := nd.earth
	for i := range nd.children {
		num := nd.children[i]
		if child, ok := doc.flat[num]; ok {
			if child.name == tag {
				return child
			}
		}
	}
	return nil
}

func (nd *Node) getNodeById(id string, idz []string) *Node {
	//inc 0:Node and Element All, 1:Element Only

	childs := nd.getChildren()
	for i := range childs {
		if childs[i].attrCheck(idz, id) {
			return childs[i]
		}

		node := childs[i].getNodeById(id, idz)
		if node != nil {
			return node
		}
	}
	return nil
}

func (nd *Node) GetNodeById(id string) *Node {
	err := nd.check()
	if err != nil {
		return nil
	}
	doc := nd.earth
	idz := doc.GetAlias("id", true)

	childs := nd.getChildren()
	for i := range childs {
		if childs[i].attrCheck(idz, id) {
			return childs[i]
		}

		node := childs[i].getNodeById(id, idz)
		if node != nil {
			return node
		}
	}
	return nil
}

func (nd *Node) getElementById(id string, idz []string) *Node {
	//inc 0:Node and Element All, 1:Element Only

	childs := nd.getChildren()
	for i := range childs {
		if !nd.nodeCheck(childs[i]) {
			if childs[i].attrCheck(idz, id) {
				return childs[i]
			}
		}

		node := childs[i].getElementById(id, idz)
		if node != nil {
			return node
		}
	}
	return nil
}

func (nd *Node) GetElementById(id string) *Node {
	err := nd.check()
	if err != nil {
		return nil
	}
	doc := nd.earth
	idz := doc.GetAlias("id", true)

	childs := nd.getChildren()
	for i := range childs {
		if !nd.nodeCheck(childs[i]) {
			if childs[i].attrCheck(idz, id) {
				return childs[i]
			}
		}

		node := childs[i].getElementById(id, idz)
		if node != nil {
			return node
		}
	}
	return nil
}

func (nd *Node) getNodesByTagName(tag string) []*Node {
	var list []*Node

	childs := nd.getChildren()
	for i := range childs {
		if childs[i].name == tag {
			list = append(list, childs[i])
		}

		chList := childs[i].getNodesByTagName(tag)
		list = append(list, chList...)
	}
	return list
}

func (nd *Node) GetNodesByTagName(tag string) []*Node {
	err := nd.check()
	if err != nil {
		return nil
	}
	var list []*Node

	if len(tag) == 0 {
		return nil
	}

	childs := nd.getChildren()
	for i := range childs {
		if childs[i].name == tag {
			list = append(list, childs[i])
		}

		chList := childs[i].getNodesByTagName(tag)
		list = append(list, chList...)
	}
	return list
}

func (nd *Node) getElementsByTagName(tag string) []*Node {
	var list []*Node

	childs := nd.getChildren()
	for i := range childs {
		if childs[i].name == tag {
			list = append(list, childs[i])
		}

		chList := childs[i].getElementsByTagName(tag)
		list = append(list, chList...)
	}
	return list
}

func (nd *Node) GetElementsByTagName(tag string) []*Node {
	err := nd.check()
	if err != nil {
		return nil
	}
	var list []*Node

	if len(tag) > 0 {
		if tag[:1] == "$" {
			return nil
		}
	} else {
		return nil
	}

	childs := nd.getChildren()
	for i := range childs {
		if childs[i].name == tag {
			list = append(list, childs[i])
		}

		chList := childs[i].getElementsByTagName(tag)
		list = append(list, chList...)
	}
	return list
}

func (nd *Node) getNodesByAttribute(keys []string, find any, check bool) []*Node {
	var list []*Node
	childs := nd.getChildren()
	for i := range childs {
		if !nd.nodeCheck(childs[i]) {
			for j := range keys {
				k := keys[j]
				if v, ok := childs[i].attri[k]; ok {
					if anyComp(v, find) {
						list = append(list, childs[i])
					}
				}
			}
		}

		chList := childs[i].getNodesByAttribute(keys, find, check)
		if len(chList) > 0 {
			list = append(list, chList...)
		}
	}
	return list
}

func (nd *Node) GetElementsByAttribute(key string, find any) []*Node {
	err := nd.check()
	if err != nil {
		return nil
	}

	doc := nd.earth
	keys := doc.GetAlias(key, true)

	var list []*Node
	childs := nd.getChildren()
	for i := range childs {
		if !nd.nodeCheck(childs[i]) {
			for j := range keys {
				k := keys[j]
				if v, ok := childs[i].attri[k]; ok {
					if anyComp(v, find) {
						list = append(list, childs[i])
					}
				}
			}
		}

		chList := childs[i].getNodesByAttribute(keys, find, true)
		if len(chList) > 0 {
			list = append(list, chList...)
		}
	}
	return list
}

func (nd *Node) GetNodesByAttribute(key string, find any) []*Node {
	err := nd.check()
	if err != nil {
		return nil
	}

	doc := nd.earth
	keys := doc.GetAlias(key, true)

	var list []*Node
	childs := nd.getChildren()
	for i := range childs {
		for j := range keys {
			k := keys[j]
			if v, ok := childs[i].attri[k]; ok {
				if anyComp(v, find) {
					list = append(list, childs[i])
				}
			}
		}

		chList := childs[i].getNodesByAttribute(keys, find, false)
		if len(chList) > 0 {
			list = append(list, chList...)
		}
	}
	return list
}

// TextContent, innerText, InnerHTML,  insertAdjacentHTML
func (nd *Node) setTextNode(text string) error {
	nd.removeChildAll()

	doc := nd.earth
	tx := doc.CreateTextNode(text)
	nd.AppendChild(tx)
	return nil
}

func (nd *Node) getTextNode() string {
	var text string
	childs := nd.getNodesByTagName("$text")
	for i := range childs {
		if v, ok := childs[i].attri["text"]; ok {
			switch x := v.(type) {
			case string:
				text = text + x
			}
		}
	}
	return text
}

func (nd *Node) TextContent(text ...string) string {
	err := nd.check()
	if err != nil {
		return ""
	}

	if len(text) > 0 {
		tx := text[0]
		tx = strings.ReplaceAll(tx, "\t", "")
		tx = strings.ReplaceAll(tx, "\n", "")
		tx = strings.ReplaceAll(tx, "\r", "")
		nd.setTextNode(tx)
	} else {
		return nd.getTextNode()
	}
	return ""
}

func (nd *Node) InnerText(text ...string) string {
	if len(text) > 0 {
		tx := text[0]
		nd.setTextNode(tx)
	} else {
		return nd.getTextNode()
	}
	return ""
}

func (nd *Node) getHTML(tagList []string) string {
	var text string

	childs := nd.Children()
	for i := range childs {
		closing := tagCheck(childs[i].name, tagList)

		tag := childs[i].name

		switch tag {
		case "$text":
			if v, ok := childs[i].attri["text"]; ok {
				text = text + any2str(v, false)
			}
		case "$comment":
			var kind string
			if v, ok := childs[i].attri["kind"]; ok {
				kind = any2str(v, true)
			}

			var comment string
			if v, ok := childs[i].attri["text"]; ok {
				comment = any2str(v, false)
			}

			switch kind {
			case "-":
				text = text + any2str(comment, false)
			case "!":
				text = text + any2str(comment, false)
			case "?":
				text = text + any2str(comment, false)
			}
		case "$space":
		case "$alias":
		case "$firend":
		case "$index":
		default:
			text = text + "<" + childs[i].name
			for k := range childs[i].attri {
				att := childs[i].attri[k]
				xx := any2str(att, true)
				text = text + " " + k + "=\"" + xx + "\""
			}

			if closing {
				text = text + ">"
			} else {
				text = text + ">"
				text = text + childs[i].getHTML(tagList)
				text = text + "</" + childs[i].name + ">"
			}
		}
	}
	return text
}

func (nd *Node) GetHTML() string {
	err := nd.check()
	if err != nil {
		return ""
	}

	tagList := SelfClosing
	slices.Sort(tagList)

	return nd.getHTML(tagList)
}

func (nd *Node) InnerHTML(text ...string) string {
	err := nd.check()
	if err != nil {
		return ""
	}

	if len(text) > 0 {
		tx := text[0]
		doc := nd.earth
		frag := doc.CreateDocumentFragment()
		frag.html2node(tx)

		// set .parent
		childs := frag.Children()
		for i := range childs {
			childs[i].parent = nd.num
		}

		nd.removeChildAll()
		nd.children = append(nd.children, frag.children...)
		frag.remove(false)
	} else {
		return nd.GetHTML()
	}
	return ""
}

// XPath
func (nd *Node) xpath(base *Node, style int) []string {
	var path []string
	if nd.num == 0 {
		log.Printf("(Error) parent is null")
		path = append(path, "?")
		return path
	}
	if nd.num == 1 {
		path = append(path, "")
		return path
	}

	if base != nil {
		if nd.num == base.num {
			switch style {
			case 0:

			case 1:
				var text string
				if nd.index > 0 {
					text = fmt.Sprintf("%v[%v]", nd.name, nd.index)
				} else {
					text = fmt.Sprintf("%v", nd.name)
				}
				path = append(path, text)
			case 2:
				path = append(path, "/")
			case 3:
				var text string
				if nd.index > 0 {
					text = fmt.Sprintf("%v[%v]", nd.name, nd.index)
				} else {
					text = fmt.Sprintf("%v", nd.name)
				}
				text = "//" + text
				path = append(path, text)
			}
			return path
		}
	}

	parent := nd.Parent()
	path = parent.xpath(base, style)

	var text string
	if nd.index > 0 {
		text = fmt.Sprintf("%v[%v]", nd.name, nd.index)
	} else {
		text = fmt.Sprintf("%v", nd.name)
	}
	path = append(path, text)

	return path
}

func (nd *Node) XPath(base *Node) string {
	//base is null: FullXPath

	//(ex.) full: /html/body/table/tr[2]/td[1]
	//      base: /html/body
	//   style 0:            table/tr[2]/td[1]
	//   style 1:       body/table/tr[2]/td[1]
	//   style 2:          //table/tr[2]/td[1]
	//   style 3:     //body/table/tr[2]/td[1]

	path := nd.xpath(base, 0)
	return strings.Join(path, "/")
}

func (nd *Node) getElementByXPath(xpaths []string, level int) *Node {
	var node *Node

	count := len(xpaths)
	if count < level {
		log.Printf("(Error) level:%v", level)
		return nil
	}

	tag := xpaths[level]
	var index int64
	if strings.Contains(tag, "[") {
		tag = strings.ReplaceAll(tag, "]", "")
		ts := strings.Split(tag, "[")
		tag = ts[0]
		x, err := strconv.Atoi(ts[1])
		if err != nil {
			log.Printf("(Error) not found XPath:%v", tag)
			return nil
		}
		index = int64(x)
	}

	childs := nd.getChildren()
	for i := range childs {
		child := childs[i]
		if child.name == tag {
			if child.index == index {
				x := level + 1
				if count == x {
					return child
				}
				node = child.getElementByXPath(xpaths, x)
			}
		}
	}
	return node
}

func (nd *Node) GetElementByXPath(xpath string) *Node {
	err := nd.check()
	if err != nil {
		return nil
	}
	var node *Node
	xpaths := strings.Split(xpath, "/")

	if xpath[:1] == "/" {
		// FullXPath
		doc := nd.earth
		root := doc.Root()
		node = root.getElementByXPath(xpaths, 1)
	} else {
		node = nd.getElementByXPath(xpaths, 0)
	}

	return node
}

// xml
func (nd *Node) getXML(ns map[string]string) string {
	var text string

	childs := nd.Children()
	for i := range childs {
		tag := childs[i].name
		if len(childs[i].space) > 0 {
			space := childs[i].space
			if short, ok := ns[space]; ok {
				if len(short) > 0 {
					tag = short + ":" + tag
				}
			}
		}

		switch tag {
		case "$text":
			if v, ok := childs[i].attri["text"]; ok {
				text = text + any2str(v, true)
			}
		case "$comment":
			var kind string
			if v, ok := childs[i].attri["kind"]; ok {
				kind = any2str(v, true)
			}

			var comment string
			if v, ok := childs[i].attri["text"]; ok {
				comment = any2str(v, true)
			}

			switch kind {
			case "-":
				text = text + comment
			case "!":
				text = text + comment
			case "?":
				text = text + "<?xml " + comment + "?>"
			}
		case "$space":
		case "$alias":
		case "$friend":
		case "$index":
		case "$link":
		default:
			text = text + "<" + tag
			for k := range childs[i].attri {
				kk := k
				att := childs[i].attri[kk]
				xx := any2str(att, true)

				if strings.Contains(kk, ";") {
					sn := strings.Split(kk, ";")
					space := sn[0]
					name := sn[1]

					if short, ok := ns[space]; ok {
						if len(short) == 0 {
							kk = name
						} else {
							kk = short + ":" + name
						}
					} else {
						kk = space + ":" + name
					}
				}

				text = text + " " + kk + "=\"" + xx + "\""
			}

			if len(childs[i].children) == 0 {
				text = text + "/>"
			} else {
				text = text + ">"
				text = text + childs[i].getXML(ns)
				text = text + "</" + tag + ">"
			}
		}
	}
	return text
}

func (nd *Node) GetXML() string {
	err := nd.check()
	if err != nil {
		return ""
	}
	var text string

	// get nameSpace
	var ns map[string]string
	ns = make(map[string]string)
	node := nd.firstByTagName("$space")
	if node != nil {
		for k := range node.attri {
			v := node.attri[k]
			x := any2str(v, false)
			ns[k] = x
		}
	}

	text = nd.getXML(ns)
	return text
}

func (nd *Node) InnerXML(text ...string) string {
	err := nd.check()
	if err != nil {
		return ""
	}

	if len(text) > 0 {
		tx := text[0]
		doc := nd.earth
		frag := doc.CreateDocumentFragment()
		frag.xml2node(tx)

		// set .parent
		childs := frag.Children()
		for i := range childs {
			childs[i].parent = nd.num
		}

		nd.removeChildAll()
		nd.children = append(nd.children, frag.children...)
		frag.remove(false)
	} else {
		return nd.GetXML()
	}
	return ""
}

// HTML to Node
func tagCheck(tag string, tagList []string) bool {
	tag = strings.ToLower(tag)
	for i := range tagList {
		if tagList[i] == tag {
			return true
		}
	}
	return false
}

func (nd *Node) html2node(htmlData string) error {
	tagList := SelfClosing
	slices.Sort(tagList)

	doc := nd.earth
	elm := nd
	tokenizer := html.NewTokenizer(strings.NewReader(htmlData))
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()
		switch tokenType {
		case html.DoctypeToken:
			text := token.String()
			child := doc.CreateComment(text)
			elm.AppendChild(child)
		case html.StartTagToken, html.SelfClosingTagToken:
			text := token.Data
			text = str2tag(text)
			child := doc.CreateElement(text)

			for i := range token.Attr {
				space := token.Attr[i].Namespace
				key := token.Attr[i].Key
				val := token.Attr[i].Val
				if len(space) > 0 {
					key = space + ":" + key
				}
				child.attri[key] = val
			}

			elm.AppendChild(child)
			if !tagCheck(text, tagList) {
				elm = child
			}
		case html.EndTagToken:
			elm = elm.Parent()
		case html.TextToken:
			text := token.String()
			text = str2tag(text)
			if len(text) > 0 {
				child := doc.CreateTextNode(text)
				elm.AppendChild(child)
			}
		case html.CommentToken:
			text := token.String()
			child := doc.CreateComment(text)
			elm.AppendChild(child)
		case html.ErrorToken:
			fmt.Printf("html.ErrorToken:%v \n", token.String())
		default:
			fmt.Printf("html.Token: '%v' \n", token)
			continue
		}
	}
	return nil
}

// Clone Node
func (nd *Node) cloneNode(parent *Node) {
	doc := nd.earth

	childs := nd.Children()
	for i := range childs {
		newChild := doc.CreateElementNS(childs[i].space, "dummy")
		newChild.name = childs[i].name
		newChild.index = childs[i].index
		for k := range childs[i].attri {
			newChild.attri[k] = childs[i].attri[k]
		}
		parent.AppendChild(newChild)

		childs[i].cloneNode(newChild)
	}
}

func (nd *Node) CloneNode(deep bool) *Node {
	if len(nd.children) == 0 {
		return nil
	}
	doc := nd.earth

	parent := doc.CreateElementNS(nd.space, "dummy")
	parent.name = nd.name

	childs := nd.Children()
	for i := range childs {
		newChild := doc.CreateElementNS(childs[i].space, "dummy")
		newChild.name = childs[i].name
		newChild.index = childs[i].index
		for k := range childs[i].attri {
			newChild.attri[k] = childs[i].attri[k]
		}
		parent.AppendChild(newChild)

		if deep {
			childs[i].cloneNode(newChild)
		}
	}
	return parent
}

// XML to Node
func (nd *Node) xml2node(htmlData string) error {
	doc := nd.earth

	var ns map[string]string
	ns = make(map[string]string)

	elm := nd
	btyeData := []byte(htmlData)
	decoder := xml.NewDecoder(bytes.NewReader(btyeData))
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch x := token.(type) {
		case xml.Directive:
			xx := string(x)
			node := doc.CreateComment(xx)
			node.attri["kind"] = "!"
			elm.AppendChild(node)
		case xml.ProcInst:
			xx := string(x.Inst)
			node := doc.CreateComment(xx)
			node.attri["kind"] = "?"
			elm.AppendChild(node)
		case xml.Comment:
			xx := string(x)
			node := doc.CreateComment(xx)
			elm.AppendChild(node)
		case xml.StartElement:
			if len(x.Name.Space) > 0 {
				if _, ok := ns[x.Name.Space]; !ok {
					ns[x.Name.Space] = ""
				}
			}

			node := doc.CreateElementNS(x.Name.Space, "dummy")
			node.name = x.Name.Local
			for i := range x.Attr {
				space := x.Attr[i].Name.Space
				local := x.Attr[i].Name.Local
				value := x.Attr[i].Value
				value = str2tag(value)

				var key string
				if len(space) > 0 {
					if space == "xmlns" {
						ns[value] = local
					}

					key = space + ";" + local
				} else {
					key = local
				}

				node.attri[key] = x.Attr[i].Value
			}
			elm.AppendChild(node)
			elm = node
		case xml.CharData:
			xx := string(x)
			xx = str2tag(xx)
			if len(xx) > 0 {
				node := doc.CreateTextNode(xx)
				elm.AppendChild(node)
			}
		case xml.EndElement:
			elm = elm.Parent()
		default:
			fmt.Printf("(Error) xml.Token '%v' \n", x)
			continue
		}
	}

	if len(ns) > 0 {
		node := doc.CreateElement("space")
		node.name = "$space"
		for k := range ns {
			node.attri[k] = ns[k]
		}
		nd.AppendChild(node)
	}
	return nil
}

// Import
func (nd *Node) Import(fileName string, kind string) error {
	err := nd.check()
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

	text := string(jsonData)
	switch kind {
	case ".html":
		nd.InnerHTML(text)
	case ".xml":
		nd.InnerXML(text)
	}
	return nil
}

// Export
func (nd *Node) Export(fileName string, kind string) error {
	err := nd.check()
	if err != nil {
		return err
	}

	fullPath, err := GetAbsolutePathName(fileName)
	if err != nil {
		return err
	}

	var text string
	switch kind {
	case ".html":
		text = nd.InnerHTML()
	case ".xml":
		text = nd.GetXML()
	}

	if len(text) == 0 {
		return fmt.Errorf("(Error) node is error num:%v \n", nd.num)
	}

	if FileExists(fullPath) {
		err = os.Remove(fullPath)
		if err != nil {
			return err
		}
	}

	var jsonData []byte
	jsonData = []byte(text)
	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
