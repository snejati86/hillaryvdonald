package com.inja;

/**
 * Created by sinasix on 5/7/16.
 */
public class EmojiEnv {
    public String owner;

    public EmojiEnv(String owner, String emojis) {
        this.owner = owner;
        this.emojis = emojis;
    }

    public String getEmojis() {
        return emojis;
    }

    public void setEmojis(String emojis) {
        this.emojis = emojis;
    }

    public String getOwner() {
        return owner;
    }

    public void setOwner(String owner) {
        this.owner = owner;
    }

    public String emojis;
}
