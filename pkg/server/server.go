/*
 *    Copyright 2019 Samsung SDS
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package server

import (
	"context"
	"fmt"
	"net"

	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/nexledger/accelerator/pkg/batch"
	"github.com/nexledger/accelerator/pkg/fabwrap"
	"github.com/nexledger/accelerator/protos"
)

type Server struct {
	host string
	port int
	ctx  fabwrap.Context
	conf *Config
	log  *zap.SugaredLogger

	client *batch.Client
	server *grpc.Server
}

func (s *Server) Serve() chan error {
	failed := make(chan error)

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
		if err != nil {
			s.log.Error("Failed to listen for gRPC server.", err)
			failed <- err
		}

		s.client, err = s.conf.BatchClient()
		if err != nil {
			s.log.Error("Failed to Create BatchClient.", err)
			failed <- err
		}

		protos.RegisterAcceleratorServiceServer(s.server, s)
		s.log.Info("Starting gRPC server.")
		if err := s.server.Serve(listener); err != nil {
			s.log.Error("gRPC server is terminated.", err)
			failed <- err
		}
	}()
	return failed
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

func (s *Server) Execute(ctx context.Context, req *protos.TxRequest) (*protos.TxResponse, error) {
	result, err := s.client.Execute(req.ChannelId, req.ChaincodeName, req.Fcn, req.Args)
	if err != nil {
		return nil, err
	}

	return &protos.TxResponse{
		TxId:       result.TxId,
		Validation: &protos.TransactionValidation{Code: result.ValidationCode, Description: peer.TxValidationCode_name[result.ValidationCode]},
		Payload:    result.Payload,
	}, nil
}

func (s *Server) Query(ctx context.Context, req *protos.TxRequest) (*protos.TxResponse, error) {
	result, err := s.client.Query(req.ChannelId, req.ChaincodeName, req.Fcn, req.Args)
	if err != nil {
		return nil, err
	}

	return &protos.TxResponse{
		TxId:       result.TxId,
		Validation: &protos.TransactionValidation{Code: result.ValidationCode, Description: peer.TxValidationCode_name[result.ValidationCode]},
		Payload:    result.Payload,
	}, nil
}

func New(configPath string) (*Server, error) {
	conf, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	var server *grpc.Server
	if len(conf.Tls.Certpath) > 0 && len(conf.Tls.KeyPath) > 0 {
		credential, err := credentials.NewServerTLSFromFile(conf.Tls.Certpath, conf.Tls.KeyPath)
		if err != nil {
			return nil, err
		}
		server = grpc.NewServer(grpc.MaxConcurrentStreams(0), grpc.Creds(credential))
	} else {
		server = grpc.NewServer(grpc.MaxConcurrentStreams(0))
	}

	return &Server{
		host:   conf.Host,
		port:   conf.Port,
		conf:   conf,
		log:    zap.S(),
		server: server,
	}, nil
}
