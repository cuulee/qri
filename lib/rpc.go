package lib

import (
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
)

// localHostIP is the IP address for localhost
const localHostIP = "127.0.0.1"

func init() {
	// Fields like dataset.Structure.Schema contain data of arbitrary types,
	// registering with the gob package prevents errors when sending them
	// over net/rpc calls.
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
}

// Receivers returns a slice of CoreRequests that defines the full local
// API of lib methods
func Receivers(inst *Instance) []Methods {
	node := inst.Node()
	r := inst.Repo()

	return []Methods{
		NewDatasetRequests(node, nil),
		NewRegistryRequests(node, nil),
		NewLogRequests(node, nil),
		NewExportRequests(node, nil),
		NewPeerRequests(node, nil),
		NewProfileMethods(inst),
		NewConfigMethods(inst),
		NewSearchRequests(node, nil),
		NewRenderRequests(r, nil),
		NewSelectionRequests(r, nil),
	}
}

// ErrRPCDisabled indicates a configuration that expressly disables
// remote procedure calls
var ErrRPCDisabled = fmt.Errorf("RPC is disabled")

// ServeRPC creates a server that listens for remote procedure call operations
// it listens on the
func ServeRPC(ctx context.Context, inst *Instance) (err error) {
	cfg := inst.Config()
	if !cfg.RPC.Enabled || cfg.RPC.Port == 0 {
		return ErrRPCDisabled
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", localHostIP, cfg.RPC.Port))
	if err != nil {
		log.Infof("RPC listen on port %d error: %s", cfg.RPC.Port, err)
		return
	}

	for _, rcvr := range Receivers(inst) {
		if err := rpc.Register(rcvr); err != nil {
			log.Errorf("error registering RPC receiver %s: %s", rcvr.CoreRequestsName(), err.Error())
			return err
		}
	}

	errs := make(chan error)
	// close the listener when the host context closes
	go func(l net.Listener) {
		<-ctx.Done()
		errs <- l.Close()
	}(listener)

	// start the RPC server
	go rpc.Accept(listener)

	return <-errs
}
