# Deploy a live demo

The `render.yaml` blueprint deploys the Go server (from the `Dockerfile`) **and** a
PostgreSQL database on [Render](https://render.com), wiring them together for you.

## Steps (all on render.com)

1. Sign up at <https://render.com> — use **Continue with GitHub** so it can read your repos.
2. **New → Blueprint**, pick the `practicum-score-graphql` repo. Render reads `render.yaml`.
3. Click **Apply**. The first build takes a few minutes (it builds the Docker image and
   creates the database).
4. Open the service URL. You should see the **GraphQL playground**, already seeded with
   synthetic demo data. Try:

   ```graphql
   { classSummary(className: "A") { students average attendanceRate } }
   ```

That's it — no connection string to copy, no database to configure by hand.

## Good to know

- The free web service **sleeps after ~15 minutes of no traffic** and takes about a minute
  to wake on the next request. Fine for a portfolio demo; the first hit after idle is slow.
- Render's **free PostgreSQL expires ~30 days** after creation. That covers a demo/application
  window. For a link that lasts, move the database to **Neon** (free, persistent): create a
  Neon project, then replace the `databases:` block and the `fromDatabase` wiring in
  `render.yaml` with a plain `DATABASE_URL` set to the Neon connection string (`?sslmode=require`).
- `-seed` is safe on every start: it only inserts data when the database has no students yet.
- If **Apply** complains that `region: singapore` is not available for the free database,
  change both `region:` lines in `render.yaml` to `oregon`.
- If Render rejects `runtime: docker`, change it to `env: docker` (older accounts use that key).
