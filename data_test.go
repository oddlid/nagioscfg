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

func TestAddProp(t *testing.T) {
	err := co.AddProp(keys[0], "11")
	if err != nil {
		t.Error(err)
	}
	co.Print(os.Stdout)
}

func TestSetProp(t *testing.T) {
	overwritten := co.SetProp(keys[0], "gaupe")
	if !overwritten {
		t.Errorf("%q should have been overwritten", keys[0])
	}
	ow2 := co.SetProp(keys[1], "jalla")
	if ow2 {
		t.Error("Key should not exist yet")
	}
}

func TestPrint(t *testing.T) {
	co.Align = co.LongestKey() + 2
	co.Print(os.Stdout)
}
