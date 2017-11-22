package main

import "log"

type voidFoundImageHook struct {
}

func (vfih voidFoundImageHook) found(pageNum int, data []byte) {
	log.Printf("Found an image")
}
