package nagioscfg

import (
	"testing"
	"strings"
	"os"
)

var cfgobjstr string = `# some comment
define service{
	  service_description    A service name with spaces
# embedded comment
	  a_key                  Some value
		singlekey
    }
	
define command {
	command_name gris
	gris_fest roligt
}

`

func TestRead(t *testing.T) {
	str_r := strings.NewReader(cfgobjstr)
	rdr := NewReader(str_r)
	co, err := rdr.Read()
	if err != nil {
		t.Error(err)
	}
	co.Print(os.Stdout)
}
