package cladogram

import . "github.com/hydra13142/bio/Sequence"

// 本函数用于统计进化树中各个节点的核苷酸突变的次数
// 假设所有突变都发生在树分枝的位置且无回复突变
// 返回各个位置发生突变的总计构成的序列
func Mutation(pnt Tree, seq []Sequence) []int {
	var (
		calc func(Tree, int) (int, int)
		dict = map[string][]byte{}
	)
	for _, s := range seq {
		dict[s.Name] = s.Char
	}
	calc = func(pnt Tree, num int) (int, int) {
		var x, y, z int
		if pnt.Leaf == nil {
			x := dict[pnt.Name]
			switch x[num] {
			case 'A':
				return 0x01, 0
			case 'G':
				return 0x02, 0
			case 'C':
				return 0x04, 0
			case 'T', 'U':
				return 0x08, 0
			default:
				return 0x10, 0
			}
		}
		for _, p := range pnt.Leaf {
			a, b := calc(p, num)
			x &= a
			y |= a
			z += b
		}
		if x == 0 {
			return y, z
		}
		t := ((x >> 0) & 1) + ((x >> 1) & 1) + ((x >> 2) & 1) + ((x >> 3) & 1) + ((x >> 4) & 1)
		return y &^ x, z + t
	}
	l := len(seq[0].Char)
	f := make([]int, l)
	for i := 0; i < l; i++ {
		x, z := calc(pnt, i)
		t := ((x >> 0) & 1) + ((x >> 1) & 1) + ((x >> 2) & 1) + ((x >> 3) & 1) + ((x >> 4) & 1)
		f[i] = z + t - 1
	}
	return f
}
