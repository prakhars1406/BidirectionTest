package main

import (
	"BidirectionalService/protos"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type CheckUserServices struct {
}

func (r *CheckUserServices) CheckUserService(stream protos.CheckUserService_CheckUserServiceServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		serviceName := req.GetServiceName()
		fmt.Printf("Received a new serviceName of...: %s\n", serviceName)
		sendErr := stream.Send(&protos.CheckUserServiceResponse{
			Message: serviceName + " subscribed",
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func main()  {
	fmt.Println("Go grpc server!")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	checkUserServices := CheckUserServices{}
	protos.RegisterCheckUserServiceServer(grpcServer, &checkUserServices)
	fmt.Println("Starting Server...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
