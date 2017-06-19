// 序列比对的动态规划算法，将匹配值函数、gap罚分函数和位置补正函数抽象出来，以便应用于多种不同格式的序列数据
package alignment

import "sort"

// 动态规划算法中，生成的矩阵里每一个方格的数据信息
type Cell struct {
	sum float64
	ori [2]int
	dir int
}

// 全局序列比对Needleman-Wunsch动态规划算法生成的矩阵
type MatrixNW [][]Cell

// 局部序列比对SmithWaterman动态规划算法生成的矩阵
type MatrixSW [][]Cell

// 全局序列比对的动态规划算法，mate返回两个序列指定位点的匹配值，nick是gap的罚分函数，coef为序列位置的修正系数
func NeedlemanWunsch(H, W int, mate func(int, int) float64, nick func(int) float64, coef func(int, int) float64) MatrixNW {
	if mate == nil {
		return nil
	}
	if nick == nil {
		nick = DefaultNick
	}
	if coef == nil {
		coef = DefaultCoef
	}
	matrix := make([][]Cell, H)
	for i := H - 1; i >= 0; i-- {
		matrix[i] = make([]Cell, W)
	}
	for j := 0; j < W; j++ {
		matrix[0][j] = Cell{coef(0, j) * mate(0, j), [2]int{0, j}, 0}
	}
	for i := 1; i < H; i++ {
		matrix[i][0] = Cell{coef(i, 0) * mate(i, 0), [2]int{i, 0}, 0}
	}
	for i := 1; i < H; i++ {
		for j := 1; j < W; j++ {
			k := coef(i, j)
			x := &matrix[i][j-1]
			y := &matrix[i-1][j-1]
			z := &matrix[i-1][j]
			h := x.sum - k*nick(x.dir)
			s := y.sum + k*mate(i, j)
			v := z.sum - k*nick(z.dir)
			switch {
			case h > s && h > v:
				matrix[i][j] = Cell{h, x.ori, x.dir + 1}
			case v > s:
				matrix[i][j] = Cell{v, z.ori, z.dir - 1}
			default:
				matrix[i][j] = Cell{s, y.ori, 0}
			}
		}
	}
	return matrix
}

// 返回最大匹配分值
func (this MatrixNW) First() float64 {
	H, W := len(this), len(this[0])
	max := float64(0)
	for i := 0; i < H; i++ {
		if score := this[i][W-1].sum; max < score {
			max = score
		}
	}
	for j := 0; j < W; j++ {
		if score := this[H-1][j].sum; max < score {
			max = score
		}
	}
	return max
}

// 返回前n高的匹配分值，这些分值的匹配都有不同的起点和终点
func (this MatrixNW) Limit(n int) []float64 {
	H, W := len(this), len(this[0])
	S := map[[2]int]float64{}
	for i := 0; i < H; i++ {
		p := &this[i][W-1]
		if t, ok := S[p.ori]; !ok || t < p.sum {
			S[p.ori] = p.sum
		}
	}
	for j := 0; j < W; j++ {
		p := &this[H-1][j]
		if t, ok := S[p.ori]; !ok || t < p.sum {
			S[p.ori] = p.sum
		}
	}
	L := []float64{}
	for _, v := range S {
		L = append(L, v)
	}
	sort.Float64s(L)
	if l := len(L); l > n {
		L = L[l-n : l]
	}
	for i, j := 0, len(L)-1; i < j; i, j = i+1, j-1 {
		L[i], L[j] = L[j], L[i]
	}
	return L
}

// 返回最高匹配分值下，匹配的起始点和对齐信息序列，1表示序列一移动1字符而序列二添加空位，2表示相反情况，3表示两个序列都移动1字符
func (this MatrixNW) Settle() ([2]int, []byte) {
	H, W := len(this), len(this[0])
	x, y, max := 0, 0, float64(0)
	for i := 0; i < H; i++ {
		if score := this[i][W-1].sum; max < score {
			max, x, y = score, i, W-1
		}
	}
	for j := 0; j < W; j++ {
		if score := this[H-1][j].sum; max < score {
			max, x, y = score, H-1, j
		}
	}
	s := make([]byte, 0, H+W)
	for k := H - 1; k > x; k-- {
		s = append(s, 1)
	}
	for k := W - 1; k > y; k-- {
		s = append(s, 2)
	}
	for x >= 0 && y >= 0 {
		switch t := this[x][y].dir; {
		case t > 0:
			for ; t > 0; t, y = t-1, y-1 {
				s = append(s, 2)
			}
		case t < 0:
			for ; t < 0; t, x = t+1, x-1 {
				s = append(s, 1)
			}
		default:
			s = append(s, 3)
			x, y = x-1, y-1
		}
	}
	for ; x >= 0; x-- {
		s = append(s, 1)
	}
	for ; y >= 0; y-- {
		s = append(s, 2)
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return [2]int{0, 0}, s
}

// 局部序列比对的动态规划算法，mate返回两个序列指定位点的匹配值，nick是gap的罚分函数，coef为序列位置的修正系数
func SmithWaterman(H, W int, mate func(int, int) float64, nick func(int) float64, coef func(int, int) float64) MatrixSW {
	if mate == nil {
		return nil
	}
	if nick == nil {
		nick = DefaultNick
	}
	if coef == nil {
		coef = DefaultCoef
	}
	matrix := make([][]Cell, H)
	for i := H - 1; i >= 0; i-- {
		matrix[i] = make([]Cell, W)
	}
	for j := 0; j < W; j++ {
		matrix[0][j] = Cell{coef(0, j) * mate(0, j), [2]int{0, j}, 0}
	}
	for i := 1; i < H; i++ {
		matrix[i][0] = Cell{coef(i, 0) * mate(i, 0), [2]int{i, 0}, 0}
	}
	for i := 1; i < H; i++ {
		for j := 1; j < W; j++ {
			k := coef(i, j)
			x := &matrix[i][j-1]
			y := &matrix[i-1][j-1]
			z := &matrix[i-1][j]
			h := x.sum - k*nick(x.dir)
			s := y.sum + k*mate(i, j)
			v := z.sum - k*nick(z.dir)
			o := k * mate(i, j)
			switch {
			case h > s && h > v && h > o:
				matrix[i][j] = Cell{h, x.ori, x.dir + 1}
			case v > s && v > o:
				matrix[i][j] = Cell{v, z.ori, z.dir - 1}
			case o > s:
				matrix[i][j] = Cell{o, [2]int{i, j}, 0}
			default:
				matrix[i][j] = Cell{s, y.ori, 0}
			}
		}
	}
	return matrix
}

// 返回最大匹配分值
func (this MatrixSW) First() float64 {
	H, W := len(this), len(this[0])
	max := float64(0)
	for i := 0; i < H; i++ {
		for j := 0; j < W; j++ {
			if score := this[i][j].sum; max < score {
				max = score
			}
		}
	}
	return max
}

// 返回前n高的匹配分值，这些分值的匹配都有不同的起点和终点
func (this MatrixSW) Limit(n int) []float64 {
	H, W := len(this), len(this[0])
	S := map[[2]int]float64{}
	for i := 0; i < H; i++ {
		for j := 0; j < W; j++ {
			p := &this[i][j]
			if t, ok := S[p.ori]; !ok || t < p.sum {
				S[p.ori] = p.sum
			}
		}
	}
	L := []float64{}
	for _, v := range S {
		L = append(L, v)
	}
	sort.Float64s(L)
	if l := len(L); l > n {
		L = L[l-n : l]
	}
	for i, j := 0, len(L)-1; i < j; i, j = i+1, j-1 {
		L[i], L[j] = L[j], L[i]
	}
	return L
}

// 返回最高匹配分值下，匹配的起始点和对齐信息序列，1表示序列一移动1字符而序列二添加空位，2表示相反情况，3表示两个序列都移动1字符
func (this MatrixSW) Settle() ([2]int, []byte) {
	H, W := len(this), len(this[0])
	x, y, max := 0, 0, float64(0)
	for i := 0; i < H; i++ {
		if score := this[i][W-1].sum; max < score {
			max, x, y = score, i, W-1
		}
	}
	for j := 0; j < W; j++ {
		if score := this[H-1][j].sum; max < score {
			max, x, y = score, H-1, j
		}
	}
	s := make([]byte, 0, H+W)
	o := this[x][y].ori
	for {
		if t := this[x][y].dir; t > 0 {
			for ; t > 0; t, y = t-1, y-1 {
				s = append(s, 2)
			}
		} else if t < 0 {
			for ; t < 0; t, x = t+1, x-1 {
				s = append(s, 1)
			}
		} else {
			s = append(s, 3)
			if [2]int{x, y} == o {
				break
			}
			x, y = x-1, y-1
		}
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return o, s
}

// 默认的gap罚分函数
func DefaultNick(n int) float64 {
	switch {
	case n < 0:
		return 2 / float64(1-n)
	case n > 0:
		return 2 / float64(n+1)
	default:
		return 2
	}
}

// 默认的序列位置修正系数函数
func DefaultCoef(int, int) float64 {
	return 1
}

// 使用两条序列和字符匹配函数来生成返回指定位置匹配值的函数，string版本
func MatchString(p, q string, f func(byte, byte) float64) func(int, int) float64 {
	return func(i, j int) float64 {
		return f(p[i], q[j])
	}
}

// 使用两条序列和字符匹配函数来生成返回指定位置匹配值的函数，[]byte版本
func MatchBytes(p, q []byte, f func(byte, byte) float64) func(int, int) float64 {
	return func(i, j int) float64 {
		return f(p[i], q[j])
	}
}

//
func AlignString(o [2]int, s []byte, p, q string) (string, string) {
	i, j := o[0], o[1]
	m := make([]byte, len(s))
	n := make([]byte, len(s))
	for k, c := range s {
		switch c {
		case 1:
			m[k], n[k] = p[i], '-'
			i++
		case 2:
			m[k], n[k] = '-', q[j]
			j++
		default:
			m[k], n[k] = p[i], q[j]
			i, j = i+1, j+1
		}
	}
	return string(m), string(n)
}

//
func AlignBytes(o [2]int, s []byte, p, q []byte) ([]byte, []byte) {
	i, j := o[0], o[1]
	m := make([]byte, len(s))
	n := make([]byte, len(s))
	for k, c := range s {
		switch c {
		case 1:
			m[k], n[k] = p[i], '-'
			i++
		case 2:
			m[k], n[k] = '-', q[j]
			j++
		default:
			m[k], n[k] = p[i], q[j]
			i, j = i+1, j+1
		}
	}
	return m, n
}
