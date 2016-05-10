var amqp = require('amqp');
var couchbase = require('couchbase')
var myCluster = new couchbase.Cluster('couchbase://'+process.env.COUCHBASE_SERVICE_PORT_8091_TCP_ADDR);
var myBucket = myCluster.openBucket('feeds',function(err){
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
            connection.queue('', {durable:false,passive:false,autoDelete:true,exclusive:true},function (q) {
                // Catch all messages
                q.bind("feeds",'',function(q){
                    // Receive messages
                    q.subscribe(function (message) {
                        // Print messages to stdout
                        var json = JSON.parse(message.data.toString());
                        console.log(json)
                        myBucket.upsert(String(json.id), json, function(err, res) {
                            if ( err ){
                                console.log(err);
                            }else{
                                console.log(res)
                            }
                        });
                    });
                });
                console.log('Bound')


            });
        });

    }
});
