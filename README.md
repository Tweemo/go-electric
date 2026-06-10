# go-electric

Compare NZ electricity costs across providers using your real usage data.

**[Try it live →](https://go-electric-web.onrender.com/)**

---

## What it does

New Zealand households can choose from multiple power companies, each with several pricing plans — but comparing them with your *actual* usage is tedious. go-electric lets you upload your usage CSV (from your current provider's portal) and instantly see what you would have paid across Contact Energy and Nova Energy's available plans, with the cheapest option highlighted.

## How it works

1. Export your electricity usage CSV from your current provider's portal
2. Upload it at [go-electric-web.onrender.com](https://go-electric-web.onrender.com/)
3. See a cost breakdown for every provider and plan, sorted cheapest first

Costs are calculated from your real hourly usage data — not averages or estimates.

Your privacy is protected by design: identifying columns (ICP, meter number, name,
account) are stripped in your browser, so only anonymous date-and-usage figures are
ever uploaded. The backend processes those entirely in memory and stores nothing.

## Tech stack

- **Backend** — Go with the [Gin](https://github.com/gin-gonic/gin) HTTP framework; parses usage CSVs in-memory and applies each provider's pricing model
- **Frontend** — Next.js (App Router) + TypeScript + Tailwind CSS v4 + shadcn/ui
- **Infrastructure** — Docker (multi-stage builds), deployed on [Render](https://render.com/) via Blueprint

## Local development

The quickest way to run both services:

```bash
docker compose up --build
```

Backend available at `http://localhost:8080`, frontend at `http://localhost:3001`.

For hot-reload during development, run them separately:

```bash
# Backend
cd backend && go run .

# Frontend (separate terminal)
cd frontend && npm install && npm run dev
```

See `backend/.env.example` and `frontend/.env.example` for environment variable reference.
