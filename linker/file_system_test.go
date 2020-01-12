package linker_test

import "gsr.dev/pilgrim/linker"

type infoReturn struct {
	returnValue linker.FileInfo
	err         error
}

type readDirReturn struct {
	returnValue []string
	err         error
}

type testFileSystem struct {
	info    map[string]infoReturn
	readDir map[string]readDirReturn
}

func (fs testFileSystem) Info(name string) (linker.FileInfo, error) {
	return fs.info[name].returnValue, fs.info[name].err
}

func (fs testFileSystem) ReadDir(name string) ([]string, error) {
	return fs.readDir[name].returnValue, fs.readDir[name].err
}
