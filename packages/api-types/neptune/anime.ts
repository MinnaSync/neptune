export type Episodes = {
    id: string;
    episode: number;
    title: string | null;
    preview: string;
    streaming_link: string;    
};

export type Meta = {
    id: number;

    color: string | null;
    poster: string | null;
    background: string | null;

    title: {
        english: string;
        romaji: string;
        native: string;
    };
    description: string | null;

    year: number;
    type: string;
    rating: number | null;
    genres: string[];
    studios: string[];
    is_nsfw: boolean;

    trailer: {
        id: string;
        platform: string;
    } | null;
};

export interface AnimeInfo {
    meta: Meta;
    details: {
        hasNextPage: boolean;
        episodes: Episodes[];
    };
};

export interface AnimeSearch {
    results: {
        id: string;
        title: string;
        poster: string;
        type: string;
        episodes: number;
        year: number;
    }[];
};

export type AnimeStreams = Record<string, {
    resolution: string;
    link: string;
}[]>;