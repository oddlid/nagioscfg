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
		log.Debugf("%s.CfgMap.AddByUUID(): Attempt to add existing key %q ignored", PKGNAME, key)
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

func (cm CfgMap) SetKeys(ids UUIDs, keys, values []string) int {
	modcnt := 0
	if ids == nil || len(ids) == 0 {
		for k := range cm {
			modcnt += cm[k].SetKeys(keys, values)
		}
	} else {
		for i := range ids {
			modcnt += cm[ids[i]].SetKeys(keys, values)
		}
	}
	return modcnt
}

func (cm CfgMap) DelKeys(ids UUIDs, keys []string) int {
	delcnt := 0
	if ids == nil || len(ids) == 0 {
		return delcnt
	} else {
		for i := range ids {
			delcnt += cm[ids[i]].DelKeys(keys)
		}
	}
	return delcnt
}

func (cm CfgMap) Append(c2 CfgMap) error {
	errcnt := 0
	for k := range c2 {
		ok := cm.AddByUUID(k, c2[k])
		if !ok {
			errcnt++
		}
	}
	if errcnt > 0 {
		return fmt.Errorf("%s.CfgMap.Append(): Failed to append %d of the %d given values", PKGNAME, errcnt, len(c2))
	}
	return nil
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
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "File ID : %s\n", v.FileID)
		fmt.Fprintf(w, "Key     : %s\n", k)
		//fmt.Fprintf(w, "UUID    : %s\n", v.UUID.String())
		fmt.Fprintf(w, "Type    : %s\n", v.Type.String())
		fmt.Fprintf(w, "Indent  : %d\n", v.Indent)
		fmt.Fprintf(w, "Align   : %d\n", v.Align)
		v.Print(w, false) // as this is mostly for debugging, don't sort by default
	}
	w.Flush()
	return buf.String()
}

func (cm CfgMap) divertMatchAllKeys(ids UUIDs, rx *regexp.Regexp, keys []string) UUIDs {
	var matches UUIDs
	if ids == nil || len(ids) == 0 {
		matches = make(UUIDs, 0, len(cm))
		for k := range cm {
			if cm[k].MatchAllKeys(rx, keys...) {
				matches = append(matches, k)
			}
		}
	} else {
		matches = make(UUIDs, 0, len(ids))
		for i := range ids {
			if cm[ids[i]].MatchAllKeys(rx, keys...) {
				matches = append(matches, ids[i])
			}
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

// MatchAllKeys matches on the keys of each CfgObj, NOT on the CfgMap keys. Returns slice of UUIDS
func (cm CfgMap) MatchAllKeys(rx *regexp.Regexp, keys ...string) UUIDs {
	return cm.divertMatchAllKeys(nil, rx, keys)
}

func (cm CfgMap) MatchAllKeysSubSet(ids UUIDs, rx *regexp.Regexp, keys ...string) UUIDs {
	return cm.divertMatchAllKeys(ids, rx, keys)
}

func (cm CfgMap) divertMatchAnyKeys(ids UUIDs, rx *regexp.Regexp, keys []string) UUIDs {
	var matches UUIDs
	if ids == nil || len(ids) == 0 {
		matches = make(UUIDs, 0, len(cm))
		for k := range cm {
			if cm[k].MatchAnyKeys(rx, keys...) {
				matches = append(matches, k)
			}
		}
	} else {
		matches = make(UUIDs, 0, len(ids))
		for i := range ids {
			if cm[ids[i]].MatchAnyKeys(rx, keys...) {
				matches = append(matches, ids[i])
			}
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

func (cm CfgMap) MatchAnyKeys(rx *regexp.Regexp, keys ...string) UUIDs {
	return cm.divertMatchAnyKeys(nil, rx, keys)
}

func (cm CfgMap) MatchAnyKeysSubSet(ids UUIDs, rx *regexp.Regexp, keys ...string) UUIDs {
	return cm.divertMatchAnyKeys(ids, rx, keys)
}

func (cm CfgMap) divertMatchAny(ids UUIDs, rx *regexp.Regexp) UUIDs {
	var matches UUIDs
	if ids == nil || len(ids) == 0 {
		matches = make(UUIDs, 0, len(cm))
		for k := range cm {
			if cm[k].MatchAny(rx) {
				matches = append(matches, k)
			}
		}
	} else {
		matches = make(UUIDs, 0, len(ids))
		for i := range ids {
			if cm[ids[i]].MatchAny(rx) {
				matches = append(matches, ids[i])
			}
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

func (cm CfgMap) MatchAny(rx *regexp.Regexp) UUIDs {
	return cm.divertMatchAny(nil, rx)
}

func (cm CfgMap) MatchAnySubSet(rx *regexp.Regexp, ids UUIDs) UUIDs {
	return cm.divertMatchAny(ids, rx)
}

func (cm CfgMap) divertSearch(subset UUIDs, q *CfgQuery) UUIDs {
	klen := len(q.Keys)
	rlen := len(q.RXs)

	// no RXs given
	if rlen == 0 {
		log.Debugf("%s.CfgMap.Search(): No regular expressions given", PKGNAME)
		return nil
	}

	ss := func() bool {
		return subset != nil
	}

	// no keys, but one or more RXs
	if klen == 0 {
		var m UUIDs
		if ss() {
			m = cm.MatchAnySubSet(q.RXs[0], subset)
		} else {
			m = cm.MatchAny(q.RXs[0])
		}
		for i := 1; i < rlen; i++ {
			m = cm.MatchAnySubSet(q.RXs[i], m)
		}
		return m // may be nil or zero length here
	}
	// one or more keys, and one or more RXs, but maybe not the same amount of each
	// ... more keys than RXs
	if klen > rlen {
		var m UUIDs
		if ss() {
			m = cm.MatchAnyKeysSubSet(subset, q.RXs[0], q.Keys...)
		} else {
			m = cm.MatchAnyKeys(q.RXs[0], q.Keys...)
		}
		for i := 1; i < rlen; i++ {
			m = cm.MatchAnyKeysSubSet(m, q.RXs[i], q.Keys...)
		}
		return m
	}
	// ... more RXs than keys
	if rlen > klen {
		var m UUIDs
		if ss() {
			m = cm.MatchAllKeysSubSet(subset, q.RXs[0], q.Keys...)
		} else {
			m = cm.MatchAllKeys(q.RXs[0], q.Keys...)
		}
		for i := 1; i < rlen; i++ {
			m = cm.MatchAllKeysSubSet(m, q.RXs[i], q.Keys...)
		}
		return m
	}
	// one or more, and the same amount of keys and RXs
	var matches UUIDs
	if !ss() {
		matches = make(UUIDs, 0, len(cm))
		for k := range cm {
			if cm[k].MatchSet(q) {
				matches = append(matches, k)
			}
		}
	} else {
		matches = make(UUIDs, 0, len(subset))
		for k := range subset {
			if cm[subset[k]].MatchSet(q) {
				matches = append(matches, subset[k])
			}
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

// Search allows more complex searches for matching CfgObjs in the CfgMap
// The strategy chosen depends on how the keys and regexes provided are balanced.
// Given any keys but no RXs, it will return nil.
// Given no keys and any RXs, it will return all objects that match all RXs on the value of any key.
// Given more keys than RXs, it will return all objects that match all RXs on any of the keys.
// Given more RXs than keys, it will return all objects that match all RXs on all of the keys.
// Given an equal amount of keys and RXs, it will return all objects that match RX on the value of the corresponding key, in given order.
func (cm CfgMap) Search(q *CfgQuery) UUIDs {
	return cm.divertSearch(nil, q)
}

// SearchSubSet searches only the CgObjs with the given UUIDs for matches
// Same underlying logic as for Search
func (cm CfgMap) SearchSubSet(q *CfgQuery, ids UUIDs) UUIDs {
	return cm.divertSearch(ids, q)
}

func (cm CfgMap) FilterType(ts ...CfgType) UUIDs {
	matches := make(UUIDs, 0, len(cm))
	for k := range cm {
		if cm[k].Type.In(ts) {
			matches = append(matches, k)
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
}

func (cm CfgMap) SplitByFileID(sort bool) map[string]UUIDs {
	fmap := make(map[string]UUIDs)
	var keys UUIDs
	if sort {
		keys = cm.Keys().Sorted()
	} else {
		keys = cm.Keys()
	}
	for k := range keys {
		fid := cm[keys[k]].FileID
		// skip etries without fileID
		if fid == "" {
			continue
			// or we could do something like this:
			//fid = "nagios-config-objects-without-fileid.cfg"
		}
		// initialize if first match of this ID
		_, ok := fmap[fid]
		if !ok {
			fmap[fid] = make(UUIDs, 0, 1)
		}
		fmap[fid] = append(fmap[fid], cm[keys[k]].UUID)
	}
	return fmap
}

func (cm CfgMap) Keys() UUIDs {
	keys := make(UUIDs, len(cm))
	i := 0
	for k := range cm {
		keys[i] = k
		i++
	}
	return keys
}
