var ws = require("nodejs-websocket")
var amqp = require('amqp');
var radio = require('radio')
var connection = amqp.createConnection({ host: 'guest:guest@'+process.env.RABBITMQ_SERVICE_PORT_5672_TCP_ADDR+':5672' });



var lineReader = require('readline').createInterface({
    input: require('fs').createReadStream('locations')
});
var stateMap = {}
lineReader.on('line', function (line) {
    stateMap[line]=true;
});

// Wait for connection to become established.
connection.on('ready', function () {

    console.log('Rabbit connected')
    // Use the default 'amq.topic' exchange
    connection.queue('tweets', {passive:true,autoDelete:false},function (q) {
        // Catch all messages
        q.bind('#');
        console.log('bound1')
        // Receive messages
        q.subscribe(function (message) {
            // Print messages to stdout
            var latest = message.data.toString();
            var tweet = JSON.parse(latest)
            var message = {}
            //console.log(tweet.coordinates);
            //console.log(tweet.user.location)
            var location = tweet.user.location;
            message['topic']='tweet';
            message['body'] = tweet.text;
            var state = getState(location);
            if (state != null )
            {
                message['loc'] = state;
            }
            //message['location'] =
            radio('new').broadcast(message);
        });
    });

    connection.queue('profanity', {passive:false,autoDelete:false},function (q) {
        // Catch all messages
        q.bind('#');
        console.log('bound2')
        // Receive messages
        q.subscribe(function (message) {
            // Print messages to stdout
            var latest = message.data.toString();
            var profanityEnvelope = JSON.parse(latest);
            var message = {}
            message['topic']='curse';
            message['body'] = profanityEnvelope;
            radio('new').broadcast(message);
        });
    });

    connection.queue('emojis', {passive:false,autoDelete:false},function (q) {
        // Catch all messages
        q.bind('#');
        console.log('bound3')
        // Receive messages
        q.subscribe(function (message) {
            // Print messages to stdout
            var latest = message.data.toString();
            console.log(latest)
            var emoji = JSON.parse(latest);
            var message = {}
            message['topic']='emoji';
            message['body'] = emoji;
            radio('new').broadcast(message);
        });
    });

    connection.queue('reddit', {passive:true,autoDelete:false},function (q) {
        // Catch all messages
        q.bind('#');
        console.log('bound4')
        // Receive messages
        q.subscribe(function (message) {
            // Print messages to stdout
            var latest = message.data.toString();
            console.log(latest)
            var reddit = JSON.parse(latest);
            var message = {}
            message['topic']='reddit';
            if ( reddit.body.body.length < 150){
                message['body'] = reddit.body.body;
                radio('new').broadcast(message);
            }
        });
    });
});
connection.on('error', function(e) {
    console.log("connection error...", e);
    process.exit(1)
});

// Scream server example: "hi" -> "HI!!!"
 ws.createServer(function (conn) {
    console.log("New connection")
    var callback = function (data,context){
            console.log(data)
            conn.sendText(JSON.stringify(data))
    };
    radio('new').subscribe(callback)

    //var timer = setInterval(function(){
    //
    //    conn.sendText(JSON.stringify({ trump : "3", hillary :"4", current :"Hello mom "}));
    //},2000);
    conn.on("close", function (code, reason) {
        //clearInterval(timer)
        radio('new').unsubscribe(callback);
        console.log("Connection closed")
    })
}).listen(8001)

function getState(text){
    var state = null;
    for ( var i in stateMap){
        if ( text.indexOf(i) !== -1 ){
            state = i;
            break;
        }
    }
    return state;
}
//function getState(text){
//    var state = null;
//    var abb = containsAbb(text);
//    if ( abb !== null){
//        state = abb;
//    }
//    abb = containsFullName(text);
//    if ( abb != null){
//        state = convertFullName(state)
//    }
//    return state;
//
//}
//
//function containsAbb(text){
//    var state = null;
//    for ( var i in abblist ){
//        if ( text.indexOf())
//    }
//}