package p2p

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	ipfsp2p "github.com/ipfs/go-ipfs/p2p"
	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/pnet"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

var resolveTimeout = 10 * time.Second

// NewRoutedHost create a p2p routing client
func newRoutedHost(listenPort int, privstr string, swarmkey []byte, peers []string) (*host.Host, *rhost.RoutedHost, *dht.IpfsDHT, error) {
	ctx := context.Background()

	skbytes, err := base64.StdEncoding.DecodeString(privstr)
	if err != nil {
		logrus.Error(err)
		return nil, nil, nil, err
	}
	priv, err := crypto.UnmarshalPrivateKey(skbytes)
	if err != nil {
		logrus.Error(err)
		return nil, nil, nil, err
	}
	bootstrapPeers := convertPeers(peers)

	// load private key swarm.key

	psk, err := pnet.DecodeV1PSK(bytes.NewReader(swarmkey))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to configure private network: %s", err)
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.PrivateNetwork(psk),
		libp2p.ConnectionManager(connmgr.NewConnManager(
			100,         // Lowwater
			400,         // HighWater,
			time.Minute, // GracePeriod
		)),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err := dht.New(ctx, h)
			return idht, err
		}),
		libp2p.EnableAutoRelay(),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
	}

	basicHost, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, nil, nil, err
	}

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	// Make the DHT
	DHT := dht.NewDHT(ctx, basicHost, dstore)

	// Make the routed host
	routedHost := rhost.Wrap(basicHost, DHT)

	cfg := DefaultBootstrapConfig
	cfg.BootstrapPeers = func() []peer.AddrInfo {
		return bootstrapPeers
	}

	id, err := peer.IDFromPrivateKey(priv)
	_, err = Bootstrap(id, routedHost, DHT, cfg)

	// connect to the chosen ipfs nodes
	if err != nil {
		return nil, nil, nil, err
	}

	// Bootstrap the host
	err = DHT.Bootstrap(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	addrs := routedHost.Addrs()
	log.Println("I can be reached at:")
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

	return &basicHost, routedHost, DHT, nil
}

// MakeIpfsP2p create ipfs p2p object
func newIpfsP2p(h *host.Host) *ipfsp2p.P2P {
	return ipfsp2p.New((*h).ID(), *h, (*h).Peerstore())
}

// P2pClient p2p client
type P2pClient struct {
	Host       *host.Host
	P2P        *ipfsp2p.P2P
	DHT        *dht.IpfsDHT
	RoutedHost *rhost.RoutedHost
}

func NewP2pClient(listenPort int, privstr string, swarmkey string, peers []string) (*P2pClient, error) {
	host, routedHost, DHT, err := newRoutedHost(listenPort, privstr, []byte(swarmkey), peers)
	if err != nil {
		return nil, err
	}
	P2P := newIpfsP2p(host)
	return &P2pClient{
		Host:       host,
		P2P:        P2P,
		DHT:        DHT,
		RoutedHost: routedHost,
	}, nil
}

// P2PListenerInfoOutput  p2p monitoring or mapping information
type P2PListenerInfoOutput struct {
	Protocol      string
	ListenAddress string
	TargetAddress string
}

// P2PLsOutput p2p monitor or map information output
type P2PLsOutput struct {
	Listeners []P2PListenerInfoOutput
}

// List p2p monitor message list
func (c *P2pClient) List() *P2PLsOutput {
	output := &P2PLsOutput{}

	c.P2P.ListenersLocal.Lock()
	for _, listener := range c.P2P.ListenersLocal.Listeners {
		output.Listeners = append(output.Listeners, P2PListenerInfoOutput{
			Protocol:      string(listener.Protocol()),
			ListenAddress: listener.ListenAddress().String(),
			TargetAddress: listener.TargetAddress().String(),
		})
	}
	c.P2P.ListenersLocal.Unlock()

	c.P2P.ListenersP2P.Lock()
	for _, listener := range c.P2P.ListenersP2P.Listeners {
		output.Listeners = append(output.Listeners, P2PListenerInfoOutput{
			Protocol:      string(listener.Protocol()),
			ListenAddress: listener.ListenAddress().String(),
			TargetAddress: listener.TargetAddress().String(),
		})
	}
	c.P2P.ListenersP2P.Unlock()

	return output
}

// Listen map local ports to p2p networks
func (c *P2pClient) Listen(proto, targetOpt string) error {
	log.Println("listening for connections")

	//targetOpt := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)
	protoId := protocol.ID(proto)

	target, err := ma.NewMultiaddr(targetOpt)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = c.P2P.ForwardRemote(context.Background(), protoId, target, false)
	logrus.Info("local port" + targetOpt + ",mapping to p2p network succeeded")
	return err
}

// Forward map p2p network remote nodes to local ports
func (c *P2pClient) Forward(proto string, port int, targetOpt string) error {
	listenOpt := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)
	listen, err := ma.NewMultiaddr(listenOpt)

	if err != nil {
		log.Println(err)
		return err
	}

	targets, err := parseIpfsAddr(targetOpt)
	protoId := protocol.ID(proto)

	err = forwardLocal(context.Background(), c.P2P, (*c.Host).Peerstore(), protoId, listen, targets)
	if err != nil {
		log.Println(err)
		return err
	}
	logrus.Info("remote node" + targetOpt + ",forward to " + listenOpt + "success")
	return err
}

// CheckForwardHealth check if the remote node is connected
func (c *P2pClient) CheckForwardHealth(proto, peerId string) error {
	targetOpt := fmt.Sprintf("/p2p/%s", peerId)
	targets, err := parseIpfsAddr(targetOpt)
	protoId := protocol.ID(proto)
	if err != nil {
		return err
	}
	cctx, cancel := context.WithTimeout(context.Background(), time.Second*30) //TODO: configurable?
	defer cancel()
	stream, err := (*c.Host).NewStream(cctx, targets.ID, protoId)
	if err != nil {
		return err
	} else {
		stream.Close()
		return nil
	}
}

// Close turn off p2p listening connection
func (c *P2pClient) Close(target string) (int, error) {
	targetAddress, err := ma.NewMultiaddr(target)
	if err != nil {
		return 0, err
	}
	match := func(listener ipfsp2p.Listener) bool {

		if !targetAddress.Equal(listener.TargetAddress()) {
			return false
		}
		return true
	}

	done := c.P2P.ListenersLocal.Close(match)
	done += c.P2P.ListenersP2P.Close(match)

	return done, nil

}

// Destroy: destroy and close the p2p client, including all subordinate listeners, stream objects
func (c *P2pClient) Destroy() error {
	for _, stream := range c.P2P.Streams.Streams {
		c.P2P.Streams.Close(stream)
	}
	match := func(listener ipfsp2p.Listener) bool {
		return true
	}
	c.P2P.ListenersP2P.Close(match)
	c.P2P.ListenersLocal.Close(match)
	err := (*c.Host).Close()
	c.P2P = nil
	c.Host = nil
	return err
}

// forwardLocal forwards local connections to a libp2p service
func forwardLocal(ctx context.Context, p *ipfsp2p.P2P, ps pstore.Peerstore, proto protocol.ID, bindAddr ma.Multiaddr, addr *peer.AddrInfo) error {

	ps.AddAddrs(addr.ID, addr.Addrs, pstore.TempAddrTTL)
	// TODO: return some info
	_, err := p.ForwardLocal(ctx, addr.ID, proto, bindAddr)
	return err
}

// parseIpfsAddr is a function that takes in addr string and return ipfsAddrs
func parseIpfsAddr(addr string) (*peer.AddrInfo, error) {
	multiaddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	pi, err := peer.AddrInfoFromP2pAddr(multiaddr)
	if err == nil {
		return pi, nil
	}

	// resolve multiaddr whose protocol is not ma.P_IPFS
	ctx, cancel := context.WithTimeout(context.Background(), resolveTimeout)
	defer cancel()
	addrs, err := madns.Resolve(ctx, multiaddr)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, errors.New("fail to resolve the multiaddr:" + multiaddr.String())
	}
	var info peer.AddrInfo
	for _, addr := range addrs {
		taddr, id := peer.SplitAddr(addr)
		if id == "" {
			// not an ipfs addr, skipping.
			continue
		}
		switch info.ID {
		case "":
			info.ID = id
		case id:
		default:
			return nil, fmt.Errorf(
				"ambiguous multiaddr %s could refer to %s or %s",
				multiaddr,
				info.ID,
				id,
			)
		}
		info.Addrs = append(info.Addrs, taddr)
	}
	return &info, nil
}
