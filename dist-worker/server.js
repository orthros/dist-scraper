require('dotenv').config();
const cheerio = require('cheerio');
var rp = require('request-promise');
var open = require('amqplib').connect(process.env.QUEUE_LOCATION);//'amqp://localhost');

const queueName = process.env.QUEUE_NAME;

async function main() {
  //TODO, get the bookNameChapter and baseURL from the message
  const bookName = 'psyren';
  const chapter = 146;  
  const baseUrl = 'http://www.mangareader.net';  

  // open.then(function (conn) {
  //   return conn.createChannel();
  // }).then(function (ch) {
  //   return ch.assertQueue(queueName).then(function (ok) {
  //     return ch.consume(queueName, function (msg) {
  //       if (msg !== null) {
  //         // await scrape(baseUrl, bookNameChapter, onFound);
  //         ch.ack(msg);
  //       }
  //     });
  //   });
  // }).catch(console.warn);

  var bookID = await getBookID(bookName);
  var chapterID = await getChapterID(bookID,  chapter);

  const finder = {
    bookId: bookID,
    chapterId: chapterID,    
    found: async function(imgNum, img){
      return await foundImage(this.bookId, this.chapterId, imgNum, img)
    }
  };

  const finder2 = {
    found: async function(imgNum, img){
      console.log('Found an image ' + imgNum);
    }
  };

  await scrape(baseUrl, '/' + bookName + '/' + chapter, finder);
}

async function getBookID(bookName) {
  const endpoint = process.env.BOOKS_ENDPOINT;
  var options = {
    method: 'POST',
    uri: endpoint + '/books',
    headers: {
      'User-Agent': 'Request-Promise'
    },
    json: true,
    body: {
      id: 0,
      title: bookName
    }
  };
  var res = await rp(options);
  return res;
}

async function getChapterID(bookID, chapterNumber){
  const endpoint = process.env.BOOKS_ENDPOINT;
  var options = {
    method: 'POST',
    uri: endpoint + '/chapters/',
    headers: {
      'User-Agent': 'Request-Promise'
    },
    json: true,
    body: {
      id: 0,
      bookID : bookID,
      chapterNumber: chapterNumber,
      chapterTitle: ''
    }
  };
  var res = await rp(options);
  return res; 
}

async function foundImage(bookID, chapterId, imgNum, img) {
  const endpoint = process.env.BOOKS_ENDPOINT;
  var options = {
    method: 'POST',
    uri: endpoint + '/pages',
    headers: {
      'User-Agent': 'Request-Promise'
    },
    json: true,
    body: {
      id:0,
      chapterId: chapterId,
      pageNumber: imgNum,
      data: Array.prototype.slice.call(Buffer.from(img),0)
    }
  }
  var response = await rp(options);
}

async function scrape(baseUrl, bookNameChapter, finder) {
  const imgSelecter = '#img';
  const totalSelecter = '#selectpage';
  const nextSelecter = '#imgholder > a:nth-child(1)';
  const totalRegex = /of (\d+)/g;
  
  //When we start scraping there is ALWAYS at least one page
  var totalPages = 1;
  var nextLocation = baseUrl + bookNameChapter;
  
  for (i = 0; i < totalPages; i++) {
    //Navigate and load
    var html = await rp(nextLocation);
    var $ = cheerio.load(html);
    
    //Determine total number of pages
    totalRegex.lastIndex = 0;
    var totalPages = +totalRegex.exec($(totalSelecter).text())[1];
    
    var imagePath = $(imgSelecter).attr('src');

    //We have the image's path 
    var options = {
      uri: imagePath,
      headers: {
        'User-Agent': 'Request-Promise'
      },
      json: true
    }
    var img = await rp(options).then(function (image) {
      return image;
    });

    //Now that we found it, let's go
    await finder.found(i, img);

    //Now we need to grab the next one
    nextLocation = baseUrl + $(nextSelecter).attr('href');
  }
}

main();