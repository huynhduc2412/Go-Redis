package core

import datastructure "Go-Redis/internal/data_structure"

var dicStore *datastructure.Dict

func init() {
	dicStore = datastructure.CreateDict()
}