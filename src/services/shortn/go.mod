module github.com/Menschomat/bly.li/services/shortn

go 1.22

toolchain go1.22.5

replace github.com/Menschomat/bly.li/shared => ../../shared

require (
	github.com/Menschomat/bly.li/shared v0.0.0
	github.com/coreos/go-oidc/v3 v3.11.0
	github.com/go-chi/chi/v5 v5.1.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-jose/go-jose/v4 v4.0.4 // indirect
	github.com/redis/go-redis/v9 v9.6.1 // indirect
	github.com/sethvargo/go-envconfig v1.1.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
)