package nagioscfg

/*
IO-related stuff for nagioscfg
Much of the stuff here is taken from Golangs encoding/json source and modified to the specific needs of this package.
See: https://golang.org/LICENSE
*/

import (
	"bufio"
	//"bytes"
	"fmt"
	"io"
	"unicode"
)

// A ParseError is returned for parsing errors.
// The first line is 1.  The first column is 0.
type ParseError struct {
	Line   int   // Line where the error occurred
	Column int   // Column (rune index) where the error occurred
	Err    error // The actual error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Err)
}

// These are the errors that can be returned in ParseError.Error
var (
	ErrNoValue("only key given where key/value expected")
	ErrUnknown("unknown parsing error")
)


type Reader struct {
	Comment rune
	line    int
	column  int
	field   bytes.Buffer
	r       *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Comment: '#',
		r:       bufio.NewReader(r),
	}
}

func (r *Reader) error(err error) error {
	return &ParseError{
		Line:   r.line,
		Column: r.column,
		Err:    err,
	}
}

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

// Read reads from a Nagios config stream and returns the next config object. 
// Should be called repeatedly. Returns err = io.EOF when done
func (r *Reader) Read() (*CfgObj, error) {
	return nil, nil
}

// ReadAll calls Read repeatedly and returns all config objects it collects
func (r *Reader) ReadAll() ([]CfgObj, error) {
	return nil, nil
}


// *** Stuff below is temporary kept for reference to what I was first thinking. For later removal. ***

// IsComment and IsBlankLine could possibly be replaced by something doing the same checks at once, to avoid looping through
// the same line twice, which will be the case when encountering blank lines when the check is "IsComment || IsBlankLine"

// IsComment loops through a line (byte slice) and looks for '#'.
// If it is found, and only, optionally, preceeded by whitespace, it returns true, otherwise false.
/*
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
*/
/*
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
*/

/*
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
*/
