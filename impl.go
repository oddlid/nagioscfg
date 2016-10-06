/*
Function/method implementations for types from data.go
*/

package nagioscfg

import (
	"fmt"
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

// Set adds the given key/value to CfgObj.Props, returning true if the key was overwritten, and false if it was added fresh
func (co *CfgObj) Set(key, val string) bool {
	_, exists := co.Props[key]
	co.Props[key] = val
	return exists // true = key was overwritten, false = key was added
}

// Add adds the given key/value to CfgObj.Props only if the key does not already exist. Returns true if added, false otherwise.
func (co *CfgObj) Add(key, val string) bool {
	_, exists := co.Props[key]
	if exists {
		return false
	}
	return !co.Set(key, val) // Set should return false, as the key doesn't exist yet, so we inverse the result
}

// Get returns the value for the given key, if it exists. "found" will be false if no such key exists.
func (co *CfgObj) Get(key string) (val string, found bool) {
	val, found = co.Props[key]
	return val, found
}

// Del deletes the entry with the given key. It returns true if anything was deleted, false otherwise.
func (co *CfgObj) Del(key string) bool {
	_, exists := co.Props[key]
	delete(co.Props, key)
	return exists // just signals if there was anything there to be deleted in the first place
}

// LongestKey returns the length of the longest key in CfgObj.Props at the time of calling
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

// GetList gets a value from CfgObj.Props and returns a string slice after splitting the value on the separator given
func (co *CfgObj) GetList(key, sep string) []string {
	val, exists := co.Get(key)
	if !exists {
		return nil
	}
	return strings.Split(val, sep)
}

// SetList takes a slice and joins it using the given separator, then sets it as the value for the given key
func (co *CfgObj) SetList(key, sep string, list ...string) bool {
	lstr := strings.Join(list, sep)
	return co.Set(key, lstr)
}

// AddList does the same as SetList, but only if the key does not already exist
func (co *CfgObj) AddList(key, sep string, list ...string) bool {
	_, exists := co.Props[key]
	if exists {
		return false
	}
	return !co.SetList(key, sep, list...) // SetList should return false as key does not exist, so invert the result
}

// GetCheckCommand returns the list value for check_command in a service object
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

// GetCheckCommandCmd returns the command name part from GetCheckCommand
func (co *CfgObj) GetCheckCommandCmd() (string, bool) {
	lst := co.GetCheckCommand()
	if lst == nil {
		return "", false
	}
	return lst[0], true
}

// GetCheckCommandArgs returns the argument list part from GetCheckCommand
func (co *CfgObj) GetCheckCommandArgs() []string {
	lst := co.GetCheckCommand()
	if lst == nil {
		return nil
	}
	return lst[1:]
}

// GetName tries to return the name for the given object, if set
func (co *CfgObj) GetName() (string, bool) {
	key := co.Type.String() + "_name"
	name, found := co.Get(key)
	if !found {
		return co.Get("name")
	}
	return name, found
}

// GetDescription tries to get the description for the given object, if set
func (co *CfgObj) GetDescription() (string, bool) {
	key := co.Type.String() + "_description"
	return co.Get(key)
}

// generateComment is set as private, as it makes "unsafe" assumptions about the existing format of the comment
func (co *CfgObj) generateComment() bool {
	var name string
	var success bool
	var is_template bool
	if co.Type == T_SERVICE {
		name, success = co.GetDescription()
		if !success {
			name, success = co.GetName() // in case it's a template, not a real service
			is_template = success
		}
	} else {
		name, success = co.GetName()
	}
	if success && strings.Index(co.Comment, "%") > -1 {
		if is_template {
			co.Comment = fmt.Sprintf("# %s template '%s'", co.Type.String(), name)
		} else {
			co.Comment = fmt.Sprintf(co.Comment, name)
		}
	}
	return success
}

// AutoAlign sets the CfgObj alignment/spacing to LongestKey + 2
func (co *CfgObj) AutoAlign() int {
	co.Align = co.LongestKey() + 2
	return co.Align
}

// LongestKey returns the length of the longest key in a collection of CfgObj
func (cos CfgObjs) LongestKey() int {
	max := 0
	for i := range cos {
		curmax := cos[i].LongestKey()
		if curmax > max {
			max = curmax
		}
	}
	return max
}

// AutoAlign sets the alignment for a collection of CfgObj
func (cos CfgObjs) AutoAlign() int {
	align := cos.LongestKey() + 2
	for i := range cos {
		cos[i].Align = align
	}
	return align
}

func (cos CfgObjs) Add(co *CfgObj) {
	// Should have some duplicate checking here
	cos = append(cos, *co)
}

// Del deletes an object from CfgObjs based on index
// See also: https://github.com/golang/go/wiki/SliceTricks
func (cos CfgObjs) Del(index int) {
	cos = append(cos[:index], cos[index+1:]...)
	// Should this have memory leak problems, try this instead:
	/*
	copy(cos[i:], cos[i+1:])
	cos[len(cos)-1] = nil // or CfgObj{} instead of nil
	cos = cos[:len(cos)-1]
	*/
}

// Find returns a collection of CfgObj based on a string match
func (cos CfgObjs) Find(match string) (CfgObjs, error) {
	return nil, nil
}
