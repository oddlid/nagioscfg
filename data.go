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
//type CfgKey int

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

var CfgKeys = [...]string{
	"active_checks_enabled",
	"address",
	"alias",
	"can_submit_commands",
	"check_command",
	"check_freshness",
	"check_interval",
	"check_period",
	"command_line",
	"command_name",
	"contact_groups",
	"contact_name",
	"contactgroup_name",
	"contactgroups",
	"contacts",
	"display_name",
	"email",
	"escalation_options",
	"escalation_period",
	"event_handler_enabled",
	"first_notification",
	"flap_detection_enabled",
	"flap_detection_options",
	"friday",
	"host_name",
	"host_notification_commands",
	"host_notification_options",
	"host_notification_period",
	"host_notifications_enabled",
	"hostgroup_members",
	"hostgroup_name",
	"hostgroups",
	"icon_image",
	"is_volatile",
	"last_notification",
	"max_check_attempts",
	"monday",
	"name",
	"notes",
	"notes_url",
	"notification_interval",
	"notification_options",
	"notification_period",
	"notifications_enabled",
	"obsess",
	"pager",
	"parallelize_check",
	"parents",
	"passive_checks_enabled",
	"process_perf_data",
	"register",
	"retain_nonstatus_information",
	"retain_status_information",
	"retry_interval",
	"saturday",
	"service_description",
	"service_notification_commands",
	"service_notification_options",
	"service_notification_period",
	"service_notifications_enabled",
	"servicegroup_name",
	"servicegroups",
	"stalking_options",
	"statusmap_image",
	"sunday",
	"thursday",
	"timeperiod_name",
	"tuesday",
	"use",
	"wednesday",
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
	Type    CfgType
	Props   map[string]string
	Indent  int
	Align   int
	Comment string
}

func NewCfgObj(ct CfgType) *CfgObj {
	p := make(map[string]string)
	return &CfgObj{
		Type:    ct,
		Props:   p,
		Indent:  DEF_INDENT,
		Align:   DEF_ALIGN,
		Comment: "# " + ct.String() + "'%s'",
	}
}


// methods - move to separate file when it grows

func (ct CfgType) String() string {
	return CfgTypes[ct]
}

func (co *CfgObj) Add(key, val string) bool {
	_, exists := co.Props[key]
	if exists {
		return false
	}
	co.Props[key] = val
	return true
}

func (co *CfgObj) Set(key, val string) bool {
	_, exists := co.Props[key]
	co.Props[key] = val
	return exists // true = key was overwritten, false = key was added
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
	fmt.Fprintf(w, "# %s\n", co.Comment)
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

func (co *CfgObj) GetCheckCommand() []string {
	if co.Type != T_SERVICE {
		return nil
	}
	lst := co.GetList(CfgKeys[4], "!") // make sure to update index here if CfgKeys is updated
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
	if success {
		co.Comment = fmt.Sprintf(co.Comment, name)
	}
	return success
}
