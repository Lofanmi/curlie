package formatter

import (
	"bytes"
	"io"
)

// JSON is a writer that formats/colorizes JSON without decoding it.
// If the stream of bytes does not start with {, the formatting is disabled.
type JSON struct {
	Out              io.Writer
	Scheme           ColorScheme
	ParseJsonUnicode bool
	inited           bool
	disabled         bool
	last             byte
	lastQuote        byte
	isValue          bool
	level            int
	buf              []byte
}

var indent = []byte(`    `)

func (j *JSON) Write(p []byte) (n int, err error) {
	if !j.inited && len(p) > 0 {
		// Only JSON object are supported.
		j.disabled = (p[0] != '{' && p[0] != '[')
		j.inited = true
	}
	if j.disabled {
		return j.Out.Write(p)
	}

	// 检查是否包含 curl 的时间指标输出
	if bytes.Contains(p, []byte("TimingMetrics")) {
		return j.Out.Write(p)
	}

	cs := j.Scheme
	cp := j.buf
	for i := 0; i < len(p); i++ {
		b := p[i]
		if j.last == '\\' {
			cp = append(cp, b)
			j.last = b
			continue
		}
		switch b {
		case '\'', '"':
			switch j.lastQuote {
			case 0:
				j.lastQuote = b
				c := cs.Field
				if j.isValue {
					c = cs.Value
				}
				cp = append(cp, c...)
				cp = append(cp, b)
			case b:
				j.lastQuote = 0
				cp = append(cp, b)
				cp = append(cp, cs.Default...)
			default:
				cp = append(cp, b)
			}
		case '{', '[':
			j.isValue = false
			j.level++
			cp = append(cp, cs.Default...)
			cp = append(cp, b)
			cp = append(cp, '\n')
			cp = append(cp, bytes.Repeat(indent, j.level)...)
		case '}', ']':
			j.level--
			if j.level < 0 {
				j.level = 0
			}
			cp = append(cp, '\n')
			cp = append(cp, bytes.Repeat(indent, j.level)...)
			cp = append(cp, cs.Default...)
			cp = append(cp, b)
		case ':':
			j.isValue = true
			cp = append(cp, cs.Default...)
			cp = append(cp, b, ' ')
		case ',':
			j.isValue = false
			cp = append(cp, cs.Default...)
			cp = append(cp, b)
			cp = append(cp, '\n')
			cp = append(cp, bytes.Repeat(indent, j.level)...)
		default:
			if j.last == ':' {
				switch b {
				case 'n', 't', 'f':
					// null, true, false
					cp = append(cp, cs.Literal...)
				default:
					// unquoted values like numbers
					cp = append(cp, cs.Value...)
				}
			}
			cp = append(cp, b)
		}
		j.last = b
	}
	n, err = j.Out.Write(cp)
	if err != nil || n != len(cp) {
		return
	}
	return len(p), nil
}
