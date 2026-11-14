package db

import "fmt"

// AddVote records that an author voted for a post. Voting twice is a no-op.
func (d DB) AddVote(authorID, postID float32) error {
	if _, err := d.client.Exec("INSERT OR IGNORE INTO votes(author, post) VALUES(?, ?)", authorID, postID); err != nil {
		return fmt.Errorf("failed to add vote: %w", err)
	}

	return nil
}

// CountVotes returns how many authors voted for a post.
func (d DB) CountVotes(postID float32) (int, error) {
	var n int
	if err := d.client.QueryRow("SELECT COUNT(*) FROM votes WHERE post = ?", postID).Scan(&n); err != nil {
		return 0, fmt.Errorf("failed to count votes: %w", err)
	}

	return n, nil
}
