import { serve } from 'bun';
import { Hono } from "hono";
import 'dotenv/config';
import './util/cache';

import middleware from "./routes/middleware";
import anime from "./routes/anime";
import logger from './util/logger';

const app = new Hono();

app.use(middleware.logRequest);

app.route("/anime", anime);

const port = process.env.PORT ? parseInt(process.env.PORT) : 8444;

logger.info(`Starting server on port ${port}.`);

serve({
    port: port,
    fetch: app.fetch,
})