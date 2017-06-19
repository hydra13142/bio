package crisper

import (
	"fmt"

	"github.com/hydra13142/bio/alignment"
	"github.com/hydra13142/bio/restriction"
	"github.com/hydra13142/bio/sequence"
	"github.com/hydra13142/sma"
)

// 表示一个探针
type Detector struct {
	site [2]int
	*sequence.Seq
}

// 实现fmt.Stringer接口
func (this Detector) String() string {
	return fmt.Sprintf("{(%d, %d) %s}", this.site[0], this.site[1], this.Char)
}

// 筛选CRISPER/CAS9探针，所得结果仍需要人工复核！
func CrisperCas9(target *sequence.Seq, genome []sequence.Sequence) []Detector {
	cut := SingleCutter(target)
	can := Candidate(target)
	mdi := []Detector{}
	for _, each := range can {
		if SuitCutter(each, cut) {
			mdi = append(mdi, each)
		}
	}
	can, mdi = mdi, []Detector{}
	for _, each := range can {
		if SuitGrna(each) <= 10 {
			mdi = append(mdi, each)
		}
	}
	can, mdi = mdi, []Detector{}
	for _, each := range can {
		if SuitGenome(each, genome) {
			mdi = append(mdi, each)
		}
	}
	return mdi
}

// 根据探针的规则，要求探针长度为20bp，探针3'末端对应的目标序列的右侧应有连续两个鸟嘌呤，来生成候选探针
func Candidate(s *sequence.Seq) (o []Detector) {
	if s.Kind() != "DNA" || s.Direction() != "5 => 3" {
		return nil
	}
	l := len(s.Char)
	for i := 20; i+2 < l; i++ {
		if string(s.Char[i:i+2]) == "GG" {
			o = append(o, Detector{[2]int{i - 10, i}, s.Slice(i-20, i)})
		}
	}
	for i := 2; i+20 < l; i++ {
		if string(s.Char[i-2:i]) == "CC" {
			o = append(o, Detector{[2]int{i, i + 10}, s.Slice(i, i+20).ReverseComplement()})
		}
	}
	return
}

// 找出在序列中只有一个酶切位点的酶以及切割位点
func SingleCutter(s *sequence.Seq) (o []restriction.Cutting) {
	cutters := restriction.FindCutters(s)
	for k, v := range cutters {
		if len(v) == 1 {
			o = append(o, restriction.Cutting{v[0], k})
		}
	}
	return
}

// 探针序列上，尤其是3'端10bp应包含酶切位点，且这个酶切位点应在序列上只存在很少几个最好只有一个
func SuitCutter(d Detector, cs []restriction.Cutting) bool {
	x := d.site
	for _, c := range cs {
		y := c.Fit
		if !(x[1] < y[0] || x[0] > y[1]) {
			return true
		}
	}
	return false
}

// 我们设计的探针连接在gRNA的5'端，对gRNA的3'端序列（互补后）进行匹配，如果分值高说明容易形成二级结构
func SuitGrna(d Detector) float64 {
	p := d.Char // q采用的pYAO的gRNA序列3'端部分的互补序列
	q := "AAAAAAAGCACCGACTCGGTGCCACTTTTTCAAGTTGATAACGGACTAGCCTTATTTTAACTTGCTATTTCTAGCTCTAAAAC"
	mtx := alignment.SmithWaterman(len(p), len(q), func(a, b int) float64 {
		if p[a] == q[b] {
			return 1
		}
		if p[a] == 'G' && q[b] == 'A' {
			return 0.5
		}
		return 0
	}, func(int) float64 {
		return 1
	}, func(i, _ int) float64 {
		return 0.81 + 0.02*float64(i)
	})
	return mtx.First()
}

// 探针长度20bp，要求3'端12bp的序列不能在基因组上有两个及以上的匹配（即存在错误匹配）
func SuitGenome(d Detector, genome []sequence.Sequence) bool {
	b1 := sma.NewBM(string(d.Char[8:]))
	b2 := sma.NewBM(string(d.Seq.ReverseComplement().Char[:12]))
	ct := 0
	for _, ch := range genome {
		if ct += len(b1.Find(ch.Char)); ct > 1 {
			return false
		}
		if ct += len(b2.Find(ch.Char)); ct > 1 {
			return false
		}
	}
	return ct == 1
}
