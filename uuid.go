package nagioscfg

/*
This file is a very compressed and simplified version of:
https://github.com/satori/go.uuid/blob/master/uuid.go
adjusted to only cater for my specific needs here.
*/

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	sMutex   sync.Mutex
	sOnce    sync.Once
	clockSeq uint16
	lastTime uint64
	hwAddr   [6]byte
)

func NewUUIDv1() UUID {
	u := UUID{}

	sOnce.Do(func() {
		buf := make([]byte, 2)
		if _, err := rand.Read(buf); err != nil {
			panic(err)
		}
		clockSeq = binary.BigEndian.Uint16(buf)
		interfaces, err := net.Interfaces()
		if err == nil {
			for _, iface := range interfaces {
				if len(iface.HardwareAddr) >= 6 {
					copy(hwAddr[:], iface.HardwareAddr)
					return
				}
			}
		}
		hwAddr[0] |= 0x01
	})

	sMutex.Lock()
	timeNow := 122192928000000000 + uint64(time.Now().UnixNano()/100)
	// If clock changed backwards since last UUID generation,
	// we should increase clock sequence.
	if timeNow <= lastTime {
		clockSeq++
	}
	lastTime = timeNow
	sMutex.Unlock()

	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	copy(u[10:], hwAddr[:])

	u[6] = (u[6] & 0x0f) | (1 << 4) // set version 4
	u[8] = (u[8] & 0xbf) | 0x80     // set variant

	return u
}

func (u UUID) Equals(u2 UUID) bool {
	return bytes.Equal(u[:], u2[:])
}

// Key returns a string suitable for using as a key in a map, but not for human reading
func (u UUID) Key() string {
	return string(u[:])
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	const dash byte = '-'
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

// Bytes returns bytes slice representation of UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (u UUID) MarshalText() (text []byte, err error) {
	text = []byte(u.String())
	return
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
func (u *UUID) UnmarshalText(text []byte) (err error) {
	urnPrefix := []byte("urn:uuid:")
	byteGroups := []int{8, 4, 4, 4, 12}

	if len(text) < 32 {
		err = fmt.Errorf("uuid: UUID string too short: %s", text)
		return
	}

	t := text[:]
	braced := false

	if bytes.Equal(t[:9], urnPrefix) {
		t = t[9:]
	} else if t[0] == '{' {
		braced = true
		t = t[1:]
	}

	b := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			if t[0] != '-' {
				err = fmt.Errorf("uuid: invalid string format")
				return
			}
			t = t[1:]
		}

		if len(t) < byteGroup {
			err = fmt.Errorf("uuid: UUID string too short: %s", text)
			return
		}

		if i == 4 && len(t) > byteGroup &&
			((braced && t[byteGroup] != '}') || len(t[byteGroup:]) > 1 || !braced) {
			err = fmt.Errorf("uuid: UUID string too long: %s", text)
			return
		}

		_, err = hex.Decode(b[:byteGroup/2], t[:byteGroup])
		if err != nil {
			return
		}

		t = t[byteGroup:]
		b = b[byteGroup/2:]
	}

	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 16 {
		err = fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
		return
	}
	copy(u[:], data)

	return
}

// FromString returns UUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func UUIDFromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

func (u UUID) FromString(input string) error {
	u, err := UUIDFromString(input)
	return err
}
