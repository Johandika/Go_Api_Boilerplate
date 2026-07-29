# Project Style

Panduan singkat agar project turunan dari boilerplate ini tetap konsisten.

## Layer

- `handler.go` hanya mengurus HTTP request dan response.
- `service.go` mengurus business logic, validasi domain, dan error domain.
- `repository.go` mengurus query database dengan GORM.
- `model.go` merepresentasikan table database.
- `response.go` mengurus bentuk JSON response.

## Alur Dependency

```text
config -> database -> repository -> service -> handler -> router -> middleware -> server
```

## Database

- Gunakan PostgreSQL + GORM untuk query aplikasi.
- Gunakan Goose untuk migration schema.
- Hindari `AutoMigrate` untuk project production.
- Perubahan table ditulis di `migrations/*.sql`.

## Middleware

Urutan middleware standar:

```text
CORS -> logging -> API key -> JWT auth -> permission -> mux
```

## Module Baru

Untuk module `product`, buat:

```text
internal/modules/product/model.go
internal/modules/product/repository.go
internal/modules/product/service.go
internal/modules/product/handler.go
internal/modules/product/response.go
```

Register di `internal/http/router.go`.

## Response

Gunakan helper dari `internal/shared`:

- `shared.WriteSuccess`
- `shared.WritePaginated`
- `shared.WriteError`
- `shared.WriteInternalError`

## Error

Buat error domain eksplisit di service:

```go
var ErrProductNotFound = errors.New("product not found")
```

Map error tersebut ke HTTP status di handler.
