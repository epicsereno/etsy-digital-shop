# Etsy Digital Shop

Digital product workspace for planning Etsy listings, tracking production, organizing marketing, and testing lightweight order workflows.

Live site: https://epicsereno.github.io/etsy-digital-shop/

## What This Includes

This repository combines a static GitHub Pages workspace with supporting shop operations files and a small Go service prototype.

- `index.html` publishes the interactive workspace through GitHub Pages.
- `listings/active/` stores listing copy, tags, titles, and pricing for active products.
- `niches/` stores product research, palettes, keyword notes, and prompt ideas by niche.
- `prompts/midjourney/` stores repeatable image-generation prompt sets.
- `marketing/pinterest/` tracks Pinterest promotion planning.
- `operations/` stores SOPs, SEO research, revenue tracking, and batch logs.
- `services/orders/` contains the Go orders service prototype.

## GitHub Pages App

The published workspace is a single-page app for managing digital product work:

- Product inventory with status, category, revenue, and download metrics.
- Listing generator for title, description, and tag drafts.
- File intake area for mockups, PDFs, ZIP files, and design exports.
- Marketing calendar for campaign planning.
- Operations checklist for production and launch workflows.
- Light and dark theme support.
- CSV export for product tracking.

GitHub Pages serves the app from the `main` branch at the repository root.

Open the live app:

```text
https://epicsereno.github.io/etsy-digital-shop/
```

Open locally by loading `index.html` in a browser. No frontend build step is required.

## Shop Workflow

Use this repo as the working system for a digital Etsy shop:

1. Research a niche in `niches/<niche-name>/`.
2. Draft product ideas, palettes, keywords, and generation prompts.
3. Create product assets and collect source notes in the appropriate listing or prompt folder.
4. Draft title, description, tags, and pricing under `listings/active/`.
5. Schedule promotion in `marketing/pinterest/schedule.md`.
6. Track production runs and revenue in `operations/`.
7. Use the GitHub Pages workspace for a quick visual dashboard.

## Current Product Focus

The first active product area is the Dark Botanical bundle:

- Listing folder: `listings/active/listing-001-dark-botanical-bundle/`
- Niche research: `niches/dark-botanical/`
- Example product in the app: Dark Botanical Wallpaper Pack

The goal is to keep each product line reproducible: research, prompts, listing copy, pricing, and marketing notes should all live in version control.

## Orders Service

Small Go service for digital-product order workflows. It currently exposes a sample sale endpoint for the Dark Botanical Wallpaper Pack.

Run the service:

```sh
go run ./services/orders/cmd/ordersd
```

Create the sample sale to John Doe:

```sh
curl -X POST http://localhost:8080/orders/sample-john-doe
```

Check health:

```sh
curl http://localhost:8080/healthz
```

List recorded orders:

```sh
curl http://localhost:8080/orders
```

Run tests:

```sh
go test ./...
```

The service keeps data in memory for now. It is meant for quick workflow testing, not production order storage.

## Development Notes

Useful commands:

```sh
git status --short --branch
go test ./...
go run ./services/orders/cmd/ordersd
```

Before merging changes, confirm:

- The README still points to the live GitHub Pages URL.
- `index.html` works as a static file.
- Go tests pass for the orders service.
- New listing or operations files are placed in the correct folder.

## Publishing

Changes merged into `main` are published by GitHub Pages. Pages is configured for:

- Source branch: `main`
- Source path: `/`
- Public URL: https://epicsereno.github.io/etsy-digital-shop/
