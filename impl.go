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

/*
Function/method implementations for types from data.go
*/

package nagioscfg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/oddebug"
	"os"
	"regexp"
)


func NewNagiosCfg() *NagiosCfg {
	return &NagiosCfg{
		SessionID: NewUUIDv1(),
		Config:    make(CfgMap),
	}
}

func (nc *NagiosCfg) LoadFiles(files ...string) error {
	mfr := NewMultiFileReader(files...)
	defer mfr.Close()
	in := mfr.ReadChan(true)
	cm := make(CfgMap)
	for o := range in {
		cm[o.UUID] = o
	}
	nc.Config = cm
	nc.pipe = false
	return nil // can change later if we use another way to read to map
}

//func (nc *NagiosCfg) LoadFiles(files ...string) error {
//	// Testing a variant that does not read via channels in parallell
//	// Only for debugging duplicate entries @2017-07-24 18:58:16
//	cm := make(CfgMap)
//	tmpcm := make(CfgMap)
//	for i := range files {
//		log.Debugf("Reading file %q (in: %s)", files[i], oddebug.DebugInfoMedium(PROJECT_PREFIX))
//		rdr := NewFileReader(files[i])
//		fileID, err := rdr.AbsPath()
//		if err != nil {
//			log.Errorf("%q (in: %s)", err, oddebug.DebugInfoMedium(PROJECT_PREFIX))
//			fileID = files[i]
//		}
//		tmpcm, err = rdr.ReadAllMap(fileID)
//		rdr.Close()
//		if err != nil {
//			log.Errorf("%s (in: %s)", err.Error(), oddebug.DebugInfoMedium(PROJECT_PREFIX))
//			return err
//		}
//		for k := range tmpcm {
//			cm[k] = tmpcm[k]
//		}
//	}
//	nc.Config = cm
//	nc.pipe = false
//	return nil
//}

func (nc *NagiosCfg) LoadStdin() (err error) {
	rdr := NewReader(os.Stdin)
	nc.Config, err = rdr.ReadAllMap("")
	nc.pipe = true // indicator that all content came from stdin and that we don't have any FileIDs
	return err
}

func (nc *NagiosCfg) DumpStdout() {
	nc.Print(os.Stdout, true) // sort by default
}

func (nc *NagiosCfg) InPipe() bool {
	return nc.pipe
	// could use this instead
	//fi, _ := os.Stdin.Stat()
	//return (fi.Mode() & os.ModeCharDevice) == 0
}

func (nc *NagiosCfg) FilterType(ts ...CfgType) UUIDs {
	m := nc.Config.FilterType(ts...)
	if m == nil {
		m = make(UUIDs, 0)
	}
	nc.matches = m
	return nc.matches
}

func (nc *NagiosCfg) Search(q *CfgQuery) UUIDs {
	if !nc.matches.Empty() {
		nc.matches = nc.Config.SearchSubSet(q, nc.matches)
	} else {
		nc.matches = nc.Config.Search(q)
	}
	return nc.matches
}

func (nc *NagiosCfg) InverseResults() UUIDs {
	if nc.matches.Empty() {
		return uuidorder // if previous search yielded nothing, then everything is the inverse
	}
	inv := make(UUIDs, 0, nc.Config.Len() - nc.matches.Len())
	for _, v := range uuidorder {
		if !v.In(nc.matches) { // this is probably slow
			inv = append(inv, v)
		}
	}
	nc.matches = inv
	return nc.matches
}

func (nc *NagiosCfg) Len() int {
	return nc.Config.Len()
}

func (nc *NagiosCfg) GetMatches() UUIDs {
	return nc.matches
}

func (nc *NagiosCfg) ClearMatches() {
	nc.matches = nil
}

func (nc *NagiosCfg) DeleteMatches() CfgMap {
	if nc.matches.Empty() {
		return nil
	}
	cm := make(CfgMap)
	for i := range nc.matches {
		cm[nc.matches[i]] = nc.Config.DelByUUID(nc.matches[i])
	}
	nc.ClearMatches()
	return cm
}

func (nc *NagiosCfg) DelKeys(keys []string) int {
	return nc.Config.DelKeys(nc.matches, keys)
}

func (nc *NagiosCfg) SetKeys(keys, values []string) int {
	return nc.Config.SetKeys(nc.matches, keys, values)
}

func (nc *NagiosCfg) HasServiceDups() (bool, map[string]UUIDs) {
	return nc.Config.hasDups()
}

// Valid checks if the given CfgType is within valid range
func (ct CfgType) Valid() bool {
	return ct >= T_COMMAND && ct < T_INVALID
}

// String returns the string representation of the CfgType
func (ct CfgType) String() string {
	if !ct.Valid() {
		return "INVALID_TYPE"
	}
	return string(CfgTypes[ct])
}

func (ct CfgType) In(types []CfgType) bool {
	for i := range types {
		if types[i] == ct {
			return true
		}
	}
	return false
}

// Type returns the int (CfgType) value for the given CfgName, or -1 if not valid
func (cn CfgName) Type() CfgType {
	for i := range CfgTypes {
		//log.Debugf("%s.CfgName.Type(): trying index #%d", PKGNAME, i)
		if CfgTypes[i] == cn {
			//log.Debugf("%s.CfgName.Type(): match at index #%d", PKGNAME, i)
			return CfgType(i)
		}
	}
	return T_INVALID
}

func (cn CfgName) Valid() bool {
	return cn.Type() != T_INVALID
}

func IsValidProperty(key string) bool {
	_, ok := CfgKeySortOrder[key]
	return ok
}

func ValidCfgNames() []string {
	l := len(CfgTypes)
	s := make([]string, l)
	for i := range CfgTypes {
		s[i] = string(CfgTypes[i])
	}
	return s
}

//func (cp CfgProp) Valid() bool {
//	_, ok := CfgKeySortOrder[string(cp)]
//	return ok
//}

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

func (cq *CfgQuery) AddRX(re string) bool {
	rx, err := regexp.Compile(re)
	if err != nil {
		log.Errorf("%q (in: %s)", err, oddebug.DebugInfoMedium(PROJECT_PREFIX))
		return false
	}
	cq.RXs = append(cq.RXs, rx)
	return true
}

func (cq *CfgQuery) AddKey(key string) bool {
	if key != "" { // won't accept empty keys
		if IsValidProperty(key) { // only accept defined keys/properties
			cq.Keys = append(cq.Keys, key)
			return true
		}
	}
	log.Errorf("Invalid key: %q (in: %s)", key, oddebug.DebugInfoMedium(PROJECT_PREFIX))
	return false
}

func (cq *CfgQuery) AddKeyRX(key, re string) bool {
	if key == "" {
		log.Error("CfgQuery.AddFilter(): Error: Empty key")
		return false
	}

	if !IsValidProperty(key) {
		log.Errorf("Invalid key: %q (in: %s)", key, oddebug.DebugInfoMedium(PROJECT_PREFIX))
		return false
	}

	return cq.AddRX(re) && cq.AddKey(key)
}

// for debugging only
//func findDups(u UUIDs) UUIDs {
//	var ret UUIDs
//	dmap := make(map[UUID]int)
//	for i := range u {
//		dmap[u[i]] += 1
//	}
//	for k := range dmap {
//		if dmap[k] > 1 {
//			if ret == nil {
//				ret = make(UUIDs, 0, len(u))
//			}
//			ret = append(ret, k)
//		}
//	}
//	return ret
//}
