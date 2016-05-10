package com.inja;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.rabbitmq.client.*;
import emoji4j.EmojiUtils;


import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

class Emoji{


    public static void main (String[] args) throws Exception{
        ConnectionFactory factory = new ConnectionFactory();
        String rabbit = System.getenv("RABBITMQ_SERVICE_PORT_5672_TCP_ADDR");
        if ( rabbit == null ){
            System.out.println("Can not find system property");
            System.exit(-1);
        }
        factory.setHost(rabbit);
        Connection connection = factory.newConnection();
        final Channel channel = connection.createChannel();
        String queueName = channel.queueDeclare().getQueue();
        channel.queueBind(queueName,"tweets","");

        channel.exchangeDeclare("emoji", "fanout", true);

        Consumer consumer = new DefaultConsumer(channel) {
            @Override
            public void handleDelivery(String consumerTag, Envelope envelope, AMQP.BasicProperties properties, byte[] body)
                    throws IOException {
                String message = new String(body, "UTF-8");
                System.out.println(message);
                ObjectMapper objectMapper = new ObjectMapper();
                Tweet tweet = objectMapper.readValue(body,Tweet.class);
                String owner = "trump";

                if ( !tweet.text.contains("@realDonaldTrump")){
                    owner = "clinton";
                }

                //System.out.println(tweet.text);
                //String regexString = Pattern.quote("text:") + "(.*?)" + Pattern.quote(pattern2);

                Pattern pattern = Pattern.compile("[\ud83c\udc00-\ud83c\udfff]|[\ud83d\udc00-\ud83d\udfff]|[\u2600-\u27ff]",
                        Pattern.UNICODE_CASE | Pattern.CASE_INSENSITIVE);
                //Pattern pattern = Pattern.compile(regexPattern);
                Matcher matcher = pattern.matcher(tweet.text);

                List<String> matchList = new ArrayList<String>();

                while (matcher.find()) {
                    matchList.add(matcher.group());
                }
                if ( matchList.size() > 0){
                    for(String s : matchList) {
                        EmojiEnv env = new EmojiEnv(owner,EmojiUtils.htmlify(s));
                        String json = objectMapper.writeValueAsString(env);
                        //String json = "{toward:\""+owner+"\",emoji:\""+emojis+"\"}";
                        channel.basicPublish("emoji","", null, json.getBytes());

                        System.out.println(json);
                    }
                }
            }
        };

        channel.basicConsume(queueName, true, consumer);

    }
}