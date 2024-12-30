package models

import "errors"

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (n *Note) Validate() error {
	if n.Title == "" {
		return errors.New("title is required")
	}

	return nil
}
