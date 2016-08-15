package nagioscfg

/*
Defines the data structures that maps to Nagios config items
Odd, 2016-08-10 17:51
*/


import (
)

type Key interface {
	Name() string
}

type Value interface {
	String() string
	Int() int
	Set(val string)
}

type ConfigItem interface {
	Key() string
	Value() string
	String(pad, align int) string
	Set(key,  val string)
}

type Cfg struct {
}

type CfgItem struct {
	Key Key
	Val Value
}

type CfgObj struct {
	Type string
	Items []ConfigItem
}
