package model

import "io"

type FileStream interface {
	io.ReadCloser

	ContentType() string
	Size() int64
}

type fileStream struct {
	reader      io.ReadCloser
	contentType string
	size        int64
}

func NewFileStream(reader io.ReadCloser, contentType string, size int64) FileStream {
	return &fileStream{
		reader:      reader,
		contentType: contentType,
		size:        size,
	}
}

func (f *fileStream) Read(p []byte) (int, error) {
	return f.reader.Read(p)
}

func (f *fileStream) Close() error {
	return f.reader.Close()
}

func (f *fileStream) ContentType() string {
	return f.contentType
}

func (f *fileStream) Size() int64 {
	return f.size
}
