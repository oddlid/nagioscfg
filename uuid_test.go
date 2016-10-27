package nagioscfg

import (
	"testing"
)

//func BenchmarkNewUUIDv4(b *testing.B) {
//	for i := 0; i <= b.N; i++ {
//		NewUUIDv4()
//	}
//}

func BenchmarkNewUUIDv1(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		NewUUIDv1()
	}
}

func BenchmarkUUIDString(b *testing.B) {
	u := NewUUIDv1()
	for i := 0; i <= b.N; i++ {
		u.String()
	}
}
