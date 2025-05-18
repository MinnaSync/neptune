import { okAsync, errAsync } from "neverthrow";
import * as cheerio from "cheerio";
import { request } from "../util/request";
import kwik from "./extractors/kwik";
import cache from "../util/cache";

type AnimeInfo = {
    externalLinks: {
        type: string;
        url: string;
    }[];
    episodes: {
        id: string;
        preview: string;
        episode: number;
        url: string;
    }[];
};

type AnimeEpisodes = {
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
    filter: number;
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
    const res = await request(new URL(`/api`, BASE_URL).toString(), 'get', HEADERS(query.id),
        { m: 'release', ...query }
    );

    if (res.isErr()) {
        return res;
    }

    const result: APIResponse<AnimeEpisodes[]> = res.value;

    return okAsync(result);
}

/**
 * Get the animepahe information for an anime.
 * @param id The anime id.
 */
async function getAnime(id: string) {
    const cached = await cache.jsonGet<AnimeInfo>(`animepahe:info:${id}`);
    if (cached) return okAsync(cached);

    const res = await request(new URL(`/anime/${id}`, BASE_URL).toString(), 'get', HEADERS(id), {});

    if (res.isErr()) {
        return res;
    }

    const html = res.value;
    const $ = cheerio.load(html);

    const externalLinks: { type: string; url: string; }[] = [];
    for (const el of $('.external-links > a').toArray()) {
        const type = $(el).text().trim();
        const link = $(el).attr('href')?.trim().replace(/^\/+/, ''); // links for some reason have 2 leading slashes.

        if (!link) continue;

        externalLinks.push({ type, url: `https://${link}` });
    }

    const episodes = await getEpisodes({ id: id, sort: 'episode_asc', page: '1' });
    if (episodes.isErr()) {
        return okAsync({
            externalLinks,
            episodes: [],
        });
    }

    const info: AnimeInfo = {
        externalLinks,
        episodes: episodes.value.data.map((ep) => ({
            id: `${id}/${ep.session}`,
            preview: ep.snapshot,
            episode: ep.episode,
            url: `https://animepahe.ru/play/${id}/${ep.session}`,
        })),
    };

    await cache.jsonSet(`animepahe:info:${id}`, info, 60 * 30);

    return okAsync(info);
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