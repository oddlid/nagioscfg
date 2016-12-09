package nagioscfg

import (
	"bufio"
	"bytes"
	"fmt"
)

func (cm CfgMap) SetByUUID(key UUID, val *CfgObj) bool {
	_, exists := cm[key]
	cm[key] = val
	return exists
}

func (cm CfgMap) Set(key string, val *CfgObj) bool {
	u, err := UUIDFromString(key)
	if err != nil {
		return false
	}
	return cm.SetByUUID(u, val)
}

func (cm CfgMap) AddByUUID(key UUID, val *CfgObj) bool {
	_, exists := cm[key]
	if exists {
		return false
	}
	return !cm.SetByUUID(key, val)
}

func (cm CfgMap) Add(key string, val *CfgObj) bool {
	u, err := UUIDFromString(key)
	if err != nil {
		return false
	}
	return cm.AddByUUID(u, val)
}

func (cm CfgMap) GetByUUID(key UUID) (val *CfgObj, found bool) {
	val, found = cm[key]
	return
}

func (cm CfgMap) Get(key string) (val *CfgObj, found bool) {
	u, err := UUIDFromString(key)
	if err != nil {
		return nil, false
	}
	return cm.GetByUUID(u)
}

func (cm CfgMap) DelByUUID(key UUID) *CfgObj {
	val := cm[key]
	delete(cm, key)
	return val // might be nil
}

func (cm CfgMap) Del(key string) *CfgObj {
	u, err := UUIDFromString(key)
	if err != nil {
		return nil
	}
	return cm.DelByUUID(u)
}

func (cm CfgMap) LongestKey() int {
	max := 0
	curmax := 0
	for _, v := range cm {
		curmax = v.LongestKey()
		if curmax > max {
			max = curmax
		}
	}
	return max
}

func (cm CfgMap) Dump() string {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	for k, v := range cm {
		fmt.Fprintf(w, "Key     : %q\n", k)
		fmt.Fprintf(w, "UUID    : %s\n", v.UUID.String())
		fmt.Fprintf(w, "Type    : %s\n", v.Type.String())
		fmt.Fprintf(w, "Indent  : %d\n", v.Indent)
		fmt.Fprintf(w, "Align   : %d\n", v.Align)
		fmt.Fprintf(w, "File ID : %s\n", v.FileID)
		v.Print(w)
	}
	w.Flush()
	return buf.String()
}
