package board

import (
	"bufio"
	"fmt"
	"strings"
)

type write2log struct {
	log func(s string)
}

func (w2l *write2log) Write(p []byte) (int, error) {
	w2l.log(string(p))
	return len(p), nil
}

func NewWriter(log func(s string)) *bufio.Writer {
	// 1024 is large enough per board row or other log line, each line is flushed separately
	return bufio.NewWriterSize(&write2log{log}, 1024)
}

func Write(b Board, bw *bufio.Writer, fmt string) {
	if len(fmt) == 0 {
		fmt = "v"
	}

	for _, f := range fmt {
		switch f {
		case 'v', 'V':
			WriteValues(b, bw)
		case 'r', 'R':
			WriteRowSets(b, bw)
		case 'c', 'C':
			WriteColumnSets(b, bw)
		case 's', 'S':
			WriteSquareSets(b, bw)
		case 't', 'T':
			bw.WriteString("Serialized: ")
			writeSerialized(b, bw)
			bw.WriteRune('\n')
			bw.Flush()
		default:
			bw.WriteString("Unsupported format: ")
			bw.WriteRune(f)
			bw.WriteRune('\n')
			bw.Flush()
		}
	}
}

func WriteValues(b BoardBase, bw *bufio.Writer) {
	bw.WriteString("╔═══════╦═══════╦═══════╗\n")
	bw.Flush()
	for row := 0; row < SequenceSize; row++ {
		if row != 0 && (row%3) == 0 {
			bw.WriteString("╠═══════╬═══════╬═══════╣\n")
			bw.Flush()
		}
		for col := 0; col < SequenceSize; col++ {
			if col%3 == 0 {
				bw.WriteString("║ ")
			}
			i := IndexFromCoordinates(row, col)
			c := rune('0' + b.Get(i))
			bw.WriteRune(c)

			switch b := b.(type) {
			case Board:
				if b.IsReadOnly(i) {
					if b.IsValidCell(i) {
						c = ' '
					} else {
						c = 'X'
					}
				} else {
					if b.IsValidCell(i) {
						c = '.'
					} else {
						c = '!'
					}
				}
			default:
				if b.IsReadOnly(i) {
					c = ' '
				} else {
					c = '.'
				}
			}

			bw.WriteRune(c)
		}
		bw.WriteString("║\n")
		bw.Flush()
	}

	bw.WriteString("╚═══════╩═══════╩═══════╝\n")
	bw.Flush()
}

func WriteRowSets(b Board, bw *bufio.Writer) {
	bw.WriteString("Rows:")
	writeSets(func(row int) ValueSet { return b.RowSet(row) }, bw)
}

func WriteColumnSets(b Board, bw *bufio.Writer) {
	bw.WriteString("Columns:")
	writeSets(func(col int) ValueSet { return b.ColumnSet(col) }, bw)
}

func WriteSquareSets(b Board, bw *bufio.Writer) {
	bw.WriteString("Squares:")
	writeSets(func(square int) ValueSet { return b.SquareSet(square) }, bw)
}

func writeSets(fs func(si int) ValueSet, bw *bufio.Writer) {
	for si := 0; si < SequenceSize; si++ {
		bw.WriteString(" [")
		bw.WriteString(fs(si).String())
		bw.WriteRune(']')
	}
	bw.WriteRune('\n')
	bw.Flush()
}

func writeEmpty(count int, bw *bufio.Writer) {
	for count >= 27 {
		count -= 27
		bw.WriteRune('Z')
	}
	if count == 0 {
		return
	}

	if count == 1 {
		bw.WriteRune('0')
	} else {
		// A - 2 empty cells; B -> 3; ... Z - 27
		c := rune('A' + (count - 2))
		bw.WriteRune(c)
	}
}

func writeSerialized(b BoardBase, bw *bufio.Writer) {
	empty := 0
	for i := 0; i < BoardSize; i++ {
		v := b.Get(i)
		if v.IsEmpty() {
			empty++
		} else {
			writeEmpty(empty, bw)
			empty = 0
			c := rune('0' + v)
			bw.WriteRune(c)
			if !b.IsReadOnly(i) {
				// provided by the player
				bw.WriteRune('_')
			}
		}
	}
	writeEmpty(empty, bw)
}

func Serialize(b BoardBase) string {
	var sb strings.Builder
	bw := bufio.NewWriter(&sb)
	writeSerialized(b, bw)
	bw.Flush()
	return sb.String()
}

// Deserialize accepts more than one consecutive zero
// (while Serialize replaces them with letters, starting from 2)
func Deserialize(s string) (Board, error) {
	b := NewBoard()
	i := 0
	for c := range s {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}

		if c >= 'a' && c <= 'z' {
			i += 1 + c - 'a'
		} else if c >= 'A' && c <= 'Z' {
			i += 1 + c - 'A'
		} else if c == '0' {
			i++
		} else if c >= '0' && c <= '9' {
			v := Value(c - '0')
			if !v.IsEmpty() {
				b.SetReadOnly(i, v)
			}
			i++
		} else if c == '_' {
			v := b.Get(i)
			if v.IsEmpty() {
				return nil, fmt.Errorf("misplaced _ sign")
			} else if !b.IsReadOnly(i) {
				return nil, fmt.Errorf("duplicate _ sign")
			}
			b.Set(i, Empty)
			b.Set(i, v)
		} else {
			return nil, fmt.Errorf("invalid board character '%v', at index: %v", c, i)
		}
	}

	if i != BoardSize {
		return nil, fmt.Errorf("incomplete board, stopped at index: %v", i)
	}

	return b, nil
}

// DeserializeSolution only accepts exactly 81 digits of [1, 9], and whitespaces
func DeserializeSolution(s string) (Solution, error) {
	var sol solutionImpl
	// start in edit mode
	sol.init(Edit)
	i := 0
	for c := range s {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}

		if c >= '1' && c <= '9' {
			v := Value(c - '0')
			sol.setInternal(i, v, true)
			i++
		} else if c == '_' {
			v := sol.Get(i)
			if v.IsEmpty() {
				return nil, fmt.Errorf("misplaced _ sign")
			} else if !sol.IsReadOnly(i) {
				return nil, fmt.Errorf("duplicate _ sign")
			}
			sol.setInternal(i, Empty, false)
			sol.setInternal(i, v, false)
		} else {
			return nil, fmt.Errorf("invalid board character '%v', at index: %v", c, i)
		}
	}

	if i != BoardSize {
		return nil, fmt.Errorf("incomplete board, stopped at index: %v", i)
	}

	err := sol.validateAndLock()
	if err != nil {
		return nil, err
	}
	return &sol, nil
}
