package com.inja;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Created by sinasix on 5/9/16.
 */
public class EmojiUtils
{
    private static final Pattern EMOJI_EXTRACTION_PATTERN = Pattern.compile("[\ud83c\udc00-\ud83c\udfff]|[\ud83d\udc00-\ud83d\udfff]|[\u2600-\u27ff]",
            Pattern.UNICODE_CASE | Pattern.CASE_INSENSITIVE);

    public static List<String> getEmojis ( String input )
    {
        Matcher matcher = EMOJI_EXTRACTION_PATTERN.matcher(input);

        List<String> matchList = new ArrayList<String>();

        while (matcher.find())
        {
            matchList.add(matcher.group());
        }
        return matchList;
    }

}
