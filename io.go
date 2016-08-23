package nagioscfg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func ReadKeyVal() {
}

func Read(r io.Reader) {
	scanner := bufio.NewScanner(r)
	var buf []byte
	for scanner.Scan() {
		buf = scanner.Bytes()
		if bytes.IndexByte(buf, '#') == 0 {
			//fmt.Printf("Seems we have a comment: %q\n", buf)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v+\n", err)
	}
}
