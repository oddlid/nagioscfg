package nagioscfg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

func IsComment(buf []byte) bool {
	// We might need to check for and skip whitespace before #, but for now, this will do
	return bytes.IndexByte(buf, '#') == 0
}

func IsBlankLine(buf []byte) bool {
	return unicode.IsSpace(buf[0]) // bogus, change me
}

func ReadKeyVal() {
}

// - check for comments and discard
// - check for blank lines and discard
// - check for beginning of a definition and enter a mode/set a flag
// - check for end of definition if "within-flag" set, unset flag if '}' encountered
// - split lines into key/value if within object definition
// - feed result to a consumer that creates object instances
func Read(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if IsComment(scanner.Bytes()) {
			//fmt.Printf("Seems we have a comment: %q\n", buf)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v+\n", err)
	}
}
