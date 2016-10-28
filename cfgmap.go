package nagioscfg

func (cm CfgMap) Set(key string, val *CfgObj) bool {
	_, exists := cm[key]
	cm[key] = val
	return exists
}

func (cm CfgMap) Add(key string, val *CfgObj) bool {
	_, exists := cm[key]
	if exists {
		return false
	}
	return !cm.Set(key, val)
}

func (cm CfgMap) Get(key string) (val *CfgObj, found bool) {
	val, found = cm[key]
	return
}

func (cm CfgMap) Del(key string) *CfgObj {
	val := cm[key]
	delete(cm, key)
	return val // might be nil
}
