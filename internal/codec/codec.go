package codec

import (
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

type Decoder interface {
	Decode(any) error
}

func NewDecoder(data io.Reader) Decoder {
	return msgpack.NewDecoder(data)
}

func Marshal(data any) ([]byte, error) {
	return msgpack.Marshal(data)
}
