package main

type FoundImageHook interface {
	found(pageNum int, data []byte)
}
