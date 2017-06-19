package tre

import (
	"errors"
	"fmt"
	"io"

	. "github.com/hydra13142/bio/cladogram"
	"github.com/hydra13142/pattern/token"
)

// cluster tree格式表示的进化树的“表述单元”
type Token struct {
	Kind  int
	Value interface{}
}

// 从io.Reader中读取一个cluster tree格式表示的树并解析为token序列
func scan(r io.Reader) (list []Token, err error) {
	word, _ := token.TokenDFA(`\w+`)
	scanner := token.NewScanner(r,
		func(data []byte, end bool) (i int, r interface{}) { // i<0 匹配失败，i==0匹配不完全，i>0匹配成功
			switch data[0] {
			case '(', ')', ',', ';':
				return 1, data[0]
			default:
				return -1, nil
			}
		},
		func(data []byte, end bool) (i int, r interface{}) { // :f
			if data[0] != ':' {
				return -1, nil
			}
			if len(data) <= 1 {
				return 0, nil
			}
			i, r = token.Float(data[1:], end)
			if i <= 0 {
				return i, nil
			}
			return i + 1, r
		},
		func(data []byte, end bool) (i int, r interface{}) { // [d]
			if data[0] != '[' {
				return -1, nil
			}
			if len(data) <= 1 {
				return 0, nil
			}
			i, r = token.Integer(data[1:], end)
			if i <= 0 {
				return i, nil
			}
			if len(data) <= 1+i {
				return 0, nil
			}
			if data[1+i] != ']' {
				return -1, nil
			}
			return i + 2, r
		}, word, token.Space)
	list = make([]Token, 0, 200)
	for scanner.Next() {
		i, r := scanner.Token()
		switch i {
		case 0: // 操作符
			switch r.(byte) {
			case '(':
				list = append(list, Token{0, '('})
			case ',':
				list = append(list, Token{4, ','})
			case ')':
				list = append(list, Token{5, ')'})
			case ';':
				list = append(list, Token{6, ';'})
			}
		case 1: // 浮点数表示枝长度
			list = append(list, Token{2, r})
		case 2: // 整数表示置信率（bootstrap中分枝合并的次数）
			list = append(list, Token{3, r})
		case 3: // 序列名称
			list = append(list, Token{1, string(r.([]byte))})
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return
}

// 将Token序列解析为一个进化树
func build(list []Token) (tkn Tree, err error) {
	medi := make([]Tree, 50)
	quot := make([]int, 50)
	t, i, j := 6, -1, -1
	for _, one := range list {
		switch t {
		case 0, 4, 6:
			switch one.Kind {
			case 0:
				j++
				quot[j] = i + 1
				t = 0
			case 1:
				i++
				medi[i] = Tree{Name: one.Value.(string)}
				t = 1
			default:
				err = errors.New("syntax error")
				return
			}
		case 1, 5:
			switch one.Kind {
			case 2:
				medi[i].Value = one.Value.(float64)
			case 3:
				medi[i].Times = one.Value.(int64)
			case 4:
				t = 4
			case 5:
				if j < 0 {
					err = errors.New("unfold quote")
					return
				}
				n := i + 1
				i = quot[j]
				l := make([]Tree, (n - i))
				copy(l, medi[i:n])
				medi[i] = Tree{Leaf: l}
				j--
				t = 5
			case 6:
				goto exit
			default:
				err = errors.New("syntax error")
				return
			}
		}
	}
exit:
	if j >= 0 {
		err = errors.New("unfold quote")
		return
	}
	if i != 0 {
		err = errors.New("not one trees")
		return
	}
	tkn = medi[0]
	return
}

// 从cluster tree格式的文件中读取进化树
func Read(r io.Reader) (Tree, error) {
	tkn, err := scan(r)
	if err != nil {
		return Tree{}, err
	}
	shu, err := build(tkn)
	if err != nil {
		return Tree{}, err
	}
	return shu, nil
}

// 将进化树写入cluster tree格式的文件中
func Write(w io.Writer, t Tree) (err error) {
	var write func(Tree) error
	write = func(t Tree) error {
		l := len(t.Leaf)
		if l != 0 {
			_, err = w.Write([]byte("("))
			if err != nil {
				return err
			}
			for i := 0; i < l-1; i++ {
				err = write(t.Leaf[i])
				if err != nil {
					return err
				}
				_, err = w.Write([]byte(","))
				if err != nil {
					return err
				}
			}
			err = write(t.Leaf[l-1])
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(w, "):%g", t.Value)
			if err != nil {
				return err
			}
			if t.Times != 0 {
				_, err = fmt.Fprintf(w, "[%d]", t.Times)
				if err != nil {
					return err
				}
			}
		} else if t.Name != "" {
			_, err = fmt.Fprintf(w, "%s:%g", t.Name, t.Value)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if err = write(t); err != nil {
		return err
	}
	_, err = w.Write([]byte{';'})
	return err
}
