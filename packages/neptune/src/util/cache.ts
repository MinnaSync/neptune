import ioredis from "ioredis";
import logger from "./logger";

const cache =
    (process.env.REDIS_HOST && process.env.REDIS_PORT
    && new ioredis({
        port: parseInt(process.env.REDIS_PORT),
        host: process.env.REDIS_HOST,
        password: process.env.REDIS_PASSWORD,
    })) as ioredis | null;

if (!cache) {
    logger.warn("Unable to initialize Redis, caching will be disabled.");
}

async function jsonGet<T>(key: string): Promise<T | null> {
    try {
        const json = await cache?.call("JSON.GET", key) as string;
        return JSON.parse(json) as T;
    } catch (e) {
        return null;
    }
}

async function jsonSet<T>(key: string, value: T, ttl?: number): Promise<boolean> {
    try {
        await cache?.call("JSON.SET", key, "$", JSON.stringify(value));

        if (ttl) {
            await cache?.call("EXPIRE", key, ttl);
        }

        return true;
    } catch {
        return false;
    }
}

export default {
    jsonGet,
    jsonSet,
}