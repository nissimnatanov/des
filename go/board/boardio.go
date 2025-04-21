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
	for row := range SequenceSize {
		if row != 0 && (row%3) == 0 {
			bw.WriteString("╠═══════╬═══════╬═══════╣\n")
			bw.Flush()
		}
		for col := range SequenceSize {
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
	for si := range SequenceSize {
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
	for i := range BoardSize {
		v := b.Get(i)
		if v == 0 {
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
	b := New().(*boardImpl)
	if err := deserializeInternal(s, &b.boardBase); err != nil {
		return nil, err
	}
	return b, nil
}

func deserializeInternal(s string, b boardBaseInternal) error {
	i := 0
	for _, c := range s {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}

		if i >= BoardSize {
			return fmt.Errorf("unexpected board character '%v', at board index: %v", c, i)
		}

		switch {
		case c >= 'a' && c <= 'z':
			i += int(1 + c - 'a')
		case c >= 'A' && c <= 'Z':
			i += int(1 + c - 'A')
		case c == '0':
			i++
		case c >= '1' && c <= '9':
			v := Value(c - '0')
			if v != 0 {
				b.setInternal(i, v, false)
			}
			i++
		case c == '_':
			v := b.Get(i)
			if v == 0 {
				return fmt.Errorf("misplaced _ sign")
			} else if !b.IsReadOnly(i) {
				return fmt.Errorf("duplicate _ sign")
			}
			b.setInternal(i, 0, false)
			b.setInternal(i, v, false)
		default:
			return fmt.Errorf("invalid board character %q, at board index: %v", c, i)
		}
	}

	if i != BoardSize {
		return fmt.Errorf("final board index is incorrect: %v", i)
	}

	return nil
}

// DeserializeSolution only accepts 81 digits of [1, 9], '_' to indicate originally editable cells and whitespaces
func DeserializeSolution(s string) (Solution, error) {
	var sol solutionImpl
	// start in edit mode
	sol.init(Edit)
	if err := deserializeInternal(s, &sol.boardBase); err != nil {
		return nil, err
	}

	err := sol.validateAndLock()
	if err != nil {
		return nil, err
	}
	return &sol, nil
}
