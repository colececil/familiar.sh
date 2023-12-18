package test

import (
	"errors"
)

// FileDouble is an implementation of io.WriteCloser to use as a test double for os.File.
type FileDouble struct {
	path     string
	content  string
	isClosed bool
}

const fileAlreadyClosed = "file already closed"

// NewFileDouble returns a new instance of FileDouble with the given path.
func NewFileDouble(path string) *FileDouble {
	return &FileDouble{
		path: path,
	}
}

// Close implements io.Closer by setting the isClosed flag to true.
func (double *FileDouble) Close() error {
	if double.isClosed {
		return errors.New(fileAlreadyClosed)
	}
	double.isClosed = true
	return nil
}

// Write implements io.Writer by writing the given bytes to a string representing the file content.
func (double *FileDouble) Write(bytes []byte) (bytesWritten int, err error) {
	if double.isClosed {
		return 0, errors.New(fileAlreadyClosed)
	}
	double.content += string(bytes)
	return len(bytes), nil
}

// GetPath returns the path of the file.
func (double *FileDouble) GetPath() string {
	return double.path
}

// GetContent returns the string representing the file content.
func (double *FileDouble) GetContent() string {
	return double.content
}

// IsClosed returns whether the file is closed.
func (double *FileDouble) IsClosed() bool {
	return double.isClosed
}
