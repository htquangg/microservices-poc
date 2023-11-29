package serdes

import (
	"fmt"
	"reflect"

	"github.com/htquangg/microservices-poc/internal/registry"

	"google.golang.org/protobuf/proto"
)

type ProtoSerde struct {
	r registry.Registry
}

var (
	_      registry.Serde = (*ProtoSerde)(nil)
	protoT                = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

func NewProtoSerde(r registry.Registry) *ProtoSerde {
	return &ProtoSerde{r: r}
}

func (s ProtoSerde) Register(v registry.Registrable, options ...registry.BuildOption) error {
	if !reflect.TypeOf(v).Implements(protoT) {
		return fmt.Errorf("%T does not implement proto.Message", v)
	}
	return registry.Register(s.r, v, s.serialize, s.deserialize, options)
}

func (s ProtoSerde) RegisterKey(key string, v interface{}, options ...registry.BuildOption) error {
	if !reflect.TypeOf(v).Implements(protoT) {
		return fmt.Errorf("%T does not implement proto.Message", v)
	}
	return registry.RegisterKey(s.r, key, v, s.serialize, s.deserialize, options)
}

func (s ProtoSerde) RegisterFactory(key string, fn func() interface{}, options ...registry.BuildOption) error {
	v := fn()

	if v == nil {
		return fmt.Errorf("%s factory return a nil value", key)
	}

	if _, ok := v.(proto.Message); !ok {
		return fmt.Errorf("%s does not implement  proto.Message", key)
	}

	return registry.RegisterFactory(s.r, key, fn, s.serialize, s.deserialize, options)
}

func (s ProtoSerde) serialize(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (s ProtoSerde) deserialize(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
