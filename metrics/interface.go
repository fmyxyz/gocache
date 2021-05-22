package metrics

import "github.com/fmyxyz/gocache/codec"

// MetricsInterface represents the metrics interface for all available providers
type MetricsInterface interface {
	RecordFromCodec(codec codec.CodecInterface)
}
