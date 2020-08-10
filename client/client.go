package main

import (
	"context"
	"encoding/json"
	"flag"
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
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewGrpcSampleClient(conn)
	ctx := context.Background()
	token := Token{
		Access:     "john access",
		Subject:    "john subject",
		Permission: []string{"perm1", "perm2"},
	}
	tokenBytes, _ := json.Marshal(token)
	ctx = metadata.AppendToOutgoingContext(ctx, "token", string(tokenBytes))
	result, err := client.Ping(ctx, &pb.PongIn{Message: "john"})
	if err != nil {
		panic(err)
	}
	spew.Dump(result)
}
