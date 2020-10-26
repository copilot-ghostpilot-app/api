// Package tweets provides functionality around reading and writing tweets.
package tweets

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	rangeMin1 = 127744
	rangeMax1 = 129750
	rangeMin2 = 126980
	rangeMax2 = 127569
	rangeMin3 = 169
	rangeMax3 = 174
	rangeMin4 = 8205
	rangeMax4 = 12953
)

type TweetsController interface {
	StoreTweet(t Tweet) error
}

func NewTweetsController(db DB) TweetsController {
	return &tweetsController{
		db: db,
	}
}

type tweetsController struct {
	db DB
}

type TweetMetdata struct {
	Media       string `json:"media"`
	HashTags    string `json:"hashtags"`
	CreatedAt   string `json:"created_date"`
	RetweetData string `json:"retweet_data"`
}

type Tweet struct {
	ID           string       `json:"id"`
	Username     string       `json:"username"`
	TweetContent string       `json:"tweet_content"`
	Metadata     TweetMetdata `json:"metadata"`
}

func (tc *tweetsController) StoreTweet(tweet Tweet) error {
	var metadata string
	tm, err := json.Marshal(tweet.Metadata)
	if err != nil {
		log.Printf("INFO: tweets: Unable to marshal tweet metadata, id=%s, metadata=%s", tweet.ID, tweet.Metadata)
	} else {
		metadata = string(tm)
	}

	if err := tc.db.StoreTweet(tweet.ID, tweet.Username, tweet.TweetContent, tweet.Metadata.CreatedAt, metadata); err != nil {
		return err
	}
	emojis, err := tc.convertTweetToEmojisList(tweet.TweetContent)
	if err != nil {
		// TODO: We need some kind of recon that reconciles missing data between tables.
		return err
	}
	for _, emoji := range emojis {
		err := tc.db.StoreEmoji(tweet.ID, emoji)
		if err != nil {
			log.Printf("ERROR: tweets: id=%s emoji=%s: %v\n", tweet.ID, emoji, err)
		}
	}
	return nil
}

// convertTweetToEmojisList extracts emoji from a tweet and store in decimal.
func (tc *tweetsController) convertTweetToEmojisList(tweetContent string) ([]string, error) {
	var emoji []string
	r, err := regexp.Compile(`\\U\w{8}`)
	if err != nil {
		return nil, err
	}
	for _, e := range r.FindAllString(tweetContent, -1) {
		dec, err := hexUnicodeToDecUnicode(e)
		if err != nil {
			return nil, err
		}
		emoji = append(emoji, dec)
	}
	replaced := string(r.ReplaceAll([]byte(tweetContent), []byte("")))
	for _, r := range replaced {
		if isEmoji(r) {
			emoji = append(emoji, fmt.Sprint(r))
		}
	}
	return emoji, nil
}

// Example: \U0001f9c3 -> 129475
func hexUnicodeToDecUnicode(s string) (string, error) {
	hex := strings.TrimPrefix(s, `\U`)
	dec, err := strconv.ParseInt(hex, 16, 32)
	if err == nil {
		return fmt.Sprint(dec), nil
	}
	return "", err
}

func isEmoji(r rune) bool {
	code := int(r)
	switch {
	case code >= rangeMin1 && code <= rangeMax1:
		return true
	case code >= rangeMin2 && code <= rangeMax2:
		return true
	case code >= rangeMin3 && code <= rangeMax3:
		return true
	case code >= rangeMin4 && code <= rangeMax4:
		return true
	default:
		return false
	}
}
