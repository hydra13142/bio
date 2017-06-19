package phy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	. "github.com/hydra13142/bio/sequence"
)

// 从phy文件读取序列数据
func Read(r io.Reader) (ans []Sequence, err error) {
	var t, l int
	_, err = fmt.Fscanf(r, "    %d    %d\r\n", &t, &l)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(r)
	mdi := make([][][]byte, t)
	ans = make([]Sequence, t)
	for i := 0; ; {
		line, err := buf.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		line = bytes.TrimRight(line, "\r\n\t\v\f\x20")
		if len(line) == 0 {
			i = 0
		} else {
			parts := bytes.Fields(line)
			switch line[0] {
			case '\x20', '\r', '\n', '\t', '\v':
			default:
				ans[i].Name = string(parts[0])
				parts = parts[1:]
			}
			mdi[i] = append(mdi[i], bytes.Join(parts, nil))
			i++
		}
		if err != nil {
			break
		}
	}
	for i := 0; i < t; i++ {
		ans[i].Char = bytes.Join(mdi[i], nil)
	}
	return ans, nil
}

// 将序列数据写入phy文件
func Write(w io.Writer, ans []Sequence) error {
	if len(ans) < 2 {
		return fmt.Errorf("Need at least 2 sequences")
	}
	x, t, l := 0, len(ans), len(ans[0].Char)
	fmt.Fprintf(w, "    %d    %d\r\n", t, l)

	for i := 0; i < t; i++ {
		name := ans[i].Name
		if len(name) > 10 {
			fmt.Fprintf(w, "%s", name[:10])
		} else {
			fmt.Fprintf(w, "%-10s", name)
		}
		for j, y := 0, x; j < 5; j = j + 1 {
			if y+10 < l {
				fmt.Fprintf(w, " %s", ans[i].Char[y:y+10])
				y += 10
			} else {
				fmt.Fprintf(w, " %-10s", ans[i].Char[y:l])
				break
			}
		}
		fmt.Fprintf(w, "\r\n")
	}
	for x += 50; x < l; x += 50 {
		fmt.Fprintf(w, "\r\n")
		for i := 0; i < t; i++ {
			fmt.Fprintf(w, "          ")
			for j, y := 0, x; j < 5; j = j + 1 {
				if y+10 < l {
					fmt.Fprintf(w, " %s", ans[i].Char[y:y+10])
					y += 10
				} else {
					fmt.Fprintf(w, " %-10s", ans[i].Char[y:l])
					break
				}
			}
			fmt.Fprintf(w, "\r\n")
		}
	}
	return nil
}
