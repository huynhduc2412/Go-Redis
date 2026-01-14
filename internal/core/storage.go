package core

import data_structure "Go-Redis/internal/data_structure"

var dicStore *data_structure.Dict
var setStore map[string]*data_structure.SimpleSet
var zsetStore map[string]*data_structure.ZSet
var cmsStore map[string]*data_structure.CMS
func init() {
	dicStore = data_structure.CreateDict()
	setStore = make(map[string]*data_structure.SimpleSet)
	zsetStore = make(map[string]*data_structure.ZSet)
	cmsStore = make(map[string]*data_structure.CMS)
}