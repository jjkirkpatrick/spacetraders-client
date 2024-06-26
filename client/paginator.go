package client

import (
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

// Paginator is a generic struct for pagination, where T is the type of data being paginated.
// It abstracts away the pagination logic, allowing users to fetch pages of data seamlessly.
type Paginator[T any] struct {
	Data          []T
	Meta          models.Meta // Use value instead of pointer to simplify usage
	fetchPageFunc func(meta models.Meta) ([]T, models.Meta, error)
	Error         error
}

// NewPaginator creates a new Paginator instance with default pagination parameters.
// fetchFunc is a function that knows how to fetch a page of data given pagination metadata.
func NewPaginator[T any](fetchFunc func(models.Meta) ([]T, models.Meta, error)) *Paginator[T] {
	// Initialize with default pagination parameters, e.g., page 1 and limit 5
	defaultMeta := models.Meta{Page: 1, Limit: 20}
	return &Paginator[T]{
		Meta:          defaultMeta,
		fetchPageFunc: fetchFunc,
		Error:         nil,
	}
}

// FetchFirstPage fetches the first page based on the current Meta settings.
func (p *Paginator[T]) fetchFirstPage() (*Paginator[T], error) {
	return p.fetchPage(1)
}

// FetchPage allows fetching a specific page by its number.
func (p *Paginator[T]) fetchPage(page int) (*Paginator[T], error) {
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
func (p *Paginator[T]) getNextPage() (*Paginator[T], error) {
	p.Meta.Page++
	newData, newMeta, err := p.fetchPageFunc(p.Meta) // Use the updated Meta
	if err != nil {
		return nil, err
	}
	p.Data = newData // Update the paginator's data
	p.Meta = newMeta // Update the paginator's meta information
	return p, nil    // Return the updated paginator instance
}

// GetPreviousPage fetches the previous page.
func (p *Paginator[T]) getPreviousPage() (*Paginator[T], error) {
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

// FetchAllPages fetches all data at once without needing to loop over paginators.
func (p *Paginator[T]) FetchAllPages() ([]T, error) {
	var allData []T

	currentPage, err := p.fetchFirstPage()
	if err != nil {
		return nil, err
	}

	allData = append(allData, currentPage.Data...)

	// Loop until no more data is available.
	for len(currentPage.Data) > 0 {
		nextPage, err := currentPage.getNextPage()
		if err != nil {
			// If an error occurs, stop fetching and return what we have so far along with the error.
			return allData, err
		}
		// Check if nextPage is empty, indicating no more pages.
		if len(nextPage.Data) == 0 {
			break
		}
		allData = append(allData, nextPage.Data...)
		currentPage = nextPage
	}

	return allData, nil
}
