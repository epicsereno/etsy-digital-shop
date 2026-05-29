# Etsy Digital Shop

Master strategy doc and quick-start SOP for digital product production, listing, marketing, and tracking.

## Orders Service

Small Go service for digital-product order workflows. It currently exposes a sample sale endpoint for the Dark Botanical Wallpaper Pack.

```sh
go run ./services/orders/cmd/ordersd
```

Create the sample sale to John Doe:

```sh
curl -X POST http://localhost:8080/orders/sample-john-doe
```

List recorded orders:

```sh
curl http://localhost:8080/orders
```
