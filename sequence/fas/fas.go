package fas

import (
	"bufio"
	"bytes"
	"io"

	. "github.com/hydra13142/bio/sequence"
)

func Read(r io.Reader) (ans []Sequence, err error) {
	buf := bufio.NewReader(r)
	mdi := make([][]byte, 0)
	ans = make([]Sequence, 0, 10)
	i := -1
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		line = bytes.TrimSpace(line)
		if len(line) != 0 {
			if line[0] == '>' {
				if i >= 0 {
					ans[i].Char = bytes.Join(mdi, nil)
					mdi = mdi[0:0]
				}
				ans = append(ans, Sequence{string(line[1:]), Seq{}})
				i++
			} else {
				mdi = append(mdi, line)
			}
		}
		if err != nil {
			break
		}
	}
	ans[i].Char = bytes.Join(mdi, nil)
	for ; i >= 0; i-- {
		l := len(ans[i].Char) - 1
		if l >= 0 && ans[i].Char[l] == '*' {
			ans[i].Char = ans[i].Char[:l]
		}
	}
	return ans, nil
}

func Write(w io.Writer, seq []Sequence) error {
	for i, l := 0, len(seq); i < l; i++ {
		_, err := w.Write([]byte(">" + seq[i].Name + "\r\n"))
		if err != nil {
			return err
		}
		_, err = w.Write(seq[i].Char)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte("\r\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
