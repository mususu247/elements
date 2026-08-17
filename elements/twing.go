package elements

type Twig struct {
	Num      int64
	Space    string
	Name     string
	Index    int64
	Attr     map[string]any
	Children []*Twig
}

func (nd *Node) toTwig() *Twig {
	var tw Twig
	tw.Num = nd.num
	tw.Space = nd.space
	tw.Name = nd.name
	tw.Index = nd.index
	tw.Attr = make(map[string]any)
	for k := range nd.attri {
		tw.Attr[k] = nd.attri[k]
	}

	childs := nd.getChildren()
	for i := range childs {
		twChild := childs[i].toTwig()
		tw.Children = append(tw.Children, twChild)
	}

	return &tw
}

func (nd *Node) fromTwig(tw *Twig, numList map[int64]int64) error {
	doc := nd.earth
	for i := range tw.Children {
		twChild := tw.Children[i]
		child := doc.CreateElementNS(twChild.Space, "dummy")
		child.name = twChild.Name
		child.index = twChild.Index

		for k := range child.attri {
			child.attri[k] = twChild.Attr[k]
		}
		nd.AppendChild(child)

		for k := range twChild.Attr {
			child.attri[k] = twChild.Attr[k]
		}

		child.fromTwig(twChild, numList)
		numList[twChild.Num] = child.num
	}
	return nil
}

func (tw *Twig) nothing() {
	for i := range tw.Children {
		tw.Children[i].nothing()
		tw.Children[i].Attr = nil
		tw.Children[i].Children = nil
		tw.Children[i] = nil
	}
}

func (tw *Twig) Nothing() {
	for i := range tw.Children {
		tw.Children[i].nothing()
		tw.Children[i].Attr = nil
		tw.Children[i].Children = nil
		tw.Children[i] = nil
	}
	tw.Attr = nil
	tw.Children = nil
	tw = nil
}
