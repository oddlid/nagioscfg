package nagioscfg

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"net"
	"sync"
	"time"
)

const epochStart = 122192928000000000

var (
	storageMutex  sync.Mutex
	storageOnce   sync.Once
	clockSequence uint16
	lastTime      uint64
	hardwareAddr  [6]byte
	epochFunc     = unixTimeFunc
)

func unixTimeFunc() uint64 {
	return epochStart + uint64(time.Now().UnixNano()/100)
}

func initClockSequence() {
	buf := make([]byte, 2)
	safeRandom(buf)
	clockSequence = binary.BigEndian.Uint16(buf)
}

func safeRandom(dest []byte) {
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
}

func initHardwareAddr() {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				copy(hardwareAddr[:], iface.HardwareAddr)
				return
			}
		}
	}

	// Initialize hardwareAddr randomly in case
	// of real network interfaces absence
	safeRandom(hardwareAddr[:])

	// Set multicast bit as recommended in RFC 4122
	hardwareAddr[0] |= 0x01
}

func initStorage() {
	initClockSequence()
	initHardwareAddr()
}

func getStorage() (uint64, uint16, []byte) {
	storageOnce.Do(initStorage)

	storageMutex.Lock()
	defer storageMutex.Unlock()

	timeNow := epochFunc()
	// Clock changed backwards since last UUID generation.
	// Should increase clock sequence.
	if timeNow <= lastTime {
		clockSequence++
	}
	lastTime = timeNow

	return timeNow, clockSequence, hardwareAddr[:]
}

func NewUUIDv1() UUID {
	u := UUID{}

	timeNow, clockSeq, hardwareAddr := getStorage()

	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	copy(u[10:], hardwareAddr)

	u.SetVersion(1)
	u.SetVariant()

	return u
}

func (u *UUID) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

func (u *UUID) SetVariant() {
	u[8] = (u[8] & 0xbf) | 0x80
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
