// Package tweets provides functionality around reading and writing tweets.
package tweets

type EmojisController interface {
	EmojiResults() ([]EmojiCount, error)
}

func NewEmojisController(db DB) EmojisController {
	return &emojisController{
		db: db,
	}
}

type emojisController struct {
	db DB
}

// EmojiCount is a pair of a emoji and the count of occurrences of emoji.
type EmojiCount struct {
	Emoji string `json:"emoji"`
	Count int `json:"count"`
}

// EmojiResults returns the pair of emojis and counts.
func (ec *emojisController) EmojiResults() ([]EmojiCount, error) {
	var results []EmojiCount
	ec.db.EmojiResults()(func(row scanFn) error {
		var ec EmojiCount
		if err := row(&ec.Emoji, &ec.Count); err != nil {
			return err
		}
		results = append(results, ec)
		return nil
	})
	return results, nil
}
