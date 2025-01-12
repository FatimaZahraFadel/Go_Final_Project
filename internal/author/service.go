package author

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type InMemoryAuthorStore struct {
	sync.Mutex
	authors map[int]Author
	nextID  int
}

func (store *InMemoryAuthorStore) CreateAuthor(ctx context.Context, author Author) (int, error) {
	store.Lock()
	defer store.Unlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		for _, existingAuthor := range store.authors {
			if existingAuthor.FirstName == author.FirstName && existingAuthor.LastName == author.LastName {
				log.Printf("author with name %s already exists", author.FirstName)
				return 0, fmt.Errorf("author with name %s already exists", author.FirstName)
			}
		}

		author.ID = store.nextID
		store.authors[author.ID] = author
		store.nextID++
		return author.ID, nil
	}
}

func (store *InMemoryAuthorStore) GetAuthorByID(ctx context.Context, id int) (Author, error) {
	store.Lock()
	defer store.Unlock()

	select {
	case <-ctx.Done():
		return Author{}, ctx.Err()
	default:
		author, found := store.authors[id]
		if !found {
			log.Printf("author with ID %d not found", id)
			return Author{}, fmt.Errorf("author with ID %d not found", id)
		}
		return author, nil
	}
}

func (store *InMemoryAuthorStore) ListAuthors(ctx context.Context) ([]Author, error) {
	store.Lock()
	defer store.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		all := make([]Author, 0, len(store.authors))
		for _, author := range store.authors {
			all = append(all, author)
		}
		return all, nil
	}
}

func (store *InMemoryAuthorStore) DeleteAuthor(ctx context.Context, id int) error {
	store.Lock()
	defer store.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if _, found := store.authors[id]; found {
			delete(store.authors, id)
			return nil
		}
		log.Printf("author with ID %d not found", id)
		return fmt.Errorf("author with ID %d not found", id)
	}
}

func (store *InMemoryAuthorStore) UpdateAuthor(ctx context.Context, id int, author Author) error {
	store.Lock()
	defer store.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if _, found := store.authors[id]; found {
			author.ID = id
			store.authors[id] = author
			return nil
		}
		log.Printf("author with ID %d not found", id)
		return fmt.Errorf("author with ID %d not found", id)
	}
}
