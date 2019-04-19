package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/qri-io/ioes"
	"github.com/qri-io/qri/api"
	"github.com/qri-io/qri/lib"
	"github.com/spf13/cobra"
)

// NewConnectCommand creates a new `qri connect` cobra command for connecting to
// the d.web, local api, rpc server, and webapp
func NewConnectCommand(f Factory, ioStreams ioes.IOStreams) *cobra.Command {
	o := ConnectOptions{IOStreams: ioStreams}
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to the distributed web by spinning up a Qri node",
		Annotations: map[string]string{
			"group": "network",
		},
		Long: `
While it’s not totally accurate, connect is like starting a server. Running 
connect will start a process and stay there until you exit the process 
(ctrl+c from the terminal, or killing the process using tools like activity 
monitor on the mac, or the aptly-named “kill” command). Connect does three main 
things:
- Connect to the qri distributed network
- Connect to IPFS
- Start a local API server

When you run connect you are connecting to the distributed web, interacting with
peers & swapping data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, args); err != nil {
				return err
			}
			return o.Run()
		},
	}

	cmd.Flags().IntVarP(&o.APIPort, "api-port", "", 0, "port to start api on")
	cmd.Flags().IntVarP(&o.RPCPort, "rpc-port", "", 0, "port to start rpc listener on")
	cmd.Flags().IntVarP(&o.WebappPort, "webapp-port", "", 0, "port to serve webapp on")
	cmd.Flags().IntVarP(&o.DisconnectAfter, "disconnect-after", "", 0, "duration to keep connected in seconds, 0 means run indefinitely")

	cmd.Flags().BoolVarP(&o.DisableAPI, "disable-api", "", false, "disables api, overrides the api-port flag")
	cmd.Flags().BoolVarP(&o.DisableRPC, "disable-rpc", "", false, "disables rpc, overrides the rpc-port flag")
	cmd.Flags().BoolVarP(&o.DisableWebapp, "disable-webapp", "", false, "disables webapp, overrides the webapp-port flag")
	cmd.Flags().BoolVarP(&o.DisableP2P, "disable-p2p", "", false, "disable peer-2-peer networking")

	cmd.Flags().BoolVarP(&o.Setup, "setup", "", false, "run setup if necessary, reading options from environment variables")
	cmd.Flags().BoolVarP(&o.ReadOnly, "read-only", "", false, "run qri in read-only mode, limits the api endpoints")
	cmd.Flags().BoolVarP(&o.RemoteMode, "remote-mode", "", false, "run qri in remote mode")
	cmd.Flags().Int64VarP(&o.RemoteAcceptSizeMax, "remote-accept-size-max", "", -1, "when running as a remote, max size of dataset to accept, -1 for any size")
	cmd.Flags().StringVarP(&o.Registry, "registry", "", "", "specify registry to setup with. only works when --setup is true")

	return cmd
}

// ConnectOptions encapsulates state for the connect command
type ConnectOptions struct {
	ioes.IOStreams
	inst *lib.Instance

	stopRPC  context.CancelFunc
	stopAPI  context.CancelFunc
	stopP2P  context.CancelFunc
	stopCron context.CancelFunc

	APIPort         int
	RPCPort         int
	WebappPort      int
	DisconnectAfter int

	DisableAPI    bool
	DisableRPC    bool
	DisableWebapp bool
	DisableP2P    bool

	Registry            string
	Setup               bool
	ReadOnly            bool
	RemoteMode          bool
	RemoteAcceptSizeMax int64
}

// Complete adds any missing configuration that can only be added just before calling Run
func (o *ConnectOptions) Complete(f Factory, args []string) (err error) {
	qriPath := f.QriRepoPath()

	if o.Setup && !QRIRepoInitialized(qriPath) {
		so := &SetupOptions{
			IOStreams: o.IOStreams,
			IPFS:      true,
			Registry:  o.Registry,
			Anonymous: true,
		}
		if err = so.Complete(f, args); err != nil {
			return err
		}
		if err = so.DoSetup(f); err != nil {
			return err
		}
	} else if !QRIRepoInitialized(qriPath) {
		return fmt.Errorf("no qri repo exists")
	}

	if err = f.Init(); err != nil {
		return err
	}
	o.inst = f.Instance()
	return
}

// Run executes the connect command with currently configured state
func (o *ConnectOptions) Run() (err error) {
	cfg := o.inst.Config().Copy()

	if o.APIPort != 0 {
		cfg.API.Port = o.APIPort
	}
	if o.RPCPort != 0 {
		cfg.RPC.Port = o.RPCPort
	}
	if o.WebappPort != 0 {
		cfg.Webapp.Port = o.WebappPort
	}

	if o.DisconnectAfter != 0 {
		cfg.API.DisconnectAfter = o.DisconnectAfter
	}

	if o.ReadOnly {
		cfg.API.ReadOnly = true
	}
	if o.RemoteMode {
		cfg.API.RemoteMode = true
	}

	cfg.API.RemoteAcceptSizeMax = o.RemoteAcceptSizeMax

	if err := o.inst.ChangeConfig(cfg); err != nil {
		log.Error(err)
	}

	o.maybeStartRPC()
	o.maybeStartP2P()
	o.maybeStartAPI()
	o.maybeStartCron()
	o.maybeTimeoutProcess()

	go func() {
		for change := range o.inst.ConfigChanges() {
			cfg := change.New

			if !cfg.RPC.Enabled && o.stopRPC != nil {
				log.Info("stopping RPC")
				o.stopRPC()
				o.stopRPC = nil
			} else if cfg.RPC.Enabled && o.stopRPC == nil {
				o.maybeStartRPC()
			}

			if !cfg.P2P.Enabled && o.stopP2P != nil {
				log.Info("stopping P2P")
				o.stopP2P()
				o.stopP2P = nil
			} else if cfg.P2P.Enabled && o.stopP2P == nil {
				o.maybeStartP2P()
			}

			if !cfg.API.Enabled && o.stopAPI != nil {
				log.Info("stopping API")
				o.stopAPI()
				o.stopAPI = nil
			} else if cfg.API.Enabled && o.stopAPI == nil {
				o.maybeStartAPI()
			}

			// TODO (b5): wire up cron checking once cron config exists
		}
	}()

	<-o.inst.Context().Done()
	return
}

func (o *ConnectOptions) maybeStartRPC() {
	if o.DisableRPC {
		log.Info("RPC is disabled by command-line flag")
		return
	}
	if !o.inst.Config().RPC.Enabled || o.inst.Config().RPC.Port == 0 {
		return
	}

	var ctx context.Context
	ctx, o.stopRPC = context.WithCancel(o.inst.Context())

	go func() {
		if err := lib.ServeRPC(ctx, o.inst); err != nil {
			log.Errorf("RPC closed: %s", err)
		}
	}()
}

func (o *ConnectOptions) maybeStartAPI() {
	if o.DisableAPI {
		log.Info("API is disabled by command-line flag")
		return
	}
	if !o.inst.Config().API.Enabled || o.inst.Config().API.Port == 0 {
		return
	}

	var ctx context.Context
	ctx, o.stopAPI = context.WithCancel(o.inst.Context())

	go func() {
		s := api.New(o.inst)
		if err := s.ServeHTTP(ctx); err != nil && err.Error() != "http: Server closed" {
			log.Errorf("API closed: %s", err)
		}
	}()
}

func (o *ConnectOptions) maybeStartP2P() {
	if o.DisableP2P {
		log.Info("P2P is disabled by command-line flag")
		return
	}
	if !o.inst.Config().P2P.Enabled {
		return
	}
	log.Error("connecting P2P")

	node := o.inst.Node()
	cfg := o.inst.Config()

	var ctx context.Context
	ctx, o.stopP2P = context.WithCancel(o.inst.Context())

	if err := node.GoOnline(ctx); err != nil {
		fmt.Println("serving error", err)
		return
	}

	info := "\n📡  Success! You are now connected to the d.web. Here's your connection details:\n"
	info += cfg.SummaryString()
	info += "IPFS Addresses:"
	for _, a := range node.EncapsulatedAddresses() {
		info = fmt.Sprintf("%s\n  %s", info, a.String())
	}
	info += "\n\n"

	node.LocalStreams.Print(info)
}

func (o *ConnectOptions) maybeStartCron() {
	// TODO (b5): implement cron
}

func (o *ConnectOptions) maybeTimeoutProcess() {
	if o.DisconnectAfter != 0 {
		log.Infof("disconnecting after %d seconds", o.DisconnectAfter)
		go func(inst *lib.Instance, t int) {
			<-time.After(time.Second * time.Duration(t))
			log.Infof("disconnecting")
			o.inst.Teardown()
		}(o.inst, o.DisconnectAfter)
	}
}
