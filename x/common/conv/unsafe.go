package conv

import (
	"reflect"
	"unsafe"
)

// UnsafeStrToBytes converts string to byte slice without allocation.
// ⚠️ Use with caution: the returned slice must not be modified.
func UnsafeStrToBytes(s string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHeader := &reflect.SliceHeader{
		Data: strHeader.Data,
		Len:  strHeader.Len,
		Cap:  strHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(sliceHeader))
}

// UnsafeBytesToStr converts byte slice to string without allocation.
// ⚠️ Use with caution: the input slice must not be modified after conversion.
func UnsafeBytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
