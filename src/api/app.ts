import { Hono } from "hono";
import middleware from "./routes/middleware";
import anime from "./routes/anime";
import { cors } from "hono/cors";

const app = new Hono();
app.use(middleware.logRequest);
app.use(cors({
    origin: "https://sync.minna.now",
    allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
}))

const routes = app
    .route('/anime', anime);

export default app;
export type AppType = typeof routes;
