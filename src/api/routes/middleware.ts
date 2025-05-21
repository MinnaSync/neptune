import type { Context, Next } from "hono";
import logger from "../../util/logger";

async function logRequest(c: Context, next: Next) {
    const method = c.req.method;

    const url = new URL(c.req.url);
    const path = url.pathname;
    const query = url.searchParams.toString();

    const start = Date.now();
    await next();
    const end = Date.now();

    const duration = `${end - start}ms`.padEnd(8);
    const status = c.res.status;

    logger.info(`${status} | ${duration} | ${method} ${path}${query ? `?${query}` : ''}`);
}

export default {
    logRequest,
};