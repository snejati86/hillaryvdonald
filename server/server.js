var ws = require("nodejs-websocket")
var amqp = require('amqp');
var radio = require('radio')
var LRU = require("lru-cache")
var uuid = require('node-uuid');

var options = { max: 100
    , length: function (n, key) { return 1}
    , dispose: function (key, n) { }
    , maxAge: 1000 * 60 * 60 }

var tweetCache = LRU(options);
var profanityCache = LRU(options);
var emojiCache = LRU(options);

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
    connection.queue('', {durable:false,passive:false,autoDelete:true,exclusive:true},function (q) {

        // Catch all messages
        q.bind('tweets','',function(q){
            console.log('bound1')
            // Receive messages
            q.subscribe(function (message) {
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
                tweetCache.set(uuid.v1(),message)
            });
        });
    });

    connection.queue('', {durable:false,passive:false,autoDelete:true,exclusive:true},function (q) {
        // Catch all messages
        q.bind('profanity','',function(q){
            console.log('bound2')
            // Receive messages
            q.subscribe(function (message) {
                // Print messages to stdout
                var latest = message.data.toString();
                var profanityEnvelope = JSON.parse(latest);
                var message = {}
                message['topic']='curse';
                message['body'] = profanityEnvelope;
                profanityCache.set(uuid.v1(),message);
                radio('new').broadcast(message);
            });
        })

    });

    connection.queue('', {durable:false,passive:false,autoDelete:true,exclusive:true},function (q) {
        // Catch all messages
        q.bind('emoji','',function(){
            console.log('bound3')
            // Receive messages
            q.subscribe(function (message) {
                var latest = message.data.toString();
                console.log(latest)
                var emoji = JSON.parse(latest);
                var message = {}
                message['topic']='emoji';
                message['body'] = emoji;
                emojiCache.set(uuid.v1(),message);
                radio('new').broadcast(message);
            });
        });

    });

    connection.queue('', {durable:false,passive:false,autoDelete:true,exclusive:true},function (q) {
        // Catch all messages
        q.bind('feelings','',function(){
            console.log('bound4')
            // Receive messages
            q.subscribe(function (message) {
                var latest = message.data.toString();
                console.log(latest)
                var feelings = JSON.parse(latest);
                message['topic']='feelings';
                message['body'] = feelings;
                radio('new').broadcast(message);
            });
        });

    });


    connection.queue('', {durable:false,passive:false,autoDelete:true,exclusive:true},function (q) {
        // Catch all messages
        q.bind('reddit','',function(){
            console.log('bound5')
            // Receive messages
            q.subscribe(function (message) {
                // Print messages to stdout
                var latest = message.data.toString();
                var reddit = JSON.parse(latest);
                console.log(reddit)
                var message = {}
                message['topic']='reddit';
                if ( reddit.data.body.length < 150){
                    message['body'] = reddit.data.body;
                    radio('new').broadcast(message);
                }
            });
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
     tweetCache.forEach(function(value,key,cache){
        conn.sendText(JSON.stringify(value));
     });
     console.log(tweetCache.itemCount)
     profanityCache.forEach(function(value,key,cache){
         console.log('Sending from cache');
         conn.sendText(JSON.stringify(value));
     });
     console.log(profanityCache.itemCount)
     emojiCache.forEach(function(value,key,cache){
         console.log('Sending from cache');
         conn.sendText(JSON.stringify(value));
     });
     console.log(emojiCache.itemCount)
    radio('new').subscribe(callback)

     conn.on("error", function (code, reason) {
         //clearInterval(timer)
         radio('new').unsubscribe(callback);
         console.log("ERROR closed")
     })
    conn.on("close", function (code, reason) {
        //clearInterval(timer)
        radio('new').unsubscribe(callback);
        console.log("Connection closed")
    })
}).listen(8001)

/**
 * Listen I know this doesn't method doesn't belong here, I know it's a long way and you're ready to go to work... all I'm saying is wait, just wait,
 * please hear me out because this is not an episode, relapse, fuck-up, it's... I'm begging you . I'm begging you.
 * Try and make believe this is not just madness because this is not just madness.
 * Two weeks ago I came out of the building, okay, I'm running across Sixth Avenue, there's a car waiting,
 * I got exactly 38 minutes to get to the airport and I'm dictating. There's this, this panicked associate sprinting along beside me,
 * scribbling in a notepad, and suddenly she starts screaming, and I realize we're standing in the middle of the street, the light's changed,
 * there's this wall of traffic, serious traffic speeding towards us, and I freeze,
 * I can't move, and I'm suddenly consumed with the overwhelming sensation that I'm covered with some sort of film.
 * It's in my hair, my face... it's like a glaze... like a... a coating, and... at first I thought,
 * oh my god, I know what this is, this is some sort of amniotic - embryonic - fluid. I'm drenched in afterbirth,
 * I've-I've breached the chrysalis, I've been reborn. But then the traffic, the stampede, the cars, the trucks, the horns,
 * the screaming and I'm thinking no-no-no-no, reset, this is not rebirth,
 * this is some kind of giddy illusion of renewal that happens in the final moment before death.
 * And then I realize no-no-no, this is completely wrong because I look back at the building and I had the most stunning moment of clarity.
 * I realized , that I had emerged not from the doors of our office, not through the portals of our vast and powerful tech company,
 * but from the asshole of an organism whose sole function is to excrete the poison,
 * the ammo, the defoliant necessary for other, larger, more powerful organisms to destroy the miracle of humanity.
 * And that I had been coated in this patina of shit for the best part of my life.
 * The stench of it and the stain of it would in all likelihood take the rest of my life to undo.
 * And you know what I did? I took a deep cleansing breath and I set that notion aside. I tabled it.
 * I said to myself as clear as this may be, as potent a feeling as this is, as true a thing as I believe that I have witnessed today,
 * it must wait.
 * It must stand the test of time.
 * And , the time is now.
 */
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
