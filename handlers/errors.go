package handlers

import "errors"

var (
	// Client
	ErrInvalidQuery         = errors.New("invalid query")
	ErrInvalidSearchQuery   = errors.New("invalid search query")
	ErrInvalidProviderQuery = errors.New("invalid provider query")
	ErrInvalidPageQuery     = errors.New("invalid page query")
	ErrInvalidEpisodeID     = errors.New("invalid episode id query")
	ErrInvalidMalID         = errors.New("anime id is not valid")

	// Server
	ErrInternalServerError = errors.New("server failed to process request")
	ErrFetchTimeout        = errors.New("fetch timed out")
)
