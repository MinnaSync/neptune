import type { AnimeInfo, AnimeSearch, AnimeStreams } from "api-types";
import { Hono } from "hono";
import { zValidator } from "@hono/zod-validator";
import { z } from "zod";

import anilist from "../../resources/anilist";
import animepahe from "../../providers/animepahe";
import client from "../client";

const app = new Hono()
// GET /meta
.get(
    "/meta",
    zValidator('query', z.object({
        id: z.string(),
        resource: z.enum(['anilist']),
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
        }
        
        return c.json(meta!);
    }
)
// GET /info
.get(
    "/info",
    zValidator('query', z.object({
        id: z.string(),
        resource: z.enum(['anilist']),
        provider: z.enum(['animepahe']),
        page: z.string().optional(),
    })),
    async (c) => {
        const { id, resource, provider, page } = c.req.query();

        let meta: AnimeInfo['meta'];
        let details: AnimeInfo['details'];

        if (provider === "animepahe") {
            const animepaheInfo = await animepahe.getAnime(id, { page: page });

            if (animepaheInfo.isErr()) {
                c.status(500);
                return c.json({});
            } 

            if (resource === "anilist") {
                const url = animepaheInfo.value.externalLinks.find(el => el.type === "AniList")?.url;
                if (!url) {
                    c.status(404);
                    return c.json({});
                };

                const res = await client.anime.meta.$get({
                    query: { id: url.split('/').pop()!, resource: "anilist" },
                });

                if (!res.ok) {
                    c.status(500);
                    return c.json({});
                }

                meta = await res.json();
            }

            details = {
                hasNextPage: animepaheInfo.value.episodes.hasNextPage,
                episodes: animepaheInfo.value.episodes.list.map((e) => ({
                    id: e.id,
                    title: `Episode ${e.episode}`,
                    episode: e.episode,
                    preview: e.preview,
                    streaming_link: e.url,
                })),
            };
        }

        return c.json({
            meta: meta!,
            details: details!,
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