package nagioscfg

import (
	"testing"
	"os"
)

var co = NewCfgObj(T_SERVICE)
var keys = [...]string {
	"max_check_attempts",
	"active_checks_enabled",
	"retain_nonstatus_information",
	"service_description",
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

func TestPrint(t *testing.T) {
	co.Add(keys[3], "Disk usage /var")
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
