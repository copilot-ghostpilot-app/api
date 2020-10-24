// Package tweets provides functionality around reading and writing tweets.
package tweets

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	// DBName is the name of the database for the tweets API.
	DBName = "tweets"
)

type scanFn func(dest ...interface{}) error

type eachRowFn func(row scanFn) error

type partialQuery func(row eachRowFn) error

// DB is the interface for all the operations allowed on tweets.
type DB interface {
	StoreTweet(id, username, tweetContent, createdAt, metadata string) error
	StoreEmoji(id, emoji string) error
	EmojiResults() partialQuery
}

// NewSQLDB creates a sql database to read and store tweets.
func NewSQLDB(db *sql.DB) DB {
	return &sqlDB{
		conn: db,
	}
}

type execQuerier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type sqlDB struct {
	conn execQuerier
}

// Store a tweet in the database.
func (db *sqlDB) StoreTweet(id, username, tweetContent, createdAt, metadata string) error {
	mdLen := len(metadata)
	if mdLen >= 1000 {
		log.Printf("INFO: db: skipping storing metadata because len of metadata=%d exceeds 1000 for id=%s, metadata=%s", mdLen, id, metadata)
		metadata = " "
	}
	parsedDate, err := time.Parse(time.RubyDate, createdAt)
	if err != nil {
		return fmt.Errorf("parse date %s: %w", createdAt, err)
	}
	parsedTime := parsedDate.Unix()
	_, err = db.conn.Exec(`INSERT INTO tweets (id, username, tweet_content, created_at, metadata) VALUES ($1, $2, $3, to_timestamp($4), $5)`, id, username, tweetContent, parsedTime, metadata)
	if err != nil {
		return fmt.Errorf("tweets: store tweet id=%s, username=%s, tweet_content=%s, created_at=%d, err=%v", id, username, tweetContent, parsedTime, err)
	}
	return nil
}

func (db *sqlDB) StoreEmoji(id, emoji string) error {
	_, err := db.conn.Exec(`INSERT INTO emojis (id, emoji) VALUES ($1, $2)`, id, emoji)
	if err != nil {
		return fmt.Errorf("emojis: store emoji for tweet id %s and emoji %s: %v", id, emoji, err)
	}
	return nil
}

func (db *sqlDB) EmojiResults() partialQuery {
	return func(row eachRowFn) error {
		rows, err := db.conn.Query(`SELECT emoji, COUNT(id) AS count FROM emojis GROUP BY emoji`)
		if err != nil {
			return fmt.Errorf("emojis: retrieve emoji results: %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := row(rows.Scan); err != nil {
				return fmt.Errorf("emojis: scan row to emoji count pair: %v", err)
			}
		}
		return rows.Err()
	}
}

// CreateTweetsTableIfNotExist creates the "tweets" table if it does not exist already.
func CreateTweetsTableIfNotExist(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tweets (
				id VARCHAR(255) NOT NULL UNIQUE,
				username VARCHAR(255) NOT NULL,
				tweet_content VARCHAR(1000) NOT NULL,
				created_at TIMESTAMP NOT NULL,
				metadata VARCHAR(1000) NOT NULL)`)
	if err != nil {
		return fmt.Errorf(`tweet: create "tweets" table: %v\n`, err)
	}
	return nil
}

// CreateEmojisTableIfNotExist creates the "emojis" table if it does not exist already.
func CreateEmojisTableIfNotExist(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS emojis (id VARCHAR(255), emoji VARCHAR(255), PRIMARY KEY (id, emoji))`)
	if err != nil {
		return fmt.Errorf(`tweet: create "emojis" table: %v\n`, err)
	}
	return nil
}
