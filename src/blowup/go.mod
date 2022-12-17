module github.com/Menschomat/bly.li/blowup

go 1.19

replace github.com/Menschomat/bly.li/shared => ../shared

require github.com/go-chi/chi/v5 v5.0.8 // direct

require github.com/Menschomat/bly.li/shared v0.0.0-00010101000000-000000000000

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.2 // indirect
)
