package relay

import (
	"math"
	"time"

	mdata "github.com/wenewzhang/core/metadata"
	mdutil "github.com/wenewzhang/core/metadata/util"
)

type metadata struct {
	readTimeout   time.Duration
	enableBind    bool
	udpBufferSize int
	noDelay       bool
	hash          string
	entryPoint    string
	directTunnel  bool
}

func (h *relayHandler) parseMetadata(md mdata.Metadata) (err error) {
	// fmt.Fprintf("lite-gost x handle relay \n")
	const (
		readTimeout   = "readTimeout"
		enableBind    = "bind"
		udpBufferSize = "udpBufferSize"
		noDelay       = "nodelay"
		hash          = "hash"
		entryPoint    = "entryPoint"
	)

	h.md.readTimeout = mdutil.GetDuration(md, readTimeout)
	h.md.enableBind = mdutil.GetBool(md, enableBind)
	h.md.noDelay = mdutil.GetBool(md, noDelay)

	if bs := mdutil.GetInt(md, udpBufferSize); bs > 0 {
		h.md.udpBufferSize = int(math.Min(math.Max(float64(bs), 512), 64*1024))
	} else {
		h.md.udpBufferSize = 4096
	}

	h.md.hash = mdutil.GetString(md, hash)

	h.md.entryPoint = mdutil.GetString(md, entryPoint)

	h.md.directTunnel = mdutil.GetBool(md, "tunnel.direct")

	return
}
