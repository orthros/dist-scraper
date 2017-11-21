using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace books_service.Models
{
    public class Chapter
    {
        public int Id { get; set; }
        public int BookID { get; set; }
        public int ChapterNumber { get; set; }
        public string ChapterTitle { get; set; }
    }
}
