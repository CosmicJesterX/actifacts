//go:build !windows
// +build !windows

package trace2receiver

import (
	"context"
	"errors"
	"net"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

var (
	errNilNextConsumer = errors.New("nil next Consumer")
)

func createTraces(_ context.Context,
	params receiver.CreateSettings,
	baseCfg component.Config,
	consumer consumer.Traces) (receiver.Traces, error) {

	if consumer == nil {
		return nil, errNilNextConsumer
	}

	trace2Cfg := baseCfg.(*Config)

	rcvr := &Rcvr_UnixSocket{
		Base: &Rcvr_Base{
			Settings:       params,
			Logger:         params.Logger,
			TracesConsumer: consumer,
			RcvrConfig:     trace2Cfg,
		},
		SocketPath: trace2Cfg.UnixSocketPath,
	}
	return rcvr, nil
}

// Gather up any requested PII from the machine or
// possibly the connection from the client process.
// Add any requested PII data to `tr2.pii[]`.
func (tr2 *trace2Dataset) pii_gather(cfg *Config, conn *net.UnixConn) {
	if cfg.piiSettings != nil && cfg.piiSettings.Include.Hostname {
		if h, err := os.Hostname(); err == nil {
			tr2.pii[string(Trace2PiiHostname)] = h
		}
	}

	if cfg.piiSettings != nil && cfg.piiSettings.Include.Username {
		if u, err := getPeerUsername(conn); err == nil {
			tr2.pii[string(Trace2PiiUsername)] = u
		}
	}
}
