package aln

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	. "github.com/hydra13142/bio/sequence"
)

func Read(r io.Reader) (ans []Sequence, err error) {
	buf := bufio.NewReader(r)
	mdi := make([][][]byte, 0)
	ans = make([]Sequence, 0)
	buf.ReadLine()
	for i, t := -1, true; ; {
		line, err := buf.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		line = bytes.TrimRight(line, " \r\n")
		if len(line) == 0 || line[0] == '\x20' {
			if i >= 0 {
				i, t = -1, false
			}
		} else {
			parts := bytes.Fields(line)
			if len(parts) != 2 {
				return nil, fmt.Errorf("Need exact two parts:%q", line)
			}
			i++
			if t {
				ans = append(ans, Sequence{string(parts[0]), Seq{}})
				mdi = append(mdi, [][]byte{parts[1]})
			} else {
				mdi[i] = append(mdi[i], parts[1])
			}
		}
		if err != nil {
			break
		}
	}
	for i, l := 0, len(ans); i < l; i++ {
		ans[i].Char = bytes.Join(mdi[i], nil)
	}
	return ans, nil
}

func Write(w io.Writer, seq []Sequence) error {
	if len(seq) < 2 {
		return fmt.Errorf("Need at least 2 sequences")
	}
	l := len(seq[0].Char)
	for i := 1; i < len(seq); i++ {
		if l != len(seq[i].Char) {
			return fmt.Errorf("Sequences' lengths are different")
		}
	}
	n := make([]string, len(seq))
	s := make([][]byte, len(seq))
	m := make([]byte, l)
	for i := 0; i < l; i++ {
		m[i] = '*'
		for j := 1; j < len(seq); j++ {
			if seq[0].Char[i] != seq[j].Char[i] {
				m[i] = ' '
				break
			}
		}
	}
	for i := 0; i < len(seq); i++ {
		if len(seq[i].Name) > 15 {
			n[i] = seq[i].Name[:15]
		} else {
			n[i] = seq[i].Name
		}
		s[i] = seq[i].Char
	}
	w.Write([]byte("CLUSTAL 2.1 multiple sequence alignment\r\n\r\n\r\n"))
	for {
		if len(m) > 60 {
			for i := 0; i < len(seq); i++ {
				fmt.Fprintf(w, "%-16s%s\r\n", n[i], s[i][:60])
				s[i] = s[i][60:]
			}
			fmt.Fprintf(w, "%-16s%s\r\n\r\n", " ", m[:60])
			m = m[60:]
		} else {
			for i := 0; i < len(seq); i++ {
				fmt.Fprintf(w, "%-16s%s\r\n", n[i], s[i])
			}
			fmt.Fprintf(w, "%-16s%s\r\n", " ", m)
			break
		}
	}
	return nil
}
