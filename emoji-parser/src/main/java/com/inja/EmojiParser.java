package com.inja;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.inja.filters.ElectionFilter;
import com.inja.models.EmojiModel;
import com.inja.models.TweetModel;
import com.rabbitmq.client.*;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;



import java.io.IOException;
import java.util.List;

/**
 * Main driver for the service.
 */
class EmojiParser
{

    /**
     * Name used for AMQP exchange to publish messages.
     */
    private static final String EMOJI_PUBLISH_NAME = "emoji";

    /**
     * Name used for AMQP exchange to subscribe to messages.
     */
    private static final String INBOUND_SUBSCRIBE_NAME = "tweets";

    /**
     * Exchange should fan out to multiple consumers.
     */
    private static final String EXCHANGE_STRATEGY = "fanout";

    private static final Logger logger = LogManager.getLogger(EmojiParser.class);


    public static void main (String[] args) throws Exception
    {

        String rabbitMqIP = System.getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR");

        if ( rabbitMqIP == null )
        {
            throw new RuntimeException("Unable to find the system property for rabbitmq.");
        }

        // Connect.
        ConnectionFactory factory = new ConnectionFactory();
        factory.setHost(rabbitMqIP);
        Connection connection = factory.newConnection();
        final Channel channel = connection.createChannel();

        //Temp queue.
        String queueName = channel.queueDeclare().getQueue();
        channel.queueBind(queueName,INBOUND_SUBSCRIBE_NAME,"");

        logger.trace("Bound to the exchange,");
        //Outbound exchange.
        channel.exchangeDeclare(EMOJI_PUBLISH_NAME, EXCHANGE_STRATEGY, true);
        final ObjectMapper objectMapper = new ObjectMapper();

        //I'll fix it later.
        final ElectionFilter electionFilter = new ElectionFilter();

        channel.basicConsume(queueName, true, new DefaultConsumer(channel)
        {
            @Override
            public void handleDelivery(String consumerTag, Envelope envelope, AMQP.BasicProperties properties, byte[] body)
                    throws IOException {

                String message = new String(body, "UTF-8");
                TweetModel tweetModel = objectMapper.readValue(body, TweetModel.class);
                logger.trace("Go message = {}"+tweetModel.text);
                List<String> emojis = EmojiUtils.getEmojis(tweetModel.text);
                if  ( emojis.size() > 0 )
                {
                    for(String emoji : emojis)
                    {
                        EmojiModel emojiModel = electionFilter.filterTweet(tweetModel);
                        emojiModel.setEmojis(emoji);
                        String jsonModel = objectMapper.writeValueAsString(emojiModel);
                        logger.trace("Publishing model = {}"+jsonModel);
                        channel.basicPublish(EMOJI_PUBLISH_NAME,"",null,jsonModel.getBytes());
                    }

                }

            }
        });


    }
}