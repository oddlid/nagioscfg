package nagioscfg

import (
	"testing"
	"strings"
)

var cfgobjstr string = `# some comment
define service{
	  service_description    A service name with spaces
# embedded comment
	  a_key                  Some value
    }
	
define command{
	command_name gris
	gris_fest roligt
}

`

func TestRead(t *testing.T) {
	str_r := strings.NewReader(cfgobjstr)
	rdr := NewReader(str_r)
	rdr.Read()
	rdr.Read()
	rdr.Read()
	rdr.Read()
	rdr.Read()
}
