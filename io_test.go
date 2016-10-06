package nagioscfg

import (
	"os"
	"strings"
	"testing"
)

var cfgobjstr string = `# some comment
define service {
	  service_description    A service name with spaces
# embedded comment
	  a_key                  Some value
		singlekey
    }
	
define command {
	command_name gris
	gris_fest roligt
}

# Bla bla, some comment crap
# I'm really too tired now

define service{
	service_description Disk usage /my/ass
	contact_group toilet
}
`

func TestRead(t *testing.T) {
	str_r := strings.NewReader(cfgobjstr)
	rdr := NewReader(str_r)
	co, err := rdr.Read()
	if err != nil {
		t.Fatal(err)
	}
	if co == nil {
		t.Fatal("CfgObj is nil")
	}
	co.AutoAlign()
	co.Print(os.Stdout)
}

func TestReadAll(t *testing.T) {
	//t.Skip("Not implemented yet")
	str_r := strings.NewReader(cfgobjstr)
	rdr := NewReader(str_r)
	cos, err := rdr.ReadAll()
	if err != nil {
		t.Error(err)
	} else {
		cos.AutoAlign()
		cos.Print(os.Stdout)
	}
}

func TestReadFile(t *testing.T) {
	path := "../op5_automation/cfg/etc/services.cfg"
	objs, err := ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(objs))
}

func TestWriteFile(t *testing.T) {
	src := "../op5_automation/cfg/etc/services.cfg"
	dst := "/tmp/services.cfg"
	objs, err := ReadFile(src)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(objs))
	err = WriteFile(dst, objs)
	if err != nil {
		t.Error(err)
	}
}
