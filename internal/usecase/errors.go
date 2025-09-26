package usecase

import "errors"

var (
	// Client
	ErrAnimeNotFound        = errors.New("anime not found")
	ErrNoSearchResultsFound = errors.New("no search results found")

	// Server
	ErrFetchFailed = errors.New("failed to fetch data")
)
