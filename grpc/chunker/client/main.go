package main

import (
	"context"
	"io"
	"log"

	"github.com/0x5010/example/grpc/chunker/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewChunkerClient(conn)
	stream, err := client.Chunker(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	var blob []byte
	for {
		c, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Printf("Transfer of %d bytes successful", len(blob))
				return
			}
			log.Fatal(err)
		}
		blob = append(blob, c.Chunk...)
	}
}
