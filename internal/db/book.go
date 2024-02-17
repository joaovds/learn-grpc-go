package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type Book struct {
	db          *sql.DB
	ID          string
	Name        string
  ISBN        string
	AuthorID  string
}

func NewBook(db *sql.DB) *Book {
	return &Book{db: db}
}

func (b *Book) Create(name, isbn, authorID string) (Book, error) {
	newId := uuid.New().String()

	_, err := b.db.Exec(
		"INSERT INTO books (id, name, isbn, author_id) VALUES ($1, $2, $3, $4)",
		newId, name, isbn, authorID,
	)
	if err != nil {
		return Book{}, err
	}

	return Book{ID: newId, Name: name, ISBN: isbn, AuthorID: authorID}, nil
}

func (b *Book) GetAll() ([]Book, error) {
	rows, err := b.db.Query("SELECT id, name, isbn, author_id FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []Book{}

	for rows.Next() {
		var book Book

		err := rows.Scan(&book.ID, &book.Name, &book.ISBN, &book.AuthorID)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (b *Book) GetByAuthorID(authorID string) ([]Book, error) {
	rows, err := b.db.Query("SELECT id, name, isbn, author_id FROM books WHERE author_id = $1", authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []Book{}

	for rows.Next() {
		var book Book

		err := rows.Scan(&book.ID, &book.Name, &book.ISBN, &book.AuthorID)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}
