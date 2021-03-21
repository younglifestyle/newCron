module myCron

go 1.14

require (
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/sony/sonyflake v1.0.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/genproto v0.0.0-20210319143718-93e7006c17a6 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

//replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
