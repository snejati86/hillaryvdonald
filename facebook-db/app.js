var amqp = require('amqp');
var couchbase = require('couchbase')
var myCluster = new couchbase.Cluster('couchbase://'+process.env.COUCHBASE_SERVICE_PORT_8091_TCP_ADDR);
var myBucket = myCluster.openBucket('facebook',function(err){
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
            // Use the default 'amq.topic' exchange
            connection.queue('', {durable:false,passive:false,autoDelete:true,exclusive:true},function (q) {
                var exc = connection.exchange('facebook-feed', {type:'fanout',durable:true},function (exchange) {
                    q.bind("facebook",'',function(q){
                        // Receive messages
                        q.subscribe(function (message) {
                            // Print messages to stdout
                            var json = JSON.parse(message.data.toString());
                            console.log(json);
                            myBucket.insert(String(json.Id), json, function(err, res) {
                                console.log(err);
                                if ( err !== null){
                                    console.log('cached');
                                }else{
                                    exc.publish('',JSON.stringify(json),null,function(err){
                                        if ( err ){
                                            console.log('Error occured');
                                            console.log(err);
                                        }else{
                                            console.log('published');
                                        }
                                    })

                                }
                            });
                        });
                    });

                });
                    // Catch all messages

                console.log('Bound')


            });
        });

    }
});
