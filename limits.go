package libp2p

import (
	"github.com/libp2p/go-libp2p/p2p/host/autonat"
	circuitv1 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv1/relay"
	circuitv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/libp2p/go-libp2p/p2p/protocol/holepunch"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"

	rcmgr "github.com/libp2p/go-libp2p-resource-manager"
)

// SetDefaultServiceLimits sets the default limits for bundled libp2p services
//
// More specifically this sets the following limits:
// - identify:
//   128 streams in, 128 streams out, 256 streams total, 4MB min, 64MB max svc memory
//   16/16/32 streams per peer
// - ping:
//   128 streams in, 128 sreams out, 256 streasms total, 4MB min, 64MB max svc memory
//   2/3/4 streams per peer
// - autonat
//   128 streams in, 128 streams out, 128 streams total, 4MB min, 64MB max svc memory
//   2/2/2 streams per peer
// - holepunch
//   128 streams in, 128 streams out, 128 streams total, 4MB min, 64MB max svc memory
//   2/2/2 streams per peer
// - relay v1 and v2 (separate services)
//   1024 streams in, 1024 streams out, 1024 streams total, 4MB min, 64MB max svc memory
//   64/64/64 streams per peer
func SetDefaultServiceLimits(limiter *rcmgr.BasicLimiter) {
	peerSvcLimit := func(numStreamsIn, numStreamsOut, numStreamsTotal int) rcmgr.Limit {
		return &rcmgr.StaticLimit{
			// memory: 256kb for window buffers plus some change for message buffers per stream
			Memory: int64(numStreamsTotal * (256<<10 + 16384)),
			BaseLimit: rcmgr.BaseLimit{
				StreamsInbound:  numStreamsIn,
				StreamsOutbound: numStreamsOut,
				Streams:         numStreamsTotal,
			},
		}
	}

	if limiter.ServiceLimits == nil {
		limiter.ServiceLimits = make(map[string]rcmgr.Limit)
	}
	if limiter.ServicePeerLimits == nil {
		limiter.ServicePeerLimits = make(map[string]rcmgr.Limit)
	}
	// identify
	if _, ok := limiter.ServiceLimits[identify.ServiceName]; !ok {
		limiter.ServiceLimits[identify.ServiceName] = limiter.DefaultServiceLimits.
			WithMemoryLimit(1, 4<<20, 64<<20). // max 64MB service memory
			WithStreamLimit(128, 128, 256)     // max 256 streams -- symmetric
		limiter.ServicePeerLimits[identify.ServiceName] = peerSvcLimit(16, 16, 32)
	}
	// ping
	if _, ok := limiter.ServiceLimits[ping.ServiceName]; !ok {
		limiter.ServiceLimits[ping.ServiceName] = limiter.DefaultServiceLimits.
			WithMemoryLimit(1, 4<<20, 64<<20). // max 64MB service memory
			WithStreamLimit(128, 128, 128)     // max 128 streams - asymmetric
		limiter.ServicePeerLimits[ping.ServiceName] = peerSvcLimit(2, 3, 4)
	}
	// autonat
	if _, ok := limiter.ServiceLimits[autonat.ServiceName]; !ok {
		limiter.ServiceLimits[autonat.ServiceName] = limiter.DefaultServiceLimits.
			WithMemoryLimit(1, 4<<20, 64<<20). // max 64MB service memory
			WithStreamLimit(128, 128, 128)     // max 128 streams - asymmetric
		limiter.ServicePeerLimits[autonat.ServiceName] = peerSvcLimit(2, 2, 2)
	}
	// holepunch
	if _, ok := limiter.ServiceLimits[holepunch.ServiceName]; !ok {
		limiter.ServiceLimits[holepunch.ServiceName] = limiter.DefaultServiceLimits.
			WithMemoryLimit(1, 4<<20, 64<<20). // max 64MB service memory
			WithStreamLimit(128, 128, 256)     // max 256 streams - symmetric
		limiter.ServicePeerLimits[autonat.ServiceName] = peerSvcLimit(2, 2, 2)
	}
	// relay/v1
	if _, ok := limiter.ServiceLimits[circuitv1.ServiceName]; !ok {
		limiter.ServiceLimits[circuitv1.ServiceName] = limiter.DefaultServiceLimits.
			WithMemoryLimit(1, 4<<20, 64<<20). // max 64MB service memory
			WithStreamLimit(1024, 1024, 1024)  // max 1024 streams - asymmetric
		limiter.ServicePeerLimits[circuitv1.ServiceName] = peerSvcLimit(64, 64, 64)
	}
	// relay/v2
	if _, ok := limiter.ServiceLimits[circuitv2.ServiceName]; !ok {
		limiter.ServiceLimits[circuitv2.ServiceName] = limiter.DefaultServiceLimits.
			WithMemoryLimit(1, 4<<20, 64<<20). // max 64MB service memory
			WithStreamLimit(1024, 1024, 1024)  // max 1024 streams - asymmetric
		limiter.ServicePeerLimits[circuitv2.ServiceName] = peerSvcLimit(64, 64, 64)
	}
}
