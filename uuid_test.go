package nagioscfg

import (
	"testing"
)

func BenchmarkNewUUIDv4(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		NewUUIDv4()
	}
}

func BenchmarkNewUUIDv1(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		NewUUIDv1()
	}
}
