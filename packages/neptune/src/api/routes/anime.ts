import type { AnimeInfo, Episodes, AnimeSearch, AnimeStreams } from "api-types";
import { Hono } from "hono";
import { zValidator } from "@hono/zod-validator";
import { z } from "zod";

import anilist from "../../resources/anilist";
import animepahe from "../../providers/animepahe";
import client from "../client";

const app = new Hono()
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
                title: anilistData.title.english,
                description: anilistData.description,
                type: anilistData.format,
                rating: anilistData.meanScore,
                genres: anilistData.genres,
                is_nsfw: anilistData.isAdult,
            };
        }
        
        return c.json(meta!);
    }
)
.get(
    "/info",
    zValidator('query', z.object({
        id: z.string(),
        resource: z.enum(['anilist']),
        provider: z.enum(['animepahe']),
    })),
    async (c) => {
        const { id, resource, provider } = c.req.query();

        let meta: AnimeInfo['meta'];
        let episodes: AnimeInfo['episodes'];

        if (provider === "animepahe") {
            const animepaheInfo = await animepahe.getAnime(id);

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

                meta = await res.json();
            }

            episodes = animepaheInfo.value.episodes.map((e) => ({
                id: e.id,
                title: `Episode ${e.episode}`,
                episode: e.episode,
                preview: e.preview,
                streaming_link: e.url,
            }));

            // episodes = animepaheInfo.value.episodes.reduce((acc, curr) => {
            //     acc.push({
            //         id: curr.id,
            //         title: `Episode ${curr.episode}`,
            //         episode: curr.episode,
            //         preview: curr.preview,
            //         streaming_link: curr.url,
            //     });

            //     return acc;
            // }, [] as Episodes[]);
        }

        return c.json({
            meta: meta!,
            episodes: episodes!,
        });
    }
)
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