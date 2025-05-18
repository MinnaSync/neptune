declare module Bun {
    interface Env {
        /**
         * The Redis host URL.
         */
        REDIS_HOST?: string;
        /**
         * The Redis port.
         */
        REDIS_PORT?: string;
        /**
         * The password for the Redis connection.
         * Leave empty if there is no password.
         */
        REDIS_PASSWORD?: string;
    }
}