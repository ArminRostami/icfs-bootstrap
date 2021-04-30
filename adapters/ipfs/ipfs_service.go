package ipfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"

	"github.com/pkg/errors"
)

type IpfsService struct {
	repoPath string
	ctx      context.Context
	ipfs     icore.CoreAPI
}

func NewService() (context.CancelFunc, *IpfsService, error) {
	pr, err := config.PathRoot()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get default config path")
	}
	if err := setupPlugins(pr); err != nil {
		return nil, nil, errors.Wrap(err, "failed to setup plugins")
	}
	ctx, cancel := context.WithCancel(context.Background())

	service := &IpfsService{ctx: ctx, repoPath: pr}
	err = service.start()
	if err != nil {
		return cancel, nil, errors.Wrap(err, "failed to start service")
	}

	return cancel, service, nil
}

func (s *IpfsService) start() error {
	err := s.setupRepo()
	if err != nil {
		return errors.Wrap(err, "failed to start ipfs service")
	}

	ipfs, err := createNode(s.ctx, s.repoPath)
	if err != nil {
		return errors.Wrap(err, "failed to spawn default node")
	}
	s.ipfs = ipfs
	return nil
}

func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open repo")
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start new node")
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func (s *IpfsService) setupRepo() error {
	if fsrepo.IsInitialized(s.repoPath) {
		cfg, err := fsrepo.ConfigAt(s.repoPath)
		if err != nil {
			return errors.Wrap(err, "failed to open config file")
		}
		log.Println(cfg.Bootstrap)
		return nil
	}

	log.Printf("setting up new repo at %s\n", s.repoPath)
	cfg, err := config.Init(io.Discard, 2048)
	if err != nil {
		return errors.Wrap(err, "failed to init config")
	}

	if err = setBootstrap(cfg); err != nil {
		return errors.Wrap(err, "failed to set bootstrap")
	}

	if err = fsrepo.Init(s.repoPath, cfg); err != nil {
		return errors.Wrap(err, "failed to init repo")
	}

	err = writeSwarmKey(swKey, s.repoPath)
	if err != nil {
		return errors.Wrap(err, "failed to copy swarm.key file")
	}

	return nil
}

func setBootstrap(cfg *config.Config) error {
	ip, err := getOutboundIP()
	if err != nil {
		return errors.Wrap(err, "failed to get ip")
	}

	bootStr := getBootstrapString(ip, cfg.Identity.PeerID)
	log.Println(bootStr)

	peers, err := config.ParseBootstrapPeers([]string{bootStr})
	if err != nil {
		return errors.Wrap(err, "failed to parse peerAddr")
	}

	cfg.SetBootstrapPeers(peers)
	return nil
}

func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", errors.Wrap(err, "failed to get ip")
	}
	ip, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "", errors.Wrap(err, "failed to get ip")
	}
	return ip.IP.String(), nil
}

func (s *IpfsService) GetConInfo() (string, string, error) {
	cfg, err := fsrepo.ConfigAt(s.repoPath)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to read config")
	}
	if len(cfg.Bootstrap) < 1 {
		return "", "", errors.Wrap(err, "no bootstrap info")
	}
	return cfg.Bootstrap[0], swKey, nil
}

func getBootstrapString(ip, id string) string {
	return fmt.Sprintf("/ip4/%s/tcp/4001/ipfs/%s", ip, id)
}

func writeSwarmKey(key, repoPath string) error {
	if err := os.WriteFile(path.Join(repoPath, "swarm.key"), []byte(key), 0644); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	return nil
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return errors.Wrap(err, "failed to load plugins")
	}

	if err := plugins.Initialize(); err != nil {
		return errors.Wrap(err, "failed to init plugins")
	}

	if err := plugins.Inject(); err != nil {
		return errors.Wrap(err, "failed to inject plugins")
	}

	return nil
}
