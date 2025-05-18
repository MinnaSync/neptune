import { serve } from 'bun';
import { Hono } from "hono";
import 'dotenv/config';
import './util/cache';

import middleware from "./routes/middleware";
import anime from "./routes/anime";

const app = new Hono();

app.use(middleware.logRequest);

app.route("/anime", anime);

serve({
    port: 8444,
    fetch: app.fetch,
})