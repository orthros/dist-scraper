package main

type ServiceFoundImageHook struct {
	Service   BookService
	ChapterID int
}

func NewServiceFoundImageHook(bookName string, chapterNumber int) ServiceFoundImageHook {
	bookService := NewBookService()

	bookID := bookService.getBookID(bookName)
	chapterID := bookService.getChapterID(bookID, chapterNumber)

	hook := ServiceFoundImageHook{
		Service:   bookService,
		ChapterID: chapterID,
	}
	return hook
}

func (sfih ServiceFoundImageHook) found(pageNum int, data []byte) {
	sfih.Service.postImage(sfih.ChapterID, pageNum, data)
}
