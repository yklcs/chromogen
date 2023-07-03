package photo

import (
	"bytes"
	"errors"
	"io"
)

func (p *Photo) Open() error {
	p.buffer = &bytes.Buffer{}
	return nil
}

func (p *Photo) Write(b []byte) (int, error) {
	if p.buffer == nil {
		return 0, errors.New("attempted to write closed photo")
	}
	return p.buffer.Write(b)
}

func (p *Photo) ReadFrom(r io.Reader) (int64, error) {
	if p.buffer == nil {
		return 0, errors.New("attempted to write closed photo")
	}
	return p.buffer.ReadFrom(r)
}

type Reader struct {
	bytes.Reader
}

func NewReader(p Photo) (*Reader, error) {
	if p.buffer == nil {
		return nil, errors.New("attempted to create reader of closed photo")
	}
	return &Reader{
		Reader: *bytes.NewReader(p.buffer.Bytes()),
	}, nil
}

func (p *Photo) Close() error {
	p.buffer = nil
	return nil
}
