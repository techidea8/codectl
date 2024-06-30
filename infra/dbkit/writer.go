package dbkit

import (
	"io"
	"os"
)

func FileWriter(path string) (file io.Writer, err error) {
	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	return
}

func StdWriter() io.Writer {
	return os.Stdout
}
