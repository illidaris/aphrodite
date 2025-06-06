package dep

import "net"

// InterfaceAddrs defines the interface used for retrieving network addresses.
type InterfaceAddrs func() ([]net.Addr, error)
