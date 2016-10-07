package nagioscfg

/*
Defines the data structures that maps to Nagios config items
Odd, 2016-08-10 17:51
*/

import (
	"io"
)

type CfgType int
type CfgName string
type IoState int
type CfgObjs []*CfgObj

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
	DEF_INDENT int    = 4
	DEF_ALIGN  int    = 32
	SEP_CMD    string = "!"
	SEP_LST    string = ","
)

const (
	IO_OBJ_OUT IoState = iota
	IO_OBJ_BEGIN
	IO_OBJ_IN
	IO_OBJ_END
)

var CfgTypes = [...]CfgName{
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

/*
type PropertyCollection interface {
	Add(key, val string) bool      // should only add if key does not yet exist. Return false if key exists
	Set(key, val string) bool      // adds or overwrites. Return true if key was overwritten
	Get(key string) (string, bool) // return val, success
	Del(key string) bool           // return true if key was present
	LongestKey() int
}
*/

type Printer interface {
	Print(w io.Writer)
}

type CfgObj struct {
	Type    CfgType
	Indent  int
	Align   int
	Comment string
	Props   map[string]string
}

type CfgFile struct {
	Path string
	Objs CfgObjs
}
