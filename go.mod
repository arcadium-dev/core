module arcadium.dev/core

go 1.26.4

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/google/go-cmp v0.7.0
	github.com/gorilla/mux v1.8.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.20.2
	github.com/rs/cors v1.11.1
	github.com/rs/zerolog v1.33.0
	golang.org/x/exp v0.0.0-20240409090435-93d18d7e34b8
	golang.org/x/vuln v1.3.0
	honnef.co/go/tools v0.7.0
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.org/x/exp/typeparams v0.0.0-20260611194520-c48552f49976 // indirect
	golang.org/x/mod v0.37.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/telemetry v0.0.0-20260611141451-d61e87d5f4a3 // indirect
	golang.org/x/tools v0.46.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

tool (
	golang.org/x/vuln
	golang.org/x/vuln/cmd/govulncheck
	honnef.co/go/tools/cmd/staticcheck
)
