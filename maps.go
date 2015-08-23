package main

var Levels [20]*Map

func initMaps() {
	Levels[0] = NewMap("level1")
	for i, lvl := range Levels {
		if lvl == nil {
			continue
		}
		lvl.ID = i
	}
}
