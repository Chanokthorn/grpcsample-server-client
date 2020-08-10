package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc/metadata"
	"log"
	"net"

	pb "github.com/Chanokthorn/grpcsample"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	pb2 "grpc-sample/grpcsample2"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")
)

type Token struct {
	Access     string   `json:"access"`
	Subject    string   `json:"subject"`
	Permission []string `json:"permission"`
}

type grpcSampleServer struct {
	pb.UnimplementedGrpcSampleServer
	pb2.UnimplementedGrpcSample2Server
}

func (g *grpcSampleServer) Ping(ctx context.Context, pongIn *pb.PongIn) (*pb.PongOut, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("cannot retrieve metadata from context")
	}
	tokenStrings := md.Get("token")
	if len(tokenStrings) != 1 {
		return nil, errors.New("invalid token")
	}
	fmt.Println("token:")
	var token Token
	err := json.Unmarshal([]byte(tokenStrings[0]), &token)
	if err != nil {
		return nil, err
	}
	spew.Dump(token)
	return &pb.PongOut{Message: pongIn.Message}, nil
}

func newServer() pb.GrpcSampleServer {
	return &grpcSampleServer{}
}

func newServer2() pb2.GrpcSample2Server {
	return &grpcSampleServer{}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(someMiddleware1, someMiddleware2)))
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGrpcSampleServer(grpcServer, newServer())
	pb2.RegisterGrpcSample2Server(grpcServer, newServer2())
	grpcServer.Serve(lis)
}
