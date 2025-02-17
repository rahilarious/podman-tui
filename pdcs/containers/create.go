package containers

import (
	"net"
	"strconv"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/volumes"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/containers/podman/v4/pkg/specgenutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// CreateOptions container create options.
type CreateOptions struct {
	Name                  string
	Labels                []string
	Image                 string
	Remove                bool
	Pod                   string
	Hostname              string
	IPAddress             string
	Network               string
	MacAddress            string
	Publish               []string
	Expose                []string
	PublishAll            bool
	DNSServer             []string
	DNSOptions            []string
	DNSSearchDomain       []string
	Volume                string
	ImageVolume           string
	SelinuxOpts           []string
	ApparmorProfile       string
	Seccomp               string
	NoNewPriv             bool
	Mask                  string
	Unmask                string
	HealthCmd             string
	HealthInterval        string
	HealthRetries         string
	HealthStartPeroid     string
	HealthTimeout         string
	HealthOnFailure       string
	HealthStartupCmd      string
	HealthStartupInterval string
	HealthStartupRetries  string
	HealthStartupSuccess  string
	HealthStartupTimeout  string
}

// Create creates a new container.
func Create(opts CreateOptions) ([]string, error) { //nolint:cyclop
	var (
		warningResponse []string
		createOptions   entities.ContainerCreateOptions
	)

	log.Debug().Msgf("pdcs: podman container create %v", opts)
	utils.DefineCreateDefaults(&createOptions)

	conn, err := registry.GetConnection()
	if err != nil {
		return warningResponse, err
	}

	if len(opts.Labels) > 0 {
		createOptions.Label = opts.Labels
	}

	createOptions.Name = opts.Name
	createOptions.Rm = opts.Remove

	createOptions.Hostname = opts.Hostname

	if len(opts.Expose) > 0 {
		createOptions.Expose = opts.Expose
	}

	createOptions.PublishAll = opts.PublishAll

	createOptions.Net, err = containerNetworkOptions(opts)
	if err != nil {
		return warningResponse, err
	}

	if opts.Pod != "" {
		createOptions.Pod = opts.Pod
		createOptions.Net.Network.NSMode = specgen.FromPod
	} else {
		createOptions.Net.Network.NSMode = specgen.Default
	}

	if opts.Volume != "" {
		// get volume path
		volFilter := make(map[string][]string)
		volFilter["name"] = []string{opts.Volume}

		volResponse, err := volumes.List(conn, new(volumes.ListOptions).WithFilters(volFilter))
		if err != nil {
			return warningResponse, err
		}

		if len(volResponse) > 0 {
			createOptions.Volume = []string{volResponse[0].Mountpoint}
		}
	}

	createOptions.ImageVolume = opts.ImageVolume

	// security options
	if opts.ApparmorProfile != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "apparmor="+opts.ApparmorProfile)
	}

	if opts.Mask != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "mask="+opts.Mask)
	}

	if opts.Unmask != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "unmask="+opts.Unmask)
	}

	if opts.Seccomp != "" {
		createOptions.SeccompPolicy = opts.Seccomp
	}

	if opts.NoNewPriv {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "no-new-privileges")
	}

	if len(opts.SelinuxOpts) > 0 {
		for _, selinuxLabel := range opts.SelinuxOpts {
			createOptions.SecurityOpt = append(createOptions.SecurityOpt, "label="+selinuxLabel)
		}
	}

	// add healthcheck options
	if err := containerHealthOptions(&createOptions, opts); err != nil {
		return warningResponse, err
	}

	s := specgen.NewSpecGenerator(opts.Name, false)
	if err := specgenutil.FillOutSpecGen(s, &createOptions, nil); err != nil {
		return warningResponse, err
	}

	s.Image = opts.Image

	// validate spec
	if err := s.Validate(); err != nil {
		return warningResponse, err
	}

	response, err := containers.CreateWithSpec(conn, s, &containers.CreateOptions{})
	if err != nil {
		return warningResponse, err
	}

	warningResponse = response.Warnings

	return warningResponse, nil
}

func containerHealthOptions(createOptions *entities.ContainerCreateOptions, opts CreateOptions) error { //nolint:cyclop
	createOptions.HealthInterval = define.DefaultHealthCheckInterval
	createOptions.StartupHCInterval = define.DefaultHealthCheckInterval
	createOptions.HealthRetries = define.DefaultHealthCheckRetries
	createOptions.HealthStartPeriod = define.DefaultHealthCheckStartPeriod
	createOptions.HealthTimeout = define.DefaultHealthCheckTimeout
	createOptions.StartupHCTimeout = define.DefaultHealthCheckTimeout
	createOptions.HealthOnFailure = opts.HealthOnFailure

	if opts.HealthCmd == "" {
		createOptions.HealthCmd = "none"

		return nil
	}

	createOptions.HealthCmd = opts.HealthCmd

	if opts.HealthInterval != "" {
		createOptions.HealthInterval = opts.HealthInterval
	}

	if opts.HealthStartPeroid != "" {
		createOptions.HealthStartPeriod = opts.HealthStartPeroid
	}

	if opts.HealthTimeout != "" {
		createOptions.HealthTimeout = opts.HealthTimeout
	}

	if opts.HealthStartupCmd != "" {
		createOptions.StartupHCCmd = opts.HealthStartupCmd
	}

	if opts.HealthStartupInterval != "" {
		createOptions.StartupHCInterval = opts.HealthStartupInterval
	}

	if opts.HealthStartupTimeout != "" {
		createOptions.StartupHCTimeout = opts.HealthStartupTimeout
	}

	if opts.HealthRetries != "" {
		retries, err := strconv.ParseUint(opts.HealthRetries, 10, 32)
		if err != nil {
			return err
		}

		retriesWd := uint(retries)
		createOptions.HealthRetries = retriesWd
	}

	if opts.HealthStartupRetries != "" {
		startupRetries, err := strconv.ParseUint(opts.HealthStartupRetries, 10, 32)
		if err != nil {
			return err
		}

		startupRetriesWd := uint(startupRetries)
		createOptions.StartupHCRetries = startupRetriesWd
	}

	if opts.HealthStartupSuccess != "" {
		startupSuccess, err := strconv.ParseUint(opts.HealthStartupSuccess, 10, 32)
		if err != nil {
			return err
		}

		startupSuccessWd := uint(startupSuccess)
		createOptions.StartupHCRetries = startupSuccessWd
	}

	return nil
}

func containerNetworkOptions(opts CreateOptions) (*entities.NetOptions, error) { //nolint:cyclop
	var (
		err           error
		perNetworkOpt types.PerNetworkOptions
	)

	netOptions := &entities.NetOptions{}
	netOptions.Networks = make(map[string]types.PerNetworkOptions)

	var dnsServers []net.IP

	for _, d := range opts.DNSServer {
		addr := net.ParseIP(d)
		if addr != nil {
			dnsServers = append(dnsServers, addr)

			continue
		}

		return nil, errors.Wrap(utils.ErrInvalidDNSAddress, d)
	}

	if len(dnsServers) > 0 {
		netOptions.DNSServers = dnsServers
		netOptions.DNSOptions = opts.DNSOptions
		netOptions.DNSSearch = opts.DNSSearchDomain
	}

	if len(opts.Publish) > 0 {
		netOptions.PublishPorts, err = specgenutil.CreatePortBindings(opts.Publish)
		if err != nil {
			return nil, err
		}
	}

	if opts.Network != "" { //nolint:nestif
		if opts.MacAddress != "" {
			mac, err := net.ParseMAC(opts.MacAddress)
			if err != nil {
				return nil, err
			}

			perNetworkOpt.StaticMAC = types.HardwareAddr(mac)
		}

		if opts.IPAddress != "" {
			addr := net.ParseIP(opts.IPAddress)

			if addr == nil {
				return nil, errors.Wrap(utils.ErrInvalidIPAddress, opts.IPAddress)
			}

			perNetworkOpt.StaticIPs = []net.IP{addr}
		}

		netOptions.Networks[opts.Network] = perNetworkOpt
	}

	return netOptions, nil
}
