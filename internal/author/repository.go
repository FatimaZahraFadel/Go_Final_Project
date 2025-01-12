package author

import "context"

type Author struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
}

type AuthorStore interface {
	CreateAuthor(ctx context.Context, author Author) (int, error)
	GetAuthorByID(ctx context.Context, id int) (Author, error)
	UpdateAuthor(ctx context.Context, author Author) (Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	ListAuthors(ctx context.Context) ([]Author, error)
}
