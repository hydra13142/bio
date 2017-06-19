package restriction

import "github.com/hydra13142/bio/sequence"

// 表示一次酶切反应，包含切割位置和切割酶两个信息
type Cutting struct {
	Site
	Cutter string
}

// 内部函数，用于将相对位置定位在绝对位置，此为正向
func inc(s Site, n int) Site {
	return Site{[2]int{n + s.Fit[0], n + s.Fit[1]}, [2]int{n + s.Cut[0], n + s.Cut[1]}}
}

// 内部函数，用于将相对位置定位在绝对位置，此为确定反向的位置
func dec(n int, s Site) Site {
	return Site{[2]int{n - s.Fit[1], n - s.Fit[0]}, [2]int{n - s.Cut[1], n - s.Cut[0]}}
}

// 搜索酶切位点，返回所有可能的切割反应构成的slice
func FindSites(seq *sequence.Seq) []Cutting {
	if seq.Kind() != "DNA" {
		return nil
	}
	l := len(seq.Char)
	if seq.Direction() == "3 => 5" {
		s := FindSites(seq.Reverse())
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = Cutting{dec(l, s[j].Site), s[j].Cutter}, Cutting{dec(l, s[i].Site), s[i].Cutter}
		}
		return s
	}
	s := make([]Cutting, 0, 10)
	for i := 0; i < l; i++ {
		for _, t := range TwoWay.Match(seq.Char[i:]) {
			println(t)
			s = append(s, Cutting{inc(Cutters[t].Site, i), t})
		}
		for _, t := range OneWay.Match(seq.Char[i:]) {
			println(t)
			s = append(s, Cutting{inc(Cutters[t].Site, i), t})
		}
	}
	seq = seq.ReverseComplement()
	for i := 0; i < l; i++ {
		for _, t := range OneWay.Match(seq.Char[i:]) {
			println(t)
			s = append(s, Cutting{dec(l-i, Cutters[t].Site), t})
		}
	}
	return s
}

// 搜索酶切位点，返回一个map，键为酶的名称，而值为该酶可以切割的位点构成的slice
func FindCutters(seq *sequence.Seq) map[string][]Site {
	if seq.Kind() != "DNA" {
		return nil
	}
	l := len(seq.Char)
	if seq.Direction() == "3 => 5" {
		s := FindCutters(seq.Reverse())
		for k, v := range s {
			for i, t := range v {
				v[i] = dec(l, t)
			}
			s[k] = v
		}
		return s
	}
	s := make(map[string][]Site, 10)
	for i := 0; i < l; i++ {
		for _, t := range TwoWay.Match(seq.Char[i:]) {
			v, _ := s[t]
			s[t] = append(v, inc(Cutters[t].Site, i))
		}
		for _, t := range OneWay.Match(seq.Char[i:]) {
			v, _ := s[t]
			s[t] = append(v, inc(Cutters[t].Site, i))
		}
	}
	seq = seq.ReverseComplement()
	for i := 0; i < l; i++ {
		for _, t := range OneWay.Match(seq.Char[i:]) {
			v, _ := s[t]
			s[t] = append(v, dec(l-i, Cutters[t].Site))
		}
	}
	return s
}
