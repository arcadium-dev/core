module arcadium.dev/core

go 1.17

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/go-kit/log v0.2.1
	github.com/gorilla/mux v1.8.0
	github.com/jackc/pgx/v4 v4.15.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.12.2
	google.golang.org/grpc v1.46.2
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.11.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.10.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.34.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

// CVE-2022-28948: gopkg.in/yaml.v3
exclude (
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/prometheus/common v0.32.1
	github.com/stretchr/testify v1.4.0
	github.com/stretchr/testify v1.5.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)
