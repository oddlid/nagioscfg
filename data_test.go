package nagioscfg

import (
	"os"
	"reflect"
	"testing"
)

var co = NewCfgObj(T_SERVICE)
var keys = [...]string {
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
	//co.Print(os.Stdout)
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


func TestPrint(t *testing.T) {
	co.Add(keys[3], "Graphite DLQ") // just to get a description/comment as well
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
