package resp

import (
	"bytes"
	"fmt"
	"log"
)

// RESPSerializer structure
type RESPSerializer struct {
	byteBuffer *bytes.Buffer
}

// NewRespSerializer creates a new RESPSerializer instance
func NewRespSerializer() *RESPSerializer {
	return &RESPSerializer{byteBuffer: &bytes.Buffer{}}
}

func (r *RESPSerializer) serializeAll(value any) {
	switch v := value.(type) {
	case nil:
		r.serializeNull()
	case bool:
		r.serializeBool(v)
	case string:
		r.serializeBulkString(v)
	case int:
		r.serializeInteger(int64(v))
	case int64:
		r.serializeInteger(v)
	case []any:
		r.serializeArray(v)
	case error:
		r.serializeError(v)
	default:
		log.Println("invalid type found for resp serialization")
		r.byteBuffer.WriteString("-ERR unknown type\r\n")
	}
}

func (r *RESPSerializer) serializeNull() {
	r.byteBuffer.WriteString("_\r\n")
}

func (r *RESPSerializer) serializeBool(value bool) {
	if value {
		r.byteBuffer.WriteString("#t\r\n")
	} else {
		r.byteBuffer.WriteString("#f\r\n")
	}
}

func (r *RESPSerializer) serializeError(err error) {
	r.byteBuffer.WriteString("-" + err.Error() + "\r\n")
}

func (r *RESPSerializer) serializeInteger(value int64) {
	r.byteBuffer.WriteString(":")
	r.writeIntegerWithDeliminatorToByte(value)
}

func (r *RESPSerializer) writeIntegerWithDeliminatorToByte(value int64) {
	r.byteBuffer.WriteString(fmt.Sprintf("%d\r\n", value))
}

func (r *RESPSerializer) serializeBulkString(value string) {
	r.byteBuffer.WriteString("$")
	r.writeIntegerWithDeliminatorToByte(int64(len(value)))
	r.byteBuffer.WriteString(value + "\r\n")
}

func (r *RESPSerializer) serializeArray(value []any) {
	r.byteBuffer.WriteString("*")
	lenght_of_passed_array := len(value)
	r.writeIntegerWithDeliminatorToByte(int64(lenght_of_passed_array))
	for i := 0; i < lenght_of_passed_array; i++ {
		// Recursively serialize each element in the array
		r.serializeAll(value[i])
	}
}

// GetSerializedData returns the serialized data as a byte slice
func (r *RESPSerializer) Serialize(data any) []byte {
	defer r.byteBuffer.Reset()
	r.serializeAll(data)
	return r.byteBuffer.Bytes()
}

func (r *RESPSerializer) SerializeSimpleString(value string) []byte {
	return []byte("+" + value + "\r\n")
}
