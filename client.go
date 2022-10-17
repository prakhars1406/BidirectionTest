package main

import (
	"BidirectionalService/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	checkUserServiceClient := protos.NewCheckUserServiceClient(conn)
	stream,err:=checkUserServiceClient.CheckUserService(context.Background())
	go func(stream protos.CheckUserService_CheckUserServiceClient){
		for i:=0;i<100000;i++{
			stream.Send(&protos.CheckUserServiceRequest{
				ServiceName: "service"+strconv.Itoa(i),
			})
		}
	}(stream)
	go func(stream protos.CheckUserService_CheckUserServiceClient){
		for{
			response,err:=stream.Recv()
			if err!=nil{
				fmt.Println("error in receiving "+err.Error())
			}else{
				fmt.Println("response received "+response.Message)
			}
		}
	}(stream)

	b1:=make(chan bool)
	<-b1
}
