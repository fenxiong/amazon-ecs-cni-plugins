// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package engine

import (
	"net"
	"syscall"

	"github.com/aws/amazon-ecs-cni-plugins/pkg/netlinkwrapper"
	"github.com/containernetworking/cni/pkg/ns"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

const (
	instanceMetadataEndpoint = "169.254.169.254/32"
)

var linkWithMACNotFoundError = errors.New("engine: device with mac address not found")

// setupNamespaceClosureContext wraps the parameters and the method to configure the container's namespace
type setupNamespaceClosureContext struct {
	netLink     netlinkwrapper.NetLink
	ifName      string
	deviceName  string
	macAddress  string
	ipv4Addr    *netlink.Addr
	ipv6Addr    *netlink.Addr
	ipv4Gateway net.IP
	ipv6Gateway net.IP
	blockIMDS   bool
}

// teardownNamespaceClosureContext wraps the parameters and the method to teardown the
// container's namespace
type teardownNamespaceClosureContext struct {
	netLink      netlinkwrapper.NetLink
	hardwareAddr net.HardwareAddr
}

// newSetupNamespaceClosureContext creates a new setupNamespaceClosure object
func newSetupNamespaceClosureContext(netLink netlinkwrapper.NetLink,
	ifName string, deviceName string, macAddress string, ipv4Address string, ipv6Address string,
	ipv4Gateway string, ipv6Gateway string, blockIMDS bool) (*setupNamespaceClosureContext, error) {
	nlIPV4Addr, err := netLink.ParseAddr(ipv4Address)
	if err != nil {
		return nil, errors.Wrap(err,
			"setupNamespaceClosure engine: unable to parse ipv4 address for the interface")
	}

	ipv4GatewayIP := net.ParseIP(ipv4Gateway)
	if ipv4GatewayIP == nil {
		return nil, errors.New(
			"setupNamespaceClosure engine: unable to parse address of the ipv4 gateway")
	}

	nsClosure := &setupNamespaceClosureContext{
		netLink:     netLink,
		ifName:      ifName,
		deviceName:  deviceName,
		macAddress:  macAddress,
		ipv4Addr:    nlIPV4Addr,
		ipv4Gateway: ipv4GatewayIP,
		blockIMDS:   blockIMDS,
	}
	if ipv6Address != "" {
		nlIPV6Addr, err := netLink.ParseAddr(ipv6Address)
		if err != nil {
			return nil, errors.Wrap(err,
				"setupNamespaceClosure engine: unable to parse ipv6 address for the interface")
		}
		ipv6GatewayIP := net.ParseIP(ipv6Gateway)
		if ipv6GatewayIP == nil {
			return nil, errors.New(
				"setupNamespaceClosure engine: unable to parse address of the ipv6 gateway")
		}
		nsClosure.ipv6Addr = nlIPV6Addr
		nsClosure.ipv6Gateway = ipv6GatewayIP
	}

	return nsClosure, nil
}

// newTeardownNamespaceClosureContext creates a new teardownNamespaceClosure object
func newTeardownNamespaceClosureContext(netLink netlinkwrapper.NetLink,
	mac string) (*teardownNamespaceClosureContext, error) {
	hardwareAddr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, errors.Wrapf(err,
			"newTeardownNamespaceClosure engine: malformatted mac address specified")
	}

	return &teardownNamespaceClosureContext{
		netLink:      netLink,
		hardwareAddr: hardwareAddr,
	}, nil
}

// run defines the closure to execute within the container's namespace to configure it
// appropriately
func (closureContext *setupNamespaceClosureContext) run(_ ns.NetNS) error {
	// Get the link for the ENI device
	eniLink, err := closureContext.netLink.LinkByName(closureContext.deviceName)
	if err != nil {
		return errors.Wrapf(err,
			"setupNamespaceClosure engine: unable to get link for device '%s'",
			closureContext.deviceName)
	}

	err = closureContext.netLink.LinkSetName(eniLink, closureContext.ifName)
	if err != nil {
		return errors.Wrap(err, "setupNamespaceClosure engine: unable to change interface name")
	}
	// Add the IPV4 Address to the link
	err = closureContext.netLink.AddrAdd(eniLink, closureContext.ipv4Addr)
	if err != nil {
		return errors.Wrap(err,
			"setupNamespaceClosure engine: unable to add ipv4 address to the interface")
	}

	if closureContext.ipv6Addr != nil {
		// Add the IPV6 Address to the link
		err = closureContext.netLink.AddrAdd(eniLink, closureContext.ipv6Addr)
		if err != nil {
			return errors.Wrap(err,
				"setupNamespaceClosure engine: unable to add ipv6 address to the interface")
		}
	}

	// Bring it up
	err = closureContext.netLink.LinkSetUp(eniLink)
	if err != nil {
		return errors.Wrap(err,
			"setupNamespaceClosure engine: unable to bring up the device")
	}

	// Add a blackhole route for IMDS endpoint if required
	if closureContext.blockIMDS {
		_, imdsNetwork, err := net.ParseCIDR(instanceMetadataEndpoint)
		if err != nil {
			// This should never happen because we always expect
			// 169.254.169.254/32 to be parsed without any errors
			return errors.Wrapf(err, "setupNamespaceClosure engine: unable to parse instance metadata endpoint")
		}
		if err = closureContext.netLink.RouteAdd(&netlink.Route{
			Dst:  imdsNetwork,
			Type: syscall.RTN_BLACKHOLE,
		}); err != nil {
			return errors.Wrapf(err, "setupNamespaceClosure engine: unable to add route to block instance metadata")
		}
	}

	// Setup ipv4 route for the gateway
	err = closureContext.netLink.RouteAdd(&netlink.Route{
		Gw: closureContext.ipv4Gateway,
	})
	if err != nil {
		return errors.Wrap(err,
			"setupNamespaceClosure engine: unable to add the route for the ipv4 gateway")
	}

	if closureContext.ipv6Addr != nil {
		// Setup ipv6 route for the gateway
		err = closureContext.netLink.RouteAdd(&netlink.Route{
			LinkIndex: eniLink.Attrs().Index,
			Gw:        closureContext.ipv6Gateway,
		})
		if err != nil && !isRouteExistsError(err) {
			return errors.Wrap(err,
				"setupNamespaceClosure engine: unable to add the route for the ipv6 gateway")
		}
	}

	return nil
}

// isRouteExistsError returns true if the error type is syscall.EEXIST
// This helps us determine if we should ignore this error as the route
// that we want to add already exists in the routing table
func isRouteExistsError(err error) bool {
	if errno, ok := err.(syscall.Errno); ok {
		return errno == syscall.EEXIST
	}

	return false
}

// run defines the closure to execute within the container's namespace to tear it down
func (closureContext *teardownNamespaceClosureContext) run(_ ns.NetNS) error {
	return nil
}

// getLinkByHardwareAddress gets the link device based on the mac address
func getLinkByHardwareAddress(netLink netlinkwrapper.NetLink, hardwareAddr net.HardwareAddr) (netlink.Link, error) {
	links, err := netLink.LinkList()
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		// TODO: Evaluate if reflect.DeepEqual is a better alternative here
		if link.Attrs().HardwareAddr.String() == hardwareAddr.String() {
			return link, nil
		}
	}

	return nil, linkWithMACNotFoundError
}
