package restriction

import (
	"fmt"
	"strings"
)

// 将限制性内切酶的匹配序列生成一个字典树，用于匹配酶切位点。本类型为字典树的节点类型，同时也用来代表一个字典树
type Node struct {
	Name []string
	Leaf map[int]*Node
}

// 用来将碱基和简并碱基表示为二进制位的集合构成的整数
var baseBinary = map[byte]int{
	'A': 0x1, 'G': 0x2, 'C': 0x4, 'T': 0x8, 'R': 0x3,
	'Y': 0xC, 'M': 0x5, 'K': 0xA, 'S': 0x6, 'W': 0x9,
	'H': 0xD, 'B': 0xE, 'V': 0x7, 'D': 0xB, 'N': 0xF}

// 用来深层拷贝一棵树，对拷贝的任何操作都不影响原本的树
func (this *Node) Copy() *Node {
	that := new(Node)
	if this.Name != nil {
		that.Name = make([]string, len(this.Name))
		copy(that.Name, this.Name)
	}
	if this.Leaf != nil {
		that.Leaf = make(map[int]*Node, len(this.Leaf))
		for k, v := range this.Leaf {
			that.Leaf[k] = v.Copy()
		}
	}
	return that
}

// 像Trie树添加一个限制性内切酶（匹配序列，名称）
func (this *Node) Add(seq, name string) {
	if seq == "" {
		this.Name = append(this.Name, name)
		return
	}
	head := baseBinary[seq[0]]
	tail := seq[1:]
	if this.Leaf != nil {
		if leaf, ok := this.Leaf[head]; ok {
			leaf.Add(tail, name)
			return
		}
		ks := make([]int, 0, len(this.Leaf))
		for k := range this.Leaf {
			ks = append(ks, k)
		}
		for _, k := range ks {
			t := k & head
			if t != 0 {
				if k != t {
					v := this.Leaf[k]
					delete(this.Leaf, k)
					this.Leaf[k-t] = v.Copy()
					this.Leaf[t] = v
				}
				this.Leaf[t].Add(tail, name)
				head -= t
			}
		}
	} else {
		this.Leaf = make(map[int]*Node, 3)
	}
	if head != 0 {
		p := new(Node)
		p.Add(tail, name)
		this.Leaf[head] = p
	}
}

// Trie树测试一条序列并返回所有可以匹配的限制性内切酶
func (this *Node) Match(seq []byte) (name []string) {
	if len(seq) != 0 && this.Leaf != nil {
		head := baseBinary[seq[0]]
		tail := seq[1:]
		if leaf, ok := this.Leaf[head]; ok {
			return append(leaf.Match(tail), this.Name...)
		}
		for k := range this.Leaf {
			if k&head == head {
				return append(this.Leaf[k].Match(tail), this.Name...)
			}
		}
	}
	return this.Name
}

// 将Trie树表示为可被golang语言识别的文本格式
func (this *Node) GoString() string {
	name := "nil"
	if this.Name != nil {
		s := make([]string, 0, len(this.Name))
		for _, v := range this.Name {
			s = append(s, v)
		}
		name = `[]string{"` + strings.Join(s, `","`) + `"}`
	}
	leaf := "nil"
	if this.Leaf != nil {
		s := make([]string, 0, len(this.Leaf))
		for k, v := range this.Leaf {
			s = append(s, fmt.Sprintf("%#x:%s", k, v.GoString()))
		}
		leaf = `map[int]*Node{` + strings.Join(s, `,`) + `}`
	}
	return `&Node{` + name + `,` + leaf + `}`
}
