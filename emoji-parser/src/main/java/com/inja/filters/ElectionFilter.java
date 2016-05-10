package com.inja.filters;

import com.inja.models.EmojiModel;
import com.inja.models.TweetModel;

/**
 * Created by sinasix on 5/9/16.
 */
public class ElectionFilter extends TweetFilter<EmojiModel> {

    final static String TRUMP = "trump";

    final static String HILLARY = "hillary";

    final static String TRUMP_HANDLE = "@realDonaldTrump";


    @Override
    public EmojiModel filterTweet(TweetModel tweet) {
        EmojiModel emojiModel = null;
        if (!tweet.text.contains("@realDonaldTrump")) {
            emojiModel = new EmojiModel(TRUMP);
        }else{
            emojiModel = new EmojiModel(HILLARY);
        }
        return emojiModel;
    }
}
