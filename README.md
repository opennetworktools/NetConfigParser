# NetConfigParser

NetConfigParser is a tool that parses your network device's startup or running configurations. It currently supports Cisco IOS-XE, with more features under development. Support for the Arista EOS is also in queue.

## Getting Started

Add the `NetConfigParser` package in your project,

```
go get github.com/opennetworktools/NetConfigParser
```

Import the package in your project,

```
import (
	netconfigparser "github.com/opennetworktools/NetConfigParser"
)
```

Refer [example/main.go](example/main.go) to know how to parse your running/startup configs.

> [!NOTE]
> NetConfigParser is under development. Feel free to go through the open issues and contribute to help us speed up the development.

## Design

For alpha release the plan is to focus only on the Cisco IOS-XE with support for bare minimum protocols. For now I am going to stick with the Cisco IOS-XE and the Arista EOS (which will be the focus for beta and v1 release).

### Supported features in the Cisco IOS-XE:

1. Interfaces
    - FastEthernet, GigabitEthernet, TenGigGigabitEthernet
    - IP Address and Subnet Mask
    - Description
    - Shutdown Status
    - Attached ACL's
2. BGP
    - ASN
    - RouterID
    - Neighbors
        - ASN
        - IP Address
        - Description
        - Timers
3. Route Map
    - Name
    - Action
    - Sequence Number
    - Attached IP Prefix List (TODO)
4. IP Prefix-list
    - Name
    - Rules
        - Sequence Number
        - Action
        - IP Address and Subnet Mask
5. IP Access-list
    - Name
    - Type
    - Rules
        - Action
        - Type
        - SrcIP, SrcMask
        - DstIP, DstMask
        - Protocol, Port Number




