package db

import "google.golang.org/protobuf/reflect/protoreflect"

type Codec interface {
	Decode(kind protoreflect.Kind, data interface{}) (interface{}, error) // db结构转内存结构
	Encode(kind protoreflect.Kind, data interface{}) (interface{}, error) // 内存结构转db结构
}

type CacheDBCodec struct {
}

func (*CacheDBCodec) Decode(kind protoreflect.Kind, data interface{}) (interface{}, error) {
	return data, nil
}

func (*CacheDBCodec) Encode(kind protoreflect.Kind, data interface{}) (interface{}, error) {
	return data, nil
}
