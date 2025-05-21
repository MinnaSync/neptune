export type Episodes = {
    /**
     * The ID given by the provider for the episode.
     */
    id: string;
    /**
     * The episode number.
     */
    episode: number;
    /**
     * The title of the episode.
     */
    title: string | null;
    /**
     * A snapshot preview of the episode's contents.
     */
    preview: string;
    /**
     * The direct link to play the episode on the provider's platform.
     */
    streaming_link: string;    
};

export type Meta = {
    /**
     * The ID for the anime given by the resource.
     */
    id: number;

    /**
     * An accent color for the anime.
     * Only Anilist has this. Usually it will be null.
     */
    color: string | null;
    /**
     * The poster image for the anime.
     */
    poster: string | null;
    /**
     * A background cover image for the anime.
     */
    background: string | null;

    /**
     * Different titles that the anime goes by.
     */
    title: {
        /**
         * The title in English.
         */
        english: string;
        /**
         * The title in Romaji.
         */
        romaji: string;
        /**
         * The title in native language.
         */
        native: string;
    };
    /**
     * A description of the anime.
     */
    description: string | null;

    /**
     * The year that the anime was released.
     */
    year: number;
    /**
     * The type of anime
     */
    type: string;
    /**
     * The rating that the anime has.
     */
    rating: number | null;
    /**
     * The genres that the anime belongs in.
     */
    genres: string[];
    /**
     * The studios that created the anime.
     */
    studios: string[];
    /**
     * Whether or not the anime is meant for adult audiences.
     */
    is_nsfw: boolean;

    /**
     * A trailer for the anime.
     */
    trailer: {
        id: string;
        platform: string;
    } | null;
};

export interface AnimeInfo {
    /**
     * The meta given by the resource.
     */
    meta: Meta;
    /**
     * The details given by the provider.
     */
    details: {
        /**
         * Whether or not there's another page of episodes.
         */
        hasNextPage: boolean;
        /**
         * The list of episodes for the anime.
         */
        episodes: Episodes[];
    };
};

export interface AnimeSearch {
    /**
     * The search results.
     */
    results: {
        /**
         * The ID given by the provider for the anime.
         */
        id: string;
        /**
         * The title of the anime.
         */
        title: string;
        /**
         * The poster image for the anime.
         */
        poster: string;
        /**
         * The type of anime (TV, OVA, ONA, etc.)
         */
        type: string;
        /**
         * How many episodes that the anime has.
         */
        episodes: number;
        /**
         * The year that the anime was released in.
         */
        year: number;
    }[];
};

export type AnimeStreams = Record<string, {
    /**
     * The resolution of the stream.
     */
    resolution: string;
    /**
     * The link to the stream file.
     */
    link: string;
}[]>;