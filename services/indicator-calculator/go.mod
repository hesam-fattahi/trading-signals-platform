module github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator

go 1.27.1

require (
	github.com/TA-Lib/ta-lib-cgo v0.2.0
	github.com/hesam-fattahi/trading-signals-platform/gen/go v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.51
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

replace github.com/hesam-fattahi/trading-signals-platform/gen/go => ../../gen/go
