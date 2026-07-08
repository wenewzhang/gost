package relay

import (
	"time"

	mdata "github.com/wenewzhang/core/metadata"
	mdutil "github.com/wenewzhang/core/metadata/util"
	"github.com/go-gost/relay"
	"github.com/google/uuid"
)

type metadata struct {
	connectTimeout time.Duration
	noDelay        bool
	tunnelID       relay.TunnelID
}

func (c *relayConnector) parseMetadata(md mdata.Metadata) (err error) {
	const (
		connectTimeout = "connectTimeout"
		noDelay        = "nodelay"
	)

	c.md.connectTimeout = mdutil.GetDuration(md, connectTimeout)
	c.md.noDelay = mdutil.GetBool(md, noDelay)

	if s := mdutil.GetString(md, "tunnelID", "tunnel.id"); s != "" {
		uuid, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		c.md.tunnelID = relay.NewTunnelID(uuid[:])
	}

	return
}
