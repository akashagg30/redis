package resp

import (
	"errors"
	"log"
	"strings"
	"time"
)

const WAITING_TIME_FOR_INPUT = 50 * time.Millisecond

type RESPDeserializer struct {
	byteReader *waitByteReader
	isClosed   bool
}

func (r *RESPDeserializer) readByte() (byte, bool) {
	b, err := r.byteReader.ReadByte()
	if err != nil {
		log.Println(err)
		return 0, false
	}
	if b == '\r' { // end of line
		r.byteReader.ReadByte() // skipping \n char
		return 0, false
	}
	return b, true
}

func (r *RESPDeserializer) deserializeAll() any {
	firstLiteral, ok := r.readByte()
	if !ok {
		return ""
	} else {
		switch firstLiteral {
		case '_':
			return r.deserializeNull()
		case '#':
			return r.deserializeBool()
		case '-':
			// log.Println("processing error")
			return r.deserializeError()
		case '+':
			// log.Println("processing simple string")
			return r.deserializeSimpleString()
		case '$':
			// log.Println("processing bulk string")
			return r.deserializeBulkString()
		case ':':
			// log.Println("processing integer")
			return r.deserializeInteger()
		case '*':
			// log.Println("processing array")
			return r.deserializeArray()
		default:
			println("some issue", string(firstLiteral))
		}

	}
	return ""
}

func (r *RESPDeserializer) deserializeError() error {
	return errors.New(r.deserializeSimpleString())
}

func (r *RESPDeserializer) deserializeNull() any {
	r.readByte()
	return nil
}

func (r *RESPDeserializer) deserializeBool() bool {
	defer r.readByte() // skipping \r\n
	b, _ := r.readByte()
	if b == 't' {
		return true
	} else {
		return false
	}
}

func (r *RESPDeserializer) deserializeSimpleString() string {
	var stringBuilder strings.Builder
	for {
		b, ok := r.readByte()
		if !ok {
			break
		}
		stringBuilder.WriteByte(b)
	}
	return stringBuilder.String()
}

func (r *RESPDeserializer) deserializeBulkString() string {
	length := int(r.deserializeInteger())
	if length == -1 {
		return ""
	}
	var stringBuilder strings.Builder
	for i := 0; i < length; i++ {
		b, _ := r.readByte()
		stringBuilder.WriteByte(b)
	}
	r.readByte() // skipping \r\n
	return stringBuilder.String()
}

func (r *RESPDeserializer) deserializeInteger() int64 {
	var data, signMuliplier int64
	signMuliplier = 1
	firstByte, _ := r.readByte()
	switch firstByte {
	case '+':
		signMuliplier = 1
	case '-':
		signMuliplier = -1
	default:
		data = int64(firstByte - '0')
	}
	for {
		b, ok := r.readByte()
		if !ok { // reached \r\n
			break
		}
		data = (data * 10) + int64(b-'0')
	}
	return data * signMuliplier
}

func (r *RESPDeserializer) deserializeArray() []any {
	length := int(r.deserializeInteger())
	if length <= 0 {
		return make([]any, 0)
	}
	data := make([]any, length)
	for i := 0; i < length; i++ {
		data[i] = r.deserializeAll()
	}
	return data
}

func (r *RESPDeserializer) Deserialize() (any, bool) {
	return r.deserializeAll(), !r.isClosed
}

func (r *RESPDeserializer) AddData(data []byte) {
	r.byteReader.writeNewData(data)
}

func (r *RESPDeserializer) Close() {
	r.byteReader.close()
	r.isClosed = true
}

func NewRESPDeserializer(data []byte) *RESPDeserializer {
	return &RESPDeserializer{byteReader: newWaitByteReader(data), isClosed: false}
}
