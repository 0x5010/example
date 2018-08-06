package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/0x5010/example/grpc/range/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewRangeChunkerClient(conn)
	var blob1, blob2 []byte
	{
		stream, err := client.Range(context.Background(), &pb.Res{
			R: []*pb.Range{
				{0, 99},
				{100, 199},
				{200, -1},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		for {
			c, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Printf("Transfer of %d bytes successful", len(blob1))
					break
				}
				log.Fatal(err)
			}
			blob1 = append(blob1, c.Chunk...)
		}
	}
	{
		stream, err := client.Range(context.Background(), &pb.Res{})
		if err != nil {
			log.Fatal(err)
		}

		for {
			c, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Printf("Transfer of %d bytes successful", len(blob2))
					break
				}
				log.Fatal(err)
			}
			blob2 = append(blob2, c.Chunk...)
		}
	}
	fmt.Println(bytes.Equal(blob1, blob2))
}
