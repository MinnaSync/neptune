import { gql, GraphQLClient } from "graphql-request";

/**
 * Use https://studio.apollographql.com/sandbox/explorer to test queries.
 * Set sandbox to https://graphql.anilist.co
 */
export const client = new GraphQLClient("https://graphql.anilist.co");

/**
 * Data ths is returned by the {@link ANIME_MEDIA_QUERY} query.
 */
export type AnimeMedia = {
    id: number;

    coverImage: {
        color: string | null;
        large: string | null;
    };
    bannerImage: string | null;

    title: {
        english: string;
		romaji: string;
		native: string;
    };
    description: string | null;

	seasonYear: number;
    format: string;
    meanScore: number;
	popularity: number;
    genres: string[];
	studios: {
		nodes: { name: string }[];
	};
	trailer: {
		id: string;
		site: string;
		thumbnail: string;
	};
    isAdult: boolean;

	streamingEpisodes: {
		title: string;
	}[];
}

/**
 * Query that returns data matching the {@link AnimeMedia} type.
 */
export const ANIME_MEDIA_QUERY =  gql`
	query Media($mediaId: Int, $type: MediaType) {
		Media(id: $mediaId, type: $type) {
			id

			coverImage {
				color
				large
			}
			bannerImage

			title {
				english
				romaji
				native
			}
			description

			seasonYear
			nextAiringEpisode {
			airingAt
			episode
			}

			format
			meanScore
			popularity
			genres
			studios {
				nodes {
					name
				}
			}
			trailer {
				id
				site
				thumbnail
			}
			isAdult
		}
	}
`;