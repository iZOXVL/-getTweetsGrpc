package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	pb "grpc-go-example/helloworld"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTweetServiceServer
}

type ApiResponse struct {
	Replies []Tweet `json:"replies"`
}

type Tweet struct {
	Text string `json:"text"`
}

func (s *server) GetTweets(ctx context.Context, req *pb.TweetRequest) (*pb.TweetTextResponse, error) {
	log.Printf("Received request for Tweet ID: %s", req.GetTweetId())

	url := fmt.Sprintf("https://twitter154.p.rapidapi.com/tweet/replies?tweet_id=%s", req.GetTweetId())
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("X-RapidAPI-Key", "b4f37cd7a4msh1f85c44bfe216d1p11fce5jsnfe8d5879e5f5")
	request.Header.Add("X-RapidAPI-Host", "twitter154.p.rapidapi.com")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Error making request to Twitter API: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	log.Printf("Response from Twitter API: %s", string(body))

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, err
	}

	var tweetTexts []string
	for _, t := range apiResponse.Replies {
		tweetTexts = append(tweetTexts, t.Text)
	}

	log.Printf("Returning %d tweet texts", len(tweetTexts))
	return &pb.TweetTextResponse{TweetTexts: tweetTexts}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTweetServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
