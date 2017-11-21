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
        public IEnumerable<string> Get()
        {
            return new string[] { "value1", "value2" };
        }

        // GET: api/Books/5
        [HttpGet("{id}", Name = "Get")]
        public Book Get(int id)
        {
            return new Book()
            {
                Id = id,
                Title = "value"
            };
        }

        // POST: api/Books
        [HttpPost]
        public int Post([FromBody]Book value)
        {
            using (var session = Store.OpenSession())
            {
                var existing = session
                    .Query<Book>()
                    .Where(x => x.Title == value.Title)
                    .FirstOrDefault();
                if(existing == null)
                {
                    existing = new Book() { Title = value.Title };
                    session.Store(existing);
                    session.SaveChanges();
                }

                return existing.Id;
            }
        }
        
        // PUT: api/Books/5
        [HttpPut("{id}")]
        public void Put(int id, [FromBody]string value)
        {
        }
        
        // DELETE: api/ApiWithActions/5
        [HttpDelete("{id}")]
        public void Delete(int id)
        {
        }
    }
}
