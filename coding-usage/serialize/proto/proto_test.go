package serialize

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/luvx21/coding-go/coding-usage/serialize/proto/proto_gen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// protoc --go_out=./proto_gen --go_opt=paths=source_relative person.proto
func Test_serialize_00(t *testing.T) {
	person := &proto_gen.Person{
		Name:    "Alice",
		Age:     30,
		Hobbies: []string{"reading", "hiking", "hiking", "hiking", "hiking"},
	}

	data, _ := proto.Marshal(person)
	fmt.Printf("原始Protobuf数据大小: %d bytes\n", len(data))

	// 3. 使用Zstd压缩
	compressed, _ := compressZstd(data)
	fmt.Printf("压缩后数据大小: %d bytes (%.1f%%)\n",
		len(compressed),
		float64(len(compressed))/float64(len(data))*100)

	// 4. 解压缩
	decompressed, _ := decompressZstd(compressed)

	newPerson := &proto_gen.Person{}
	_ = proto.Unmarshal(decompressed, newPerson)
	fmt.Printf("反序列化结果: %v\n", newPerson)
}

// 使用 JSON 转换（间接使用 protobuf）
func Test_serialize_01(t *testing.T) {
	data := map[string]any{
		"name": "foobar",
		"age":  25,
		"tags": []any{"go", "protobuf"},
		"ext": map[string]any{
			"foo": "bar",
		},
	}
	msg, _ := structpb.NewStruct(data)

	// ---------------------------------------------------------------------------
	protoData, _ := proto.Marshal(msg)
	fmt.Printf("Proto 二进制: %x\n", protoData)

	newMsg0 := &structpb.Struct{}
	_ = proto.Unmarshal(protoData, newMsg0)
	fmt.Printf("反序列化结果: name=%s, age=%.0f\n",
		newMsg0.Fields["name"].GetStringValue(),
		newMsg0.Fields["age"].GetNumberValue(),
	)

	// ---------------------------------------------------------------------------
	jsonData, _ := protojson.Marshal(msg)
	fmt.Printf("JSON 格式: %s\n", jsonData)

	newMsg := &structpb.Struct{}
	_ = protojson.Unmarshal(jsonData, newMsg)
	fmt.Printf("反序列化结果: name=%s, age=%.0f\n",
		newMsg.Fields["name"].GetStringValue(),
		newMsg.Fields["age"].GetNumberValue(),
	)

	// ---------------------------------------------------------------------------
	jsonBytes, _ := protojson.Marshal(newMsg0)

	var restoredData map[string]any
	_ = json.Unmarshal(jsonBytes, &restoredData)
	fmt.Println("反序列化结果:", restoredData)
}

func Test_serialize_02(t *testing.T) {
	data := map[string]any{
		"message": "Hello",
		"number":  42,
	}
	dynamicData, _ := structpb.NewStruct(data)
	// 2. 用 Any 包装动态数据
	anyData, _ := anypb.New(dynamicData)

	// 3. 序列化-反序列化
	protoData, _ := proto.Marshal(anyData)
	newAny := &anypb.Any{}
	_ = proto.Unmarshal(protoData, newAny)

	// 5. 从 Any 中提取原始数据
	unpackedData := &structpb.Struct{}
	_ = anypb.UnmarshalTo(newAny, unpackedData, proto.UnmarshalOptions{})

	fmt.Printf("获取数据: %v\n", unpackedData)
}

func Test_serialize_03(t *testing.T) {
	data := map[string]any{
		"name": "foobar",
		"age":  25,
		"tags": []any{"go", "protobuf"},
		"ext": map[string]any{
			"foo": "bar",
		},
	}
	bytes, _ := serialize(data)
	r := &structpb.Struct{}
	_ = proto.Unmarshal(bytes, r)

	jsonBytes, _ := protojson.Marshal(r)
	var restoredData map[string]any
	_ = json.Unmarshal(jsonBytes, &restoredData)
	fmt.Println("反序列化结果:", restoredData)
}
