package linker_test

type testFileInfo struct {
	exists   bool
	linkname string
}

func (fi testFileInfo) Exists() bool     { return fi.exists }
func (fi testFileInfo) Linkname() string { return fi.linkname }
