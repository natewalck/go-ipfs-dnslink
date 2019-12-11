# go-ipfs-dnslink
Quick and messy binary to update dnslink TXT records on namecheap. This uses https://github.com/billputer/go-namecheap

To use this binary, setup a config.yaml with the following keys (or use env vars). See --help for more information.

    --- 
    user: "USERNAMEHERE"
    api_user: "USERNAMEHERE"
    api_token: "API_KEY_HERE"
    domain: "your.domain.com"

Note: You may optionally add an int in config.yaml for ttl.

And run the binary:

    ./ipfs-dnslink --config ./config.yaml --cid IPFS_HASH_HERE 
