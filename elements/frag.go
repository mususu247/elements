package elements

type frag struct {
	self *Node
}

func (nd *Node) Frag() *frag {
	err := nd.check()
	if err != nil {
		return nil
	}

	if nd.name != "$frag" {
		return nil
	}
	var fg frag
	fg.self = nd
	return &fg
}

func (fg *frag) Claer() error {
	fg.self.children = nil
	return nil
}

func (fg *frag) Add(node *Node) error {
	fg.self.children = append(fg.self.children, node.num)
	return nil
}

func (fg *frag) Import(list []*Node) error {
	fg.Claer()

	for i := range list {
		fg.self.children = append(fg.self.children, list[i].num)
	}
	return nil
}

func (fg *frag) Nothing() {
	fg.self = nil
	fg = nil
}
