# Deploy (free, persistent)

The live demo runs on two free services: **Neon** for PostgreSQL (permanent free tier)
and **Render** for the Go server (free web service, built from the `Dockerfile`).

## 1. Database — Neon

1. Sign up at <https://neon.tech> (you can use GitHub).
2. Create a project. Pick the **Singapore** region (closest to Indonesia).
3. Copy the **connection string** it shows — it looks like
   `postgresql://user:password@ep-xxx.ap-southeast-1.aws.neon.tech/neondb?sslmode=require`.
   Keep it; you paste it into Render in the next step.

## 2. Server — Render

1. Sign up at <https://render.com> and connect your GitHub.
2. **New → Blueprint**, pick the `practicum-score-graphql` repo. Render reads `render.yaml`.
3. When it asks for `DATABASE_URL`, paste the Neon connection string from step 1.
4. **Apply**. The first build takes a few minutes (it builds the Docker image).
5. Open the service URL. You should see the **GraphQL playground**, already seeded with
   synthetic demo data. Try:

   ```graphql
   { classSummary(className: "A") { students average attendanceRate } }
   ```

## Good to know

- The free Render web service **sleeps after ~15 minutes of no traffic** and takes about a
  minute to wake on the next request. Fine for a portfolio demo; the first hit after idle is slow.
- The Neon free database is **permanent** (it does not expire), unlike Render's own free Postgres.
- `-seed` is safe to run on every start: it only inserts data when the database has no students yet.
- Once it is live, put the URL in this repo's **About → homepage** and in your portfolio.
