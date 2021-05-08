package runtime

import "encoding/json"

const (
	ContentTypeJSON = "application/json"
)

var (
	knownDecoders = map[string]Decoder{}
)

func init() {
	jsonD := newJsonDecoder()
	RegisterDecoder(jsonD)
}

type Decoder interface {
	// decoder name. Typically: ContentType
	Name() string
	Decode(data []byte, into Object) (err error)
}

func RegisterDecoder(d Decoder) {
	knownDecoders[d.Name()] = d
}

func NegotiateDecoder(mineType string) Decoder {
	return knownDecoders[mineType]
}

type jsonDecoder struct {
}

func newJsonDecoder() *jsonDecoder {
	return &jsonDecoder{}
}

func (d *jsonDecoder) Name() string {
	return ContentTypeJSON
}

func (d *jsonDecoder) Decode(data []byte, into Object) error {
	return json.Unmarshal(data, into)
}
