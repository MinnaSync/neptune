declare module Bun {
    interface Env {
        /**
         * The port to run the server on.
         * Leave blank to use the default port.
         * Default is 8443.
         */
        PORT?: string;

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