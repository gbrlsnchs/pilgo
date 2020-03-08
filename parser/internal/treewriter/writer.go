package treewriter

import (
	"fmt"
	"io"
)

// Writer is a buffer that prints a tree.
type Writer struct {
	w    io.Writer
	root Node
}

// NewWriter builds the tree and stores it in the Writer.
func NewWriter(w io.Writer, root Node) *Writer {
	return &Writer{w: w, root: root}
}

func (w *Writer) Write(prelude []byte) (int, error) {
	var (
		total, n int
		err      error
		ww       = w.w
	)
	if len(prelude) > 0 {
		n, err = fmt.Fprintln(ww, prelude)
		if err != nil {
			return total + n, err
		}
		total += n
	}
	n, err = fmt.Fprint(ww, ".")
	if err != nil {
		return total + n, err
	}
	if n, err := write(ww, w.root, make([]bool, 0)); err != nil {
		return total + n, err
	}
	total += n
	return total, nil
}

func write(w io.Writer, n Node, lastlist []bool) (int, error) {
	var total int
	for i, isLast := range lastlist {
		deepest := i == len(lastlist)-1
		if deepest {
			char := "├"
			if isLast {
				char = "└"
			}
			n, err := fmt.Fprintf(w, "%s──", char)
			if err != nil {
				return total + n, err
			}
			continue
		}
		char := "│"
		if isLast {
			char = " "
		}
		fmt.Fprintf(w, "%s   ", char)
	}
	nstr := fmt.Sprint(n)
	if nstr != "" {
		n, err := fmt.Fprintf(w, " %s", nstr)
		if err != nil {
			return total + n, err
		}
		total += n
	}
	nn, err := fmt.Fprintln(w)
	if err != nil {
		return total + nn, err
	}
	total += nn
	nlen := n.Len()
	for i := 0; i < nlen; i++ {
		isLast := i == nlen-1
		n, err := write(w, n.At(i), append(lastlist, isLast))
		if err != nil {
			return total + n, err
		}
		total += n
	}
	return total, nil
}
