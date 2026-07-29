# Go API Boilerplate

Boilerplate Go API dengan gaya folder seperti project TMS Universal.

Stack default:

- Go `net/http` `ServeMux`
- PostgreSQL
- GORM
- Goose migration
- JWT middleware
- API key middleware
- Handler-Service-Repository per module

## Inisiasi Project Baru

Copy folder ini ke nama project baru, lalu jalankan:

```powershell
go mod edit -module github.com/username/nama-project
go mod tidy
Copy-Item .env.example .env
```

Update import path:

```powershell
.\scripts\rename-module.ps1 -OldModule go-api-boilerplate -NewModule github.com/username/nama-project
```

Jalankan aplikasi:

```powershell
go run ./cmd/api
```

Jalankan PostgreSQL lokal:

```powershell
docker compose up -d
```

## Struktur

```text
cmd/api                 entry point HTTP API
internal/config         load env config
internal/database       koneksi PostgreSQL via GORM
internal/http           router dan middleware
internal/jwt            helper JWT
internal/modules        feature modules
internal/shared         response, validation, context helper
migrations              Goose migrations embedded
```

## Menambah Module Baru

Untuk module `product`, buat:

```text
internal/modules/product/model.go
internal/modules/product/repository.go
internal/modules/product/service.go
internal/modules/product/handler.go
internal/modules/product/response.go
```

Lalu register di `internal/http/router.go`:

```go
productRepo := product.NewRepository(gormDB)
productService := product.NewService(productRepo)
product.NewHandler(productService).Register(mux)
```

Tambahkan migration baru:

```powershell
goose -dir migrations create create_products sql
```

## Catatan Migration

Gunakan Goose untuk schema database. GORM dipakai untuk query aplikasi, bukan `AutoMigrate` production.

## Panduan Style

Lihat [docs/project-style.md](docs/project-style.md) untuk aturan layer dan pola module.
