// Binary dmcn-bridge is the SMTP-DMCN bridge node that allows legacy
// email clients to exchange messages with DMCN users.
//
// See PRD Section 6.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mertenvg/logr/v2"

	"github.com/mertenvg/dmcn/internal/bridge"
)

var log logr.Logger

func main() {
	logr.AddWriter(os.Stderr, logr.WithFormatter(logr.FormatWithColours), logr.WithFilter(logr.Verbose))
	log = logr.With(logr.M("component", "bridge-cli"))

	if len(os.Args) < 2 || os.Args[1] != "start" {
		printUsage()
		os.Exit(1)
	}

	args := os.Args[2:]

	if err := cmdStart(args); err != nil {
		log.Errorf("%v", err)
		logr.Wait()
		os.Exit(1)
	}

	logr.Wait()
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `usage: dmcn-bridge start [options]

options:
  --node <multiaddr>         Multiaddr of running dmcn-node (required)
  --smtp-listen <addr>       SMTP listen address (default: 0.0.0.0:2525)
  --libp2p-listen <multiaddr> libp2p listen address (default: /ip4/127.0.0.1/tcp/0)
  --bridge-domain <domain>   Bridge email domain (default: bridge.localhost)
  --dmcn-domain <domain>     DMCN address domain (default: dmcn.localhost)
  --bridge-address <addr>    Bridge's own DMCN address (default: bridge@bridge.localhost)
  --keystore <path>          Keystore file path (default: bridge-keystore.enc)
  --passphrase <pass>        Keystore passphrase (default: dmcn-bridge-dev)
  --poll-interval <duration> Relay poll interval (default: 5s)
  --smtp-relay <host:port>   Outbound SMTP relay (default: stub mode)`)
}

func parseFlag(args []string, name string) string {
	flag := "--" + name
	for i, arg := range args {
		if arg == flag && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(arg, flag+"=") {
			return arg[len(flag)+1:]
		}
	}
	return ""
}

func cmdStart(args []string) error {
	nodeAddr := parseFlag(args, "node")
	if nodeAddr == "" {
		return fmt.Errorf("--node is required (multiaddr of running dmcn-node)")
	}

	smtpListen := parseFlag(args, "smtp-listen")
	if smtpListen == "" {
		smtpListen = "0.0.0.0:2525"
	}
	libp2pListen := parseFlag(args, "libp2p-listen")
	if libp2pListen == "" {
		libp2pListen = "/ip4/127.0.0.1/tcp/0"
	}
	bridgeDomain := parseFlag(args, "bridge-domain")
	if bridgeDomain == "" {
		bridgeDomain = "bridge.localhost"
	}
	dmcnDomain := parseFlag(args, "dmcn-domain")
	if dmcnDomain == "" {
		dmcnDomain = "dmcn.localhost"
	}
	bridgeAddress := parseFlag(args, "bridge-address")
	if bridgeAddress == "" {
		bridgeAddress = "bridge@bridge.localhost"
	}
	keystorePath := parseFlag(args, "keystore")
	if keystorePath == "" {
		keystorePath = "bridge-keystore.enc"
	}
	passphrase := parseFlag(args, "passphrase")
	if passphrase == "" {
		passphrase = "dmcn-bridge-dev"
	}
	pollStr := parseFlag(args, "poll-interval")
	pollInterval := 5 * time.Second
	if pollStr != "" {
		d, err := time.ParseDuration(pollStr)
		if err != nil {
			return fmt.Errorf("invalid --poll-interval: %w", err)
		}
		pollInterval = d
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := bridge.New(ctx, bridge.Config{
		NodeAddr:       nodeAddr,
		SMTPListenAddr: smtpListen,
		LibP2PAddr:     libp2pListen,
		BridgeDomain:   bridgeDomain,
		DMCNDomain:     dmcnDomain,
		BridgeAddress:  bridgeAddress,
		KeystorePath:   keystorePath,
		Passphrase:     passphrase,
		PollInterval:   pollInterval,
	})
	if err != nil {
		return err
	}
	defer b.Stop()

	if err := b.Start(); err != nil {
		return err
	}

	log.Infof("DMCN bridge started")
	log.Infof("SMTP listening on %s", b.SMTPAddr())
	log.Infof("bridge domain: %s", bridgeDomain)
	log.Infof("DMCN domain: %s", dmcnDomain)

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Info("shutting down...")

	return nil
}
