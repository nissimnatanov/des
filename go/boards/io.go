package boards

import (
	"fmt"
	"strings"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

func Format(b Board, fmt string) string {
	if len(fmt) == 0 {
		fmt = "v"
	}
	var sb strings.Builder
	for _, f := range fmt {
		switch f {
		case 'v', 'V':
			writeValues(b, &sb)
		case 't', 'T':
			sb.WriteString("Serialized: ")
			writeSerialized(b, false, &sb)
			sb.WriteByte('\n')
		default:
			sb.WriteString("Unsupported format: ")
			sb.WriteRune(f)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func writeValues(b Board, w writer) {
	w.WriteString("╔═══════╦═══════╦═══════╗\n")
	for row := range SequenceSize {
		if row != 0 && (row%3) == 0 {
			w.WriteString("╠═══════╬═══════╬═══════╣\n")
		}
		for col := range SequenceSize {
			if col%3 == 0 {
				w.WriteString("║ ")
			}
			i := indexes.IndexFromCoordinates(row, col)
			c := byte('0' + b.Get(i))
			w.WriteByte(c)

			if b.IsReadOnly(i) {
				if b.isValidCell(i) {
					c = ' '
				} else {
					c = 'X'
				}
			} else {
				if b.isValidCell(i) {
					c = '.'
				} else {
					c = '!'
				}
			}

			w.WriteByte(c)
		}
		w.WriteString("║\n")
	}

	w.WriteString("╚═══════╩═══════╩═══════╝\n")
}

func writeEmpty(count int, w writer) {
	for count >= 26 {
		count -= 26
		w.WriteByte('Z')
	}
	if count == 0 {
		return
	}

	// A -> 1 empty cell; B -> 2; ... Z -> 26
	c := byte('A' + (count - 1))
	w.WriteByte(c)
}

type writer interface {
	WriteByte(c byte) error
	WriteString(s string) (int, error)
}

func writeSerialized(b Board, asKey bool, w *strings.Builder) {
	empty := 0
	for i := range Size {
		v := b.Get(i)
		if v == 0 {
			disallowed := b.getDisallowedByUser(i)
			if disallowed.IsEmpty() {
				empty++
				continue
			}
			// we have to treat each disallowed value as a separate empty cell
			writeEmpty(empty, w)
			empty = 0
			// we do not need to write '0', presence of [] indicates that the cell is empty
			w.WriteByte('[')
			for _, d := range disallowed.Values() {
				w.WriteByte('0' + byte(d))
			}
			w.WriteByte(']')
			continue
		}

		writeEmpty(empty, w)
		empty = 0
		c := byte('0' + v)
		w.WriteByte(c)
		if !asKey && !b.IsReadOnly(i) {
			// provided by the player
			w.WriteByte('_')
		}
	}
	writeEmpty(empty, w)
}

type writerCounter int

func (wc *writerCounter) WriteByte(c byte) error {
	*wc++
	return nil
}

func (wc *writerCounter) WriteString(s string) (int, error) {
	*wc += writerCounter(len(s))
	return len(s), nil
}

func Serialize(b Board) string {
	var sb strings.Builder
	// regular serialization needs more space since it also marks read-write cells
	sb.Grow(Size * 2)
	writeSerialized(b, false, &sb)
	return sb.String()
}

func SerializeAsKey(b Board) string {
	// Serialize is also used for cache key calculation
	var sb strings.Builder
	// optimizing for the board generators here - they won't need more than 81 bytes
	sb.Grow(Size)
	writeSerialized(b, true, &sb)
	return sb.String()
}

// Deserialize accepts more than one consecutive zero
// (while Serialize replaces them with letters, starting from 2)
func Deserialize(s string) (*Game, error) {
	b := New()
	if err := deserializeBase(s, &b.base); err != nil {
		return nil, err
	}
	b.recalculateAllStats()
	return b, nil
}

func deserializeBase(s string, b *base) error {
	i := 0
	for _, c := range s {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}

		if i >= Size {
			return fmt.Errorf("unexpected board character '%v', at board index: %v", c, i)
		}

		switch {
		case c >= 'a' && c <= 'z':
			i += int(1 + c - 'a')
		case c >= 'A' && c <= 'Z':
			i += int(1 + c - 'A')
		case c == '0':
			// same as 'A'
			i++
		case c >= '1' && c <= '9':
			v := values.Value(c - '0')
			if v != 0 {
				b.setInternal(i, v, true)
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

	if i != Size {
		return fmt.Errorf("final board index is incorrect: %v", i)
	}

	return nil
}

// DeserializeSolution only accepts 81 digits of [1, 9], '_' to indicate originally editable cells and whitespaces
func DeserializeSolution(s string) (*Solution, error) {
	sol := &Solution{}
	// start in edit mode
	sol.base.init(Edit)
	if err := deserializeBase(s, &sol.base); err != nil {
		return nil, err
	}

	err := sol.validateAndLock()
	if err != nil {
		return nil, err
	}
	return sol, nil
}
