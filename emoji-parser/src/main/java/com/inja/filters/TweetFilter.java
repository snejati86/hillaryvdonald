package com.inja.filters;

import com.inja.models.FilterResult;
import com.inja.models.TweetModel;

/**
 * Created by sinasix on 5/9/16.
 */
public abstract class TweetFilter<T extends FilterResult>
{
    abstract public T filterTweet(TweetModel tweet);
}
