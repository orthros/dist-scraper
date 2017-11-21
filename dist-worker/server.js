require('dotenv').config();
const cheerio = require('cheerio');
var rp = require('request-promise');
var open = require('amqplib').connect(process.env.QUEUE_LOCATION);//'amqp://localhost');

const queueName = process.env.QUEUE_NAME;

async function main() {
  const bookNameChapter = '/psyren/146';
  const baseUrl = 'http://www.mangareader.net';
  const onFound = function (img) { console.log('Found one!'); }
  await scrape(baseUrl, bookNameChapter, onFound);
}

async function scrape(baseUrl, bookNameChapter, onFound) {
  const imgSelecter = '#img';
  const totalSelecter = '#selectpage';
  const nextSelecter = '#imgholder > a:nth-child(1)';
  const totalRegex = /of (\d+)/g;

  var nextLocation = baseUrl + bookNameChapter;
  var html = await rp(nextLocation);
  //Determine total number of pages
  var $ = cheerio.load(html);
  var totalPages = +totalRegex.exec($(totalSelecter).text())[1];

  for (i = 0; i < totalPages; i++) {
    var html = await rp(nextLocation);
    //Determine total number of pages
    var $ = cheerio.load(html);
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
    onFound(img);

    //Now we need to grab the next one
    nextLocation = baseUrl + $(nextSelecter).attr('href');
  }
}

main();
// open.then(function (conn) {
//   return conn.createChannel();
// }).then(function (ch) {
//   return ch.assertQueue(queueName).then(function (ok) {
//     return ch.consume(queueName, function (msg) {
//       if (msg !== null) {
//         console.log(msg.content.toString());
//         ch.ack(msg);
//       }
//     });
//   });
// }).catch(console.warn);


// var scrape = function (browser, url, chapterExtension, totalExpression, nextExpression, foundImage) {
//   const page = await browser.newPage();
//   await page.goto(url + chapterExtension);

//   //We are at the root for our page, find the image
//   //download it to memory
//   //call the "foundImage" function
//   //foundImage(img);
//   //Move to the next

//   await browser.close();
// }

// var foundImage = function (bookID, img, pageNum) {
//   const url = process.env.BOOKS_ENDPOINT;
//   request.post(url + 'books/' + bookID + '/pages', { form: { image: img, pageNumber: pageNum } });
// }