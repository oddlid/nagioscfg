package nagioscfg

/*
Defines the data structures that maps to Nagios config items
Odd, 2016-08-10 17:51
*/


import (
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
	Type string
	Items map[string]string
}

// methods - move to separate file when it grows

func (ct CfgType) String() string {
	return CfgTypes[ct]
}
