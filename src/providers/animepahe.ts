import { okAsync, errAsync } from "neverthrow";
import * as cheerio from "cheerio";
import { request } from "../util/request";
import kwik from "./extractors/kwik";
import cache from "../util/cache";

type AnimePage = {
    externalLinks: {
        type: string;
        url: string;
    }[];
};

type AnimeInfo = {
    page: AnimePage;
    episodes: {
        hasNextPage: boolean;
        from: number;
        to: number;
        total: number;
        list: {
            id: string;
            preview: string;
            episode: number;
            url: string;
        }[];
    };
};

type AnimeEpisode = {
    id: number;
    anime_id: number;
    episode: number;
    episode2: number; // not sure if this is even used ever.
    edition: string;
    title: string;
    snapshot: string;
    disc: string;
    audio: string;
    duration: string;
    session: string;
    filler: number;
    created_at: string;
};

type StreamingSources = {
    resolution: string;
    link: string;
};

type SearchResults = {
    id: number;
    title: string;
    type: string;
    episodes: number;
    status: string;
    season: string;
    year: number;
    score: number;
    poster: string;
    session: string;
}

type APIResponse<T> = {
    total: number;
    per_page: number;
    current_page: number;
    last_page: number;
    next_page_url: string;
    prev_page_url: string;
    from: number;
    to: number;
    data: T;
};

const BASE_URL = 'https://animepahe.ru';

const HEADERS = (session?: string) => ({
    accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
    'accept-encoding': 'gzip, deflate, br',
    'accept-language': 'en-US,en;q=0.9',
    'connection': 'keep-alive',
    cookie: '__ddg2_=;',
    dnt: '1',
    host: 'animepahe.ru',
    referer: session ? `${BASE_URL}/anime/${session}` : BASE_URL,
    'set-fetch-dest': 'empty',
    'set-fetch-mode': 'cors',
    'set-fetch-site': 'same-origin',
    'x-requested-with': 'XMLHttpRequest',
    'user-agent': "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:138.0) Gecko/20100101 Firefox/138.0",
});

/**
 * Fetch the streaming urls for an episode.
 * @param id The anime id.
 * @param session The session.
 */
async function getEpisodeStreams(id: string, session: string) {
    const cached = await cache.jsonGet<Record<string, StreamingSources[]>>(`animepahe:streams:${id}:${session}`);
    if (cached) return okAsync(cached);

    const res = await request(new URL(`/play/${id}/${session}`, BASE_URL).toString(), 'get', HEADERS(id), {});

    if (res.isErr()) {
        return errAsync(res.error);
    }

    const html = res.value;
    const $ = cheerio.load(html);

    const sources: Record<string, StreamingSources[]> = {};

    for (const el of $('#resolutionMenu > button').toArray()) {
        const link = $(el).attr('data-src')?.trim();
        const resolution = $(el).attr('data-resolution')?.trim();
        const lang = $(el).attr('data-audio')?.trim();

        const streamingUrl = await kwik.extract(link!);

        if (!sources[lang!]) {
            sources[lang!] = [];
        }

        sources[lang!].push({
            resolution: resolution!,
            link: streamingUrl,
        });
    }

    await cache.jsonSet(`animepahe:streams:${id}:${session}`, sources, 60 * 30);

    return okAsync(sources);
}

/**
 * Get a list of episodes for an anime.
 * @param query The URL query to make.
 */
async function getEpisodes(query: { id: string; sort?: "episode_asc" | "episode_desc", page?: string }) {
    const cached = await cache.jsonGet<AnimeInfo['episodes']>(`animepahe:episodes:${query.id}:page_${query?.page || "1"}:sort_${query?.sort || "episode_asc"}`);
    if (cached) return okAsync(cached);

    const res = await request(new URL(`/api`, BASE_URL).toString(), 'get', HEADERS(query.id),
        { m: 'release', ...query }
    );

    if (res.isErr()) {
        return res;
    }

    const result: APIResponse<AnimeEpisode[]> = res.value;
    const episodes: AnimeInfo['episodes'] =  {
        hasNextPage: result.current_page < result.last_page,
        from: result.from,
        to: result.to,
        total: result.total,
        list: result.data.map((ep) => ({
            id: `${query.id}/${ep.session}`,
            preview: ep.snapshot,
            episode: ep.episode,
            url: `https://animepahe.ru/play/${query.id}/${ep.session}`,
        }))
    }

    await cache.jsonSet(`animepahe:episodes:${query.id}:page_${query?.page || "1"}:sort_${query?.sort || "episode_asc"}`, episodes, 60 * 30);

    return okAsync(episodes);
}

/**
 * Get the page information for an anime.
 * @param id The ID of the anime.
 */
async function getAnimePage(id: string) {
    const cached = await cache.jsonGet<AnimePage>(`animepahe:info:${id}`);
    if (cached) return okAsync(cached);

    const res = await request(new URL(`/anime/${id}`, BASE_URL).toString(), 'get', HEADERS(id), {});

    if (res.isErr()) {
        return res;
    }

    const html = res.value;
    const $ = cheerio.load(html);

    /**
     * Gets the external links relating to the anime.
     */
    const externalLinks: AnimePage['externalLinks'] = [];
    for (const el of $('.external-links > a').toArray()) {
        const type = $(el).text().trim();
        const link = $(el).attr('href')?.trim().replace(/^\/+/, ''); // links for some reason have 2 leading slashes.

        if (!link) continue;

        externalLinks.push({ type, url: `https://${link}` });
    }
    
    const info: AnimePage = {
        externalLinks
    };

    await cache.jsonSet(`animepahe:info:${id}`, info, 60 * 60 * 7);

    return okAsync(info);
}

/**
 * Get the animepahe information for an anime.
 * @param id The anime id.
 */
async function getAnime(id: string, query?: { page?: string }) {
    const page = await getAnimePage(id);
    if (page.isErr()) {
        return page;
    }

    const episodes = await getEpisodes({ id: id, sort: 'episode_desc', page: query?.page || "1" });
    if (episodes.isErr()) {
        return episodes
    }

    const animeInfo: AnimeInfo = {
        page: page.value,
        episodes: episodes.value,
    }

    return okAsync(animeInfo);
}

async function searchAnime(query: { q: string; }) {
    const res = await request(new URL(`/api`, BASE_URL).toString(), 'get', HEADERS(),
        { m: 'search', ...query }
    );

    if (res.isErr()) {
        return res;
    }

    const result: APIResponse<SearchResults[]> = res.value;

    return okAsync(result);
}

export default {
    getEpisodes,
    getAnime,
    getEpisodeStreams,
    searchAnime,
}