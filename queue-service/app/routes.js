// var amqp = require('amqplib/callback_api');

// amqp.connect('amqp://localhost', function(err, conn) {
//   conn.createChannel(function(err, ch) {
//     var q = 'task_queue';
//     var msg = process.argv.slice(2).join(' ') || "Hello World!";

//     ch.assertQueue(q, {durable: true});
//     ch.sendToQueue(q, new Buffer(msg), {persistent: true});
//     console.log(" [x] Sent '%s'", msg);
//   });
//   setTimeout(function() { conn.close(); process.exit(0) }, 500);
// });


module.exports = function (app, ch, queueName) {

    // server routes ===========================================================
    // handle things like api calls
    // authentication routes
    app.get('/api/words', function (req, res) {
        res.status(200).json({
            word: ["Hello world", "We are here", "How are you"]
        });
    })

    app.get('/api/queue', function (req, res) {

        var msg = {
            BaseUrl: "http://mangareader.net",
            BookName: "Psyren",
            ChapterNumber: 146
        }

        var message = JSON.stringify(msg)

        ch.assertQueue(queueName, {
            durable: true
        });
        ch.sendToQueue(queueName, new Buffer(message), {
            persistent: true
        });
        res.status(200)
    })

    // frontend routes =========================================================
    // route to handle all angular requests
    app.get('*', function (req, res) {
        res.sendfile('./public/index.html');
    });

};