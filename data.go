package nagioscfg

/*
Defines the data structures that maps to Nagios config items
Odd, 2016-08-10 17:51
*/

import (
	//"io"
	"regexp"
)

//type WriteMap map[string]CfgMap // used to sort/write out according to FileID


// UUID representation compliant with specification
// described in RFC 4122.
type UUID [16]byte

type CfgType int
type CfgName string
type IoState int
type CfgObjs []*CfgObj
type CfgMap map[UUID]*CfgObj

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

const (
	T_COMMAND CfgType = iota
	T_CONTACT
	T_CONTACTGROUP
	T_HOST
	T_HOSTDEPENDENCY
	T_HOSTESCALATION
	T_HOSTEXTINFO
	T_HOSTGROUP
	T_SERVICE
	T_SERVICEDEPENDENCY
	T_SERVICEESCALATION
	T_SERVICEEXTINFO
	T_SERVICEGROUP
	T_TIMEPERIOD
)

var CfgTypes = [...]CfgName{
	"command",
	"contact",
	"contactgroup",
	"host",
	"hostdependency",
	"hostescalation",
	"hostextinfo",
	"hostgroup",
	"service",
	"servicedependency",
	"serviceescalation",
	"serviceextinfo",
	"servicegroup",
	"timeperiod",
}

var CfgKeys = map[int]string{
	0:  "2d_coords",
	1:  "3d_coords",
	2:  "action_url",
	3:  "active_checks_enabled",
	4:  "address",
	5:  "addressx",
	6:  "alias",
	7:  "can_submit_commands",
	8:  "check_command",
	9:  "check_freshness",
	10: "check_interval",
	11: "check_period",
	12: "command_line",
	13: "command_name",
	14: "contact_groups",
	15: "contact_name",
	16: "contactgroup_members",
	17: "contactgroup_name",
	18: "contactgroups",
	19: "contacts",
	20: "dependency_period",
	21: "dependent_host_name",
	22: "dependent_hostgroup_name",
	23: "dependent_service_description",
	24: "dependent_servicegroup_name",
	25: "display_name",
	26: "email",
	27: "escalation_options",
	28: "escalation_period",
	29: "event_handler",
	30: "event_handler_enabled",
	31: "exclude",
	32: "execution_failure_criteria",
	33: "first_notification",
	34: "first_notification_delay",
	35: "flap_detection_enabled",
	36: "flap_detection_options",
	37: "freshness_threshold",
	38: "friday",
	39: "high_flap_threshold",
	40: "host_name",
	41: "host_notification_commands",
	42: "host_notification_options",
	43: "host_notification_period",
	44: "host_notifications_enabled",
	45: "hostgroup_members",
	46: "hostgroup_name",
	47: "hostgroups",
	48: "icon_image",
	49: "icon_image_alt",
	50: "inherits_parent",
	51: "initial_state",
	52: "is_volatile",
	53: "last_notification",
	54: "low_flap_threshold",
	55: "max_check_attempts",
	56: "members",
	57: "monday",
	58: "notes",
	59: "notes_url",
	60: "notification_failure_criteria",
	61: "notification_interval",
	62: "notification_options",
	63: "notification_period",
	64: "notifications_enabled",
	65: "obsess_over_host",
	66: "obsess_over_service",
	67: "pager",
	68: "parents",
	69: "passive_checks_enabled",
	70: "process_perf_data",
	71: "retain_nonstatus_information",
	72: "retain_status_information",
	73: "retry_interval",
	74: "saturday",
	75: "service_description",
	76: "service_notification_commands",
	77: "service_notification_options",
	78: "service_notification_period",
	79: "service_notifications_enabled",
	80: "servicegroup_members",
	81: "servicegroup_name",
	82: "servicegroups",
	83: "stalking_options",
	84: "statusmap_image",
	85: "sunday",
	86: "thursday",
	87: "timeperiod_name",
	88: "tuesday",
	89: "vrml_image",
	90: "wednesday",
}

// Key order for each type defined here:
// https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/objectdefinitions.html
var CfgKeySortOrder = map[string]map[CfgType]int{
	CfgKeys[0]: map[CfgType]int{ // 2d_coords
		T_HOST:              41,
		T_HOSTEXTINFO:       8,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[1]: map[CfgType]int{ // 3d_coords
		T_HOST:              42,
		T_HOSTEXTINFO:       9,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[2]: map[CfgType]int{ // action_url
		T_HOST:              36,
		T_HOSTEXTINFO:       3,
		T_HOSTGROUP:         6,
		T_SERVICE:           36,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEEXTINFO:    4,
		T_SERVICEGROUP:      6,
	},
	CfgKeys[3]: map[CfgType]int{ // active_checks_enabled
		T_HOST:              11,
		T_SERVICE:           11,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[4]: map[CfgType]int{ // address
		T_HOST:              3,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[5]: map[CfgType]int{ // addressx
		T_CONTACT:           13,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[6]: map[CfgType]int{ // alias
		T_HOST:              1,
		T_HOSTGROUP:         1,
		T_SERVICEGROUP:      1,
		T_CONTACT:           1,
		T_CONTACTGROUP:      1,
		T_TIMEPERIOD:        1,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[7]: map[CfgType]int{ // can_submit_commands
		T_CONTACT:           14,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[8]: map[CfgType]int{ // check_command
		T_HOST:              6,
		T_SERVICE:           6,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[9]: map[CfgType]int{ // check_freshness
		T_HOST:              15,
		T_SERVICE:           15,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[10]: map[CfgType]int{ // check_interval
		T_HOST:              9,
		T_SERVICE:           9,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[19]: map[CfgType]int{ // contacts
		T_HOST:              26,
		T_HOSTESCALATION:    2,
		T_SERVICE:           31,
		T_SERVICEESCALATION: 3,
	},
	CfgKeys[40]: map[CfgType]int{ // host_name
		T_HOST:              0,
		T_HOSTDEPENDENCY:    2,
		T_HOSTESCALATION:    0,
		T_HOSTEXTINFO:       0,
		T_SERVICE:           0,
		T_SERVICEDEPENDENCY: 5,
		T_SERVICEESCALATION: 0,
		T_SERVICEEXTINFO:    0,
	},
	CfgKeys[75]: map[CfgType]int{ // service_description
		T_SERVICE:           2,
		T_SERVICEDEPENDENCY: 7,
		T_SERVICEESCALATION: 2,
		T_SERVICEEXTINFO:    1,
	},
}


type CfgObj struct {
	Type    CfgType
	UUID    UUID
	Indent  int
	Align   int
	FileID  string
	Comment string
	Props   map[string]string
}

type CfgQuery struct {
	Keys []string
	RXs  []*regexp.Regexp
}

// Top level struct for managing collections of CfgObj
//type NagiosCfg struct {
//	Objs map[string]CfgMap // key by FileID
//}

//type CfgFile struct {
//	Path string
//	Objs CfgObjs
//}
/*
type PropertyCollection interface {
	Add(key, val string) bool      // should only add if key does not yet exist. Return false if key exists
	Set(key, val string) bool      // adds or overwrites. Return true if key was overwritten
	Get(key string) (string, bool) // return val, success
	Del(key string) bool           // return true if key was present
	LongestKey() int
}
*/

//type CfgObjCollection interface {
//	Add(key string, val *CfgObj) bool
//	Set(key string, val *CfgObj) bool
//	Get(key string) (*CfgObj, bool)
//	Del(key string) *CfgObj
//	LongestKey() int
//	MatchKeys(rx *regexp.Regexp, keys ...string) CfgObjCollection
//	MatchAny(rx *regexp.Regexp)  CfgObjCollection
//}

//type Printer interface {
//	Print(w io.Writer)
//}

