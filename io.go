package nagioscfg

import (
	"bufio"
	//"bytes"
	"fmt"
	"io"
	"unicode"
)


type Reader struct {
	Comment rune
	r       *bufio.Reader
	field   bytes.Buffer
	line    int
	column  int
}

// copy of https://golang.org/src/encoding/csv/reader.go
func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.r.ReadRune()
	if r1 == '\r' {
		r1, _, err = r.r.ReadRune()
		if err == nil {
			if r1 != '\n' {
				r.r.UnreadRune()
				r1 = '\r'
			}
		}
	}
	r.column++
	return r1, err
}

func (r *Reader) skip(delim rune) error {
	for {
		r1, err := r.readRune()
		if err != nil {
			return err
		}
		if r1 == delim {
			return nil
		}
	}
}

func (r *Reader) parseLine() {
	r.line++
	r.column = -1

	r1, _, err := r.r.ReadRune()

	if err != nil {
		//return nil, err
	}

	if r.Comment != 0 && r1 == r.Comment {
		//return nil, r.skip('\n')
	}
	r.r.UnreadRune()
	//...
}

// IsComment and IsBlankLine could possibly be replaced by something doing the same checks at once, to avoid looping through
// the same line twice, which will be the case when encountering blank lines when the check is "IsComment || IsBlankLine"

// IsComment loops through a line (byte slice) and looks for '#'.
// If it is found, and only, optionally, preceeded by whitespace, it returns true, otherwise false.
func IsComment(buf []byte) bool {
	for i := range buf {
		if buf[i] == '#' {
			return true
		}
		if !unicode.IsSpace(rune(buf[i])) {
			return false
		}
	}
	return false
}

func IsBlankLine(buf []byte) bool {
	for i := range buf {
		if !unicode.IsSpace(rune(buf[i])) {
			return false
		}
	}
	return true
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
		if IsComment(scanner.Bytes()) || IsBlankLine(scanner.Bytes()) {
			//fmt.Printf("Seems we have a comment: %q\n", buf)
			continue
		}
		// ...
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v+\n", err)
	}
}
