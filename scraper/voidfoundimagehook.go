package main

import "log"

type VoidFoundImageHook struct {
}

func (vfih VoidFoundImageHook) found(pageNum int, data []byte) {
	log.Printf("Found an image")
}
