package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/roguesoftware/tla-proto"
)

const port = ":50507"

var initialVotes []*pb.VoteItem

type server struct {
	pb.UnimplementedVoteServiceServer
}

func (s *server) GetVotes(ctx context.Context, in *pb.VoteRequest) (*pb.VoteReply, error) {
	contextID := in.GetContextId()
	userID := in.GetUserId()

	log.Printf("Received: %v %v", contextID, userID)

	var votes []*pb.VoteItem
	votes = initialVotes[0:1]

	return &pb.VoteReply{Votes: votes}, nil
}

func main() {
	// load initial votes
	fileName := "votes.json"
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening %v: %v", fileName, err)
	}
	defer jsonFile.Close()

	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonBytes, &initialVotes)

	// create listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// start server
	s := grpc.NewServer()
	pb.RegisterVoteServiceServer(s, &server{})
	log.Printf("Registered vote server with %v votes", len(initialVotes))
	if err := s.Serve(lis); err != nil {
		log.Fatal(s.Serve(lis))
	}
}
