package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	uc "github.com/minna-sync/neptune/internal/usecase"
	"github.com/minna-sync/neptune/pkg/httpx"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type AnimeStreamingRouter struct {
	uc uc.AnimeLookupUsecase
}

func NewAnimeRouter(r *chi.Mux, uc uc.AnimeLookupUsecase) {
	router := AnimeStreamingRouter{uc: uc}

	r.Route("/v1/anime", func(r chi.Router) {
		r.Get("/search", router.SearchAnime)

		r.Route("/{malId}", func(r chi.Router) {
			r.Get("/", router.GetAnime)
			r.Get("/episodes", router.GetEpisodes)
			r.Get("/streams", router.GetEpisodesStreaming)
		})
	})
}

type AnimeSearchQuery struct {
	Query string `schema:"q"`
}

func (q *AnimeSearchQuery) Validate() error {
	if q.Query == "" {
		return ErrInvalidSearchQuery
	}

	return nil
}

type AnimeSearch APIResponse[[]uc.AnimeSearchEntry]

func (h *AnimeStreamingRouter) SearchAnime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var query AnimeSearchQuery
	if err := httpx.DecodeURLQuery(r, &query); err != nil {
		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidQuery.Error(),
		}, http.StatusBadRequest)
		if err != nil {
			log.Error(err)
		}

		return
	}
	if err := query.Validate(); err != nil {
		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: err.Error(),
		}, http.StatusBadRequest)
		if err != nil {
			log.Error(err)
		}

		return
	}

	results, err := h.uc.SearchAnime(ctx, query.Query)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		log.WithError(err).Error("Failed to fetch anime search results.")

		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInternalServerError.Error(),
		}, http.StatusInternalServerError)

		if err != nil {
			log.Error(err)
		}

		return
	}

	err = httpx.WriteJSON(w, AnimeSearch{
		Success: true,
		Data:    results,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

type AnimeInfo APIResponse[uc.AnimeInfo]

func (h *AnimeStreamingRouter) GetAnime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	malId := chi.URLParam(r, "malId")
	id, err := strconv.Atoi(malId)
	if err != nil {
		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidMalID.Error(),
		}, http.StatusBadRequest)
		if err != nil {
			log.Error(err)
		}

		return
	}

	anime, err := h.uc.GetAnime(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, uc.ErrAnimeNotFound) {
			err = httpx.WriteJSON(w, APIError{
				Success: false,
				Message: uc.ErrAnimeNotFound.Error(),
			}, http.StatusNotFound)
			if err != nil {
				log.Error(err)
			}

			return
		}

		log.WithError(err).Error("Failed to fetch anime info.")

		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInternalServerError.Error(),
		}, http.StatusInternalServerError)
		if err != nil {
			log.Error(err)
		}

		return
	}

	err = httpx.WriteJSON(w, AnimeInfo{
		Success: true,
		Data:    *anime,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

type AnimeEpisodesQuery struct {
	Provider string `schema:"provider"`
	Page     int    `schema:"page"`
}

func (q *AnimeEpisodesQuery) Validate() error {
	if q.Provider == "" {
		return ErrInvalidProviderQuery
	}

	if q.Page < 0 {
		return ErrInvalidPageQuery
	}

	if q.Page == 0 {
		q.Page = 1 // default to page 1
	}

	return nil
}

type AnimeEpisodes APIResponse[[]uc.AnimeEpisode]

func (h *AnimeStreamingRouter) GetEpisodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	malId := chi.URLParam(r, "malId")
	id, err := strconv.Atoi(malId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var query AnimeEpisodesQuery
	if err := httpx.DecodeURLQuery(r, &query); err != nil {
		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: err.Error(),
		}, http.StatusBadRequest)
		if err != nil {
			log.Error(err)
		}

		return
	}
	if err := query.Validate(); err != nil {
		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: err.Error(),
		}, http.StatusBadRequest)
		if err != nil {
			log.Error(err)
		}

		return
	}

	episodes, err := h.uc.GetEpisodes(ctx, id, query.Provider, query.Page)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, uc.ErrAnimeNotFound) {
			err = httpx.WriteJSON(w, APIError{
				Success: false,
				Message: uc.ErrAnimeNotFound.Error(),
			}, http.StatusNotFound)
			if err != nil {
				log.Error(err)
			}

			return
		}

		log.WithField("error", err.Error()).Error("Failed to fetch anime episodes.")

		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInternalServerError.Error(),
		}, http.StatusInternalServerError)

		if err != nil {
			log.Error(err)
		}

		return
	}

	err = httpx.WriteJSON(w, AnimeEpisodes{
		Success: true,
		Data:    episodes,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

type AnimeEpisodeLinksQuery struct {
	Provider  string `schema:"provider"`
	EpisodeID string `schema:"episode_id"`
}

func (q *AnimeEpisodeLinksQuery) Validate() error {
	if q.Provider == "" {
		return ErrInvalidProviderQuery
	}

	if q.EpisodeID == "" {
		return ErrInvalidEpisodeID
	}

	return nil
}

type AnimeEpisodeLinks APIResponse[[]uc.AnimeEpisodeLinks]

func (h *AnimeStreamingRouter) GetEpisodesStreaming(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5) // shouldn't take more than 5 seconds to fetch
	defer cancel()

	malId := chi.URLParam(r, "malId")
	id, err := strconv.Atoi(malId)
	if err != nil {
		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidMalID.Error(),
		}, http.StatusBadRequest)
		if err != nil {
			log.Error(err)
		}

		return
	}

	var query AnimeEpisodeLinksQuery
	if err := httpx.DecodeURLQuery(r, &query); err != nil {
		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidQuery.Error(),
		}, http.StatusBadRequest)
		if err != nil {
			log.Error(err)
		}

		return
	}
	if err := query.Validate(); err != nil {
		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: err.Error(),
		}, http.StatusBadRequest)
		if err != nil {
			log.Error(err)
		}

		return
	}

	links, err := h.uc.GetStreamingLinks(ctx, id, query.Provider, query.EpisodeID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			err = httpx.WriteJSON(w, APIError{
				Success: false,
				Message: ErrFetchTimeout.Error(),
			}, http.StatusGatewayTimeout)
			if err != nil {
				log.Error(err)
			}

			return
		}

		log.WithError(err).Error("Failed to fetch anime episode links.")

		err = httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInternalServerError.Error(),
		}, http.StatusInternalServerError)

		if err != nil {
			log.Error(err)
		}

		return
	}

	err = httpx.WriteJSON(w, AnimeEpisodeLinks{
		Success: true,
		Data:    links,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}
