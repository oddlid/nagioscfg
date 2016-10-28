/*
Function/method implementations for types from data.go
*/

package nagioscfg

//import (
//)

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
