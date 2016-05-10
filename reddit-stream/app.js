var request = require('request')
var amqp = require('amqp');
var couchbase = require('couchbase')
var connection = amqp.createConnection({ host: 'guest:guest@'+process.env.RABBITMQ_SERVICE_PORT_5672_TCP_ADDR+':5672' });
connection.on('error', function(e) {
    console.log("connection error...", e);
    process.exit(1)
});
connection.on('ready', function () {
    console.log('Rabbit connected')

    // Use the default 'amq.topic' exchange
    var myCluster = new couchbase.Cluster('couchbase://'+process.env.COUCHBASE_SERVICE_PORT_8091_TCP_ADDR);
    var myBucket = myCluster.openBucket('reddit',function(err) {
        if (err) {
            console.log(err);
            process.exit(1)
        } else {

            console.log('Connected to couchbase');
            setInterval(function(){
                request("https://www.reddit.com/r/all/comments/.json?limit=100", function(error, response, body) {
                    if ( error ){ process.exit(1);}
                    var json = JSON.parse(body);
                    var child = json.data.children;
                    for ( var i in child ){
                        var id = child[i].data.id;
                        var text = child[i].data.body.toLowerCase();
                        if ( text.indexOf('trump') !== -1 || text.indexOf('hillary') !== -1 ) {
                            var envelope = {};
                            envelope['topic']='reddit';
                            envelope['body']=child[i].data;
                            //console.log(envelope);
                            myBucket.insert(id, envelope, function(err, res) {
                                console.log(err);
                                if ( err !== null){
                                    console.log('cached');
                                }else{
                                            connection.publish('reddit',JSON.stringify(envelope),function(err){
                                                if ( err ){
                                                    console.log(err);
                                                }else{
                                                    console.log('published');
                                                }
                                            })

                                }
                            });
                        }
                    }
                });
            },3000);

        }
    });





});
