import type { AnimeInfo, AnimeSearch, AnimeStreams } from "@minnasync/api-types";
import { Hono } from "hono";
import { zValidator } from "@hono/zod-validator";
import { z } from "zod";

import client from "../client";
import anilist from "../../resources/anilist";
import jikan from "../../resources/jikan";
import animeIds, { type AnimeIds, Resource } from "../../resources/animeIDs";
import animepahe from "../../providers/animepahe";

const app = new Hono()
// GET /meta
.get(
    "/meta",
    zValidator('query', z.object({
        id: z.string(),
        resource: z.enum(['anilist', 'mal']),
    })),
    async (c) => {
        const { id, resource } = c.req.query();

        let meta: AnimeInfo['meta'];
        if (resource === "anilist") {
            const anilistData = await anilist.getAnime(parseInt(id));

            meta = {
                id: anilistData.id,
                color: anilistData.coverImage.color,
                poster: anilistData.coverImage.large,
                background: anilistData.bannerImage,
                title: {
                    english: anilistData.title.english,
                    romaji: anilistData.title.romaji,
                    native: anilistData.title.native,
                },
                description: anilistData.description,
                year: anilistData.seasonYear,
                type: anilistData.format,
                rating: anilistData.meanScore,
                genres: anilistData.genres,
                studios: anilistData.studios.nodes.map((s) => s.name),
                is_nsfw: anilistData.isAdult,

                trailer: anilistData.trailer ? {
                    id: anilistData.trailer.id,
                    platform: anilistData.trailer.site,
                } : null,
            };
        } else if (resource === "mal") {
            const malInfo = await jikan.getAnime(parseInt(id));

            if (malInfo.isErr()) {
                c.status(500);
                return c.json({
                    success: false,
                    message: 'failed to fetch anime info.',
                    data: null,
                });
            }

            const { data } = malInfo.value;

            meta = {
                id: data.mal_id,
                color: null,
                poster:
                    data.images.jpg.large_image_url ||
                    data.images.jpg.image_url,
                background: null,
                title: {
                    english: data.title_english || null,
                    romaji: data.title,
                    native: data.title_japanese,
                },
                description: data.synopsis || null,
                year: data.year,
                type: data.type,
                rating: data.score,
                genres: data.genres.map((g) => g.name),
                studios: data.studios.map((s) => s.name),
                is_nsfw: data.rating.startsWith("R"),
                trailer: data.trailer ? {
                    /**
                     * I think MAL only uses youtube as a trailer site.
                     */
                    id: data.trailer.youtube_id,
                    platform: 'youtube',
                } : null,
            };
        } else {
            c.status(403);
            return c.json({
                success: false,
                message: "invalid resource.",
                data: null,
            });
        }
        
        return c.json(meta!);
    }
)
// GET /info
.get(
    "/info",
    zValidator('query', z.object({
        id: z.string(),
        provider: z.enum(['animepahe']),
        page: z.string().optional(),
    })),
    async (c) => {
        const { id, provider, page } = c.req.query();

        let meta: AnimeInfo['meta'];
        let details: AnimeInfo['details'];

        if (provider === "animepahe") {
            const animepaheInfo = await animepahe.getAnime(id, { page: page });

            if (animepaheInfo.isErr()) {
                c.status(500);
                return c.json({
                    success: false,
                    message: animepaheInfo.error.message,
                    data: null,
                });
            }

            const { page: info, episodes } = animepaheInfo.value;

            let ids: { mal: number; anilist: number; } = { mal: 0, anilist: 0 };
            for (const url of info.externalLinks) {
                let mappedIds: AnimeIds | undefined = undefined;

                switch (url.type) {
                    case "AniList":
                        const anilistId = url.url.split('/').pop()!;
                        const anilistLinked = await animeIds.fromResource(Resource.ANILIST, anilistId);

                        if (anilistLinked.isErr()) {
                            continue;
                        }

                        mappedIds = anilistLinked.value;

                        break;
                    case "MAL":
                        const malId = url.url.split('/').pop()!;
                        const malLinked = await animeIds.fromResource(Resource.MAL, malId);

                        if (malLinked.isErr()) {
                            continue;
                        }

                        mappedIds = malLinked.value;

                        break;
                    case "AniDB":
                        const anidbId = url.url.split('/').pop()!;
                        const anidbLinked = await animeIds.fromResource(Resource.ANIDB, anidbId);

                        if (anidbLinked.isErr()) {
                            continue;
                        }

                        mappedIds = anidbLinked.value;

                        break;
                }

                if (mappedIds === undefined) {
                    c.status(503);

                    return c.json({
                        success: false,
                        message: 'failed to map ids.',
                        data: null,
                    });
                }

                if (mappedIds.anilist_id) {
                    ids.anilist = mappedIds.anilist_id;
                }

                if (mappedIds.mal_id) {
                    ids.mal = mappedIds.mal_id;
                }

                break;
            }

            const metaRes = await client.anime.meta.$get({
                query: { id: ids.anilist.toString(), resource: "anilist" },
            });

            if (!metaRes.ok) {
                c.status(500);
                return c.json({
                    success: false,
                    message: 'failed to fetch anime info.',
                    data: null,
                });
            }
            meta = await metaRes.json() as AnimeInfo['meta'];

            /**
             * This is done to make sure all pages necessary are fetched.
             * In the event that an episode is on another page, it will fetch that page.
             */
            const neededPages = new Set(episodes.list.map((e) => Math.ceil(e.episode / 100)));
            const pages = await Promise.all(Array.from(neededPages).map(async (page) => {
                const episodes = await jikan.getEpisodes(ids.mal, { page: page });

                if (episodes.isErr()) {
                    return {
                        page,
                        episodes: [],
                    };
                }

                return {
                    page,
                    episodes: episodes.value.data,
                };
            }));

            /**
             * All of this is to ensure episode titles are mapped properly.
             * We use the total and from to determine if the session is a separate season.
             * If it was, we'll use the total and index to determine the actual episode number.
             */
            const startingEp = episodes.list[0].episode;
            const episodesList = episodes.list;
            details = {
                hasNextPage: episodes.hasNextPage,
                episodes: episodesList.reduce((acc, e, index) => {
                    const epNumber = startingEp === episodes.total - (episodes.from - 1)
                        ? e.episode
                        : episodes.total - index;
                    const page = pages.find((p) => p.page === Math.ceil(epNumber / 100));

                    if (!page) return [
                        ...acc,
                        {
                            id: e.id,
                            title: `Episode ${epNumber}`,
                            episode: epNumber,
                            preview: e.preview,
                            streaming_link: e.url,
                        }
                    ];

                    const episode = page.episodes.find((ep) => ep.mal_id === epNumber);
                    if (!episode) return [
                        ...acc,
                        {
                            id: e.id,
                            title: `Episode ${epNumber}`,
                            episode: epNumber,
                            preview: e.preview,
                            streaming_link: e.url,
                        }
                    ];

                    return [
                        ...acc,
                        {
                            id: e.id,
                            title: episode.title,
                            episode: epNumber,
                            preview: e.preview,
                            streaming_link: e.url,
                        },
                    ];

                }, [] as AnimeInfo['details']['episodes']),
            };
        } else {
            c.status(403);
            return c.json({
                success: false,
                message: "invalid provider.",
                data: null,
            });
        }

        return c.json({
            meta: meta,
            details: details,
        });
    }
)
// GET /search/:query
.get(
    '/search/:query',
    zValidator('query', z.object({
        provider: z.enum(['animepahe']),
    })),
    async (c) => {
        const query = c.req.param('query');
        const { provider } = c.req.query();
        
        let results: AnimeSearch = {
            results: []
        };
        if (provider === "animepahe") {
            const queryResults = await animepahe.searchAnime({ q: query });

            if (queryResults.isErr()) {
                c.status(500);
                return c.json({});
            }

            results.results = queryResults.value.data.map((r) => ({
                id: r.session,
                title: r.title,
                poster: r.poster,
                type: r.type,
                episodes: r.episodes,
                year: r.year,
            }));
        };

        return c.json(results);
    }
)
// GET /streams/animepahe/:id/:session
.get(
    "/streams/animepahe/:id/:session",
    async (c) => {
        const { id, session } = c.req.param();

        let resources: AnimeStreams = {}
        const streams = await animepahe.getEpisodeStreams(id, session);
        if (streams.isErr()) {
            c.status(500);
            return c.json({});
        }

        resources = streams.value;

        return c.json(resources);
    }
);

export default app;