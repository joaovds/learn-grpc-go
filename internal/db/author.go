package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type Author struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewAuthor(db *sql.DB) *Author {
	return &Author{db: db}
}

func (a *Author) Create(name, description string) (Author, error) {
	newId := uuid.New().String()

	_, err := a.db.Exec(
		"INSERT INTO authors (id, name, description) VALUES ($1, $2, $3)",
		newId, name, description,
	)
	if err != nil {
		return Author{}, err
	}

	return Author{ID: newId, Name: name, Description: description}, nil
}

func (a *Author) GetAll() ([]Author, error) {
	rows, err := a.db.Query("SELECT id, name, description FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	authors := []Author{}

	for rows.Next() {
		var author Author

		err := rows.Scan(&author.ID, &author.Name, &author.Description)
		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (a *Author) GetByBookID(bookID string) (Author, error) {
	var author Author

	err := a.db.QueryRow(
		"SELECT au.id, au.name, au.description FROM authors au JOIN books bo ON au.id = bo.author_id WHERE bo.id = $1",
		bookID,
	).Scan(&author.ID, &author.Name, &author.Description)
	if err != nil {
		return Author{}, err
	}

	return author, nil
}

func (a *Author) GetById(id string) (Author, error) {
	var author Author
	err := a.db.QueryRow(
		"SELECT id, name, description FROM authors WHERE id = $1",
		id,
	).Scan(&author.ID, &author.Name, &author.Description)
	if err != nil {
		return Author{}, err
	}

	return author, nil
}
