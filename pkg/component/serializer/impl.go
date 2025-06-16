package serializer

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"google.golang.org/protobuf/proto"
)

type SerializerType string

const (
	JSON     SerializerType = "json"
	XML      SerializerType = "xml"
	Protobuf SerializerType = "protobuf"
)

type _JSONSerializer struct{}

func (j _JSONSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j _JSONSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j _JSONSerializer) Type() SerializerType {
	return JSON
}

func (j _JSONSerializer) Clone() Serializer {
	return _JSONSerializer{}
}

type _XMLSerializer struct{}

func (X _XMLSerializer) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (X _XMLSerializer) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (X _XMLSerializer) Type() SerializerType {
	return XML
}

func (X _XMLSerializer) Clone() Serializer {
	return _XMLSerializer{}
}

type _ProtobufSerializer struct {
	strict bool
}

func (p _ProtobufSerializer) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok && p.strict {
		return nil, errors.New("value is not a proto.Message")
	}
	return proto.Marshal(msg)
}

func (p _ProtobufSerializer) Unmarshal(data []byte, v interface{}) error {

	msg, ok := v.(proto.Message)
	if !ok && p.strict {
		return errors.New("value is not a proto.Message")
	}
	return proto.Unmarshal(data, msg)
}

func (p _ProtobufSerializer) Type() SerializerType {
	return Protobuf
}

func (p _ProtobufSerializer) Clone() Serializer {
	return _ProtobufSerializer{
		strict: p.strict,
	}
}
