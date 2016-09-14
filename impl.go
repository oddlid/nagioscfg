/*
Function/method implementations for types from data.go
*/

package nagioscfg

import (
	"fmt"
	"io"
	"strings"
)

func NewCfgObj(ct CfgType) *CfgObj {
	p := make(map[string]string)
	return &CfgObj{
		Type:    ct,
		Props:   p,
		Indent:  DEF_INDENT,
		Align:   DEF_ALIGN,
		Comment: "# " + ct.String() + " '%s'",
	}
}

func (ct CfgType) String() string {
	return string(CfgTypes[ct])
}

// Type returns the int value for the given CfgName, or -1 it not valid
func (cn CfgName) Type() CfgType {
	for i := range CfgTypes {
		if CfgTypes[i] == cn {
			return CfgType(i)
		}
	}
	return -1
}

func (co *CfgObj) Set(key, val string) bool {
	_, exists := co.Props[key]
	co.Props[key] = val
	return exists // true = key was overwritten, false = key was added
}

func (co *CfgObj) Add(key, val string) bool {
	_, exists := co.Props[key]
	if exists {
		return false
	}
	return !co.Set(key, val) // Set should return false, as the key doesn't exist yet, so we inverse the result
}

func (co *CfgObj) Get(key string) (string, bool) {
	val, exists := co.Props[key]
	return val, exists
}

func (co *CfgObj) Del(key string) bool {
	_, exists := co.Props[key]
	delete(co.Props, key)
	return exists // just signals if there was anything there to be deleted in the first place
}

func (co *CfgObj) LongestKey() int {
	max := 0
	for k, _ := range co.Props {
		l := len(k)
		if l > max {
			max = l
		}
	}
	return max
}

func (co *CfgObj) Print(w io.Writer) {
	prefix := strings.Repeat(" ", co.Indent)
	//co.Align = co.LongestKey() + 1
	fstr := fmt.Sprintf("%s%s%d%s", prefix, "%-", co.Align, "s%s\n")
	co.generateComment() // this might fail, but don't care yet
	fmt.Fprintf(w, "%s\n", co.Comment)
	fmt.Fprintf(w, "define %s{\n", co.Type.String())
	for k, v := range co.Props {
		fmt.Fprintf(w, fstr, k, v)
	}
	fmt.Fprintf(w, "%s}\n", prefix)
}

func (co *CfgObj) GetList(key, sep string) []string {
	val, exists := co.Get(key)
	if !exists {
		return nil
	}
	return strings.Split(val, sep)
}

func (co *CfgObj) SetList(key, sep string, list ...string) bool {
	lstr := strings.Join(list, sep)
	return co.Set(key, lstr)
}

func (co *CfgObj) AddList(key, sep string, list ...string) bool {
	_, exists := co.Props[key]
	if exists {
		return false
	}
	return !co.SetList(key, sep, list...) // SetList should return false as key does not exist, so invert the result
}

func (co *CfgObj) GetCheckCommand() []string {
	if co.Type != T_SERVICE {
		return nil
	}
	lst := co.GetList(CfgKeys[4], SEP_CMD) // make sure to update index here if CfgKeys is updated
	if lst == nil {
		return nil
	}
	return lst
}

func (co *CfgObj) GetCheckCommandCmd() (string, bool) {
	lst := co.GetCheckCommand()
	if lst == nil {
		return "", false
	}
	return lst[0], true
}

func (co *CfgObj) GetCheckCommandArgs() []string {
	lst := co.GetCheckCommand()
	if lst == nil {
		return nil
	}
	return lst[1:]
}

func (co *CfgObj) GetName() (string, bool) {
	key := co.Type.String() + "_name"
	return co.Get(key)
}

func (co *CfgObj) GetDescription() (string, bool) {
	key := co.Type.String() + "_description"
	return co.Get(key)
}

// generateComment is set as private, as it makes "unsafe" assumptions about the existing format of the comment
func (co *CfgObj) generateComment() bool {
	var name string
	var success bool
	if co.Type == T_SERVICE {
		name, success = co.GetDescription()
	} else {
		name, success = co.GetName()
	}
	if success && strings.Index(co.Comment, "%") > -1 {
		co.Comment = fmt.Sprintf(co.Comment, name)
	}
	return success
}

//func (co *CfgObj) Align(row int) int {
//	if row == 0 { // "auto"
//		co.Align = co.LongestKey() + 2
//	} else {
//		co.Align = row
//	}
//	return co.Align
//}
