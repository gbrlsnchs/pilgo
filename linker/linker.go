package linker

import (
	"errors"
	"fmt"

	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/parser"
)

var (
	// ErrLinkExists means an unrelated file exists in place of a link.
	ErrLinkExists = errors.New("file exists in place of link")
	// ErrLinkNotExpands means a file in place of a link can't be expanded.
	ErrLinkNotExpands = errors.New("file exists in place of link and is not expandable")
	// ErrTargetNotExists means a target doesn't exist and thus can't be symlinked.
	ErrTargetNotExists = errors.New("target doesn't exist")
	// ErrTargetNotExpands means a target is not a directory and therefore can't be expanded.
	ErrTargetNotExpands = errors.New("target can't be expanded")
)

// Linker is a file symlinker.
type Linker struct {
	fs fs.FileSystem
}

// New creates a new linker with a given file system.
func New(fs fs.FileSystem) *Linker { return &Linker{fs} }

// Link creates every symlink needed in tr. Before creating any symlinks,
// it resolves nodes and checks for conflicts. If any conflicts or errors
// are found, it aborts the operation.
//
// Also, if needed, it creates parent directories if those don't already exist.
func (ln *Linker) Link(tr *parser.Tree) error {
	var (
		links   [][2]parser.File
		prepare = func(n *parser.Node) error {
			if err := ln.Resolve(n); err != nil {
				return err
			}
			if n.Status == parser.StatusReady {
				links = append(links, [2]parser.File{n.Target, n.Link})
			}
			return nil
		}
	)
	if err := tr.Walk(prepare); err != nil {
		return err
	}
	for _, link := range links {
		tgpath := link[0]
		lnpath := link[1]
		parent := lnpath.Dir()
		if err := ln.fs.MkdirAll(parent); err != nil {
			return err
		}
		if err := ln.fs.Symlink(tgpath.FullPath(), lnpath.FullPath()); err != nil {
			return err
		}
	}
	return nil
}

// Resolve checks and resolves nodes in a parsed tree.
func (ln *Linker) Resolve(n *parser.Node) error {
	tgpath := n.Target.FullPath()
	target, err := ln.fs.Stat(tgpath)
	if err != nil {
		return err
	}
	if !target.Exists() {
		n.Status = parser.StatusError
		return newLinkErr(tgpath, ErrTargetNotExists)
	}
	if len(n.Children) > 0 || len(n.Link.Path) == 0 {
		n.Status = parser.StatusSkip
		return nil
	}
	lnpath := n.Link.FullPath()
	link, err := ln.fs.Stat(lnpath)
	if err != nil {
		return err
	}
	if !link.Exists() {
		n.Status = parser.StatusReady
		return nil
	}
	if linkname := link.Linkname(); linkname != "" {
		if linkname == tgpath {
			n.Status = parser.StatusDone
			return nil
		}
		n.Status = parser.StatusConflict
		return newLinkErr(lnpath, ErrLinkExists)
	}
	if !target.IsDir() {
		n.Status = parser.StatusConflict
		return newLinkErr(tgpath, ErrTargetNotExpands)
	}
	if !link.IsDir() {
		n.Status = parser.StatusConflict
		return newLinkErr(lnpath, ErrLinkNotExpands)
	}
	children, err := ln.fs.ReadDir(tgpath)
	if err != nil {
		return err
	}
	expand(n, children)
	n.Status = parser.StatusExpand
	return nil
}

func expand(n *parser.Node, children []fs.FileInfo) {
	if len(children) <= 0 {
		return
	}
	n.Children = make([]*parser.Node, len(children))
	for i, c := range children {
		n.Children[i] = &parser.Node{
			Target: parser.File{
				BaseDir: n.Target.BaseDir,
				Path:    append(n.Target.Path, c.Name()),
			},
			Link: parser.File{
				BaseDir: n.Link.BaseDir,
				Path:    append(n.Link.Path, c.Name()),
			},
			Children: nil,
		}
	}
}

func newLinkErr(path string, err error) error { return fmt.Errorf("linker: %s: %w", path, err) }
