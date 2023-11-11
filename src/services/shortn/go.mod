module github.com/Menschomat/bly.li/services/shortn

go 1.21

toolchain go1.21.3

replace github.com/Menschomat/bly.li/shared => ../../shared

require (
	github.com/Menschomat/bly.li/shared v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.10
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/redis/go-redis/v9 v9.3.0 // indirect
)
