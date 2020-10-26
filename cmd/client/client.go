package client

import (
	"context"
	"flag"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	v1 "test/grpc/api/proto"
	"time"
)

const (
	apiVersion = "v1"
)

func main() {
	address := flag.String("server", "", "gRPC server in fotmat host:port")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if nil != err {
		log.Fatal("服务器连不上：%v", err)
	}
	defer conn.Close()

	c := v1.NewToDoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)

	remider, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)
	req1 := v1.CreateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Id:          0,
			Title:       "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
		},
	}

	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatal("创建失败")
	}
	log.Printf("Create result%v", res1)
	id := res1.Id

	req2 := v1.ReadRequest{Api: apiVersion, Id: id}
	res2, err := c.Read(ctx, &req2)
	if nil != err {
		log.Fatal("Read failed %v", err)
	}
	log.Printf("Read result %v", res2)

	req3 := v1.UpdateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Id:          res2.ToDo.Id,
			Title:       res2.ToDo.Title,
			Description: res2.ToDo.Description + "updated",
		},
	}
	res3, err := c.Update(ctx, &req3)
	if nil != err {
		log.Fatal("Update result %v", res3)
	}
	log.Printf("Update result %v", res3)

}
