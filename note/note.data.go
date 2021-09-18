package note

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/GrandTaho/noto/database"
)

func getNotes() ([]Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	id, 
	title, 
	body, 
	author, 
	created, 
	updated, 
	tag FROM note`)

	if err != nil {
		return nil, err
	}
	defer results.Close()
	notes := make([]Note, 0)
	for results.Next() {
		var note Note
		results.Scan(&note.Id,
			&note.Title,
			&note.Body,
			&note.Author,
			&note.Created,
			&note.Updated,
			&note.Tag)
		notes = append(notes, note)
	}
	return notes, nil
}

func insertNote(note Note) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO note 
	(title, body, author, created, updated, tag) VALUES 
	(?, ?, ?, NOW(), NOW(), ?)`, note.Author, note.Body, note.Author, note.Tag)
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	return int(insertId), nil
}

func getNote(id int) (*Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT 
	id, 
	title, 
	body, 
	author, 
	created, 
	updated, 
	tag 
	FROM note WHERE id = ?`, id)
	note := &Note{}
	err := row.Scan(
		&note.Id,
		&note.Title,
		&note.Body,
		&note.Author,
		&note.Created,
		&note.Updated,
		&note.Tag,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return note, nil
}

func updateNote(note Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if note.Id == nil || *note.Id == 0 {
		log.Println("Invalid note id when trying to update", err)
		return errors.New("note have invalid id")
	}

	_, err := database.DbConn.ExecContext(ctx, `UPDATE note SET 
	title=?, 
	body=?, 
	author=?, 
	updated=NOW(), 
	tag=? WHERE id=?`,
		note.Title,
		note.Body,
		note.Author,
		note.Tag,
		note.Id)
	if err != nil {
		log.Println("Error updating note", err)
		return err
	}
	return nil
}

func removeNote(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM note WHERE id=?`, id)
	if err != nil {
		log.Printf("Error deleting note with ID = ?. Error: %v", id, err)
		return err
	}
	return nil
}
