/*
   Copyright 2017 Odd Eivind Ebbesen

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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
type UUIDs []UUID

type CfgType int
type CfgName string
type CfgProp string
type IoState int
type CfgObjs []*CfgObj
type CfgMap map[UUID]*CfgObj

const PKGNAME string = "nagioscfg"
const VERSION string = "2017-07-21"
const PROJECT_PREFIX string = "github.com/vgtmnm/"

const (
	DEF_INDENT int    = 4
	DEF_ALIGN  int    = 31
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
	T_INVALID
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
	89: "use",
	90: "vrml_image",
	91: "wednesday",
	// Newly discovered keys, not in Nagios core, only op5, in unsorted order
	92: "name",
	93: "obsess",
	94: "parallelize_check",
	95: "register",
	96: "hourly_value",
}

// Key order for each type defined here:
// https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/objectdefinitions.html
// 2017-07-21 14:26:39 - Just discovered that op5 has several extra keys not defined by Nagios. Fuuuuuuck....
//   parallelize_check, obsess, register, name, hourly_value
//   It also seems that check_command should have a lower sort order than servicegroups, according to diffs.
//   This is not according to Nagios specs, but op5's output
var CfgKeySortOrder = map[string]map[CfgType]int{
	CfgKeys[0]: map[CfgType]int{ // 2d_coords
		T_HOST:              42,
		T_HOSTEXTINFO:       8,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[1]: map[CfgType]int{ // 3d_coords
		T_HOST:              43,
		T_HOSTEXTINFO:       9,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[2]: map[CfgType]int{ // action_url
		T_HOST:              37,
		T_HOSTEXTINFO:       3,
		T_HOSTGROUP:         6,
		T_SERVICE:           37,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEEXTINFO:    4,
		T_SERVICEGROUP:      6,
	},
	CfgKeys[3]: map[CfgType]int{ // active_checks_enabled
		T_HOST:              12,
		T_SERVICE:           12,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[4]: map[CfgType]int{ // address
		T_HOST:              4,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[5]: map[CfgType]int{ // addressx
		T_CONTACT:           14,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[6]: map[CfgType]int{ // alias
		T_HOST:              2,
		T_HOSTGROUP:         1,
		T_SERVICEGROUP:      1,
		T_CONTACT:           2,
		T_CONTACTGROUP:      1,
		T_TIMEPERIOD:        1,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[7]: map[CfgType]int{ // can_submit_commands
		T_CONTACT:           15,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[8]: map[CfgType]int{ // check_command
		T_HOST:              7,
		T_SERVICE:           7,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[9]: map[CfgType]int{ // check_freshness
		T_HOST:              16,
		T_SERVICE:           16,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[10]: map[CfgType]int{ // check_interval
		T_HOST:              10,
		T_SERVICE:           10,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[11]: map[CfgType]int{ // check_period
		T_HOST:              14,
		T_SERVICE:           14,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[12]: map[CfgType]int{ // command_line
		T_COMMAND:           1,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[13]: map[CfgType]int{ // command_name
		T_COMMAND:           0,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[14]: map[CfgType]int{ // contact_groups
		T_HOST:              28,
		T_HOSTESCALATION:    3,
		T_SERVICE:           33,
		T_SERVICEESCALATION: 4,
	},
	CfgKeys[15]: map[CfgType]int{ // contact_name
		T_CONTACT:           1,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[16]: map[CfgType]int{ // contactgroup_members
		T_CONTACTGROUP:      3,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[17]: map[CfgType]int{ // contactgroup_name
		T_CONTACTGROUP:      0,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[18]: map[CfgType]int{ // contactgroups
		T_CONTACT:           3,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[19]: map[CfgType]int{ // contacts
		T_HOST:              27,
		T_HOSTESCALATION:    2,
		T_SERVICE:           32,
		T_SERVICEESCALATION: 3,
	},
	CfgKeys[20]: map[CfgType]int{ // dependency_period
		T_HOSTDEPENDENCY:    7,
		T_SERVICEDEPENDENCY: 11,
	},
	CfgKeys[21]: map[CfgType]int{ // dependent_host_name
		T_HOSTDEPENDENCY:    0,
		T_SERVICEDEPENDENCY: 0,
	},
	CfgKeys[22]: map[CfgType]int{ // dependent_hostgroup_name
		T_HOSTDEPENDENCY:    1,
		T_SERVICEDEPENDENCY: 1,
	},
	CfgKeys[23]: map[CfgType]int{ // dependent_service_description
		T_SERVICEDEPENDENCY: 4,
	},
	CfgKeys[24]: map[CfgType]int{ // dependent_servicegroup_name
		T_SERVICEDEPENDENCY: 3,
	},
	CfgKeys[25]: map[CfgType]int{ // display_name
		T_HOST:              3,
		T_SERVICE:           4,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[26]: map[CfgType]int{ // email
		T_CONTACT:           12,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[27]: map[CfgType]int{ // escalation_options
		T_HOSTESCALATION:    8,
		T_SERVICEESCALATION: 9,
	},
	CfgKeys[28]: map[CfgType]int{ // escalation_period
		T_HOSTESCALATION:    7,
		T_SERVICEESCALATION: 8,
	},
	CfgKeys[29]: map[CfgType]int{ // event_handler
		T_HOST:              18,
		T_SERVICE:           18,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[30]: map[CfgType]int{ // event_handler_enabled
		T_HOST:              19,
		T_SERVICE:           19,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[31]: map[CfgType]int{ // exclude
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_TIMEPERIOD:        4,
	},
	CfgKeys[32]: map[CfgType]int{ // execution_failure_criteria
		T_HOSTDEPENDENCY:    5,
		T_SERVICEDEPENDENCY: 9,
	},
	CfgKeys[33]: map[CfgType]int{ // first_notification
		T_HOSTESCALATION:    4,
		T_SERVICEESCALATION: 5,
	},
	CfgKeys[34]: map[CfgType]int{ // first_notification_delay
		T_HOST:              30,
		T_SERVICE:           28,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[35]: map[CfgType]int{ // flap_detection_enabled
		T_HOST:              22,
		T_SERVICE:           22,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[36]: map[CfgType]int{ // flap_detection_options
		T_HOST:              23,
		T_SERVICE:           23,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[37]: map[CfgType]int{ // freshness_threshold
		T_HOST:              17,
		T_SERVICE:           17,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[38]: map[CfgType]int{ // friday
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_TIMEPERIOD:        2,
	},
	CfgKeys[39]: map[CfgType]int{ // high_flap_threshold
		T_HOST:              21,
		T_SERVICE:           21,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[40]: map[CfgType]int{ // host_name
		T_HOST:              1,
		T_HOSTDEPENDENCY:    2,
		T_HOSTESCALATION:    0,
		T_HOSTEXTINFO:       0,
		T_SERVICE:           1,
		T_SERVICEDEPENDENCY: 5,
		T_SERVICEESCALATION: 0,
		T_SERVICEEXTINFO:    0,
	},
	CfgKeys[41]: map[CfgType]int{ // host_notification_commands
		T_CONTACT:           10,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[42]: map[CfgType]int{ // host_notification_options
		T_CONTACT:           8,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[43]: map[CfgType]int{ // host_notification_period
		T_CONTACT:           6,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[44]: map[CfgType]int{ // host_notifications_enabled
		T_CONTACT:           4,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[45]: map[CfgType]int{ // hostgroup_members
		T_HOSTGROUP:         3,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[46]: map[CfgType]int{ // hostgroup_name
		T_HOSTDEPENDENCY:    3,
		T_HOSTESCALATION:    1,
		T_HOSTGROUP:         0,
		T_SERVICE:           2,
		T_SERVICEDEPENDENCY: 6,
		T_SERVICEESCALATION: 1,
	},
	CfgKeys[47]: map[CfgType]int{ // hostgroups
		T_HOST:              6,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[48]: map[CfgType]int{ // icon_image
		T_HOST:              38,
		T_HOSTEXTINFO:       4,
		T_SERVICE:           38,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEEXTINFO:    5,
	},
	CfgKeys[49]: map[CfgType]int{ // icon_image_alt
		T_HOST:              39,
		T_HOSTEXTINFO:       5,
		T_SERVICE:           39,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEEXTINFO:    6,
	},
	CfgKeys[50]: map[CfgType]int{ // inherits_parent
		T_HOSTDEPENDENCY:    4,
		T_SERVICEDEPENDENCY: 8,
	},
	CfgKeys[51]: map[CfgType]int{ // initial_state
		T_HOST:              8,
		T_SERVICE:           8,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[52]: map[CfgType]int{ // is_volatile
		T_SERVICE:           6,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[53]: map[CfgType]int{ // last_notification
		T_HOSTESCALATION:    5,
		T_SERVICEESCALATION: 6,
	},
	CfgKeys[54]: map[CfgType]int{ // low_flap_threshold
		T_HOST:              20,
		T_SERVICE:           20,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[55]: map[CfgType]int{ // max_check_attempts
		T_HOST:              9,
		T_SERVICE:           9,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[56]: map[CfgType]int{ // members
		T_CONTACTGROUP:      2,
		T_HOSTGROUP:         2,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEGROUP:      2,
	},
	CfgKeys[57]: map[CfgType]int{ // monday
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_TIMEPERIOD:        2,
	},
	CfgKeys[58]: map[CfgType]int{ // notes
		T_HOST:              35,
		T_HOSTEXTINFO:       1,
		T_HOSTGROUP:         4,
		T_SERVICE:           35,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEEXTINFO:    2,
		T_SERVICEGROUP:      4,
	},
	CfgKeys[59]: map[CfgType]int{ // notes_url
		T_HOST:              36,
		T_HOSTEXTINFO:       2,
		T_HOSTGROUP:         5,
		T_SERVICE:           36,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEEXTINFO:    3,
		T_SERVICEGROUP:      5,
	},
	CfgKeys[60]: map[CfgType]int{ // notification_failure_criteria
		T_HOSTDEPENDENCY:    6,
		T_SERVICEDEPENDENCY: 10,
	},
	CfgKeys[61]: map[CfgType]int{ // notification_interval
		T_HOST:              29,
		T_HOSTESCALATION:    6,
		T_SERVICE:           27,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEESCALATION: 7,
	},
	CfgKeys[62]: map[CfgType]int{ // notification_options
		T_HOST:              32,
		T_SERVICE:           30,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[63]: map[CfgType]int{ // notification_period
		T_HOST:              31,
		T_SERVICE:           29,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[64]: map[CfgType]int{ // notifications_enabled
		T_HOST:              33,
		T_SERVICE:           31,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[65]: map[CfgType]int{ // obsess_over_host
		T_HOST:              15,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[66]: map[CfgType]int{ // obsess_over_service
		T_SERVICE:           15,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[67]: map[CfgType]int{ // pager
		T_CONTACT:           13,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[68]: map[CfgType]int{ // parents
		T_HOST:              5,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[69]: map[CfgType]int{ // passive_checks_enabled
		T_HOST:              13,
		T_SERVICE:           13,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[70]: map[CfgType]int{ // process_perf_data
		T_HOST:              24,
		T_SERVICE:           24,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[71]: map[CfgType]int{ // retain_nonstatus_information
		T_CONTACT:           17,
		T_HOST:              26,
		T_SERVICE:           26,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[72]: map[CfgType]int{ // retain_status_information
		T_CONTACT:           16,
		T_HOST:              25,
		T_SERVICE:           25,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[73]: map[CfgType]int{ // retry_interval
		T_HOST:              11,
		T_SERVICE:           11,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[74]: map[CfgType]int{ // saturday
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_TIMEPERIOD:        2,
	},
	CfgKeys[75]: map[CfgType]int{ // service_description
		T_SERVICE:           3,
		T_SERVICEDEPENDENCY: 7,
		T_SERVICEESCALATION: 2,
		T_SERVICEEXTINFO:    1,
	},
	CfgKeys[76]: map[CfgType]int{ // service_notification_commands
		T_CONTACT:           11,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[77]: map[CfgType]int{ // service_notification_options
		T_CONTACT:           9,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[78]: map[CfgType]int{ // service_notification_period
		T_CONTACT:           7,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[79]: map[CfgType]int{ // service_notifications_enabled
		T_CONTACT:           5,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[80]: map[CfgType]int{ // servicegroup_members
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_SERVICEGROUP:      3,
	},
	CfgKeys[81]: map[CfgType]int{ // servicegroup_name
		T_SERVICEDEPENDENCY: 2,
		T_SERVICEGROUP:      0,
	},
	CfgKeys[82]: map[CfgType]int{ // servicegroups
		T_SERVICE:           5,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[83]: map[CfgType]int{ // stalking_options
		T_HOST:              34,
		T_SERVICE:           34,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[84]: map[CfgType]int{ // statusmap_image
		T_HOST:              41,
		T_HOSTEXTINFO:       7,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[85]: map[CfgType]int{ // sunday
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_TIMEPERIOD:        2,
	},
	CfgKeys[86]: map[CfgType]int{ // thursday
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_TIMEPERIOD:        2,
	},
	CfgKeys[87]: map[CfgType]int{ // timeperiod_name
	},
	CfgKeys[88]: map[CfgType]int{ // tuesday
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_TIMEPERIOD:        2,
	},
	// due to this key being forgotten, I have to adjust the value of every T_HOST/T_CONTACT/T_SERVICE by +1
	CfgKeys[89]: map[CfgType]int{ // use
		T_CONTACT:           0,
		T_HOST:              0,
		T_SERVICE:           0,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[90]: map[CfgType]int{ // vrml_image
		T_HOST:              40,
		T_HOSTEXTINFO:       6,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[91]: map[CfgType]int{ // wednesday
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
		T_TIMEPERIOD:        2,
	},
	// The following crap was not specified in Nagios Core, and is op5 specific.
	// I haven't cared to figure out sorting order, so it's arbitrary, starting on the first free slot for T_SERVICE
	CfgKeys[92]: map[CfgType]int{ // name
		T_SERVICE:           43,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[93]: map[CfgType]int{ // obsess
		T_SERVICE:           41,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[94]: map[CfgType]int{ // parallelize_check
		T_SERVICE:           40,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[95]: map[CfgType]int{ // register
		T_SERVICE:           42,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
	CfgKeys[95]: map[CfgType]int{ // hourly_value
		T_SERVICE:           44,
		T_SERVICEDEPENDENCY: 99, // value outside defined range, will not be used, only here for alignment
	},
}

var uuidorder UUIDs // append to this every time an object is read

type CfgObj struct {
	Type    CfgType           `json:"-"`
	UUID    UUID              `json:"uuid"`
	Indent  int               `json:"-"`
	Align   int               `json:"-"`
	FileID  string            `json:"fileid"`
	Comment string            `json:"-"`
	Props   map[string]string `json:"props"`
}

type CfgQuery struct {
	Keys []string
	RXs  []*regexp.Regexp
}

// Top level struct for managing collections of CfgObj
type NagiosCfg struct {
	SessionID UUID
	Config    CfgMap // the full config
	pipe      bool   // indicator of whether the content came from stdin and should be written to stdout or not
	matches   UUIDs  // subset of config
	inorder   UUIDs  // uuids ordered by how they were read in
}

//type GenericReader interface {
//	Close() error
//	GetChannel() <-chan *CfgObj
//	GetMap() (CfgMap, error)
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


