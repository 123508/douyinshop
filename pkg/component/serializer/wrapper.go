package serializer

type SerializerWrapper struct {
	serializer Serializer
	Strategy   SerializerType
}

func NewSerializerWrapper(s SerializerType) *SerializerWrapper {
	return (&SerializerWrapper{
		Strategy: s,
	}).choose()
}

func (s *SerializerWrapper) choose() *SerializerWrapper {

	if s.serializer != nil {
		return s
	}

	switch s.Strategy {
	case JSON:
		s.serializer = _JSONSerializer{}
	case XML:
		s.serializer = _XMLSerializer{}
	case Protobuf:
		s.serializer = _ProtobufSerializer{}
	default:
		s.serializer = _JSONSerializer{}
	}
	return s
}

func (s *SerializerWrapper) SetSerializer(serial Serializer) {
	s.serializer = serial
}

func (s *SerializerWrapper) GetSerializer() Serializer {
	return s.serializer
}

// Serialize 业务数据序列化
func (s *SerializerWrapper) Serialize(v interface{}) ([]byte, error) {

	if s.serializer != nil {
		return s.serializer.Marshal(v)
	}

	return s.choose().serializer.Marshal(v)
}

// Deserialize 业务数据反序列化
func (s *SerializerWrapper) Deserialize(data []byte, v interface{}) error {
	if s.serializer != nil {
		return s.serializer.Unmarshal(data, v)
	}

	return s.choose().serializer.Unmarshal(data, v)
}

func (s *SerializerWrapper) Clone() *SerializerWrapper {

	if s.serializer == nil {
		s.choose()
	}

	return &SerializerWrapper{
		serializer: s.serializer.Clone(),
		Strategy:   s.Strategy,
	}
}
