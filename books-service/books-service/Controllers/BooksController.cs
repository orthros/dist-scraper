using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using books_service.Models;
using Marten;

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

        [HttpPost("{id}/pages")]
        public async Task<bool> AddPage(int id, [FromBody]Page page)
        {
            using (var session = Store.OpenSession())
            {
                var existingBook = await session.Query<Book>().FirstOrDefaultAsync(x => x.Id == id);
                if(existingBook == null) { return false; }
                session.Store<Page>(new Page()
                {
                    Data = page.Data,
                    BookID = existingBook.Id,
                    PageNumber = page.PageNumber
                });
                await session.SaveChangesAsync();
                return true;
            }
        }
        
        //// PUT: api/Books/5
        //[HttpPut("{id}")]
        //public void Put(int id, [FromBody]string value)
        //{
        //}
        
        // DELETE: api/ApiWithActions/5
        [HttpDelete("{id}")]
        public async Task<int> Delete(int id)
        {
            using (var sess = Store.OpenSession())
            {
                var existing = await sess.Query<Book>().FirstOrDefaultAsync(x => x.Id == id);
                if(existing != null)
                {
                    //Delete the pages
                    sess.DeleteWhere<Page>(page => page.BookID == existing.Id);
                    //Delete the book
                    sess.Delete(existing);
                    await sess.SaveChangesAsync();
                    return id;
                }
            }
            return -1;
        }
    }
}
