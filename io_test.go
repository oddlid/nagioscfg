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
	co, err := rdr.Read(false, "/dev/null")
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
	cos, err := rdr.ReadAll(false, "/dev/null")
	if err != nil {
		t.Error(err)
	} else {
		cos.AutoAlign()
		cos.Print(os.Stdout)
	}
}

func TestReadAllMap(t *testing.T) {
	str_r := strings.NewReader(cfgobjstr)
	rdr := NewReader(str_r)
	m, err := rdr.ReadAllMap("/dev/null")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("=== Map: ===\n%s\n", m.Dump())
	}
}

func TestReadFile(t *testing.T) {
	path := "../op5_automation/cfg/etc/services.cfg"
	objs, err := ReadFile(path, false)
	if err != nil {
		t.Error(err)
	}
	t.Log("Number of objets read: ", len(objs))
}

func BenchmarkReadFile(b *testing.B) {
	path := "../op5_automation/cfg/etc/services.cfg"
	for i := 0; i <= b.N; i++ {
		ReadFile(path, false)
	}
}

func BenchmarkReadFileSetUUID(b *testing.B) {
	path := "../op5_automation/cfg/etc/services.cfg"
	for i := 0; i <= b.N; i++ {
		ReadFile(path, true)
	}
}

func TestObjReadFile(t *testing.T) {
	path := "../op5_automation/cfg/etc/services.cfg"
	cf := NewCfgFile(path)
	err := cf.Read(false)
	if err != nil {
		t.Error(err)
	}
	t.Log("Number of objets read: ", len(cf.Objs))
}

func TestWriteFile(t *testing.T) {
	src := "../op5_automation/cfg/etc/services.cfg"
	dst := "/tmp/services.cfg"
	objs, err := ReadFile(src, false)
	if err != nil {
		t.Error(err)
	}
	t.Log("Number of objets read: ", len(objs))
	err = WriteFile(dst, objs)
	if err != nil {
		t.Error(err)
	}
}

func TestObjWriteFile(t *testing.T) {
	src := "../op5_automation/cfg/etc/services.cfg"
	dst := "/tmp/services.cfg"
	cf := NewCfgFile(src)
	err := cf.Read(false)
	if err != nil {
		t.Error(err)
	}
	t.Log("Number of objets read: ", len(cf.Objs))
	cf.Path = dst
	err = cf.Write()
	if err != nil {
		t.Error(err)
	}
}
