# Nexledger Accelerator
Accelerator is a software component designed to improve the performance of a blockchain network, e.g. Hyperledger Fabric, in terms of transaction throughput. Accelerator enables the blockchain network to deal with a huge number of transaction requests from applications. 


## Getting Started
### Prerequisites
To use Accelerator, you must install the following tools in advance.
- Go 1.12+

To run the examples, Docker should be installed to run the Fabric network.
- Docker
- Docker-compose

### Building Accelerator
Accelerator supports go module for dependency management. To build the executable, please simply execute `go build`.
```bash
$ go build cmd/accelerator.go
```

## Running ping example
You can learn how to configure Accelerator by running `ping` example. The example is placed in In `examples/ping`. 

To bootstrap Fabric network, please run `start.sh` script.
```bash
$ ./examples/ping/start.sh
```

To serve requests from clients, Accelerator should be up and running with proper configuration.
```bash
$ accelerator -f examples/ping/configs/accelerator.yaml
```

Accelerator is a gRPC server and the gRPC services are described in `protos/accelerator.proto`.
You may send transactions using `examples/ping/ping_test.go` that has gRPC client for ping example. 

You can terminate and remove the network by run `stop.sh` script.
```bash
$ ./examples/ping/stop.sh
```

### Under the hood
#### Modifying chaincode
In order to apply Accelerator to your business, you need to modify your chaincode.

`contracts/src/ping/ping.go` is an example chaincode with simple KV write/read operations.
To run aggregated transactions from Accelerator individually, `ping.go` imports `batchutil.go` and delegates the funcation invocation to `Invoke()` in `batchutil.go`.
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
Whitepaper includes:
- The key design features of Accelerator enabling high performance enterprise-wide blockchain technology
- The evaluation results that show the performance improvement of Hyperledger Fabric by Accelerator in practical scenarios
- The use cases that provide an insight for understanding industrial blockchain platforms

## Further Information
For further information please contact Samsung SDS(nexledger@samsung.com).

