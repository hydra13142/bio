package sequence

// 代表任意的未命名序列，可以是DNA、RNA、多肽或未知序列
type Seq struct {
	// 具体的序列数据，因为暴露在外，应尽量避免直接对其的操作
	Char []byte
	// 从低位算起，第0位表示方向，第1位表示是否有简并碱基，第2位表示是否有非法字符，其余表示序列类型
	kind int
}

// 代表任意的命名序列
type Sequence struct {
	// 序列的名称
	Name string
	// 实际序列数据
	Seq
}

// 创建一个正向（5'端到3'端，N端到C端）的序列，参数slice不经拷贝直接使用
func NewForwardSeq(seq []byte) *Seq {
	return &Seq{seq, 0}
}

// 创建一个反向（3'端到5'端，C端到N端）的序列，参数slice不经拷贝直接使用
func NewReverseSeq(seq []byte) *Seq {
	return &Seq{seq, 1}
}

// 将序列视为DNA并进行格式化，剔除空白字符、小写转为大写、识别非法字符，方法不申请新的slice
func (this *Seq) AsDNA() {
	s := this.Char
	l := len(s)
	i := 0
	if s[l-1] == '*' {
		l--
	}
	for j := 0; j < l; j++ {
		switch c := s[j]; c {
		case 'A', 'T', 'C', 'G', 'R', 'Y', 'M', 'K', 'S', 'W', 'H', 'B', 'V', 'D', 'N', 'X', '-':
			s[i] = c
		case 'a', 't', 'c', 'g', 'r', 'y', 'm', 'k', 's', 'w', 'h', 'b', 'v', 'd', 'n', 'x':
			s[i] = c - ('a' - 'A')
		case '\r', '\n', '\t', '\v', '\x20':
			continue
		default:
			s[i] = '?'
		}
		i++
	}
	*this = Seq{s[:i], this.kind&1 | 1<<3}
	for j := 0; j < i; j++ {
		if s[j] != 'A' && s[j] != 'T' && s[j] != 'C' && s[j] != 'G' {
			this.kind |= 2
			break
		}
	}
	for j := 0; j < i; j++ {
		if s[j] == '?' {
			this.kind |= 4
			break
		}
	}
}

// 将序列视为RNA并进行格式化，剔除空白字符、小写转为大写、识别非法字符，方法不申请新的slice
func (this *Seq) AsRNA() {
	s := this.Char
	l := len(s)
	i := 0
	if s[l-1] == '*' {
		l--
	}
	for j := 0; j < l; j++ {
		switch c := s[j]; c {
		case 'A', 'U', 'C', 'G', 'R', 'Y', 'M', 'K', 'S', 'W', 'H', 'B', 'V', 'D', 'N', 'X', '-':
			s[i] = c
		case 'a', 'u', 'c', 'g', 'r', 'y', 'm', 'k', 's', 'w', 'h', 'b', 'v', 'd', 'n', 'x':
			s[i] = c - ('a' - 'A')
		case '\r', '\n', '\t', '\v', '\x20':
			continue
		default:
			s[i] = '?'
		}
		i++
	}
	*this = Seq{s[:i], this.kind&1 | 2<<3}
	for j := 0; j < i; j++ {
		if s[j] != 'A' && s[j] != 'U' && s[j] != 'C' && s[j] != 'G' {
			this.kind |= 2
			break
		}
	}
	for j := 0; j < i; j++ {
		if s[j] == '?' {
			this.kind |= 4
			break
		}
	}
}

// 将序列视为多肽并进行格式化，剔除空白字符、小写转为大写、识别非法字符，方法不申请新的slice
func (this *Seq) AsPipetide() {
	s := this.Char
	l := len(s)
	i := 0
	for j := 0; j < l; j++ {
		switch c := s[j]; {
		case c >= 'a' && c <= 'z':
			c -= 'a' - 'A'
			fallthrough
		case c >= 'A' && c <= 'Z':
			switch c {
			case 'B', 'J', 'O', 'U', 'X', 'Z':
				s[i] = '?'
			default:
				s[i] = c
			}
		default:
			switch c {
			case '\r', '\n', '\t', '\v', '\x20':
				continue
			case '-', '*':
				s[i] = c
			default:
				s[i] = '?'
			}
		}
		i++
	}
	*this = Seq{s[:i], this.kind&1 | 3<<3}
	for j := 0; j < i; j++ {
		if s[j] == '?' {
			this.kind |= 4
			break
		}
	}
}

// 返回序列的类型：DNA、RNA、多肽或未知序列
func (this *Seq) Kind() string {
	switch this.kind >> 3 {
	case 0:
		return "Unknown" // 未知
	case 1:
		return "DNA"
	case 2:
		return "RNA"
	default:
		return "Peptide" // 多肽
	}
}

// 返回序列的方向，DNA和RNA是5'到3'或反过来；多肽是N端到C端或反过来；未知序列正向或反向
func (this *Seq) Direction() string {
	if this.kind&1 == 0 {
		switch this.kind >> 3 {
		case 0:
			return "Forward"
		case 1, 2:
			return "5 => 3"
		default:
			return "N => C"
		}
	} else {
		switch this.kind >> 3 {
		case 0:
			return "Reverse"
		case 1, 2:
			return "3 => 5"
		default:
			return "C => N"
		}
	}
}

// 返回序列的字符串表示
func (this *Seq) String() string {
	return string(this.Char)
}

// 返回序列是否包含简并碱基
func (this *Seq) Ambiguous() bool {
	return this.kind&2 != 0
}

// 返回序列是否包含非法字符
func (this *Seq) Illegal() bool {
	return this.kind&4 != 0
}

// 删除序列中的gap标识符'-'
func (this *Seq) DeleteGaps() *Seq {
	s := this.Char
	l := len(s)
	t := make([]byte, l)
	i := 0
	for _, c := range s {
		if c != '-' {
			t[i] = c
			i++
		}
	}
	return &Seq{t[:i], this.kind}
}

// 返回原序列的一个切片序列（切片进行了拷贝）
func (this *Seq) Slice(i, j int) *Seq {
	s := this.Char
	l := len(s)
	if j <= 0 {
		j += l
	}
	if i >= l || j < 0 || i >= j {
		return nil
	}
	t := make([]byte, j-i)
	copy(t, s[i:])
	return &Seq{t, this.kind}
}

// 返回反向序列
func (this *Seq) Reverse() *Seq {
	s := this.Char
	t := make([]byte, len(s))
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = s[j], s[i]
	}
	return &Seq{t, this.kind ^ 1}
}

// 返回互补序列，只适用于DNA、RNA；否则返回nil
func (this *Seq) Complement() *Seq {
	s := this.Char
	t := make([]byte, len(s))
	switch this.kind >> 3 {
	case 1:
		for i, c := range s {
			switch c {
			case 'A':
				t[i] = 'T'
			case 'T':
				t[i] = 'A'
			case 'C':
				t[i] = 'G'
			case 'G':
				t[i] = 'C'
			case 'R':
				t[i] = 'Y'
			case 'Y':
				t[i] = 'R'
			case 'M':
				t[i] = 'K'
			case 'K':
				t[i] = 'M'
			case 'H':
				t[i] = 'D'
			case 'B':
				t[i] = 'V'
			case 'V':
				t[i] = 'B'
			case 'D':
				t[i] = 'H'
			default: // 'S','W','N','-','*','?'
				t[i] = c
			}
		}
	case 2:
		for i, c := range s {
			switch c {
			case 'A':
				t[i] = 'U'
			case 'U':
				t[i] = 'A'
			case 'C':
				t[i] = 'G'
			case 'G':
				t[i] = 'C'
			case 'R':
				t[i] = 'Y'
			case 'Y':
				t[i] = 'R'
			case 'M':
				t[i] = 'K'
			case 'K':
				t[i] = 'M'
			case 'H':
				t[i] = 'D'
			case 'B':
				t[i] = 'V'
			case 'V':
				t[i] = 'B'
			case 'D':
				t[i] = 'H'
			default: // 'S','W','N','-','*','?'
				t[i] = c
			}
		}
	default:
		return nil
	}
	return &Seq{t, this.kind ^ 1}
}

// 返回反向互补序列，只适用于DNA、RNA；否则返回nil
func (this *Seq) ReverseComplement() *Seq {
	s := this.Char
	t := make([]byte, len(s))
	switch this.kind >> 3 {
	case 1:
		for i, j := 0, len(s)-1; j >= 0; i, j = i+1, j-1 {
			switch s[j] {
			case 'A':
				t[i] = 'T'
			case 'T':
				t[i] = 'A'
			case 'C':
				t[i] = 'G'
			case 'G':
				t[i] = 'C'
			case 'R':
				t[i] = 'Y'
			case 'Y':
				t[i] = 'R'
			case 'M':
				t[i] = 'K'
			case 'K':
				t[i] = 'M'
			case 'H':
				t[i] = 'D'
			case 'B':
				t[i] = 'V'
			case 'V':
				t[i] = 'B'
			case 'D':
				t[i] = 'H'
			default: // 'S','W','N','-','*','?'
				t[i] = s[j]
			}
		}
	case 2:
		for i, j := 0, len(s)-1; j >= 0; i, j = i+1, j-1 {
			switch s[j] {
			case 'A':
				t[i] = 'U'
			case 'U':
				t[i] = 'A'
			case 'C':
				t[i] = 'G'
			case 'G':
				t[i] = 'C'
			case 'R':
				t[i] = 'Y'
			case 'Y':
				t[i] = 'R'
			case 'M':
				t[i] = 'K'
			case 'K':
				t[i] = 'M'
			case 'H':
				t[i] = 'D'
			case 'B':
				t[i] = 'V'
			case 'V':
				t[i] = 'B'
			case 'D':
				t[i] = 'H'
			default: // 'S','W','N','-','*','?'
				t[i] = s[j]
			}
		}
	default:
		return nil
	}
	return &Seq{t, this.kind}
}

// 只适用于DNA（作为模板链），返回转录后的RNA；否则返回nil
func (this *Seq) Transcript() *Seq {
	if this.kind>>3 != 1 {
		return nil
	}
	s := this.Char
	t := make([]byte, len(s))
	for i, c := range s {
		switch c {
		case 'T':
			t[i] = 'U'
		default:
			t[i] = c
		}
	}
	return &Seq{t, (this.kind & 7) | (2 << 3)}
}

// 只适用于RNA，返回反转录后的DNA（作为模板链）；否则返回nil
func (this *Seq) ReverseTranscript() *Seq {
	if this.kind>>3 != 2 {
		return nil
	}
	s := this.Char
	t := make([]byte, len(s))
	for i, c := range s {
		switch c {
		case 'U':
			t[i] = 'T'
		default:
			t[i] = c
		}
	}
	return &Seq{t, (this.kind & 7) | (1 << 3)}
}

// 只适用于RNA，翻译后的多肽序列；否则返回nil
func (this *Seq) Translate() *Seq {
	if this.kind>>3 != 2 {
		return nil
	}
	s := this.Char
	l := len(s)
	t := make([]byte, l/3)
	if this.kind&1 == 0 {
		for i, j := 3, 0; i < l; i, j = i+3, j+1 {
			codon := string(s[i-3 : i])
			if c, ok := FromCodon[codon]; ok {
				t[j] = c
			} else {
				t[j] = '?'
			}
		}
	} else {
		for i, j := 3, 0; i < l; i, j = i+3, j+1 {
			codon := string([]byte{s[i-1], s[i-2], s[i-3]})
			if c, ok := FromCodon[codon]; ok {
				t[j] = c
			} else {
				t[j] = '?'
			}
		}
	}
	for _, c := range t {
		if c == '?' {
			return &Seq{t, (this.kind & 1) | 4 | (3 << 3)}
		}
	}
	return &Seq{t, (this.kind & 1) | (3 << 3)}
}
