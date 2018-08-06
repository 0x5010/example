package main

import (
	"crypto/rand"
	"log"
	"net"

	"github.com/0x5010/example/grpc/chunker/pb"
	"google.golang.org/grpc"
)

const chunkSize = 64 * 1024

type chunkerSrv []byte

func (c chunkerSrv) Chunker(_ *pb.Empty, srv pb.Chunker_ChunkerServer) error {
	chunk := &pb.Chunk{}
	n := len(c)
	for cur := 0; cur < n; cur += chunkSize {
		if cur+chunkSize > n {
			chunk.Chunk = c[cur:n]
		} else {
			chunk.Chunk = c[cur : cur+chunkSize]
		}
		if err := srv.Send(chunk); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	blob := make([]byte, 128*1024*1024) // 128M
	rand.Read(blob)
	pb.RegisterChunkerServer(s, chunkerSrv(blob))
	log.Println("serving on localhost:8888")
	log.Fatal(s.Serve(listen))
}
