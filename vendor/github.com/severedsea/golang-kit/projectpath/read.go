package projectpath

import (
	"io/ioutil"
)

// Read reads a file completely. But if the file does not exist, try to find it in the parent directory, [repeat...]
func Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(Abs(filename))
}
