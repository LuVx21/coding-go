package serialize

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func serialize(v any) ([]byte, error) {
	var msg proto.Message
	switch v := v.(type) {
	case nil:
		msg = structpb.NewNullValue()
	case bool:
		msg = wrapperspb.Bool(v)
	case int:
		msg = wrapperspb.Int64(int64(v))
	case int8:
		msg = wrapperspb.Int32(int32(v))
	case int16:
		msg = wrapperspb.Int32(int32(v))
	case int32:
		msg = wrapperspb.Int32(v)
	case int64:
		msg = wrapperspb.Int64(v)
	case uint:
		msg = wrapperspb.UInt64(uint64(v))
	case uint8:
		msg = wrapperspb.UInt32(uint32(v))
	case uint16:
		msg = wrapperspb.UInt32(uint32(v))
	case uint32:
		msg = wrapperspb.UInt32(v)
	case uint64:
		msg = wrapperspb.UInt64(v)
	case float32:
		msg = wrapperspb.Float(v)
	case float64:
		msg = wrapperspb.Double(v)
	case json.Number:
		n, err := v.Float64()
		if err != nil {
			msg = wrapperspb.String(string(v))
		} else {
			msg = wrapperspb.Double(n)
		}
	case string:
		msg = wrapperspb.String(v)
	case []byte:
		msg = wrapperspb.Bytes(v)
	case map[string]any:
		msg, _ = structpb.NewStruct(v)
	case []any:
		msg, _ = structpb.NewList(v)
	default:
		return nil, protoimpl.X.NewError("invalid type: %T", v)
	}

	return proto.Marshal(msg)
}

func Unmarshal(b []byte, m proto.Message) error {
	return proto.Unmarshal(b, m)
}
