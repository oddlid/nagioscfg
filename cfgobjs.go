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
	"regexp"
)

// MatchKeys runs MatchKeys for each obj and returns a collection of CfgObjs that match
func (cos CfgObjs) MatchKeys(rx *regexp.Regexp, keys ...string) CfgObjs {
	objlen := len(cos)
	if objlen == 0 {
		return nil
	}
	m := make(CfgObjs, 0, objlen)
	for i := range cos {
		if cos[i].MatchAllKeys(rx, keys...) {
			m = append(m, cos[i])
		}
	}
	if len(m) > 0 {
		return m
	}
	return nil
}

// MatchAny runs MatchAny for each obj and returns a collection of CfgObjs that match
func (cos CfgObjs) MatchAny(rx *regexp.Regexp) CfgObjs {
	objlen := len(cos)
	if objlen == 0 {
		return nil
	}
	m := make(CfgObjs, 0, objlen)
	for i := range cos {
		if cos[i].MatchAny(rx) {
			m = append(m, cos[i])
		}
	}
	if len(m) > 0 {
		return m
	}
	return nil
}

// LongestKey returns the length of the longest key in a collection of CfgObj
func (cos CfgObjs) LongestKey() int {
	max := 0
	curmax := 0
	for i := range cos {
		curmax = cos[i].LongestKey()
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

// Add appends an object to CfgObjs
func (cos *CfgObjs) Add(co *CfgObj) {
	// Should have some duplicate checking here
	*cos = append(*cos, co)
}

// Del deletes an object from CfgObjs based on index
func (cos *CfgObjs) Del(index int) {
	// Without preserving order:
	// This is after benchmarking the by far most efficient method, even better than using container/List,
	// which is almost as fast
	(*cos)[index] = (*cos)[len(*cos)-1]
	(*cos)[len(*cos)-1] = nil
	(*cos) = (*cos)[:len(*cos)-1]
}

// DelUUID deletes the object with a matching UUID. Does not keep slice in order.
func (cos *CfgObjs) DelUUID(u UUID) {
	for i := range *cos {
		if (*cos)[i].UUID.Equals(u) {
			(*cos)[i] = (*cos)[len(*cos)-1]
			(*cos)[len(*cos)-1] = nil
			(*cos) = (*cos)[:len(*cos)-1]
		}
	}
}
