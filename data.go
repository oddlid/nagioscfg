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

type CfgObj struct {
	Type   CfgType
	Props  map[string]string
	Indent int
	Align  int
}

func NewCfgObj(ct CfgType) *CfgObj {
	p := make(map[string]string)
	return &CfgObj{
		Type:   ct,
		Props:  p,
		Indent: DEF_INDENT,
		Align:  DEF_ALIGN,
	}
}

// methods - move to separate file when it grows

func (ct CfgType) String() string {
	return CfgTypes[ct]
}

// AddProp adds a config property to an object if it doesn't already exist. Returns error if the key already exists.
func (co *CfgObj) AddProp(k, v string) error {
	_, exists := co.Props[k]
	if exists {
		return fmt.Errorf("Key %q already exists with value %q", k, v)
	}
	co.Props[k] = v
	return nil
}

func (co *CfgObj) SetProp(k, v string) bool {
	_, exists := co.Props[k]
	co.Props[k] = v
	return exists // true = key was overwritten, false = key was added
}

func (co *CfgObj) DelProp(k string) bool {
	_, ok := co.Props[k]
	delete(co.Props, k)
	return ok // just signals if there was anything there to be deleted in the first place
}

func (co *CfgObj) Prop(k string) (string, bool) {
	v, exists := co.Props[k]
	return v, exists
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

