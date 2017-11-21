using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using books_service.Models;
using Marten;
using System.IO;

namespace books_service.Controllers
{
    [Produces("application/json")]
    [Route("api/Books")]
    public class BooksController : Controller
    {
        public IDocumentStore Store { get; }

        public BooksController(IDocumentStore store)
        {
            Store = store;
        }

        // GET: api/Books
        [HttpGet]
        public async Task<IEnumerable<Book>> Get()
        {
            using (var session = Store.LightweightSession())
            {
                return await session.Query<Book>()
                                    .ToListAsync();
            }                
        }

        // GET: api/Books/5
        [HttpGet("{id}", Name = "Get")]
        public async Task<Book> Get(int id)
        {
            using (var session = Store.LightweightSession())
            {
                return await session.Query<Book>()
                                    .FirstOrDefaultAsync(x => x.Id == id);
            }                
        }

        [HttpGet("{id}/chapters")]
        public async Task<IEnumerable<Chapter>> GetChapters(int id)
        {
            using (var session = Store.LightweightSession())
            {
                return await session.Query<Chapter>()
                                    .Where(x => x.BookID == id)
                                    .ToListAsync();
            }
        }

        [HttpPost("{id}/chapters")]
        public async Task<int> PostChapter(int id, [FromBody]Chapter chapter)
        {
            using (var session = Store.OpenSession())
            {
                var existingBook = await session.Query<Book>()
                                                .Where(x => x.Id == id)
                                                .FirstOrDefaultAsync();
                if (existingBook == null) { return -1; }

                var foundChapter = await session.Query<Chapter>()
                                                   .Where(x => x.ChapterNumber == chapter.ChapterNumber)
                                                   .FirstOrDefaultAsync();
                if (foundChapter == null)
                {
                    foundChapter = new Chapter()
                    {
                        BookID = id,
                        ChapterNumber = chapter.ChapterNumber,
                        ChapterTitle = chapter.ChapterTitle
                    };
                    session.Store(foundChapter);

                    await session.SaveChangesAsync();
                }
                return foundChapter.Id;
            }
        }

        // POST: api/Books
        [HttpPost]
        public async Task<int> Post([FromBody]Book value)
        {
            using (var session = Store.OpenSession())
            {
                var existing = await session
                    .Query<Book>()
                    .Where(x => x.Title == value.Title)
                    .FirstOrDefaultAsync();
                if(existing == null)
                {
                    existing = new Book() { Title = value.Title };                    
                    session.Store(existing);
                    await session.SaveChangesAsync();                    
                }
                return existing.Id;
            }
        }

        [HttpPost("{id}/chapters/{chapter}/pages")]
        public async Task<bool> AddPage(int id, int chapter, [FromBody]Page page)
        {
            using (var session = Store.OpenSession())
            {
                var existingChapter = await session.Query<Chapter>().FirstOrDefaultAsync(x => x.Id == chapter && x.BookID == id);
                if (existingChapter == null) { return false; }
                var existingPage = await session.Query<Page>().FirstOrDefaultAsync(x => x.PageNumber == page.PageNumber && x.ChapterID == page.ChapterID);
                if(existingPage == null)
                {
                    existingPage = new Page()
                    {
                        PageNumber = page.PageNumber,
                        ChapterID = page.ChapterID
                    };
                    session.Store<Page>(existingPage);
                }
                existingPage.Data = page.Data;

                session.Update(existingPage);
                await session.SaveChangesAsync();
                return true;
            }
        }

        [HttpGet("{bookID}/chapters/{chapterID}/pages")]
        public async Task<IEnumerable<int>> GetPages(int bookID, int chapterID)
        {
            using (var session = Store.LightweightSession())
            {                
                return await session.Query<Page>().Where(x => x.ChapterID == chapterID).Select(x=> x.Id).ToListAsync();
            }
        }

        [HttpGet("{bookID}/chapters/{chapterID}/pages/{pageID}")]
        public async Task<Page> GetPage(int bookID, int chapterID, int pageID)
        {
            using (var session = Store.LightweightSession())
            {
                return await session.Query<Page>().FirstOrDefaultAsync(x => x.Id == pageID);
            }
        }
        
        //// PUT: api/Books/5
        //[HttpPut("{id}")]
        //public void Put(int id, [FromBody]string value)
        //{
        //}
        
        // DELETE: api/ApiWithActions/5
        //[HttpDelete("{id}")]
        //public async Task<int> Delete(int id)
        //{
        //    using (var sess = Store.OpenSession())
        //    {
        //        var existing = await sess.Query<Book>().FirstOrDefaultAsync(x => x.Id == id);
        //        if(existing != null)
        //        {
        //            //Delete the pages                    
        //            sess.DeleteWhere<Page>(page => page.BookID == existing.Id);
        //            //Delete the book
        //            sess.Delete(existing);
        //            await sess.SaveChangesAsync();
        //            return id;
        //        }
        //    }
        //    return -1;
        //}
    }
}
