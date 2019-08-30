package nebula

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/nebula-api/nebulapb"
	"google.golang.org/grpc"
)

var cConn *grpc.ClientConn
var client nebulapb.NebulaClient

func NewClient() {
	address := os.Getenv("NEBULA_ADDRESS")
	if len(address) == 0 {
		address = "localhost:17200"
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Fatalf("[Nebula] Failed connect to Nebula-API")
		return
	}
	logrus.Printf("[Nebula] Connected Nebula-API")
	//defer conn.Close()

	cConn = conn
	client = nebulapb.NewNebulaClient(conn)
}

func Shutdown() {
	cConn.Close()
}

// GetServerEntry - Get All Servers
func GetServerEntry() ([]*nebulapb.ServerEntry, error) {
	logrus.Printf("[Server] Getting...")
	v, err := client.GetServerEntry(context.Background(), &nebulapb.GetServerEntryRequest{})
	return v.Entry, err
}

// AddServerEntry - Add Server to Database
func AddServerEntry(name, displayName, address string, port int32, fallback bool) error {
	logrus.Printf("[Server] Request..")
	e := nebulapb.ServerEntry{
		Name:        name,
		DisplayName: displayName,
		Address:     address,
		Port:        port,
		Motd:        "",
		Fallback:    fallback,
	}
	_, err := client.AddServerEntry(context.Background(), &nebulapb.AddServerEntryRequest{Entry: &e})

	return err
}

func RemoveServerEntry(name string) error {
	logrus.Printf("[Server] Remove....")
	_, err := client.RemoveServerEntry(context.Background(), &nebulapb.RemoveServerEntryRequest{Name: name})

	return err
}

func SetMotd(motd string) error {
	_, err := client.SetMotd(context.Background(), &nebulapb.SetMotdRequest{Motd: motd})

	return err
}

func SetFavicon(favicon string) error {
	_, err := client.SetFavicon(context.Background(), &nebulapb.SetFaviconRequest{Favicon: favicon})

	return err
}
