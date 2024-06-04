package marshaler

import "encoding/json"

type JsonMarshaler struct{}

func (JsonMarshaler) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (JsonMarshaler) Unmarshal(b []byte, i interface{}) error {
	return json.Unmarshal(b, i)
}
