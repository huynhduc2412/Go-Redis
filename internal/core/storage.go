package core

import datastructure "Go-Redis/internal/data_structure"

var dicStore *datastructure.Dict
var setStore map[string]*datastructure.SimpleSet
func init() {
	dicStore = datastructure.CreateDict()
	setStore = make(map[string]*datastructure.SimpleSet)
}