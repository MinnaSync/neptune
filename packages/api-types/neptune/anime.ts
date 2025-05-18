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

    title: string;
    description: string | null;

    type: string;
    rating: number | null;
    genres: string[];
    is_nsfw: boolean;
};

export interface AnimeInfo {
    meta: Meta;
    episodes: Episodes[];
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