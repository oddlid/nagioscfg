package nagioscfg

import (
	"testing"
	"os"
)

var co = NewCfgObj(T_COMMAND)
var keys = [...]string {
	"max_check_attempts",
	"active_checks_enabled",
	"retain_nonstatus_information",
}

func TestAdd(t *testing.T) {
	ok := co.Add(keys[0], "11")
	if !ok {
		t.Error("Failed to add first key")
	}
	ok = co.Add(keys[0], "gris")
	if ok {
		t.Error("Should not be allowed to add same key more than once")
	}
	co.Print(os.Stdout)
}

func TestSet(t *testing.T) {
	overwritten := co.Set(keys[0], "gaupe")
	if !overwritten {
		t.Errorf("%q should have been overwritten", keys[0])
	}
	ow2 := co.Set(keys[1], "jalla")
	if ow2 {
		t.Error("Key should not exist yet")
	}
}

func TestLongestKey(t *testing.T) {
	lk := co.LongestKey()
	correct_len := len(keys[1])
	if lk != correct_len {
		t.Errorf("LongestKey() returned %d when correct length is %d", lk, correct_len)
	}
}

func TestPrint(t *testing.T) {
	co.Align = co.LongestKey() + 2
	co.Print(os.Stdout)
}
