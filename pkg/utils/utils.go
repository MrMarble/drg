package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"reflect"
)

func ReadStruct(r io.Reader, structType interface{}) error {
	// get type from pointer
	v := reflect.Indirect(reflect.ValueOf(structType))

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()
		switch fieldType.Kind() {
		case reflect.Int32:
			field.SetInt(int64(ReadNextInt32(r)))
		case reflect.Uint32:
			field.SetUint(uint64(ReadNextUint32(r)))
		case reflect.Uint16:
			field.SetUint(uint64(ReadNextUint16(r)))
		case reflect.String:
			field.SetString(ReadNextString(r))
		case reflect.Struct:
			err := ReadStruct(r, field.Addr().Interface())
			if err != nil {
				return err
			}
		case reflect.Map:
			m := reflect.MapOf(fieldType.Key(), fieldType.Elem())
			field.Set(reflect.MakeMap(m))

			length := ReadNextUint32(r)
			for i := 0; i < int(length); i++ {
				key := ReadNextBytes(r, 16)
				value := ReadNextInt32(r)

				field.SetMapIndex(reflect.ValueOf(hex.EncodeToString(key)), reflect.ValueOf(value))
			}
		default:
			return fmt.Errorf("Unsupported type: %v", fieldType.Kind())
		}
	}
	return nil
}

func ReadNextBytes(r io.Reader, number int) []byte {
	bytes := make([]byte, number)

	_, err := r.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func ReadNextInt32(r io.Reader) int32 {
	var i int32

	binary.Read(r, binary.LittleEndian, &i)
	return i
}

func ReadNextUint32(r io.Reader) uint32 {
	var i uint32

	binary.Read(r, binary.LittleEndian, &i)
	return i
}

func ReadNextUint16(r io.Reader) uint16 {
	var i uint16
	binary.Read(r, binary.LittleEndian, &i)
	return i
}

func ReadNextString(r io.Reader) string {
	var buffer bytes.Buffer
	var length int32

	binary.Read(r, binary.LittleEndian, &length)
	if length > 65536 || length < 0 {
		log.Fatal("Invalid string length ", length)
	}
	if length == 0 {
		return ""
	}
	for i := 0; i < int(length); i++ {
		var b byte
		binary.Read(r, binary.LittleEndian, &b)
		buffer.WriteByte(b)

	}
	return buffer.String()[:length-1]
}

func ReadNextFloat32(r io.Reader) float32 {
	var f float32
	binary.Read(r, binary.LittleEndian, &f)
	return f
}

func Peek(r io.ReadSeeker, size int) []byte {
	bytes := make([]byte, size)
	_, err := r.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	_, err = r.Seek(-int64(size), io.SeekCurrent)

	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func ReadNextBool(r io.Reader) bool {
	var b byte
	binary.Read(r, binary.LittleEndian, &b)
	return b != 0
}

func ReadNextInt64(r io.Reader) int64 {
	var i int64
	binary.Read(r, binary.LittleEndian, &i)
	return i
}

func WriteString(w io.Writer, s string) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(s)+1))
	w.Write(buf)
	w.Write([]byte(s))
	w.Write([]byte{0}) // null byte
	return nil
}

func WriteInt64(w io.Writer, i int64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(i))
	w.Write(buf)
	return nil
}

func WriteInt32(w io.Writer, i int32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(i))
	w.Write(buf)
	return nil
}
