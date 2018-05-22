// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package swarm

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"
	"unicode"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/chequebook"
	"github.com/ethereum/go-ethereum/contracts/ens"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/swarm/api"
	httpapi "github.com/ethereum/go-ethereum/swarm/api/http"
	"github.com/ethereum/go-ethereum/swarm/fuse"
	"github.com/ethereum/go-ethereum/swarm/network"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

var (
	startTime          time.Time
	updateGaugesPeriod = 5 * time.Second
	startCounter       = metrics.NewRegisteredCounter("stack,start", nil)
	stopCounter        = metrics.NewRegisteredCounter("stack,stop", nil)
	uptimeGauge        = metrics.NewRegisteredGauge("stack.uptime", nil)
	dbSizeGauge        = metrics.NewRegisteredGauge("storage.db.chunks.size", nil)
	cacheSizeGauge     = metrics.NewRegisteredGauge("storage.db.cache.size", nil)
)

// Swarm stack
type Swarm struct {
	config      *api.Config            // swarm configuration
	api         *api.Api               // high level api layer (fs/manifest)
	dns         api.Resolver           // DNS registrar
	dbAccess    *network.DbAccess      // access to local chunk db iterator and storage counter
	storage     storage.ChunkStore     // internal access to storage, common interface to cloud storage backends
	dpa         *storage.DPA           // distributed preimage archive, the local API to the storage with document level storage/retrieval support
	depo        network.StorageHandler // remote request handler, interface between bzz protocol and the storage
	cloud       storage.CloudStore     // procurement, cloud storage backend (can multi-cloud)
	hive        *network.Hive          // the logistic manager
	backend     chequebook.Backend     // simple blockchain Backend
	privateKey  *ecdsa.PrivateKey
	corsString  string
	swapEnabled bool
	lstore      *storage.LocalStore // local store, needs to store for releasing resources after node stopped
	sfs         *fuse.SwarmFS       // need this to cleanup all the active mounts on node exit
}

type SwarmAPI struct {
	Api     *api.Api
	Backend chequebook.Backend
	PrvKey  *ecdsa.PrivateKey
}

func (sw *Swarm) API() *SwarmAPI {
	return &SwarmAPI{
		Api:     sw.api,
		Backend: sw.backend,
		PrvKey:  sw.privateKey,
	}
}

// NewSwarm creates a new swarm service instance
// implements node.Service
func NewSwarm(ctx *node.ServiceContext, backend chequebook.Backend, config *api.Config) (sw *Swarm, err error) {
	if bytes.Equal(common.FromHex(config.PublicKey), storage.ZeroKey) {
		return nil, fmt.Errorf("empty public key")
	}
	if bytes.Equal(common.FromHex(config.BzzKey), storage.ZeroKey) {
		return nil, fmt.Errorf("empty bzz key")
	}

	sw = &Swarm{
		config:      config,
		swapEnabled: config.SwapEnabled,
		backend:     backend,
		privateKey:  config.Swap.PrivateKey(),
		corsString:  config.Cors,
	}
	log.Debug(fmt.Sprintf("Setting up Swarm service components"))

	hash := storage.MakeHashFunc(config.ChunkerParams.Hash)
	sw.lstore, err = storage.NewLocalStore(hash, config.StoreParams)
	if err != nil {
		return
	}

	// setup local store
	log.Debug(fmt.Sprintf("Set up local storage"))

	sw.dbAccess = network.NewDbAccess(sw.lstore)
	log.Debug(fmt.Sprintf("Set up local db access (iterator/counter)"))

	// set up the kademlia hive
	sw.hive = network.NewHive(
		common.HexToHash(sw.config.BzzKey), // key to hive (kademlia base address)
		config.HiveParams,                    // configuration parameters
		config.SwapEnabled,                   // SWAP enabled
		config.SyncEnabled,                   // syncronisation enabled
	)
	log.Debug(fmt.Sprintf("Set up swarm network with Kademlia hive"))

	// setup cloud storage backend
	sw.cloud = network.NewForwarder(sw.hive)
	log.Debug(fmt.Sprintf("-> set swarm forwarder as cloud storage backend"))

	// setup cloud storage internal access layer
	sw.storage = storage.NewNetStore(hash, sw.lstore, sw.cloud, config.StoreParams)
	log.Debug(fmt.Sprintf("-> swarm net store shared access layer to Swarm Chunk Store"))

	// set up Depo (storage handler = cloud storage access layer for incoming remote requests)
	sw.depo = network.NewDepo(hash, sw.lstore, sw.storage)
	log.Debug(fmt.Sprintf("-> REmote Access to CHunks"))

	// set up DPA, the cloud storage local access layer
	dpaChunkStore := storage.NewDpaChunkStore(sw.lstore, sw.storage)
	log.Debug(fmt.Sprintf("-> Local Access to Swarm"))
	// Swarm Hash Merklised Chunking for Arbitrary-length Document/File storage
	sw.dpa = storage.NewDPA(dpaChunkStore, sw.config.ChunkerParams)
	log.Debug(fmt.Sprintf("-> Content Store API"))

	if len(config.EnsAPIs) > 0 {
		opts := []api.MultiResolverOption{}
		for _, c := range config.EnsAPIs {
			tld, endpoint, addr := parseEnsAPIAddress(c)
			r, err := newEnsClient(endpoint, addr, config)
			if err != nil {
				return nil, err
			}
			opts = append(opts, api.MultiResolverOptionWithResolver(r, tld))
		}
		sw.dns = api.NewMultiResolver(opts...)
	}

	sw.api = api.NewApi(sw.dpa, sw.dns)
	// Manifests for Smart Hosting
	log.Debug(fmt.Sprintf("-> Web3 virtual server API"))

	sw.sfs = fuse.NewSwarmFS(sw.api)
	log.Debug("-> Initializing Fuse file system")

	return sw, nil
}

// parseEnsAPIAddress parses string according to format
// [tld:][contract-addr@]url and returns ENSClientConfig structure
// with endpoint, contract address and TLD.
func parseEnsAPIAddress(s string) (tld, endpoint string, addr common.Address) {
	isAllLetterString := func(s string) bool {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return false
			}
		}
		return true
	}
	endpoint = s
	if i := strings.Index(endpoint, ":"); i > 0 {
		if isAllLetterString(endpoint[:i]) && len(endpoint) > i+2 && endpoint[i+1:i+3] != "//" {
			tld = endpoint[:i]
			endpoint = endpoint[i+1:]
		}
	}
	if i := strings.Index(endpoint, "@"); i > 0 {
		addr = common.HexToAddress(endpoint[:i])
		endpoint = endpoint[i+1:]
	}
	return
}

// newEnsClient creates a new ENS client for that is a consumer of
// a ENS API on a specific endpoint. It is used as a helper function
// for creating multiple resolvers in NewSwarm function.
func newEnsClient(endpoint string, addr common.Address, config *api.Config) (*ens.ENS, error) {
	log.Info("connecting to ENS API", "url", endpoint)
	client, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error connecting to ENS API %s: %s", endpoint, err)
	}
	ensClient := ethclient.NewClient(client)

	ensRoot := config.EnsRoot
	if addr != (common.Address{}) {
		ensRoot = addr
	} else {
		a, err := detectEnsAddr(client)
		if err == nil {
			ensRoot = a
		} else {
			log.Warn(fmt.Sprintf("could not determine ENS contract address, using default %s", ensRoot), "err", err)
		}
	}
	transactOpts := bind.NewKeyedTransactor(config.Swap.PrivateKey())
	dns, err := ens.NewENS(transactOpts, ensRoot, ensClient)
	if err != nil {
		return nil, err
	}
	log.Debug(fmt.Sprintf("-> Swarm Domain Name Registrar %v @ address %v", endpoint, ensRoot.Hex()))
	return dns, err
}

// detectEnsAddr determines the ENS contract address by getting both the
// version and genesis hash using the client and matching them to either
// mainnet or testnet addresses
func detectEnsAddr(client *rpc.Client) (common.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var version string
	if err := client.CallContext(ctx, &version, "net_version"); err != nil {
		return common.Address{}, err
	}

	block, err := ethclient.NewClient(client).BlockByNumber(ctx, big.NewInt(0))
	if err != nil {
		return common.Address{}, err
	}

	switch {

	case version == "1" && block.Hash() == params.MainnetGenesisHash:
		log.Info("using Mainnet ENS contract address", "addr", ens.MainNetAddress)
		return ens.MainNetAddress, nil

	case version == "3" && block.Hash() == params.TestnetGenesisHash:
		log.Info("using Testnet ENS contract address", "addr", ens.TestNetAddress)
		return ens.TestNetAddress, nil

	default:
		return common.Address{}, fmt.Errorf("unknown version and genesis hash: %s %s", version, block.Hash())
	}
}

/*
Start is called when the stack is started
* starts the network kademlia hive peer management
* (starts netStore level 0 api)
* starts DPA level 1 api (chunking -> store/retrieve requests)
* (starts level 2 api)
* starts http proxy server
* registers url scheme handlers for bzz, etc
* TODO: start subservices like sword, swear, swarmdns
*/
// implements the node.Service interface
func (sw *Swarm) Start(srv *p2p.Server) error {
	startTime = time.Now()
	connectPeer := func(url string) error {
		node, err := discover.ParseNode(url)
		if err != nil {
			return fmt.Errorf("invalid node URL: %v", err)
		}
		srv.AddPeer(node)
		return nil
	}
	// set chequebook
	if sw.swapEnabled {
		ctx := context.Background() // The initial setup has no deadline.
		err := sw.SetChequebook(ctx)
		if err != nil {
			return fmt.Errorf("Unable to set chequebook for SWAP: %v", err)
		}
		log.Debug(fmt.Sprintf("-> cheque book for SWAP: %v", sw.config.Swap.Chequebook()))
	} else {
		log.Debug(fmt.Sprintf("SWAP disabled: no cheque book set"))
	}

	log.Warn(fmt.Sprintf("Starting Swarm service"))
	sw.hive.Start(
		discover.PubkeyID(&srv.PrivateKey.PublicKey),
		func() string { return srv.ListenAddr },
		connectPeer,
	)
	log.Info(fmt.Sprintf("Swarm network started on bzz address: %v", sw.hive.Addr()))

	sw.dpa.Start()
	log.Debug(fmt.Sprintf("Swarm DPA started"))

	// start swarm http proxy server
	if sw.config.Port != "" {
		addr := net.JoinHostPort(sw.config.ListenAddr, sw.config.Port)
		go httpapi.StartHttpServer(sw.api, &httpapi.ServerConfig{
			Addr:       addr,
			CorsString: sw.corsString,
		})
		log.Info(fmt.Sprintf("Swarm http proxy started on %v", addr))

		if sw.corsString != "" {
			log.Debug(fmt.Sprintf("Swarm http proxy started with corsdomain: %v", sw.corsString))
		}
	}

	sw.periodicallyUpdateGauges()

	startCounter.Inc(1)
	return nil
}

func (sw *Swarm) periodicallyUpdateGauges() {
	ticker := time.NewTicker(updateGaugesPeriod)

	go func() {
		for range ticker.C {
			sw.updateGauges()
		}
	}()
}

func (sw *Swarm) updateGauges() {
	dbSizeGauge.Update(int64(sw.lstore.DbCounter()))
	cacheSizeGauge.Update(int64(sw.lstore.CacheCounter()))
	uptimeGauge.Update(time.Since(startTime).Nanoseconds())
}

// Stop implements the node.Service interface
// stops all component services.
func (sw *Swarm) Stop() error {
	sw.dpa.Stop()
	err := sw.hive.Stop()
	if ch := sw.config.Swap.Chequebook(); ch != nil {
		ch.Stop()
		ch.Save()
	}

	if sw.lstore != nil {
		sw.lstore.DbStore.Close()
	}
	sw.sfs.Stop()
	stopCounter.Inc(1)
	return err
}

// Protocols implements the node.Service interface
func (sw *Swarm) Protocols() []p2p.Protocol {
	proto, err := network.Bzz(sw.depo, sw.backend, sw.hive, sw.dbAccess, sw.config.Swap, sw.config.SyncParams, sw.config.NetworkId)
	if err != nil {
		return nil
	}
	return []p2p.Protocol{proto}
}

// APIs implements node.Service
// Apis returns the RPC Api descriptors the Swarm implementation offers
func (sw *Swarm) APIs() []rpc.API {
	return []rpc.API{
		// public APIs
		{
			Namespace: "bzz",
			Version:   "0.1",
			Service:   &Info{sw.config, chequebook.ContractParams},
			Public:    true,
		},
		// admin APIs
		{
			Namespace: "bzz",
			Version:   "0.1",
			Service:   api.NewControl(sw.api, sw.hive),
			Public:    false,
		},
		{
			Namespace: "chequebook",
			Version:   chequebook.Version,
			Service:   chequebook.NewApi(sw.config.Swap.Chequebook),
			Public:    false,
		},
		{
			Namespace: "swarmfs",
			Version:   fuse.Swarmfs_Version,
			Service:   sw.sfs,
			Public:    false,
		},
		// storage APIs
		// DEPRECATED: Use the HTTP API instead
		{
			Namespace: "bzz",
			Version:   "0.1",
			Service:   api.NewStorage(sw.api),
			Public:    true,
		},
		{
			Namespace: "bzz",
			Version:   "0.1",
			Service:   api.NewFileSystem(sw.api),
			Public:    false,
		},
		// {Namespace, Version, api.NewAdmin(sw), false},
	}
}

func (sw *Swarm) Api() *api.Api {
	return sw.api
}

// SetChequebook ensures that the local checquebook is set up on chain.
func (sw *Swarm) SetChequebook(ctx context.Context) error {
	err := sw.config.Swap.SetChequebook(ctx, sw.backend, sw.config.Path)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("new chequebook set (%v): saving config file, resetting all connections in the hive", sw.config.Swap.Contract.Hex()))
	sw.hive.DropAll()
	return nil
}

// NewLocalSwarm without netStore
func NewLocalSwarm(datadir, port string) (sw *Swarm, err error) {

	prvKey, err := crypto.GenerateKey()
	if err != nil {
		return
	}

	config := api.NewDefaultConfig()
	config.Path = datadir
	config.Init(prvKey)
	config.Port = port

	dpa, err := storage.NewLocalDPA(datadir)
	if err != nil {
		return
	}

	sw = &Swarm{
		api:    api.NewApi(dpa, nil),
		config: config,
	}

	return
}

// serialisable info about swarm
type Info struct {
	*api.Config
	*chequebook.Params
}

func (info *Info) Info() *Info {
	return info
}
