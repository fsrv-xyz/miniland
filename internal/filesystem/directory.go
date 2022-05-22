package filesystem

import (
	"fmt"
	"os"
)

func Mkdir(path string, mode os.FileMode) error {
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("mkdir %v: %v", path, err)
	}
	return nil
}
