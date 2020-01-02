package linker_test

import "gsr.dev/pilgrim/linker"

type infoReturn struct {
	returnValue linker.FileInfo
	err         error
}

type testFileSystem struct {
	info map[string]infoReturn
}

func (fs testFileSystem) Info(name string) (linker.FileInfo, error) {
	return fs.info[name].returnValue, fs.info[name].err
}
