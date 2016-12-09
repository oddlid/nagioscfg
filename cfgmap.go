package nagioscfg

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"regexp"
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

// MatchKeys matches on the keys of each CfgObj, NOT on the CfgMap keys. Returns slice of UUIDS
func (cm CfgMap) MatchKeys(rx *regexp.Regexp, keys ...string) []UUID {
	matches := make([]UUID, 0, len(cm))
	for k := range cm {
		if cm[k].MatchKeys(rx, keys...) {
			matches = append(matches, k)
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

func (cm CfgMap) MatchAny(rx *regexp.Regexp) []UUID {
	matches := make([]UUID, 0, len(cm))
	for k := range cm {
		if cm[k].MatchAny(rx) {
			matches = append(matches, k)
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

// Search allows more complex searches for matching CfgObjs in the CfgMap
func (cm CfgMap) Search(q *CfgQuery) []UUID {
	// Make sure there is a regexp for each given key
	if !q.Balanced() {
		log.Debug("CfgMap.Search(): number of keys and regexes in given CfgQuery does not match")
		return nil
	}
	matches := make([]UUID, 0, len(cm))
	for k := range cm {
		if cm[k].MatchAll(q) {
			matches = append(matches, k)
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

// SearchSubSet searches only the CgObjs with the given UUIDs for matches
func (cm CfgMap) SearchSubSet(q *CfgQuery, ids ...UUID) []UUID {
	// Make sure there is a regexp for each given key
	if !q.Balanced() {
		log.Debug("CfgMap.SearchSubSet(): number of keys and regexes in given CfgQuery does not match")
		return nil
	}
	matches := make([]UUID, 0, len(ids))
	for i := range ids {
		if cm[ids[i]].MatchAll(q) {
			matches = append(matches, ids[i])
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

func (cm CfgMap) FilterType(t CfgType) []UUID {
	return nil
}
