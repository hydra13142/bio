package cladogram

// 表示一个进化树（的节点）
type Tree struct {
	Name  string  // 节点的名称
	Leaf  []Tree  // 节点的子节点
	Value float64 // 节点到父节点的枝长度
	Times int64   // 节点的置信率
}

// 用来生成进化树的文字表示时储存中间数据
type middle struct {
	deep int
	half int
	text []string
}

// 树的宽度，即叶节点个数
func (tr *Tree) Width() int {
	s := 0
	if len(tr.Leaf) != 0 {
		for _, o := range tr.Leaf {
			s += o.Width()
		}
		return s
	}
	return 1
}

// 树的深度，即从根节点到叶节点的最大长度
func (tr *Tree) Depth() float64 {
	var s float64
	if len(tr.Leaf) != 0 {
		for _, o := range tr.Leaf {
			f := o.Depth()
			if s < f {
				s = f
			}
		}
	}
	return s + tr.Value
}

// 树的非零最短枝的长度
func (tr *Tree) Basic() float64 {
	var s float64 = 10000
	if tr.Value != 0 {
		s = tr.Value
	}
	if len(tr.Leaf) != 0 {
		for _, o := range tr.Leaf {
			f := o.Basic()
			if s > f && f != 0 {
				s = f
			}
		}
	}
	return s
}

// 将树表示为文本图，字符串切片的每个字符串对应一行字符
func (this *Tree) Text() []string {
	md := this.text()
	return md.text
}

// 生成进化树的文本图
func (this *Tree) text() middle {
	if this.Leaf == nil {
		return middle{1, 0, []string{this.Name}}
	}
	child := make([]middle, len(this.Leaf))
	for i, p := range this.Leaf {
		child[i] = p.text()
	}
	t := 0
	for _, c := range child {
		if t < c.deep {
			t = c.deep
		}
	}
	t++
	for i, c := range child {
		s, l := "", ""
		for j := t - c.deep; j > 0; j-- {
			s += "    "
			l += "----"
		}
		s, l = s[1:], l[1:]
		for k, v := range c.text {
			if k == c.half {
				c.text[k] = l + v
			} else {
				c.text[k] = s + v
			}
		}
		child[i] = c
	}
	merge := child[0].text
	l := len(child[0].text)
	z := 0
	for _, c := range child[1:] {
		merge = append(merge, "")
		merge = append(merge, c.text...)
		z = len(c.text)
		l += 1 + z
		z -= c.half
	}
	x := child[0].half
	y := l - z
	z = (x + y) / 2
	for i := 0; i < l; i++ {
		switch {
		case i < x || i > y:
			merge[i] = " " + merge[i]
		case i == x || i == y:
			merge[i] = "+" + merge[i]
		default:
			merge[i] = "|" + merge[i]
		}
	}
	return middle{t, (x + y) / 2, merge}
}
