package client

import (
	"fmt"
	"time"

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

// FetchAllPages fetches all data concurrently using 4 workers.
func (p *Paginator[T]) FetchAllPages() ([]T, error) {
	// Get first page to determine total pages
	firstPage, err := p.fetchFirstPage()
	if err != nil {
		return nil, err
	}

	totalPages := (firstPage.Meta.Total + firstPage.Meta.Limit - 1) / firstPage.Meta.Limit

	// Create channels for results and errors
	results := make(chan []T, totalPages)
	errors := make(chan error, totalPages)
	pages := make(chan int, totalPages)

	// Calculate number of workers based on total pages
	// Use min(totalPages, 8) to avoid creating more workers than needed
	numWorkers := totalPages
	if numWorkers > 12 {
		numWorkers = 12
	}

	// Start workers based on calculated number
	for i := 0; i < numWorkers; i++ {
		go func() {
			for page := range pages {
				paginator := &Paginator[T]{
					Meta:          p.Meta,
					fetchPageFunc: p.fetchPageFunc,
				}
				fmt.Println("Fetching page", page)

				// Try up to 3 times
				var data *Paginator[T]
				var err error
				for retries := 0; retries < 3; retries++ {
					data, err = paginator.fetchPage(page)
					if err == nil {
						break
					}
					// Wait a bit before retrying
					time.Sleep(time.Second * time.Duration(retries+1))
				}

				if err != nil {
					errors <- fmt.Errorf("failed to fetch page %d after 3 retries: %w", page, err)
					continue
				}
				results <- data.Data
			}
		}()
	}

	// Send pages to workers
	go func() {
		for page := 1; page <= totalPages; page++ {
			pages <- page
		}
		close(pages)
	}()

	// Collect results
	var allData []T
	allData = append(allData, firstPage.Data...)

	// Collect remaining pages
	for i := 1; i < totalPages; i++ {
		select {
		case data := <-results:
			allData = append(allData, data...)
		case err := <-errors:
			return allData, err
		}
	}

	return allData, nil
}
