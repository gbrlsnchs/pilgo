package linker_test

type listReturn struct {
	returnValue []string
	err         error
}

type testFileInfo struct {
	exists   bool
	isDir    bool
	linkname string
	list     listReturn
}

func (fi testFileInfo) Exists() bool            { return fi.exists }
func (fi testFileInfo) IsDir() bool             { return fi.isDir }
func (fi testFileInfo) Linkname() string        { return fi.linkname }
func (fi testFileInfo) List() ([]string, error) { return fi.list.returnValue, fi.list.err }
