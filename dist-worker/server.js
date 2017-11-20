require('dotenv').configure();
const puppeteer = require('puppeteer');
var request = require('request-promise');
var open = require('amqplib').connect(process.env.QUEUE_LOCATION);//'amqp://localhost');

const queueName = process.env.QUEUE_NAME;

var imgSelecter = "";
var totalSelecter = "";
const browser = await puppeteer.launch();
const page = await browser.newPage();
await page.goto('');
//Determine total number of pages
var totalPages = await page.$(totalSelecter);
for(i=0; i< totalPages; i++) {
  var imagePath = await page.$(imgSelecter).then(function(handle) {
      return handle.asElement().getProperty("src");
  });  
  //We have the image's path 
  var options = {
    uri:'' + "/" + imagePath,
    headers: {
      'User-Agent' : 'Request-Promise'
    },
    json: true
  }
  var img = await request(options).then(function(image) {
    return image;
  });
  
  //Now throw it off to the service

  //Now we need to grab the next one
  await page.click(totalSelecter);
  //Now we are navigating to the next page
  await page.waitForNavigation();
}

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