var amqp = require('amqp');
var couchbase = require('couchbase')
var myCluster = new couchbase.Cluster('couchbase://'+process.env.COUCHBASE_SERVICE_PORT_8091_TCP_ADDR);
var myBucket = myCluster.openBucket('tweets',function(err){
    if ( err ) {
        console.log(err);
        process.exit(1)
    }else{
        console.log('Couchbase connected')
        var connection = amqp.createConnection({ host: 'guest:guest@'+process.env.RABBITMQ_SERVICE_PORT_5672_TCP_ADDR+':5672' });
        connection.on('error', function(e) {
            console.log("connection error...", e);
        });
        connection.on('ready', function () {
            console.log('queue ready');
            // Use the default 'amq.topic' exchange
            connection.queue('', {durable:false,passive:false,autoDelete:false},function (q) {
                // Catch all messages
                q.bind('tweets');
                //q.bind('#');
                console.log('Bound')
                // Receive messages
                q.subscribe(function (message) {
                    // Print messages to stdout
                    console.log('got tweet ');
                    var json = JSON.parse(message.data.toString());
                    myBucket.upsert(String(json.id), json, function(err, res) {
                        if ( err ){
                            console.log(err);
                        }else{
                            console.log(res)
                        }
                    });
                });
            });
        });

    }
});
