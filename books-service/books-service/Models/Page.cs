using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace books_service.Models
{
    public class Page
    {
        public int Id { get; set; }
        public int BookID { get; set; }
        public byte[] Data { get; set; }
        public int PageNumber { get; set; }
    }
}
