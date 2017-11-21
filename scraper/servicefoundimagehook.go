package main

type ServiceFoundImageHook struct {
	Service   BookService
	ChapterID int
}

func (sfih ServiceFoundImageHook) found(pageNum int, data []byte) {
	sfih.Service.postImage(sfih.ChapterID, pageNum, data)
}
