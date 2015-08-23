package main

import "github.com/thinkofdeath/monstergame/html"

var (
	Images = map[string]*html.Image{
		"tiles":  nil,
		"player": nil,
		"level1": nil,
	}
	loadedImages = 0
)

func loadImages() {
	for k := range Images {
		img := html.NewImage("images/"+k+".png", imageLoad)
		Images[k] = img
	}
}

func imageLoad() {
	loadedImages++
	if loadedImages == len(Images) {
		realMain()
	}
}
