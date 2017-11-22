package main

type serviceFoundImageHook struct {
	Service   BookService
	ChapterID int
}

func newServiceFoundImageHook(bookName string, chapterNumber int) serviceFoundImageHook {
	bookService := NewBookService()

	bookID := bookService.getBookID(bookName)
	chapterID := bookService.getChapterID(bookID, chapterNumber)

	hook := serviceFoundImageHook{
		Service:   bookService,
		ChapterID: chapterID,
	}
	return hook
}

func (sfih serviceFoundImageHook) found(pageNum int, data []byte) {
	sfih.Service.postImage(sfih.ChapterID, pageNum, data)
}
