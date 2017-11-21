package main

type ServiceFoundImageHook struct {
	Service   BookService
	ChapterID int
}

func NewServiceFoundImageHook(service BookService, chapterID int) ServiceFoundImageHook {
	hook := ServiceFoundImageHook{
		Service:   service,
		ChapterID: chapterID,
	}
	return hook
}

func (sfih ServiceFoundImageHook) found(pageNum int, data []byte) {
	sfih.Service.postImage(sfih.ChapterID, pageNum, data)
}
