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

import (
	"bufio"
	"bytes"
	"encoding/json"
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
		log.Debugf("Attempt to add existing key %q ignored %s", key, dbgStr(false))
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
		return fmt.Errorf("Failed to append %d of the %d given values %s", errcnt, len(c2), dbgStr(true))
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
		log.Debugf("No regular expressions given %s", dbgStr(false))
		return nil
	}

	ss := func() bool {
		return subset != nil && len(subset) > 0
	}

	// no keys, but one or more RXs
	// This part seems wrong and needs to be reworked
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
	// Rework this part too
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
	// Rework this one as well
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
	// The remaining parts work quite good
	var matches UUIDs
	if !ss() {
		matches = make(UUIDs, 0, len(cm))
		keys := cm.Keys() // do this to get objects in original order, if possible
		for k := range keys {
			if cm[keys[k]].MatchSet(q) {
				//log.Debugf("%q matched %q (in: %s)", k, q, oddebug.DebugInfoMedium(PROJECT_PREFIX))
				matches = append(matches, keys[k])
			}
		}
	} else {
		matches = make(UUIDs, 0, len(subset))
		for k := range subset {
			if cm[subset[k]].MatchSet(q) {
				//log.Debugf("%q matched %q in subset (in: %s)", subset[k], q, oddebug.DebugInfoMedium(PROJECT_PREFIX))
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
	if uuidorder != nil {
		return cm.divertSearch(uuidorder, q) // this should make the search use the order given when config was read
	}
	return cm.divertSearch(nil, q)
}

// SearchSubSet searches only the CgObjs with the given UUIDs for matches
// Same underlying logic as for Search
func (cm CfgMap) SearchSubSet(q *CfgQuery, ids UUIDs) UUIDs {
	return cm.divertSearch(ids, q)
}

func (cm CfgMap) FilterType(ts ...CfgType) UUIDs {
	keys := cm.Keys() // do this to get objects in original order, if possible
	matches := make(UUIDs, 0, len(keys))
	for k := range keys {
		if cm[keys[k]].Type.In(ts) {
			matches = append(matches, keys[k])
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
	//if sort {
	//	//keys = cm.Keys().Sorted() // sorting this way does not produce the desired results, so we skip it until we have a working solution
	//	log.Debugf("Ignoring sorting of obj DB (in: %s)", oddebug.DebugInfoMedium(PROJECT_PREFIX))
	//	keys = cm.Keys()
	//} else {
	//	keys = cm.Keys()
	//}
	keys = cm.Keys() // automatically "sorted" if possible

	// Debug
	//dups1 := findDups(keys)
	//if dups1 != nil {
	//	log.Debugf("Duplicate keys: %q (in: %s)", dups1, oddebug.DebugInfoMedium(PROJECT_PREFIX))
	//}

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
		//fmap[fid] = append(fmap[fid], cm[keys[k]].UUID)
		fmap[fid] = append(fmap[fid], keys[k]) // less lookups than above
	}

	// debug dups
	//log.Debugf("fmap length: %d (in: %s)", len(fmap), oddebug.DebugInfoMedium(PROJECT_PREFIX))
	//for k := range fmap {
	//	dups := findDups(fmap[k])
	//	if dups != nil {
	//		log.Debugf("Dups in fmap[%s]: %q (in: %s)", k, fmap[k], oddebug.DebugInfoMedium(PROJECT_PREFIX))
	//	}
	//}

	return fmap
}

func (cm CfgMap) Len() int {
	return len(cm)
}

// Keys tries to deliver keys in the order they were read, otherwise it's random
func (cm CfgMap) Keys() UUIDs {
	ulen := len(uuidorder)
	clen := cm.Len()
	keys := make(UUIDs, clen)
	i := 0
	if uuidorder != nil && ulen == clen { // we assume nothing has been added or deleted if length is the same (naive...)
		copy(keys, uuidorder)
	} else if uuidorder != nil && ulen > clen { // objects have been deleted since input was read
		// here we just skip keys that are no longer present
		for k := range uuidorder {
			_, ok := cm.GetByUUID(uuidorder[k])
			if ok {
				keys[i] = uuidorder[k]
				i++
			}
		}
	} else { // give up and take Golangs random order
		for k := range cm {
			keys[i] = k
			i++
		}
	}

	// Can't figure this out in a reliable manner
	//} else if uuidorder != nil && clen > ulen { // objects have been added since input was read
	//	for k := range cm {
	//		idx := uuidorder.IndexOf(k)
	//		if idx > -1 {
	//		} else {
	//		}
	//	}
	//}
	// should still try to keep some order here, by comparing keys if uuidorder is not empty
	//for k := range cm {
	//	keys[i] = k
	//	i++
	//}

	// 2017-07-24 13:02:10: We have a problem with duplicates getting printed out when saving back,
	//  so trying to see if it's related to conversion between slices and maps
	// 2017-07-24 13:22:47 - that was not the problem, as dups still occured with below code
	//tmpmap := make(map[UUID]int)
	//tmpkeys := make(UUIDs, 0, len(keys))
	//for _, v := range keys {
	//	tmpmap[v] += 1
	//}
	//for k := range tmpmap {
	//	tmpkeys = append(tmpkeys, k)
	//}

	return keys
}

// debug dups
// mapDups searches via host_name + ; + service_description, not UUID
func (cm CfgMap) mapDups() map[string]UUIDs {
	dups := make(map[string]UUIDs)
	for u := range cm {
		if cm[u].Type != T_SERVICE {
			continue
		}
		udesc, uok := cm[u].GetUniqueCheckName()
		if !uok {
			//udesc = "INVALID_ENTRY"
			continue
		}
		_, ok := dups[udesc]
		if !ok {
			dups[udesc] = make(UUIDs, 0, 2)
		}
		dups[udesc] = append(dups[udesc], u)
	}
	// Clean up; delete entries that don't have duplicates
	for k := range dups {
		if len(dups[k]) == 1 {
			// might be overkill to both set to nil and delete, but it can't hurt
			dups[k] = nil
			delete(dups, k)
		} else {
			log.Debugf("Dups for key %q: %d %s", k, len(dups[k]), dbgStr(false))
		}
	}
	return dups
}

func (cm CfgMap) hasDups() (bool, map[string]UUIDs) {
	dupmap := cm.mapDups()
	for k := range dupmap {
		if dupmap[k] != nil && len(dupmap[k]) > 1 {
			return true, dupmap
		}
	}
	return false, dupmap
}

func (cm CfgMap) RemoveDuplicateServices(dups map[string]UUIDs) int {
	// We take the dup-map as an argument (may be nil), so that if mapDups has been run from somewhere
	// else earlier, one can pass the results as an argument in order to not have to run through
	// everything one more time
	var sdups map[string]UUIDs
	if dups == nil {
		sdups = cm.mapDups()
	} else {
		sdups = dups
	}

	num_deleted := 0
	for k := range sdups {
		for i := range sdups[k] {
			if i == 0 { // leave the first entry
				continue
			}
			obj := cm.DelByUUID(sdups[k][i])
			if obj != nil {
				num_deleted++
			}
		}
		delete(sdups, k) // this should affect map passed in, hopefully, just so it's not reused by accident
	}

	return num_deleted
}

// json stuff

func (cm CfgMap) MarshalJSON() ([]byte, error) {
	mlen := cm.Len()
	cnt := 0
	buf := bytes.NewBufferString("{")
	for k, v := range cm {
		jk, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		jv, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("%s:%s", string(jk), string(jv)))
		cnt++
		if cnt < mlen {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (cm CfgMap) UnmarshalJSON(b []byte) error {
	var tmp map[string]json.RawMessage
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	for k, v := range tmp {
		u, err := UUIDFromString(k)
		if err != nil {
			return fmt.Errorf("%s %s", err.Error(), dbgStr(true))
		}
		co := CfgObj{}
		err = json.Unmarshal(v, &co)
		if err != nil {
			return fmt.Errorf("%s %s", err.Error(), dbgStr(true))
		}
		cm[u] = &co
	}

	return nil
}
