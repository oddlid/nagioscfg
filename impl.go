/*
Function/method implementations for types from data.go
*/

package nagioscfg

import (
	log "github.com/Sirupsen/logrus"
	"regexp"
)

// String returns the string representation of the CfgType
func (ct CfgType) String() string {
	return string(CfgTypes[ct])
}

// Type returns the int (CfgType) value for the given CfgName, or -1 if not valid
func (cn CfgName) Type() CfgType {
	for i := range CfgTypes {
		if CfgTypes[i] == cn {
			return CfgType(i)
		}
	}
	return -1
}

// size returns the runtime bytes size for the given objects map ( to calculate objs from input file size). Only for debugging.
/*
func (co *CfgObj) size() int {
	var size int
	for k, v := range co.Props {
		size += co.Indent + (co.Align - len(k)) + len(v)
	}
	size += 64 // approx buffer for comments etc.
	return size
}
*/

func NewCfgQuery() *CfgQuery {
	return &CfgQuery{
		Keys: make([]string, 0, 2),
		RXs:  make([]*regexp.Regexp, 0, 2),
	}
}

// Balanced() verifies that there is a matching number of keys and regexes in the instance
func (cq CfgQuery) Balanced() bool {
	return len(cq.Keys) == len(cq.RXs)
}

func (cq *CfgQuery) AddFilter(key, re string) bool {
	if key == "" {
		log.Error("CfgQuery.AddFilter(): Error: Empty key")
		return false
	}
	rx, err := regexp.Compile(re)
	if err != nil {
		log.Error(err)
		return false
	}

	cq.Keys = append(cq.Keys, key)
	cq.RXs = append(cq.RXs, rx)

	return true
}
