package com.inja.models;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true) 
public final class TweetModel {
    public final Contributors contributors;
    public final Coordinates coordinates;
    public final String created_at;
    public final Current_user_retweet current_user_retweet;
    public final Entities entities;
    public final long favorite_count;
    public final boolean favorited;
    public final String filter_level;
    public final long id;
    public final String id_str;
    public final String in_reply_to_screen_name;
    public final long in_reply_to_status_id;
    public final String in_reply_to_status_id_str;
    public final long in_reply_to_user_id;
    public final String in_reply_to_user_id_str;
    public final String lang;
    public final boolean possibly_sensitive;
    public final long retweet_count;
    public final boolean retweeted;
    public final Retweeted_status retweeted_status;
    public final String source;
    public final Scopes scopes;
    public final String text;
    public final boolean truncated;
    public final User user;
    public final boolean withheld_copyright;
    public final Withheld_in_countries withheld_in_countries;
    public final String withheld_scope;
    public final Extended_entities extended_entities;
    public final long quoted_status_id;
    public final String quoted_status_id_str;
    public final Quoted_status quoted_status;

    @JsonCreator
    public TweetModel(
            @JsonProperty("contributors")
            Contributors contributors,
            @JsonProperty("coordinates")
            Coordinates coordinates,
            @JsonProperty("created_at")
            String created_at,
            @JsonProperty("current_user_retweet")
            Current_user_retweet current_user_retweet,
            @JsonProperty("entities")
            Entities entities,
            @JsonProperty("favorite_count")
            long favorite_count,
            @JsonProperty("favorited")
            boolean favorited,
            @JsonProperty("filter_level")
            String filter_level,
            @JsonProperty("id")
            long id,
            @JsonProperty("id_str")
            String id_str,
            @JsonProperty("in_reply_to_screen_name")
            String in_reply_to_screen_name,
            @JsonProperty("in_reply_to_status_id")
            long in_reply_to_status_id,
            @JsonProperty("in_reply_to_status_id_str")
            String in_reply_to_status_id_str,
            @JsonProperty("in_reply_to_user_id")
            long in_reply_to_user_id,
            @JsonProperty("in_reply_to_user_id_str")
            String in_reply_to_user_id_str,
            @JsonProperty("lang")
            String lang,
            @JsonProperty("possibly_sensitive")
            boolean possibly_sensitive,
            @JsonProperty("retweet_count")
            long retweet_count,
            @JsonProperty("retweeted")
            boolean retweeted,
            @JsonProperty("retweeted_status")
            Retweeted_status retweeted_status,
            @JsonProperty("source")
            String source,
            @JsonProperty("scopes")
            Scopes scopes,
            @JsonProperty("text")
            String text,
            @JsonProperty("truncated")
            boolean truncated,
            @JsonProperty("user")
            User user,
            @JsonProperty("withheld_copyright")
            boolean withheld_copyright,
            @JsonProperty("withheld_in_countries")
            Withheld_in_countries withheld_in_countries,
            @JsonProperty("withheld_scope")
            String withheld_scope,
            @JsonProperty("extended_entities")
            Extended_entities extended_entities,
            @JsonProperty("quoted_status_id")
            long quoted_status_id,
            @JsonProperty("quoted_status_id_str")
            String quoted_status_id_str,
            @JsonProperty("quoted_status")
            Quoted_status quoted_status
    )

    {
        this.contributors = contributors;
        this.coordinates = coordinates;
        this.created_at = created_at;
        this.current_user_retweet = current_user_retweet;
        this.entities = entities;
        this.favorite_count = favorite_count;
        this.favorited = favorited;
        this.filter_level = filter_level;
        this.id = id;
        this.id_str = id_str;
        this.in_reply_to_screen_name = in_reply_to_screen_name;
        this.in_reply_to_status_id = in_reply_to_status_id;
        this.in_reply_to_status_id_str = in_reply_to_status_id_str;
        this.in_reply_to_user_id = in_reply_to_user_id;
        this.in_reply_to_user_id_str = in_reply_to_user_id_str;
        this.lang = lang;
        this.possibly_sensitive = possibly_sensitive;
        this.retweet_count = retweet_count;
        this.retweeted = retweeted;
        this.retweeted_status = retweeted_status;
        this.source = source;
        this.scopes = scopes;
        this.text = text;
        this.truncated = truncated;
        this.user = user;
        this.withheld_copyright = withheld_copyright;
        this.withheld_in_countries = withheld_in_countries;
        this.withheld_scope = withheld_scope;
        this.extended_entities = extended_entities;
        this.quoted_status_id = quoted_status_id;
        this.quoted_status_id_str = quoted_status_id_str;
        this.quoted_status = quoted_status;
    }
    @JsonIgnoreProperties(ignoreUnknown = true) 
    public static final class Contributors {

        @JsonCreator
        public Contributors() {
        }
    }
    @JsonIgnoreProperties(ignoreUnknown = true) 
    public static final class Coordinates {

        @JsonCreator
        public Coordinates() {
        }
    }
    @JsonIgnoreProperties(ignoreUnknown = true) 
    public static final class Current_user_retweet {

        @JsonCreator
        public Current_user_retweet() {
        }
    }
    @JsonIgnoreProperties(ignoreUnknown = true) 
    public static final class Entities {
        public final Hashtag hashtags[];
        public final Media[] media;
        public final Url urls[];
        public final User_mention user_mentions[];

        @JsonCreator
        public Entities(@JsonProperty("hashtags") Hashtag[] hashtags, @JsonProperty("media") Media[] media, @JsonProperty("urls") Url[] urls, @JsonProperty("user_mentions") User_mention[] user_mentions) {
            this.hashtags = hashtags;
            this.media = media;
            this.urls = urls;
            this.user_mentions = user_mentions;
        }
        @JsonIgnoreProperties(ignoreUnknown = true) 
        public static final class Hashtag {

            @JsonCreator
            public Hashtag() {
            }
        }
        @JsonIgnoreProperties(ignoreUnknown = true)
        public static final class Media {

            @JsonCreator
            public Media() {
            }
        }
        @JsonIgnoreProperties(ignoreUnknown = true) 
        public static final class Url {

            @JsonCreator
            public Url() {
            }
        }
        @JsonIgnoreProperties(ignoreUnknown = true) 
        public static final class User_mention {
            public final int[] indices;
            public final long id;
            public final String id_str;
            public final String name;
            public final String screen_name;

            @JsonCreator
            public User_mention(@JsonProperty("indices") int[] indices, @JsonProperty("id") long id, @JsonProperty("id_str") String id_str, @JsonProperty("name") String name, @JsonProperty("screen_name") String screen_name) {
                this.indices = indices;
                this.id = id;
                this.id_str = id_str;
                this.name = name;
                this.screen_name = screen_name;
            }
        }
    }
    @JsonIgnoreProperties(ignoreUnknown = true) 
    public static final class Retweeted_status {

        @JsonCreator
        public Retweeted_status() {
        }
    }
    @JsonIgnoreProperties(ignoreUnknown = true) 
    public static final class Scopes {

        @JsonCreator
        public Scopes() {
        }
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    public static final class User {
        public final boolean contributors_enabled;
        public final String created_at;
        public final boolean default_profile;
        public final boolean default_profile_image;
        public final String description;
        public final String email;
        public final Entities entities;
        public final long favourites_count;
        public final boolean follow_request_sent;
        public final boolean following;
        public final long followers_count;
        public final long friends_count;
        public final boolean geo_enabled;
        public final long id;
        public final String id_str;
        public final boolean id_translator;
        public final String lang;
        public final long listed_count;
        public final String location;
        public final String name;
        public final boolean notifications;
        public final String profile_background_color;
        public final String profile_background_image_url;
        public final String profile_background_image_url_https;
        public final boolean profile_background_tile;
        public final String profile_banner_url;
        public final String profile_image_url;
        public final String profile_image_url_https;
        public final String profile_link_color;
        public final String profile_sidebar_border_color;
        public final String profile_sidebar_fill_color;
        public final String profile_text_color;
        public final boolean profile_use_background_image;
        public final String screen_name;
        public final boolean show_all_inline_media;
        public final Status status;
        public final long statuses_count;
        public final String time_zone;
        public final String url;
        public final long utc_offset;
        public final boolean verified;
        public final String withheld_in_countries;
        public final String withheld_scope;

        @JsonCreator
        public User(@JsonProperty("contributors_enabled") boolean contributors_enabled, @JsonProperty("created_at") String created_at, @JsonProperty("default_profile") boolean default_profile, @JsonProperty("default_profile_image") boolean default_profile_image, @JsonProperty("description") String description, @JsonProperty("email") String email, @JsonProperty("entities") Entities entities, @JsonProperty("favourites_count") long favourites_count, @JsonProperty("follow_request_sent") boolean follow_request_sent, @JsonProperty("following") boolean following, @JsonProperty("followers_count") long followers_count, @JsonProperty("friends_count") long friends_count, @JsonProperty("geo_enabled") boolean geo_enabled, @JsonProperty("id") long id, @JsonProperty("id_str") String id_str, @JsonProperty("id_translator") boolean id_translator, @JsonProperty("lang") String lang, @JsonProperty("listed_count") long listed_count, @JsonProperty("location") String location, @JsonProperty("name") String name, @JsonProperty("notifications") boolean notifications, @JsonProperty("profile_background_color") String profile_background_color, @JsonProperty("profile_background_image_url") String profile_background_image_url, @JsonProperty("profile_background_image_url_https") String profile_background_image_url_https, @JsonProperty("profile_background_tile") boolean profile_background_tile, @JsonProperty("profile_banner_url") String profile_banner_url, @JsonProperty("profile_image_url") String profile_image_url, @JsonProperty("profile_image_url_https") String profile_image_url_https, @JsonProperty("profile_link_color") String profile_link_color, @JsonProperty("profile_sidebar_border_color") String profile_sidebar_border_color, @JsonProperty("profile_sidebar_fill_color") String profile_sidebar_fill_color, @JsonProperty("profile_text_color") String profile_text_color, @JsonProperty("profile_use_background_image") boolean profile_use_background_image, @JsonProperty("screen_name") String screen_name, @JsonProperty("show_all_inline_media") boolean show_all_inline_media, @JsonProperty("status") Status status, @JsonProperty("statuses_count") long statuses_count, @JsonProperty("time_zone") String time_zone, @JsonProperty("url") String url, @JsonProperty("utc_offset") long utc_offset, @JsonProperty("verified") boolean verified, @JsonProperty("withheld_in_countries") String withheld_in_countries, @JsonProperty("withheld_scope") String withheld_scope) {
            this.contributors_enabled = contributors_enabled;
            this.created_at = created_at;
            this.default_profile = default_profile;
            this.default_profile_image = default_profile_image;
            this.description = description;
            this.email = email;
            this.entities = entities;
            this.favourites_count = favourites_count;
            this.follow_request_sent = follow_request_sent;
            this.following = following;
            this.followers_count = followers_count;
            this.friends_count = friends_count;
            this.geo_enabled = geo_enabled;
            this.id = id;
            this.id_str = id_str;
            this.id_translator = id_translator;
            this.lang = lang;
            this.listed_count = listed_count;
            this.location = location;
            this.name = name;
            this.notifications = notifications;
            this.profile_background_color = profile_background_color;
            this.profile_background_image_url = profile_background_image_url;
            this.profile_background_image_url_https = profile_background_image_url_https;
            this.profile_background_tile = profile_background_tile;
            this.profile_banner_url = profile_banner_url;
            this.profile_image_url = profile_image_url;
            this.profile_image_url_https = profile_image_url_https;
            this.profile_link_color = profile_link_color;
            this.profile_sidebar_border_color = profile_sidebar_border_color;
            this.profile_sidebar_fill_color = profile_sidebar_fill_color;
            this.profile_text_color = profile_text_color;
            this.profile_use_background_image = profile_use_background_image;
            this.screen_name = screen_name;
            this.show_all_inline_media = show_all_inline_media;
            this.status = status;
            this.statuses_count = statuses_count;
            this.time_zone = time_zone;
            this.url = url;
            this.utc_offset = utc_offset;
            this.verified = verified;
            this.withheld_in_countries = withheld_in_countries;
            this.withheld_scope = withheld_scope;
        }

        @JsonIgnoreProperties(ignoreUnknown = true) 
        public static final class Entities {

            @JsonCreator
            public Entities() {
            }
        }

        public static final class Status {

            @JsonCreator
            public Status() {
            }
        }
    }
    @JsonIgnoreProperties(ignoreUnknown = true) 
    public static final class Withheld_in_countries {

        @JsonCreator
        public Withheld_in_countries() {
        }
    }
    @JsonIgnoreProperties(ignoreUnknown = true) 
    public static final class Extended_entities {

        @JsonCreator
        public Extended_entities() {
        }
    }
    @JsonIgnoreProperties(ignoreUnknown = true)
    public static final class Quoted_status {

        @JsonCreator
        public Quoted_status() {
        }
    }
}
