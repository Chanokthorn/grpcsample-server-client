package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
)

func someMiddleware1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (interface{}, error) {
	newCtx := context.WithValue(ctx, "user_id", "john@example.com")
	fmt.Println("pre: middleware1")
	res, err := next(newCtx, req)
	fmt.Println("post: middleware1")
	return res, err
}


func someMiddleware2(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (interface{}, error) {
	newCtx := context.WithValue(ctx, "user_id", "john@example.com")
	fmt.Println("pre: middleware2")
	res, err := next(newCtx, req)
	fmt.Println("post: middleware2")
	return res, err
}