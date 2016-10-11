package nagioscfg

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var co = NewCfgObj(T_SERVICE)
var keys = [...]string{
	"max_check_attempts",
	"active_checks_enabled",
	"retain_nonstatus_information",
	"service_description",
	"check_command",
	"contact_groups",
}
var cmd = []string{
	"vgt_check_graphite_v2",
	"192.168.0.1",
	"some.prod.solution.queues.queue.DLQ",
	"4000",
	"5000",
	"gt",
}
var cgs = []string{
	"devs",
	"ops",
	"support",
	"managers",
}

var comment = []byte("					    #    lkdsglknag  \n")
var notcomment = []byte("			 define gris")
var blankline = []byte("						    \n")

func TestString(t *testing.T) {
	str := co.Type.String()
	exp := "service"
	if str != exp {
		t.Errorf("Expected String() to return %q, but got %q", exp, str)
	}
}

func TestSet(t *testing.T) {
	overwritten := co.Set(keys[0], "gaupe")
	if overwritten {
		t.Errorf("%q should not exist yet", keys[0])
	}
	ow2 := co.Set(keys[0], "11")
	if !ow2 {
		t.Errorf("Key %q should have been overwritten", keys[0])
	}
}

func TestAdd(t *testing.T) {
	ok := co.Add(keys[1], "1")
	if !ok {
		t.Error("Failed to add second key")
	}
	ok = co.Add(keys[1], "gris")
	if ok {
		t.Errorf("Should not be allowed to add same key %q more than once", keys[1])
	}
}

func TestGet(t *testing.T) {
	ret, exists := co.Get(keys[0])
	if !exists {
		t.Error("Get returned false")
	}
	if ret != "11" { // set in TestSet()
		t.Errorf("Expected %q, but got %q", "11", ret)
	}
}

func TestDel(t *testing.T) {
	k := "dkey"
	v := "dval"
	deleted := co.Del(k)
	if deleted {
		t.Error("Delete non-existing key should return false")
	}
	co.Add(k, v)
	deleted = co.Del(k)
	if !deleted {
		t.Errorf("Failed to delete key %q", k)
	}
	ret, exists := co.Get(k)
	if exists {
		t.Errorf("Key %q should be deleted, but got value %q", k, ret)
	}
}

func TestLongestKey(t *testing.T) {
	lk := co.LongestKey()
	correct_len := len(keys[1])
	if lk != correct_len {
		t.Errorf("LongestKey() returned %d when correct length is %d", lk, correct_len)
	}
}

func TestSetList(t *testing.T) {
	exists := co.SetList(keys[4], SEP_CMD, cmd...)
	if exists {
		t.Errorf("key %q should not exist yet", keys[4])
	}
}

func TestAddList(t *testing.T) {
	ok := co.AddList(keys[5], SEP_LST, cgs...)
	if !ok {
		t.Error("Failed to add contact groups")
	}
	ok = co.AddList(keys[5], SEP_LST, "gris", "hund", "katt")
	if ok {
		t.Errorf("Should not be allowed to add list to key %q more than once", keys[5])
	}
}

func TestGetList(t *testing.T) {
	lst := co.GetList(keys[5], SEP_LST)
	if lst == nil {
		t.Errorf("Should get a valid list from key %q", keys[5])
	}
	if !reflect.DeepEqual(lst, cgs) {
		t.Error("Returned list is not equal to the one we put in")
	}
}

func TestGetCheckCommand(t *testing.T) {
	lst := co.GetCheckCommand()
	if lst == nil {
		t.Error("Check command list should not be nil")
	}
	if !reflect.DeepEqual(lst, cmd) {
		t.Error("Returned command list does not equal what we put in")
	}
}

func TestGetCheckCommandCmd(t *testing.T) {
	checkcmd, ok := co.GetCheckCommandCmd()
	if !ok {
		t.Errorf("GetCheckCommandCmd() failed to return %q", cmd[0])
	}
	if checkcmd != cmd[0] {
		t.Errorf("Command %q is not equal to %q", checkcmd, cmd[0])
	}
}

func TestGetCheckCommandArgs(t *testing.T) {
	args := co.GetCheckCommandArgs()
	if args == nil {
		t.Error("GetCheckCommandArgs() returned nil")
	}
}

func TestGetName(t *testing.T) {
	o := NewCfgObj(T_COMMAND)
	key := "command_name"
	name := "testcommand"
	o.Set(key, name)
	ret, exists := o.GetName()
	if !exists {
		t.Errorf("Expected %q, but got nothing", name)
	}
	if ret != name {
		t.Errorf("Expected %q, but got %q", name, ret)
	}

	o = NewCfgObj(T_HOST)
	key = "name"
	name = "host-template-something"
	o.Set(key, name)
	ret, exists = o.GetName()
	if !exists {
		t.Errorf("Expected %q, but got nothing", name)
	}
	if ret != name {
		t.Errorf("Expected %q, but got %q", name, ret)
	}
}

func TestGetHostname(t *testing.T) {
	o := NewCfgObj(T_HOST)
	k := CfgKeys[24] //"host_name"
	v := "printserver"
	o.Set(k, v)
	ret, exists := o.GetHostname()
	if !exists {
		t.Errorf("Expected %q, but got nothing", v)
		//o.Print(os.Stdout)
	}
	if ret != v {
		t.Errorf("Expected %q, but got %q", v, ret)
	}
}

func TestGetDescription(t *testing.T) {
	o := NewCfgObj(T_SERVICE)
	key := "service_description"
	name := "testservice"
	o.Set(key, name)
	ret, exists := o.GetDescription()
	if !exists {
		t.Errorf("Expected %q, but got nothing", name)
	}
	if ret != name {
		t.Errorf("Expected %q, but got %q", name, ret)
	}
}

func TestGetUniqueCheckName(t *testing.T) {
	o := NewCfgObj(T_SERVICE)
	k1 := "host_name"
	k2 := "service_description"
	v1 := "host.domain.tld"
	v2 := "PLING_PLONG_LuftBallong"
	exp := fmt.Sprintf("%s;%s", v1, v2)
	o.Set(k1, v1)
	o.Set(k2, v2)
	ret, ok := o.GetUniqueCheckName()
	if !ok {
		t.Errorf("Expected %q but got nothing", exp)
	}
	if exp == "" {
		t.Errorf("Expected %q but got %q", exp, ret)
	}
	t.Logf("Unique name: %q", ret)
}

func TestGenerateComment(t *testing.T) {
	co.Add(keys[3], "Graphite DLQ")
	ok := co.generateComment()
	if !ok {
		t.Error("Attempt to generate comment returned false")
	}
	exp_comment := "# service 'Graphite DLQ'"
	if co.Comment != exp_comment {
		t.Errorf("Expected comment %q, but got %q", exp_comment, co.Comment)
	}

}

func TestPrint(t *testing.T) {
	co.Align = co.LongestKey() + 2
	co.Print(os.Stdout)
}

/*
func TestRead(t *testing.T) {
	t.Skip("TestRead temporarily disabled")
	fpath := "../op5_automation/cfg/etc/services.cfg"
	file, err := os.Open(fpath)
	if err != nil {
		t.Fatalf("Unable to open config file %q", fpath)
	}
	defer file.Close()

	Read(file)
}

func TestIsComment(t *testing.T) {
	if !IsComment(comment) {
		t.Error("Should be detected as a comment")
	}
	if IsComment(notcomment) {
		t.Error("Should be detected as not a comment")
	}
}

func TestIsBlankLine(t *testing.T) {
	if !IsBlankLine(blankline) {
		t.Error("Should be detected as a blank line")
	}
	if IsBlankLine(notcomment) {
		t.Error("Should be detected as not a blank line")
	}
}

func BenchmarkIsComment(b *testing.B) {
	for n:= 0; n < b.N; n++ {
		IsComment(comment)
	}
}

func BenchmarkIsBlankLine(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsBlankLine(blankline)
	}
}
*/
