package main

import (
	"strings"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/gbrlsnchs/pilgo/fs"
)

type readMode struct {
	include cliutil.MultiValueOptionSet
	exclude cliutil.MultiValueOptionSet
	hidden  bool
}

func (md *readMode) resolve(files []fs.FileInfo) []string {
	eligible := make([]string, 0, len(files))
	for _, fi := range files {
		fname := fi.Name()
		if fname == "" || !md.hidden && strings.HasPrefix(fname, ".") {
			continue
		}
		if len(md.include) > 0 {
			if _, ok := md.include[fname]; !ok {
				continue
			}
		}
		if _, ok := md.exclude[fname]; ok {
			continue
		}
		eligible = append(eligible, fname)
	}
	return eligible
}
