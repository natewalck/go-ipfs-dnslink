# go-ipfs-dnslink
Quick and messy app to update dnslinks on namecheap

To user this binary, setup a config.yaml with the following keys (or use env vars):

    --- 
    user: "USERNAMEHERE"
    api_user: "USERNAMEHERE"
    api_token: "API_KEY_HERE"

And run the binary:

    ./ipfs-dnslink --config ./config.yaml --domain your.domain.com --cid IPFS_HASH_HERE 
