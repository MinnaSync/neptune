import { serve } from 'bun';
import 'dotenv/config';
import './util/cache';
import app from './api/app';

serve({
    port: process.env.PORT ? parseInt(process.env.PORT) : 8444,
    fetch: app.fetch,
});