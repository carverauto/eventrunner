# IPv6 Routing Configuration

## Overview
Our network consists of multiple subnets with both IPv4 and IPv6 addressing:
- Primary LAN: 192.168.1.0/24 → 2001:470:c0b5:1::/64
- Server Network: 192.168.2.0/24 → 2001:470:c0b5:2::/64
- HE.net IPv6 tunnel terminating at our router (192.168.1.1)

## Initial Issues
Initially, hosts on the primary LAN (192.168.1.0/24) could not reach IPv6 hosts on the server network (2001:470:c0b5:2::/64), despite the router having connectivity to both networks.

## Resolution Steps

### 1. IPv6 Forwarding
Verified IPv6 forwarding was enabled on the router:
```bash
sysctl net.ipv6.conf.all.forwarding
# Should return 1
```

### 2. Static Routes
Added necessary static routes on the UDM Pro to ensure proper IPv6 routing between subnets:
```bash
ip -6 route add 2001:470:c0b5:2::/64 dev br2
```

Note: These static routes will need to be re-added after router reboots unless configured persistently.

### 3. BGP Configuration (FRR)
Modified FRR configuration to handle BGP peering with Calico nodes. Key changes included:
- Explicit IPv4 and IPv6 neighbor configurations
- Removed unnecessary HE.net BGP configuration (tunnel works without BGP)
- Added proper address-family configurations for both IPv4 and IPv6
- Configured route redistribution for connected and kernel routes

Key FRR configuration points:
- BGP ASN: 65000 (Router)
- Calico Node ASN: 65001
- Redistributing connected and kernel routes
- No route-maps required for basic functionality

## Known Issues and Solutions
- Capability code 71 warnings in FRR logs are benign and can be ignored
- HE.net tunnel does not require BGP configuration
- Neighbor discovery (NDP) must be allowed through firewalls

## Testing Connectivity
You can verify IPv6 connectivity using:
```bash
# From a host on 192.168.1.0/24 network
ping6 2001:470:c0b5:2::110
telnet 2001:470:c0b5:2::110 22
```

## Notes
- Monitor FRR logs for BGP session status: `tail -f /var/log/frr/bgpd.log`
- Check BGP peer status: `vtysh -c "show bgp ipv6 unicast summary"`
- Review route tables: `vtysh -c "show ipv6 route bgp"`
- Calico node status can be checked with: `calicoctl node status`