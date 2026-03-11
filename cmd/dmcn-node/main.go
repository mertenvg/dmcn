// Binary dmcn-node is the combined DMCN development node that runs a
// DHT registry and relay service in a single process.
//
// See PRD Section 5.3.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mertenvg/dmcn/internal/core/identity"
	"github.com/mertenvg/dmcn/internal/core/message"
	"github.com/mertenvg/dmcn/internal/keystore"
	"github.com/mertenvg/dmcn/internal/node"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "start":
		err = cmdStart(args)
	case "identity":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "usage: dmcn-node identity <generate|register|lookup>")
			os.Exit(1)
		}
		subcmd := args[0]
		subargs := args[1:]
		switch subcmd {
		case "generate":
			err = cmdIdentityGenerate(subargs)
		case "register":
			err = cmdIdentityRegister(subargs)
		case "lookup":
			err = cmdIdentityLookup(subargs)
		default:
			fmt.Fprintf(os.Stderr, "unknown identity subcommand: %s\n", subcmd)
			os.Exit(1)
		}
	case "message":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "usage: dmcn-node message <send|fetch>")
			os.Exit(1)
		}
		subcmd := args[0]
		subargs := args[1:]
		switch subcmd {
		case "send":
			err = cmdMessageSend(subargs)
		case "fetch":
			err = cmdMessageFetch(subargs)
		default:
			fmt.Fprintf(os.Stderr, "unknown message subcommand: %s\n", subcmd)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `usage: dmcn-node <command> [options]

commands:
  start                          Start a DMCN node
  identity generate              Generate a new identity key pair
  identity register              Register an identity in the DHT
  identity lookup                Look up an identity by address
  message send                   Send a message
  message fetch                  Fetch pending messages`)
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
	listen := parseFlag(args, "listen")
	if listen == "" {
		listen = "/ip4/0.0.0.0/tcp/7400"
	}
	keystorePath := parseFlag(args, "keystore")
	if keystorePath == "" {
		keystorePath = "keystore.enc"
	}
	passphrase := parseFlag(args, "passphrase")
	if passphrase == "" {
		passphrase = "dmcn-dev" // default dev passphrase
	}
	bootstrap := parseFlag(args, "bootstrap")

	cfg := node.Config{
		ListenAddr:   listen,
		KeystorePath: keystorePath,
		Passphrase:   passphrase,
	}
	if bootstrap != "" {
		cfg.BootstrapPeers = strings.Split(bootstrap, ",")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n, err := node.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer n.Close()

	addrs := n.Addrs()
	fmt.Printf("DMCN node started\n")
	fmt.Printf("Peer ID: %s\n", n.PeerID())
	for _, addr := range addrs {
		fmt.Printf("Listening: %s\n", addr)
	}

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\nShutting down...")

	return nil
}

func cmdIdentityGenerate(args []string) error {
	address := parseFlag(args, "address")
	if address == "" {
		return fmt.Errorf("--address is required")
	}
	keystorePath := parseFlag(args, "keystore")
	if keystorePath == "" {
		keystorePath = "keystore.enc"
	}
	passphrase := parseFlag(args, "passphrase")
	if passphrase == "" {
		passphrase = "dmcn-dev"
	}

	kp, err := identity.GenerateIdentityKeyPair()
	if err != nil {
		return fmt.Errorf("generate key pair: %w", err)
	}

	ks := newKeystore(keystorePath, passphrase)
	if err := ks.Store(address, kp); err != nil {
		return fmt.Errorf("store key pair: %w", err)
	}

	rec, err := identity.NewIdentityRecord(address, kp)
	if err != nil {
		return fmt.Errorf("create identity record: %w", err)
	}

	fmt.Printf("Identity generated for %s\n", address)
	fmt.Printf("Fingerprint: %s\n", rec.Fingerprint())
	fmt.Printf("Stored in: %s\n", keystorePath)

	return nil
}

func cmdIdentityRegister(args []string) error {
	address := parseFlag(args, "address")
	if address == "" {
		return fmt.Errorf("--address is required")
	}
	nodeAddr := parseFlag(args, "node")
	if nodeAddr == "" {
		return fmt.Errorf("--node is required (multiaddr of running node)")
	}
	keystorePath := parseFlag(args, "keystore")
	if keystorePath == "" {
		keystorePath = "keystore.enc"
	}
	passphrase := parseFlag(args, "passphrase")
	if passphrase == "" {
		passphrase = "dmcn-dev"
	}

	ctx := context.Background()

	ks := newKeystore(keystorePath, passphrase)
	kp, err := ks.Load(address)
	if err != nil {
		return fmt.Errorf("load key pair: %w", err)
	}

	// Create a temporary node to connect and register
	n, err := node.New(ctx, node.Config{
		ListenAddr:     "/ip4/127.0.0.1/tcp/0",
		BootstrapPeers: []string{nodeAddr},
	})
	if err != nil {
		return fmt.Errorf("create node: %w", err)
	}
	defer n.Close()

	rec, err := identity.NewIdentityRecord(address, kp)
	if err != nil {
		return fmt.Errorf("create identity record: %w", err)
	}
	if err := rec.Sign(kp); err != nil {
		return fmt.Errorf("sign identity: %w", err)
	}
	if err := n.Registry().Register(ctx, rec); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	fmt.Printf("Identity registered: %s\n", address)
	return nil
}

func cmdIdentityLookup(args []string) error {
	address := parseFlag(args, "address")
	if address == "" {
		return fmt.Errorf("--address is required")
	}
	nodeAddr := parseFlag(args, "node")
	if nodeAddr == "" {
		return fmt.Errorf("--node is required (multiaddr of running node)")
	}

	ctx := context.Background()

	n, err := node.New(ctx, node.Config{
		ListenAddr:     "/ip4/127.0.0.1/tcp/0",
		BootstrapPeers: []string{nodeAddr},
	})
	if err != nil {
		return fmt.Errorf("create node: %w", err)
	}
	defer n.Close()

	rec, err := n.Registry().Lookup(ctx, address)
	if err != nil {
		return fmt.Errorf("lookup: %w", err)
	}

	fmt.Printf("Address: %s\n", rec.Address)
	fmt.Printf("Fingerprint: %s\n", rec.Fingerprint())
	fmt.Printf("Created: %s\n", rec.CreatedAt)
	if !rec.ExpiresAt.IsZero() {
		fmt.Printf("Expires: %s\n", rec.ExpiresAt)
	}
	if err := rec.Verify(); err != nil {
		fmt.Printf("Signature: INVALID (%v)\n", err)
	} else {
		fmt.Printf("Signature: valid\n")
	}

	return nil
}

func cmdMessageSend(args []string) error {
	from := parseFlag(args, "from")
	to := parseFlag(args, "to")
	body := parseFlag(args, "body")
	nodeAddr := parseFlag(args, "node")
	keystorePath := parseFlag(args, "keystore")
	if keystorePath == "" {
		keystorePath = "keystore.enc"
	}
	passphrase := parseFlag(args, "passphrase")
	if passphrase == "" {
		passphrase = "dmcn-dev"
	}

	if from == "" || to == "" || body == "" || nodeAddr == "" {
		return fmt.Errorf("--from, --to, --body, and --node are required")
	}

	ctx := context.Background()

	ks := newKeystore(keystorePath, passphrase)
	senderKP, err := ks.Load(from)
	if err != nil {
		return fmt.Errorf("load sender key: %w", err)
	}

	n, err := node.New(ctx, node.Config{
		ListenAddr:     "/ip4/127.0.0.1/tcp/0",
		BootstrapPeers: []string{nodeAddr},
	})
	if err != nil {
		return fmt.Errorf("create node: %w", err)
	}
	defer n.Close()

	// Lookup recipient
	recipientRec, err := n.Registry().Lookup(ctx, to)
	if err != nil {
		return fmt.Errorf("lookup recipient: %w", err)
	}

	// Compose message
	msg, err := message.NewPlaintextMessage(from, to, "", body, senderKP.Ed25519Public)
	if err != nil {
		return fmt.Errorf("compose message: %w", err)
	}

	sm := &message.SignedMessage{Plaintext: *msg}
	if err := sm.Sign(senderKP.Ed25519Private); err != nil {
		return fmt.Errorf("sign message: %w", err)
	}

	recipients := []message.RecipientInfo{{
		DeviceID:  senderKP.DeviceID, // using sender's device ID as placeholder
		X25519Pub: recipientRec.X25519Public,
	}}

	env, err := message.Encrypt(sm, recipients)
	if err != nil {
		return fmt.Errorf("encrypt message: %w", err)
	}

	// Store on relay (connect to the node and send)
	hash, err := n.Relay().ClientStoreWithAddress(ctx, n.PeerID(), from, senderKP, env)
	if err != nil {
		return fmt.Errorf("store message: %w", err)
	}

	fmt.Printf("Message sent to %s\n", to)
	fmt.Printf("Envelope hash: %x\n", hash)

	return nil
}

func cmdMessageFetch(args []string) error {
	address := parseFlag(args, "address")
	nodeAddr := parseFlag(args, "node")
	keystorePath := parseFlag(args, "keystore")
	if keystorePath == "" {
		keystorePath = "keystore.enc"
	}
	passphrase := parseFlag(args, "passphrase")
	if passphrase == "" {
		passphrase = "dmcn-dev"
	}

	if address == "" || nodeAddr == "" {
		return fmt.Errorf("--address and --node are required")
	}

	ctx := context.Background()

	ks := newKeystore(keystorePath, passphrase)
	kp, err := ks.Load(address)
	if err != nil {
		return fmt.Errorf("load key: %w", err)
	}

	n, err := node.New(ctx, node.Config{
		ListenAddr:     "/ip4/127.0.0.1/tcp/0",
		BootstrapPeers: []string{nodeAddr},
	})
	if err != nil {
		return fmt.Errorf("create node: %w", err)
	}
	defer n.Close()

	envs, hashes, err := n.Relay().ClientFetch(ctx, n.PeerID(), kp, address)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	if len(envs) == 0 {
		fmt.Println("No pending messages")
		return nil
	}

	fmt.Printf("Fetched %d message(s)\n", len(envs))
	for i, env := range envs {
		sm, err := message.Decrypt(env, kp.X25519Private, kp.X25519Public)
		if err != nil {
			fmt.Printf("  [%d] Failed to decrypt: %v\n", i+1, err)
			continue
		}
		if err := sm.Verify(); err != nil {
			fmt.Printf("  [%d] Invalid signature: %v\n", i+1, err)
			continue
		}
		fmt.Printf("  [%d] From: %s\n", i+1, sm.Plaintext.SenderAddress)
		fmt.Printf("      Body: %s\n", string(sm.Plaintext.Body.Content))
		fmt.Printf("      Hash: %x\n", hashes[i])
	}

	return nil
}

func newKeystore(path, passphrase string) *keystore.Keystore {
	return keystore.New(path, passphrase)
}
