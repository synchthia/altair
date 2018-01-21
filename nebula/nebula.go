package nebula

import (
	"context"

	"github.com/sirupsen/logrus"
	pb "gitlab.com/Startail/Nebula-API/nebulapb"
	"google.golang.org/grpc"
)

var cConn *grpc.ClientConn
var client pb.NebulaClient

func NewClient() {
	conn, err := grpc.Dial("localhost:17200", grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("can not connect: %v", err)
	}
	//defer conn.Close()

	cConn = conn
	client = pb.NewNebulaClient(conn)
}

func Shutdown() {
	cConn.Close()
}

// GetServerEntry - Get All Servers
func GetServerEntry() ([]*pb.ServerEntry, error) {
	logrus.Printf("[Server] Getting...")
	v, err := client.GetServerEntry(context.Background(), &pb.GetServerEntryRequest{})
	return v.Entry, err
}

// AddServerEntry - Add Server to Database
func AddServerEntry(name, displayName, address string, port int32) error {
	logrus.Printf("[Server] Request..")
	e := pb.ServerEntry{
		Name:        name,
		DisplayName: displayName,
		Address:     address,
		Port:        port,
		Motd:        "",
	}
	_, err := client.AddServerEntry(context.Background(), &pb.AddServerEntryRequest{Entry: &e})

	return err
}

func RemoveServerEntry(name string) error {
	logrus.Printf("[Server] Remove....")
	_, err := client.RemoveServerEntry(context.Background(), &pb.RemoveServerEntryRequest{Name: name})

	return err
}
