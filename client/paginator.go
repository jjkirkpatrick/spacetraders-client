package client

import (
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

// Paginator is a generic struct for pagination, where T is the type of data being paginated.
// It abstracts away the pagination logic, allowing users to fetch pages of data seamlessly.
type Paginator[T any] struct {
	Data          []T
	Meta          models.Meta // Use value instead of pointer to simplify usage
	fetchPageFunc func(meta models.Meta) ([]T, models.Meta, *models.APIError)
}

// NewPaginator creates a new Paginator instance with default pagination parameters.
// fetchFunc is a function that knows how to fetch a page of data given pagination metadata.
func NewPaginator[T any](fetchFunc func(models.Meta) ([]T, models.Meta, *models.APIError)) *Paginator[T] {
	// Initialize with default pagination parameters, e.g., page 1 and limit 5
	defaultMeta := models.Meta{Page: 1, Limit: 5}
	return &Paginator[T]{
		Meta:          defaultMeta,
		fetchPageFunc: fetchFunc,
	}
}

// FetchFirstPage fetches the first page based on the current Meta settings.
func (p *Paginator[T]) FetchFirstPage() (*Paginator[T], *models.APIError) {
	return p.FetchPage(1)
}

// FetchPage allows fetching a specific page by its number.
func (p *Paginator[T]) FetchPage(page int) (*Paginator[T], *models.APIError) {
	p.Meta.Page = page
	data, meta, err := p.fetchPageFunc(p.Meta)
	if err != nil {
		return nil, err
	}
	p.Data = data
	p.Meta = meta
	return p, nil
}

// GetNextPage fetches the next page.
func (p *Paginator[T]) GetNextPage() (*Paginator[T], *models.APIError) {
	p.Meta.Page++                                    // Correctly increment the page number
	newData, newMeta, err := p.fetchPageFunc(p.Meta) // Use the updated Meta
	if err != nil {
		return nil, err
	}
	p.Data = newData // Update the paginator's data
	p.Meta = newMeta // Update the paginator's meta information
	return p, nil    // Return the updated paginator instance
}

// GetPreviousPage fetches the previous page.
func (p *Paginator[T]) GetPreviousPage() (*Paginator[T], *models.APIError) {
	if p.Meta.Page > 1 {
		p.Meta.Page-- // Decrement the page number for the previous page
	}
	newData, newMeta, err := p.fetchPageFunc(p.Meta) // Fetch the previous page using the updated meta
	if err != nil {
		return nil, err
	}
	p.Data = newData // Update the paginator's data with the previous page's data
	p.Meta = newMeta // Update the paginator's meta information
	return p, nil    // Return the same paginator instance
}
