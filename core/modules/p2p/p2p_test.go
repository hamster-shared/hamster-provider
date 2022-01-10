package p2p

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/protocol"
	"log"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestLocalForward(t *testing.T) {

	target := "QmUhYMMxwpS6CbZ6x95sC2pdTT1jriPvbJGjW3YM6MrDmC"
	ctx := context.Background()

	r := rand.Reader

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	skbytes, _ := crypto.MarshalPrivateKey(priv)
	swarmKey := []byte("/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7")
	ha, _, _, err := newRoutedHost(4001, base64.StdEncoding.EncodeToString(skbytes), swarmKey, DEFAULT_IPFS_PEERS)
	if err != nil {
		log.Fatal(err)
	}

	protoOpt := "/x/ssh"
	listenOpt := "/ip4/127.0.0.1/tcp/2222"
	targetOpt := fmt.Sprintf("/p2p/%s", target)
	listen, err := ma.NewMultiaddr(listenOpt)

	if err != nil {
		log.Fatalln(err)
	}

	targets, err := parseIpfsAddr(targetOpt)
	proto := protocol.ID(protoOpt)

	var pa = newIpfsP2p(ha)
	err = forwardLocal(ctx, pa, (*ha).Peerstore(), proto, listen, targets)
	if err != nil {
		log.Fatalln(err)
	}
	//select {}
}

func TestListenLocal(t *testing.T) {
	// Make a host that listens on the given multiaddress
	log.Println("using global bootstrap")
	ctx := context.Background()

	r := rand.Reader

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	skbytes, _ := crypto.MarshalPrivateKey(priv)
	swarmKey := []byte("/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7")
	ha, _, _, err := newRoutedHost(4001, base64.StdEncoding.EncodeToString(skbytes), swarmKey, DEFAULT_IPFS_PEERS)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening for connections")

	pa := newIpfsP2p(ha)
	protoOpt := "/x/ssh"
	targetOpt := "/ip4/127.0.0.1/tcp/22"
	proto := protocol.ID(protoOpt)

	target, err := ma.NewMultiaddr(targetOpt)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = pa.ForwardRemote(ctx, proto, target, false)
	//select {} // hang forever
}
