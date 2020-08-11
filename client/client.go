package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	pb "github.com/Chanokthorn/grpcsample"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/testdata"
	"log"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

type Token struct {
	Access     string   `json:"access"`
	Subject    string   `json:"subject"`
	Permission []string `json:"permission"`
}

func grpcClientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	token := Token{
		Access:     "john access",
		Subject:    "john subject",
		Permission: []string{"perm1", "perm2"},
	}
	tokenString, err := json.Marshal(&token)
	if err != nil {
		return err
	}
	newCtx := metadata.AppendToOutgoingContext(ctx, "token", string(tokenString))
	err = invoker(newCtx, method, req, reply, cc, opts...)
	return err
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithUnaryInterceptor(grpcClientInterceptor))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewGrpcSampleClient(conn)
	ctx := context.Background()
	result, err := client.Ping(ctx, &pb.PongIn{Message: "john"})
	if err != nil {
		fmt.Print("error from john")
	}
	spew.Dump(result)
}
