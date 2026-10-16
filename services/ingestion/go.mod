module github.com/hesam-fattahi/trading-signals-platform/services/ingestion

go 1.27.1

require (
	github.com/gorilla/websocket v1.5.3
	github.com/hesam-fattahi/trading-signals-platform/gen/go v0.1.0
	github.com/segmentio/kafka-go v0.4.51
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/klauspost/compress v1.20.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.31 // indirect
)

replace github.com/hesam-fattahi/trading-signals-platform/gen/go => ../../gen/go
