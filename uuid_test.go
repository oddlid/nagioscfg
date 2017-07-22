/*
   Copyright 2017 Odd Eivind Ebbesen

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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

func TestUUIDFromString(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	s1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

	u1, err := UUIDFromString(s1)
	if err != nil {
		t.Errorf("Error parsing UUID from string: %s", err)
	}

	t.Logf("u:  %s", u)
	t.Logf("s1: %s", s1)
	t.Logf("u1: %s", u1)
}
