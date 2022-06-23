package p2p

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/protocol"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"

	ma "github.com/multiformats/go-multiaddr"
)

const (
	swarmKey  = "/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7"
	test_peer = "12D3KooWHPbFSqWiKgh1QzuX64otKZNfYuUu1cYRmfCWnxEqjb5k"
)

func newTestP2pClient() (*P2pClient, error) {
	identity, err := config.CreateIdentity()
	if err != nil {
		return nil, err
	}
	p2pClient, err := NewP2pClient(4001, identity.PrivKey, swarmKey, DEFAULT_IPFS_PEERS)
	if err != nil {
		return nil, err
	}
	return p2pClient, nil
}

func TestList(t *testing.T) {

	p2pClient, err := newTestP2pClient()
	assert.NoError(t, err)

	output := p2pClient.List()
	assert.Empty(t, output.Listeners)

	targetOpt := "/ip4/127.0.0.1/tcp/8080"
	err = p2pClient.Listen("/x/ssh", targetOpt)
	output = p2pClient.List()
	assert.NotEmpty(t, output.Listeners)
	assert.Equal(t, 1, len(output.Listeners))
	assert.Equal(t, targetOpt, output.Listeners[0].TargetAddress)
}

func TestListen(t *testing.T) {
	p2pClient, err := newTestP2pClient()
	assert.NoError(t, err)

	targetOpt := "/ip4/127.0.0.1/tcp/8080"
	err = p2pClient.Listen("/x/ssh", targetOpt)
	assert.NoError(t, err)
	output := p2pClient.List()
	listen := output.Listeners[0]

	fmt.Println(listen.ListenAddress)

	checkClient, err := newTestP2pClient()
	assert.NoError(t, err)
	err = checkClient.Forward("/x/ssh", 8081, listen.ListenAddress)
	assert.NoError(t, err)
	assert.NotEmpty(t, checkClient.List().Listeners)
	assert.Equal(t, listen.ListenAddress, checkClient.List().Listeners[0].TargetAddress)

}

func TestForward(t *testing.T) {
	p2pClient, err := newTestP2pClient()
	assert.NoError(t, err)

	targetOpt := fmt.Sprintf("/p2p/%s", test_peer)

	err = p2pClient.Forward("/x/ssh", 8081, targetOpt)
	assert.NotEmpty(t, p2pClient.List().Listeners)

}

func TestForwardHealth(t *testing.T) {
	p2pClient, err := newTestP2pClient()
	assert.NoError(t, err)
	err = p2pClient.CheckForwardHealth("/x/ssh", test_peer)
	assert.NoError(t, err)
}

func TestClose(t *testing.T) {
	p2pClient, err := newTestP2pClient()
	assert.NoError(t, err)
	targetOpt := fmt.Sprintf("/p2p/%s", test_peer)
	err = p2pClient.Forward("/x/ssh", 8081, targetOpt)
	assert.NoError(t, err)

	count, err := p2pClient.Close(targetOpt)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	data := p2pClient.List()
	assert.Equal(t, 0, len(data.Listeners))
}

func TestDestroy(t *testing.T) {
	p2pClient, err := newTestP2pClient()
	assert.NoError(t, err)

	err = p2pClient.Destroy()
	assert.NoError(t, err)

	address := net.JoinHostPort("127.0.0.1", "4001")
	// 3 second timeout
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)

	assert.Error(t, err)
	assert.Empty(t, conn)

}

func TestForwardLocal(t *testing.T) {

	target := "12D3KooWHPbFSqWiKgh1QzuX64otKZNfYuUu1cYRmfCWnxEqjb5k"
	ctx := context.Background()

	r := rand.Reader

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	assert.NoError(t, err)
	skbytes, _ := crypto.MarshalPrivateKey(priv)

	ha, _, _, err := newRoutedHost(4001, base64.StdEncoding.EncodeToString(skbytes), []byte(swarmKey), DEFAULT_IPFS_PEERS)
	assert.NoError(t, err)

	protoOpt := "/x/ssh"
	listenOpt := "/ip4/127.0.0.1/tcp/2222"
	targetOpt := fmt.Sprintf("/p2p/%s", target)
	listen, err := ma.NewMultiaddr(listenOpt)

	assert.NoError(t, err)

	targets, err := parseIpfsAddr(targetOpt)
	assert.NoError(t, err)
	proto := protocol.ID(protoOpt)
	var pa = newIpfsP2p(ha)
	err = forwardLocal(ctx, pa, (*ha).Peerstore(), proto, listen, targets)
	assert.NoError(t, err)

	check_connect_avaliable("127.0.0.1", []string{"2222"}, t)

}

func check_connect_avaliable(host string, ports []string, t *testing.T) {
	for _, port := range ports {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			fmt.Println("Connecting error:", err)
			assert.NoError(t, err)
		}
		if conn != nil {
			defer conn.Close()
			fmt.Println("Opened", net.JoinHostPort(host, port))
		}
	}
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
		assert.NoError(t, err)
	}
	log.Println("listening for connections")

	pa := newIpfsP2p(ha)
	protoOpt := "/x/ssh"
	targetOpt := "/ip4/127.0.0.1/tcp/22"
	proto := protocol.ID(protoOpt)

	target, err := ma.NewMultiaddr(targetOpt)
	if err != nil {
		log.Error(err)
		assert.NoError(t, err)
	}
	_, err = pa.ForwardRemote(ctx, proto, target, false)
	assert.NoError(t, err)

	// check p2p port is recorded
	fmt.Println(pa.ListenersP2P.Listeners)
	assert.NotEmpty(t, pa.ListenersP2P.Listeners[protoOpt])

}
