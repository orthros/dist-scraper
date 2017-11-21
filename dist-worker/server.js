require('dotenv').config();
const cheerio = require('cheerio');
var rp = require('request-promise');
var open = require('amqplib').connect(process.env.QUEUE_LOCATION);//'amqp://localhost');

const queueName = process.env.QUEUE_NAME;

async function main() {
  //TODO, get the bookNameChapter and baseURL from the message
  const bookNameChapter = '/psyren/146';
  const baseUrl = 'http://www.mangareader.net';
  const onFound = function (img) { console.log('Found one!'); }

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

  await getBookID('psyren');
  await scrape(baseUrl, bookNameChapter, onFound);
}

async function getBookID(bookName) {
  const endpoint = process.env.BOOKS_ENDPOINT;
  var options = {
    method: 'PUT',
    uri: endpoint + '/books/',
    headers: {
      'User-Agent': 'Request-Promise'
    },
    json: true
  };
  var res = await rp(options, { title: bookName });
  return res;
}

async function foundImage(bookID, imgNum, img) {
  const endpoint = process.env.BOOKS_ENDPOINT;
  var options = {
    method: 'POST',
    uri: endpoint + '/books/' + bookID + '/images',
    headers: {
      'User-Agent': 'Request-Promise'
    },
    json: true
  };
  var body = {
    imageNumber: imgNum,
    image: img
  };
  var img = await rp(options, body).then(function (image) {
    return image;
  });
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