import { AnimeClient, type JikanResponse, type Anime, type AnimeEpisode } from "@tutkli/jikan-ts"; 
import { okAsync, errAsync } from "neverthrow";
import cache from "../util/cache";

const client = new AnimeClient();

async function getAnime(id: number) {
    const cached = await cache.jsonGet<JikanResponse<Anime>>(`jikan:anime_${id}:info`);
    if (cached) return okAsync(cached);

    try {
        const info = await client.getAnimeById(id);

        await cache.jsonSet(`jikan:anime_${id}:info`, info, 60 * 60 * 24);
        return okAsync(info);
    } catch (e) {
        return errAsync(e);
    }
}

async function getEpisodes(id: number, opts?: { page?: number }) {
    const cached = await cache.jsonGet<JikanResponse<AnimeEpisode[]>>(`jikan:anime_${id}:episodes:page_${opts?.page || "1"}`);
    if (cached) return okAsync(cached);

    try {
        const episodes = await client.getAnimeEpisodes(id, opts?.page);

        await cache.jsonSet(`jikan:anime_${id}:episodes:page_${opts?.page || "1"}`, episodes, 60 * 60 * 24);
        return okAsync(episodes);
    } catch (e) {
        return errAsync(e);
    }
}

export default {
    getAnime,
    getEpisodes,
}