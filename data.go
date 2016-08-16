package nagioscfg

/*
Defines the data structures that maps to Nagios config items
Odd, 2016-08-10 17:51
*/


import (
	"fmt"
	"io"
	"strings"
)

type CfgType int

const (
	T_COMMAND CfgType = iota
	T_CONTACTGROUP
	T_CONTACT
	T_HOSTESCALATION
	T_HOSTGROUP
	T_HOST
	T_SERVICEESCALATION
	T_SERVICEGROUP
	T_SERVICE
	T_TIMEPERIOD
)

const (
	DEF_INDENT int = 4
	DEF_ALIGN  int = 32
)

var CfgTypes = [...]string {
	"command",
	"contactgroup",
	"contact",
	"hostescalation",
	"hostgroup",
	"host",
	"serviceescalation",
	"servicegroup",
	"service",
	"timeperiod",
}

type PropertyCollection interface {
	Add(key, val string) bool      // should only add if key does not yet exist. Return false if key exists
	Set(key, val string) bool      // adds or overwrites. Return true if key was overwritten
	Get(key string) (string, bool) // return val, success
	Del(key string) bool           // return true if key was present
	LongestKey() int
}

type Printer interface {
	Print(w io.Writer)
}


type CfgObj struct {
	Type   CfgType
	Props  map[string]string
	Indent int
	Align  int
}

//type CmdItem struct {
//	Key string
//	Cmd string
//}

func NewCfgObj(ct CfgType) *CfgObj {
	p := make(map[string]string)
	return &CfgObj{
		Type:   ct,
		Props:  p,
		Indent: DEF_INDENT,
		Align:  DEF_ALIGN,
	}
}


//func SplitList(separator string) []string {
//	var list []string
//	return list
//}
//
//func JoinList(separator string, args ...string) string {
//	return ""
//}

// methods - move to separate file when it grows

func (ct CfgType) String() string {
	return CfgTypes[ct]
}

func (co *CfgObj) Add(key, val string) bool {
	_, exists := co.Props[key]
	if exists {
		return !exists
	}
	co.Props[key] = val
	return true
}

func (co *CfgObj) Set(key, val string) bool {
	_, exists := co.Props[key]
	co.Props[key] = val
	return exists // true = key was overwritten, false = key was added
}

func (co *CfgObj) Del(key string) bool {
	_, exists := co.Props[key]
	delete(co.Props, key)
	return exists // just signals if there was anything there to be deleted in the first place
}

func (co *CfgObj) Get(key string) (string, bool) {
	val, exists := co.Props[key]
	return val, exists
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
	ct := co.Type.String()
	fstr := fmt.Sprintf("%s%s%d%s", prefix, "%-", co.Align, "s%s\n")
	//fmt.Fprintf(w, "# %s '%s'\n", ct, "bogus")
	fmt.Fprintf(w, "define %s{\n", ct)
	for k, v := range co.Props {
		fmt.Fprintf(w, fstr, k, v)
	}
	fmt.Fprintf(w, "%s}\n", prefix)
}

