package linker_test

type testFileInfo struct {
	exists bool
}

func (fi testFileInfo) Exists() bool { return fi.exists }
