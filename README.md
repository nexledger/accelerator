# Nexledger Accelerator
Nexledger Accelerator is a software component designed to improve the performance of a blockchain network, e.g. Hyperledger Fabric, in terms of transaction throughput. Accelerator enables the blockchain network to deal with explosive transaction requests from applications. 

Accelerator receives transactions from clients on behalf of blockchain nodes and provides transaction acceleration in terms of TPS (Transaction Per Second) by classifying, aggregating, and routing the transactions. The current version of Accelerator is compatible with Hyperledger Fabric v1.4.


## Getting Started
### Prerequisites
- Go 1.11+
- Docker (17.06.2-ce or greater)
- Docker-compose (1.14.0 or greater)

### Building Accelerator
Accelerator supports go module for dependency management. To build the executable, please simply execute `go build`.
```bash
$ go build cmd/accelerator.go
```

## Running ping example
The ping example shows how to configure and run Accelerator. The example is placed in In `examples/ping`. 

To bootstrap Fabric network, please run `start.sh` script. It configures the Hyperledger Fabric network and installs/instantiates the example chaincode.
```bash
$ ./examples/ping/start.sh
```
  
To serve requests from clients, Accelerator should be up and running with proper configuration.
```bash
$ ./accelerator -f examples/ping/configs/accelerator.yaml
```

Accelerator is a gRPC server and the gRPC services are described in `protos/accelerator.proto`.
You may send transactions using `examples/ping/ping_test.go` that has gRPC client for ping example. 
```bash
$  cd examples/ping
$  go test
```



You can terminate and remove the network by run `stop.sh` script.
```bash
$ ./examples/ping/stop.sh
```

### Under the hood
#### Modifying chaincode
In order to apply Accelerator to your business, you need to modify your chaincode. 
Accelerator aggregates multiple transactions into a batched transaction and submits the batched transaction to the endorsers. 
So chaincodes operating with Accelerator should be modified to execute aggregated transactions individually.

`contracts/src/ping/ping.go` is the example chaincode with simple KV write/read operations.
`ping.go` imports `batchutil.go` for segregating batched transactions from Accelerator and delegates the invocations to `Invoke()` in `batchutil.go`.

```go
func (t *PingPongChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fnc := string(stub.GetArgs()[0])
	switch fnc {
	case "ping":
		return Invoke(stub, t.ping)
	case "pong":
		return Invoke(stub, t.pong)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'insert', 'query'")
}
```  

#### Configuring Accelerator
Accelerator should be configured with target chaincode functions and corresponding batch configs. 
The configuration file for `ping` example is placed at `configs/accelerator.yaml` 
```yaml
sdk: "examples/ping/configs/accelerator-sdk.yaml"
host: "localhost"
port: 8090
userName: "Admin"
organization: "peerorg1"
batch:
  - type: "execute"
    channelId: "accelerator"
    chaincodeName: "ping"
    fcn: "ping"
    queueSize: 1000
    maxWaitTimeSeconds: 5
    maxBatchItems: 10
  - type: "query"
    channelId: "accelerator"
    chaincodeName: "ping"
    fcn: "pong"
    queueSize: 1000
    maxWaitTimeSeconds: 5
    maxBatchItems: 10
```
- `sdk`: Path to the fabric SDK configuration File
- `host`: Host of Accelerator
- `port`: Port of Accelerator
- `queueSize`: The size of the in-memory queue that kept requested transactions until processing.
- `maxWaitTimeSeconds`: Maximum waiting time in seconds to create a new batch transaction.
- `maxBatchItems`: Maximum number of items for a new batch transaction.

## Whitepaper
[Whitepaper](https://github.com/nexledger/accelerator/blob/master/docs/Whitepaper-Acceleratoring%20Throughput%20in%20Permissioned%20Blockchain%20Networks.pdf) includes:
- The key design features of Accelerator enabling high performance enterprise-wide blockchain technology
- The evaluation results that show the performance improvement of Hyperledger Fabric by Accelerator in practical scenarios
- The use cases that provide an insight for understanding industrial blockchain platforms

## Further Information
For further information please contact Samsung SDS(nexledger@samsung.com).

