import { client, ANIME_MEDIA_QUERY, type AnimeMedia } from "./query";
import cache from "../util/cache";

async function getAnime(id: number) {
    const cached = await cache.jsonGet<AnimeMedia>(`anilist:${id}`);
    if (cached) return cached;

    const { Media } = await client.request<{ Media: AnimeMedia }>(ANIME_MEDIA_QUERY, {
        mediaId: id,
        type: "ANIME"
    });

    await cache.jsonSet(`anilist:${id}`, Media, 60 * 60 * 24);

    return Media;
}

export default {
    getAnime,
}