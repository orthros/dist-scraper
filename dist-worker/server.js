require('dotenv').configure();
const puppeteer = require('puppeteer');
var request = require('request');
var open = require('amqplib').connect(process.env.QUEUE_LOCATION);//'amqp://localhost');

const queueName = process.env.QUEUE_NAME;

const browser = await puppeteer.launch();

open.then(function(conn) {
    return conn.createChannel();
  }).then(function(ch) {
    return ch.assertQueue(queueName).then(function(ok) {
      return ch.consume(queueName, function(msg) {
        if (msg !== null) {
          console.log(msg.content.toString());
          ch.ack(msg);
        }
      });
    });
  }).catch(console.warn);

var scrape = function(browser, url, chapterExtension , totalExpression, nextExpression, foundImage) {
    const page = await browser.newPage();
    await page.goto(url + chapterExtension);
    
    //We are at the root for our page, find the image
    //download it to memory
    //call the "foundImage" function
    //foundImage(img);
    //Move to the next

    await browser.close();
}

var foundImage = function(bookID, img, pageNum) {
  const url = process.env.BOOKS_ENDPOINT;
  request.post(url + 'books/'+ bookID +'/pages', { form: { image:img, pageNumber:pageNum } });
}