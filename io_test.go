package nagioscfg

import (
	"testing"
	"strings"
)

var cfgobjstr string = `# some comment
define service{
	  service_description    A service name with spaces
	  a_key                  Some value
    }
`

func TestRead(t *testing.T) {
	str_r := strings.NewReader(cfgobjstr)
	rdr := NewReader(str_r)
	rdr.Read()
	rdr.Read()
	rdr.Read()
	rdr.Read()
}
