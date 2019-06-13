package main

import (
	"bytes"
	"encoding/gob"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func invoke(stub shim.ChaincodeStubInterface, target func(shim.ChaincodeStubInterface, []string) pb.Response) pb.Response {
	jobs := make([][][]byte, 0)
	if err := decode(stub.GetArgs()[1], &jobs); err != nil {
		return shim.Error("Failed to unmarshal request")
	}

	results := make([]pb.Response, 0, len(jobs))
	for _, job := range jobs {
		args := make([]string, 0, len(job))
		for _, arg := range job {
			args = append(args, string(arg))
		}

		results = append(results, target(stub, args))
	}

	response, err := encode(results)
	if err != nil {
		return shim.Error("Failed to marshal response")
	}
	return shim.Success(response)
}

func encode(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(d []byte, v interface{}) error {
	buf := bytes.NewBuffer(d)
	if err := gob.NewDecoder(buf).Decode(v); err != nil {
		return err
	}
	return nil
}
