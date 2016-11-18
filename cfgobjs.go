package nagioscfg

import (
	//"container/list"
	"regexp"
)

//func (cos CfgObjs) ToList() *list.List {
//	l := list.New()
//	for i := range cos {
//		l.PushBack(cos[i])
//	}
//	return l
//}

// MatchKeys runs MatchKeys for each obj and returns a collection of CfgObjs that match
func (cos CfgObjs) MatchKeys(rx *regexp.Regexp, keys ...string) CfgObjs {
	objlen := len(cos)
	if objlen == 0 {
		return nil
	}
	m := make(CfgObjs, 0, objlen)
	for i := range cos {
		if cos[i].MatchKeys(rx, keys...) {
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

// GetMap returns a CfgMap filtered on the given type
//func (cos CfgObjs) GetMap(typ CfgType, global bool) CfgMap {
//	if len(cos) == 0 {
//		return nil
//	}
//	objmap := make(CfgMap)
//	for i := range cos {
//		if cos[i].Type == typ {
//			var key string
//			switch typ {
//			case T_SERVICE:
//				if global {
//					ret, ok := cos[i].GetUniqueCheckName()
//					if ok {
//						key = ret
//					}
//				} else {
//					ret, ok := cos[i].GetDescription()
//					if ok {
//						key = ret
//					}
//				}
//			default:
//				ret, ok := cos[i].GetName()
//				if ok {
//					key = ret
//				}
//			}
//			if key != "" {
//				objmap[key] = cos[i]
//			}
//		}
//	}
//	return objmap
//}

// GetUUIDMap returns a CfgMap with each CfgObj's UUID as the key
//func (cos CfgObjs) GetUUIDMap() CfgMap {
//	m := make(CfgMap)
//	for i := range cos {
//		u := cos[i].GetUUID()
//		if u != nil {
//			m[u.Key()] = cos[i]
//		}
//	}
//	return m
//}

/*
// GetFilteredMap returns a map of objects matching the given filters
func (cos CfgObjs) GetFilteredMap() CfgMap {
	// I'm not taking type as an argument, as one might want to seach for stuff that can be
	// attached to several kinds of objects, like contact_groups on both hosts and services
	if len(cos) == 0 {
		return nil
	}
	//matches := make(CfgMap)
	return nil
}
*/

// GetServiceMap is a wrapper for GetMap(T_SERVICE, ...)
//func (cos CfgObjs) GetServiceMap(global bool) CfgMap {
//	return cos.GetMap(T_SERVICE, global)
//}

// GetHostMap is a wrapper for GetMap(T_HOST, ...)
//func (cos CfgObjs) GetHostMap() CfgMap {
//	return cos.GetMap(T_HOST, false)
//}

// GetCommandMap is a wrapper for GetMap(T_COMMAND, ...)
//func (cos CfgObjs) GetCommandMap() CfgMap {
//	return cos.GetMap(T_COMMAND, false)
//}

// GetContactGroupMap is a wrapper for GetMap(T_CONTACTGROUP, ...)
//func (cos CfgObjs) GetContactGroupMap() CfgMap {
//	return cos.GetMap(T_CONTACTGROUP, false)
//}

// GetContactMap is a wrapper for GetMap(T_CONTACT, ...)
//func (cos CfgObjs) GetContactMap() CfgMap {
//	return cos.GetMap(T_CONTACT, false)
//}

// GetHostGroupMap is a wrapper for GetMap(T_HOSTGROUP, ...)
//func (cos CfgObjs) GetHostGroupMap() CfgMap {
//	return cos.GetMap(T_HOSTGROUP, false)
//}

// GetServiceGroupMap is a wrapper for GetMap(T_SERVICEGROUP, ...)
//func (cos CfgObjs) GetServiceGroupMap() CfgMap {
//	return cos.GetMap(T_SERVICEGROUP, false)
//}

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
	//cos = append(cos[:index], cos[index+1:]...)

	// Should this have memory leak problems, try this instead:
	//copy((*cos)[index:], (*cos)[index+1:])
	//(*cos)[len(*cos)-1] = nil
	//(*cos) = (*cos)[:len(*cos)-1]

	// Without preserving order:
	// This is after benchmarking the by far most efficient method, even better than using container/List,
	// which is almost as fast
	(*cos)[index] = (*cos)[len(*cos)-1]
	(*cos)[len(*cos)-1] = nil
	(*cos) = (*cos)[:len(*cos)-1]

	// This version uses a lot more memory (~1000x), but is much faster if deleting from the beginning of the slice (~3x)
	//o := make(CfgObjs, 0, len(*cos)-1)
	//for i := range *cos {
	//	if i != index {
	//		o = append(o, (*cos)[i])
	//	}
	//}
	//*cos = o
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
