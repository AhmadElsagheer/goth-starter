import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { auth } from './lib/auth';
import { env } from './config';

const app = new Hono()

app.use(
  "/api/auth/*", // or replace with "*" to enable cors for all routes
  cors({
    origin: env.BETTER_AUTH_TRUSTED_ORIGINS,
    allowHeaders: ["Content-Type", "Authorization"],
    allowMethods: ["POST", "GET", "OPTIONS"],
    exposeHeaders: ["Content-Length"],
    maxAge: 600,
    credentials: true,
  }),
);


app.on(["POST", "GET"], "/api/auth/*", (c) => auth.handler(c.req.raw));


export default {
  port: env.SERVER_PORT,
  fetch: app.fetch,
}
