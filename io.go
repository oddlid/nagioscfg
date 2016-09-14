package nagioscfg

/*
IO-related stuff for nagioscfg
Much of the stuff here is taken from Golangs encoding/json source and modified to the specific needs of this package.
See: https://golang.org/LICENSE
*/

import (
	"bufio"
	"bytes"
	"errors"
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
	ErrNoValue = errors.New("only key given where key/value expected")
	ErrUnknown = errors.New("unknown parsing error")
)


type Reader struct {
	Comment rune
	line    int
	column  int
	field   bytes.Buffer
	r       *bufio.Reader
}

func _debug(args ...interface{}) {
	fmt.Println(args)
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

// this is basically "dos2unix"
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

// skip advances the reader until it reaches delim, ignoring everything it reads
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

func (r *Reader) parseFields() (haveField bool, delim rune, err error) {
	r.field.Reset() // clear buffer at each call

	r1, err := r.readRune()
	for err == nil && r1 != '\n' && unicode.IsSpace(r1) {
		//_debug("Skipping space")
		r1, err = r.readRune()
	}
	if err == io.EOF && r.column != 0 {
		return true, 0, err
	}
	if err != nil {
		return false, 0, err
	}

	switch r1 {
//	case '\n':
//		if r.column == 0 {
//			_debug("Bailing out from newline in first column")
//			return false, r1, nil
//		}
//		_debug("Bailing out from newline not in first column")
//		return false, r1, nil
//	case ' ':
//		_debug("Encountered space")
//		return false, r1, nil
//	case '{':
//		_debug("Found {")
//		return false, r1, nil
//	case '}':
//		_debug("Found }")
//		return false, r1, nil
	case '\n':
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		fallthrough
	case '{':
		fallthrough
	case '}':
		return false, r1, nil
	default:
		for {
			//_debug("Writing rune", r1)
			if !unicode.IsSpace(r1) {
				r.field.WriteRune(r1)
			}
			r1, err = r.readRune()
			if err != nil || r1 == '{' || r1 == '}' || unicode.IsSpace(r1) {
			//if err != nil || unicode.IsSpace(r1) {
				//_debug("Came across { or } or had an error")
				break
			}
			if r1 == '\n' {
				_debug("End of line, returning")
				return true, r1, nil
			}
		}
	}

	if err != nil {
		if err == io.EOF {
			return true, 0, err
		}
		return false, 0, err
	}

	return true, r1, nil
}

func (r *Reader) parseLine() (fields []string, err error) {
	r.line++
	r.column = -1

	r1, _, err := r.r.ReadRune()
	if err != nil {
		return nil, err
	}
	if r.Comment != 0 && r1 == r.Comment {
		return nil, r.skip('\n')
	}
	r.r.UnreadRune()

	for {
		haveField, delim, err := r.parseFields()
		if haveField {
			if fields == nil {
				fields = make([]string, 0, 6) 
			}
			fields = append(fields, r.field.String())
			//fmt.Printf("%#v\n", fields)
		}
		if delim == '\n' || delim == '{' || delim == '}' || err == io.EOF {
			return fields, err
		} else if err != nil {
			return nil, err
		}
	}
}

// Read reads from a Nagios config stream and returns the next config object. 
// Should be called repeatedly. Returns err = io.EOF when done
func (r *Reader) Read() (*CfgObj, error) {
	var fields []string
	var err error
	for {
		fields, err = r.parseLine()
		if fields != nil {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Fields: %#v\n", fields)

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
