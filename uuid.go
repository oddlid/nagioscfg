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

func NewUUIDv4() UUID {
	// I should probably find another way to generate IDs, as this seems unnecessary slow...
	u := UUID{}
	_, err := rand.Read(u[:]) // this step seems to be quite slow, like about 138 times slower than when the values are statically set...
	if err != nil {
		panic(err)
	}
	u[6] = (u[6] & 0x0f) | (4 << 4) // set version 4
	u[8] = (u[8] & 0xbf) | 0x80     // set variant
	return u
}

func (u UUID) Equals(u2 UUID) bool {
	return bytes.Equal(u[:], u2[:])
}
