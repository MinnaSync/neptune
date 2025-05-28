import { serve } from 'bun';
import 'dotenv/config';
import './util/cache';
import app from './api/app';
// import ids, { Resource } from './resources/ids';

// const test = await ids.fromResource(Resource.ANILIST, 179965);
// if (test.isOk()) {
//     console.log(test.value);
// }

serve({
    port: process.env.PORT ? parseInt(process.env.PORT) : 8444,
    fetch: app.fetch,
});