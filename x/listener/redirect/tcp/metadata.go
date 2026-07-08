package tcp

import (
	mdata "github.com/wenewzhang/core/metadata"
	mdutil "github.com/wenewzhang/core/metadata/util"
)

type metadata struct {
	tproxy bool
}

func (l *redirectListener) parseMetadata(md mdata.Metadata) (err error) {
	const (
		tproxy = "tproxy"
	)
	l.md.tproxy = mdutil.GetBool(md, tproxy)
	return
}
