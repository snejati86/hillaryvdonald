var amqp = require('amqp');
var connection = amqp.createConnection({ host: 'guest:guest@192.168.99.100:5672' });

var watson = require('watson-developer-cloud');

var tone_analyzer = watson.tone_analyzer({
    password: 'cAFf3ctGr7vA',
    username: '18b0c519-367d-40b3-bfd2-cc1f9d226c04',
    version: 'v3-beta',
    version_date: '2016-02-11'
});


// Wait for connection to become established.
connection.on('ready', function () {
    // Use the default 'amq.topic' exchange
    connection.queue('tweets', {passive:true,autoDelete:false},function (q) {
        // Catch all messages
        q.bind('#');

        // Receive messages
        q.subscribe(function (message) {
            // Print messages to stdout

            var latest = message.data.toString().Body;
            tone_analyzer.tone({ text: latest},
                function(err, tone) {
                    if (err)
                        console.log(err);
                    else
                        console.log(JSON.stringify(tone, null, 2));
                });

        });
    });
});

