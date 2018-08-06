package main

import (
	"crypto/rand"
	"log"
	"net"

	"github.com/0x5010/example/grpc/range/pb"
	"google.golang.org/grpc"
)

const chunkSize = 64 * 1024

type chunkerSrv []byte

func (c chunkerSrv) Range(r *pb.Res, srv pb.RangeChunker_RangeServer) error {
	chunk := &pb.Chunk{}
	ranges := c.parseRanges(r)

	for _, rr := range ranges {
		start, stop := rr[0], rr[1]
		for cur := start; cur < stop; cur += chunkSize {
			if cur+chunkSize > stop {
				chunk.Chunk = c[cur:stop]
			} else {
				chunk.Chunk = c[cur : cur+chunkSize]
			}
			if err := srv.Send(chunk); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c chunkerSrv) parseRanges(r *pb.Res) [][2]int {
	n := len(c)
	ranges := [][2]int{}
	rs := r.GetR()
	if len(rs) == 0 {
		return [][2]int{[2]int{0, n}}
	}
	for _, rr := range rs {
		start, stop := rangeLimit(rr, n)
		if start == -1 {
			return nil
		}
		ranges = append(ranges, [2]int{start, stop})
	}
	return ranges
}

func rangeLimit(r *pb.Range, llen int) (int, int) {
	start, stop := int(r.Start), int(r.Stop)+1
	if stop > llen || stop == 0 {
		stop = llen
	}
	if start < 0 || stop < 0 || start >= stop {
		return -1, -1
	}
	return start, stop
}

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	blob := make([]byte, 12*1024*1024) // 12M
	rand.Read(blob)
	pb.RegisterRangeChunkerServer(s, chunkerSrv(blob))
	log.Println("serving on localhost:8888")
	log.Fatal(s.Serve(listen))
}
