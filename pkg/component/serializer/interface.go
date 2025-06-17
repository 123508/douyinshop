package serializer

type Serializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	Type() SerializerType
	Clone() Serializer
}
