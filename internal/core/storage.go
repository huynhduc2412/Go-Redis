package core

import data_structure "Go-Redis/internal/data_structure"

var dicStore *data_structure.Dict
var setStore map[string]*data_structure.SimpleSet
var zsetStore map[string]*data_structure.ZSet
func init() {
	dicStore = data_structure.CreateDict()
	setStore = make(map[string]*data_structure.SimpleSet)
	zsetStore = make(map[string]*data_structure.ZSet)
}