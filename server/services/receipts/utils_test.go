package receipts

import (
	"path/filepath"
	"runtime"
)

func getTestDataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(
		filepath.Dir(
			filepath.Dir(
				filepath.Dir(
					filepath.Dir(filename)))), "test", "data")
}
