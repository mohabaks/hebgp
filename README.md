# BGP Information Query Tool

This is a simple command-line tool that allows you to query information about IP addresses, subnets, and Autonomous System Numbers (ASNs) using the [bgp.he.net](https://bgp.he.net) website.

## Features

- Query information about IP addresses, subnets, ASNs, and organizations.
- Retrieve data in JSON format for easy integration with other tools.
- Supports multiple command-line options to specify the type of query.

## Usage

```
# Query for ASN information
hebgp -asn AS63293

# Query for IP information
hebgp -ip 1.1.1.1

# Query for network block information
hebgp -net 41.223.111.0/22
hebgp -net 41.223.111.0/22|jq '.[]'

# Query for organization information
hebgp -org facebook

# Show help message
hebpg -h
```

## Installation

> **Dependencies**: [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery): A package for parsing HTML documents using CSS selectors.

```
go install github.com/mohabaks/hebgp@latest
```

OR

```
git clone https://github.com/yourusername/hebgp.git
cd hebgp
go build
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
