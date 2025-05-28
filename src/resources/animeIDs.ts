import { okAsync, errAsync } from "neverthrow";
import cache from "../util/cache";

export type AnimeIds = {
    livechart_id: number;
    'anime-planet_id': string;
    anisearch_id: number;
    anidb_id: number;
    kitsu_id: number;
    mal_id: number;
    'notify.moe_id': string;
    anilist_id: number;
    thetvdb_id: number;
    imdb_id: string;
    themoviedb_id: number;
    animecountdown_id: number;
};

type AnimeIdKeys = keyof AnimeIds;

export enum Resource {
    MAL,
    ANILIST,
    ANIDB,
}

async function fetchIdsList() {
    const cached = await cache.jsonGet<AnimeIds[]>("internal:anime_ids");
    if (cached) return okAsync(cached);

    const res = await fetch("https://raw.githubusercontent.com/Fribb/anime-lists/refs/heads/master/anime-list-full.json");
    if (!res.ok) {
        return errAsync(new Error("Failed to fetch anime ids."));
    }

    const json: AnimeIds[]= await res.json();

    await cache.jsonSet("internal:anime_ids", json, 60 * 60 * 24);
    return okAsync(json);
}

function findIds(list: AnimeIds[], key: AnimeIdKeys, value: string) {
    return list.find((ids) => ids[key]?.toString() === value);
}

async function fromResource(resource: Resource, id: number | string) {
    const list = await fetchIdsList();
    if (list.isErr()) {
        return list;
    }

    const ids = list.value;
    
    switch (resource) {
        case Resource.MAL:
            const malIds = findIds(ids, "mal_id", id.toString());
            if (!malIds) return errAsync(new Error(`No IDs found for MAL: ${id}`));
            
            return okAsync(malIds);
        case Resource.ANILIST:
            const anilistIds = findIds(ids, "anilist_id", id.toString());
            if (!anilistIds) return errAsync(new Error(`No IDs found for Anilist: ${id}`));
            
            return okAsync(anilistIds);
        case Resource.ANIDB:
            const anidbIds = findIds(ids, "anidb_id", id.toString());
            if (!anidbIds) return errAsync(new Error(`No IDs found for AniDB: ${id}`));

            return okAsync(anidbIds);
        default:
            return errAsync(new Error(`Unknown resource ${resource}`));
    }
}

export default {
    fromResource
}