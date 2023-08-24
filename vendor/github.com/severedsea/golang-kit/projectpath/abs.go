package projectpath

import (
	"os"
	"path/filepath"
)

// Abs returns absolute path of a file in project directory
// on errors, we return the requested `filename` for simplicity
// since eventually using that file (read/write) would error in a more convenient place to handle
// i.e. those read/write methods has `error` included in return value
func Abs(filename string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return filename
	}

	return _abs(cwd, filename)
}

// the reason we split into `_abs` is because the `filename`
// could be a relative path, e.g. `config/secret.yml` and we want to keep the `config/` hierarchy while traversing
// the parent directory of `dirname` recursively
func _abs(dirname, filename string) string {
	fullpath, err := filepath.Abs(dirname + "/" + filename)
	if err != nil {
		return filename
	}

	if _, err = os.Stat(fullpath); err != nil {
		parentdir := filepath.Dir(dirname)
		if parentdir == "/" {
			return filename
		}
		return _abs(parentdir, filename)
	}

	return fullpath
}
