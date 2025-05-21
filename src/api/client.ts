import { hc } from 'hono/client'
import type { AppType } from './app';

const client = hc<AppType>(`http://localhost:${process.env.PORT}`);

export default client;