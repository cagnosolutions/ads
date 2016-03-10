package tmp

import "os"

type data struct {
	path string
	size int
	used int
	file *os.File
	buff mmap
}

func newData(path string) (*data, error) {

	return nil, nil
}
