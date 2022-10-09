
## full-node


[doc](https://aptos.dev/nodes/full-node/fullnode-source-code-or-docker#approach-2-using-docker)


1. Install Docker.

2. Create a directory for your local public fullnode, and cd into it. For example:

```shell
mkdir aptos-fullnode && cd aptos-fullnode
```

3. Run the following script to prepare your local config and data dir for Devnet:
```shell
mkdir data && \
curl -O https://raw.githubusercontent.com/aptos-labs/aptos-core/devnet/config/src/config/test_data/public_full_node.yaml && \
curl -O https://devnet.aptoslabs.com/waypoint.txt && \
curl -O https://devnet.aptoslabs.com/genesis.blob
```


TIP
To connect to other networks, you can find genesis and waypoint here -> https://github.com/aptos-labs/aptos-genesis-waypoint

4. Finally, start the fullnode via docker:
```shell
docker run --pull=always --rm -p 8080:8080 -p 9101:9101 -p 6180:6180 -v $(pwd):/opt/aptos/etc -v $(pwd)/data:/opt/aptos/data --workdir /opt/aptos/etc --name=aptos-fullnode aptoslabs/validator:devnet aptos-node -f /opt/aptos/etc/public_full_node.yaml
```


Ensure you have opened the relevant ports - 8080, 9101 and 6180 and you may also need to update the 127.0.0.1 with 0.0.0.0 in the public_full_node.yaml - listen_address and api\address
