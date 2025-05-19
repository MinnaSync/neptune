import { Hono } from "hono";
import middleware from "./routes/middleware";
import anime from "./routes/anime";

const app = new Hono();
app.use(middleware.logRequest);

const routes = app
    .route('/anime', anime);

export default app;
export type AppType = typeof routes;
