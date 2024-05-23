package main

import (
	"context"
	"log"
	"time"

	pb "grpc-go-example/helloworld"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Connecting to gRPC server...")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTweetServiceClient(conn)

	// Aumenta el tiempo de espera del contexto a 10 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Sending request to gRPC server...")
	r, err := c.GetTweets(ctx, &pb.TweetRequest{TweetId: "1792994026976694390"})
	if err != nil {
		log.Fatalf("could not get tweets: %v", err)
	}

	for _, tweetText := range r.GetTweetTexts() {
		log.Printf("Tweet Text: %s", tweetText)
	}
	log.Println("Request completed.")
}
