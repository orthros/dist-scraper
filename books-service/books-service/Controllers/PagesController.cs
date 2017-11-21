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
    [Route("api/Pages")]
    public class PagesController : Controller
    {
        public IDocumentStore Store { get; }

        public PagesController(IDocumentStore store) {
            Store = store;
        }

        [HttpGet("{pageID}")]
        public async Task<Page> GetPage(int pageID)
        {
            using (var session = Store.LightweightSession())
            {
                return await session.Query<Page>().FirstOrDefaultAsync(x => x.Id == pageID);
            }
        }

        [HttpPost]
        public async Task<bool> AddPage([FromBody]Page page)
        {
            using (var session = Store.OpenSession())
            {
                var existingChapter = await session.Query<Chapter>().FirstOrDefaultAsync(x => x.Id == page.ChapterID );
                if (existingChapter == null) { return false; }
                var existingPage = await session.Query<Page>().FirstOrDefaultAsync(x => x.PageNumber == page.PageNumber && x.ChapterID == page.ChapterID);
                if (existingPage == null)
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
    }
}