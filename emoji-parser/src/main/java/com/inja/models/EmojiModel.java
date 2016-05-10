package com.inja.models;

/**
 * This is the Emoji Model sent to the system.
 * Getters and setters are needed for Jackson.
 */
public class EmojiModel implements FilterResult
{
    public String owner;

    public String emojis;

    public EmojiModel(String owner) {
        this.owner = owner;
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
}
