package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	uc "github.com/minna-sync/neptune/internal/usecase"
	"github.com/minna-sync/neptune/pkg/animepahe"
	jikan "github.com/minna-sync/neptune/pkg/jikan"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"golang.org/x/text/language"
)

type Services struct {
	Jikan     *jikan.Client
	AnimePahe *animepahe.Client
}

type AnimeLookupUsecase struct {
	sf singleflight.Group

	services Services
}

func NewAnimeLookupUsecase(redis *redis.Client) uc.AnimeLookupUsecase {
	uc := &AnimeLookupUsecase{}

	uc.services.Jikan = jikan.NewClient(redis)
	uc.services.AnimePahe = animepahe.NewClient(redis)

	return uc
}

func (h *AnimeLookupUsecase) SearchAnime(ctx context.Context, term string) ([]uc.AnimeSearchEntry, error) {
	query := url.Values{}
	query.Add("q", term)
	query.Add("order_by", "start_date")
	query.Add("sort", "desc")
	query.Add("limit", "25")
	query.Add("unapproved", "")

	malSearch, _, err := h.services.Jikan.AnimeService.GetAnimeSearch(ctx, query)
	if err != nil {
		return nil, err
	}

	results := []uc.AnimeSearchEntry{}
	for _, anime := range malSearch.Data {
		if anime.Type == "" || anime.Year == 0 {
			continue
		}

		err := h.services.Jikan.Cache.SetAnimeInfo(ctx, anime.MalID, &anime)
		if err != nil {
			log.WithError(err).Error("Failed to cache anime info for ID: ", anime.MalID)
		}

		results = append(results, uc.AnimeSearchEntry{
			MALId:  anime.MalID,
			Poster: anime.Images.JPG.LargeImageURL,
			Title: uc.AnimeTitle{
				Native:  &anime.Title,
				Romaji:  &anime.TitleJapanese,
				English: &anime.TitleEnglish,
			},
			Type:          anime.Type,
			Year:          anime.Year,
			TotalEpisodes: anime.Episodes,
			IsNSFW:        strings.HasPrefix(anime.Rating, "R"),
		})
	}

	return results, nil
}

func (h *AnimeLookupUsecase) GetAnime(ctx context.Context, id int) (*uc.AnimeInfo, error) {
	info, err, _ := h.sf.Do(fmt.Sprintf("anime:%d:info", id), func() (any, error) {
		info, err := h.services.Jikan.Cache.GetAnimeInfo(ctx, id)

		if info == nil {
			info, resp, err := h.services.Jikan.AnimeService.GetAnimeById(ctx, id)

			if resp.StatusCode != http.StatusOK {
				if resp.StatusCode == http.StatusNotFound {
					return nil, uc.ErrAnimeNotFound
				}

				return nil, uc.ErrFetchFailed
			}

			if err != nil {
				return nil, err
			}

			h.services.Jikan.Cache.SetAnimeInfo(ctx, id, info.Data)
			return info.Data, nil
		}

		return info, err
	})

	if err != nil {
		return nil, err
	}

	return &uc.AnimeInfo{
		Info: info.(*jikan.AnimeInfoBase),
	}, nil
}

func (h *AnimeLookupUsecase) GetEpisodes(ctx context.Context, id int, provider string, page int) ([]uc.AnimeEpisode, error) {
	anime, err := h.GetAnime(ctx, id)
	if err != nil {
		return nil, err
	}

	switch provider {
	case "kwik":
		var eps []uc.AnimeEpisode

		cachedEpisodes, err := h.services.AnimePahe.Cache.GetEpisodes(ctx, id, page)
		if cachedEpisodes == nil {
			session, err := h.services.AnimePahe.GetMatchedResult(ctx, anime.Info)
			if err != nil {
				if errors.Is(err, animepahe.ErrNoSearchResultsFound) {
					return nil, uc.ErrAnimeNotFound
				}

				return nil, err
			}

			episodes, resp, err := h.services.AnimePahe.Releases(ctx, *session, 1)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("animepahe api returned status code %d", resp.StatusCode)
			}
			if len(episodes.Data) <= 0 {
				return nil, fmt.Errorf("no results found for %s", anime.Info.TitleEnglish)
			}

			err = h.services.AnimePahe.Cache.SetEpisodes(ctx, id, &episodes.Data, page)
			if err != nil {
				fmt.Print(err)
			}

			for _, episode := range episodes.Data {
				eps = append(eps, uc.AnimeEpisode{
					ID:       episode.Session,
					Title:    fmt.Sprintf("Episode %d", episode.Episode),
					Episode:  episode.Episode,
					Snapshot: episode.Snapshot,
				})
			}

			return eps, nil
		}

		for _, episode := range *cachedEpisodes {
			eps = append(eps, uc.AnimeEpisode{
				ID:       episode.Session,
				Title:    fmt.Sprintf("Episode %d", episode.Episode),
				Episode:  episode.Episode,
				Snapshot: episode.Snapshot,
			})
		}

		return eps, err
	default:
		return nil, fmt.Errorf("provider %s not supported", provider)
	}
}

func (h *AnimeLookupUsecase) GetStreamingLinks(ctx context.Context, id int, provider string, episode string) ([]uc.AnimeEpisodeLinks, error) {
	info, _ := h.GetAnime(ctx, id)

	switch provider {
	case "kwik":
		var episodeLinks []uc.AnimeEpisodeLinks

		session, err := h.services.AnimePahe.GetMatchedResult(ctx, info.Info)
		if err != nil {
			return nil, err
		}

		links, err := h.services.AnimePahe.GetEpisodeStreamingLinks(ctx, id, *session, episode)
		if err != nil {
			if errors.Is(err, animepahe.ErrNoSearchResultsFound) {
				return nil, uc.ErrNoSearchResultsFound
			}

			return nil, err
		}

		for _, link := range links {
			lang, err := language.All.Parse(link.Language)
			if err != nil {
				log.WithError(err).Error("Failed to parse language.")
			}

			episodeLinks = append(episodeLinks, uc.AnimeEpisodeLinks{
				URL:        link.URL,
				Resolution: link.Resolution,
				Language:   lang.String(),
				Subtitles:  []uc.AnimeEpisodeSubtitles{},
			})
		}

		return episodeLinks, nil
	default:
		return nil, fmt.Errorf("provider %s not supported", provider)
	}
}
