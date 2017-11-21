using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Marten;
using books_service.Models;

namespace books_service.Controllers
{
    [Produces("application/json")]
    [Route("api/Chapters")]
    public class ChaptersController : Controller
    {
        public IDocumentStore Store { get; }
        
        public ChaptersController(IDocumentStore store)
        {
            Store = store;
        }


        // GET: api/Chapters/5
        [HttpGet("{id}")]
        public async Task<Chapter> GetChapter(int id)
        {
            using (var session = Store.LightweightSession())
            {
                return await session.Query<Chapter>()
                                    .FirstOrDefaultAsync(x => x.Id == id);
            }
        }
        
        [HttpGet("{chapterID}/pages")]
        public async Task<IEnumerable<int>> GetPages(int chapterID)
        {
            using (var session = Store.LightweightSession())
            {
                return await session.Query<Page>().Where(x => x.ChapterID == chapterID).Select(x => x.Id).ToListAsync();
            }
        }

        [HttpPost]
        public async Task<int> PostChapter([FromBody]Chapter chapter)
        {
            using (var session = Store.OpenSession())
            {
                var existingBook = await session.Query<Book>()
                                                .Where(x => x.Id == chapter.BookID)
                                                .FirstOrDefaultAsync();
                if (existingBook == null) { return -1; }

                var foundChapter = await session.Query<Chapter>()
                                                   .Where(x => x.ChapterNumber == chapter.ChapterNumber)
                                                   .FirstOrDefaultAsync();
                if (foundChapter == null)
                {
                    foundChapter = new Chapter()
                    {
                        BookID = chapter.Id,
                        ChapterNumber = chapter.ChapterNumber,
                        ChapterTitle = chapter.ChapterTitle
                    };
                    session.Store(foundChapter);

                    await session.SaveChangesAsync();
                }
                return foundChapter.Id;
            }
        }

    }
}