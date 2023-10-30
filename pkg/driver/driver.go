package driver

import (
	"log"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
)

const CSI_DRIVER_NAME = "digitalocean"

type (
	CsiDriver struct {
		name string

		grpcServer        *grpc.Server
		grpcServerAddress string

		// DigitalOcean region where volumes will be provisioned.
		region string
	}

	NewCsiDriverArgs struct {
		Region   string
		Token    string // DigitalOcean API token
		Endpoint string
	}
)

func NewCsiDriver(args NewCsiDriverArgs) *CsiDriver {
	return &CsiDriver{
		name:              CSI_DRIVER_NAME,
		grpcServerAddress: parseEndpoint(args.Endpoint),
		region:            args.Region,
	}
}

// Run starts the gRPC server.
func (c *CsiDriver) Run() error {
	listener, err := net.Listen("unix", c.grpcServerAddress)
	if err != nil {
		log.Panicf("Error listen at the given endpoint : %v", err)
	}

	c.grpcServer = grpc.NewServer()

	csi.RegisterNodeServer(c.grpcServer, c)
	csi.RegisterControllerServer(c.grpcServer, c)
	csi.RegisterIdentityServer(c.grpcServer, c)

	return c.grpcServer.Serve(listener)
}

// parseEndpoint tries to parse a CSI driver endpoint to a gRPC server address.
func parseEndpoint(endpoint string) string {
	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("Error parsing endpoint : %v", err)
	}

	// Only Unix domain sockets can be used as endpoints as mentioned in the CSI specification.
	// Unix domain sockets are a type of Inter Process Communication (IPC) mechanism that allows
	// processes running on the same machine to communicate with each other. They reside within the
	// kernel. They have a file-system like address space - processes reference them using file paths.
	// They are more efficient than network sockets as the communication stays within the kernel and
	// does not have to go through the network stack.
	if parsedEndpoint.Scheme != "unix" {
		log.Fatalf("Only unix domain sockets are supported as the endpoint")
	}

	var gRPCServerAddress string
	if parsedEndpoint.Host == "" {
		gRPCServerAddress = filepath.FromSlash(parsedEndpoint.Path)
	} else {
		gRPCServerAddress = path.Join(parsedEndpoint.Host, filepath.FromSlash(parsedEndpoint.Path))
	}

	// Suppose we start the gRPC server for the first time. If we hit Ctrl + c quitting the server and
	// try to restart it then we will get an error that address is already in use. Thats'y we will
	// remove the Unix socket file everytime before starting the gRPC server.
	if err := os.Remove(gRPCServerAddress); err != nil && !os.IsExist(err) {
		log.Panicf("Error removing listener address : %v", err)
	}

	return gRPCServerAddress
}
