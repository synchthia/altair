package systera

import (
	"context"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/Systera-API/systerapb"
	"google.golang.org/grpc"
)

var cConn *grpc.ClientConn
var client systerapb.SysteraClient

func NewClient() {
	address := os.Getenv("SYSTERA_ADDRESS")
	if len(address) == 0 {
		address = "localhost:17300"
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Fatalf("[Systera] Failed connect to Systera-API")
		return
	}
	logrus.Printf("[Systera] Connected Systera-API")
	//defer conn.Close()

	cConn = conn
	client = systerapb.NewSysteraClient(conn)
}

func Shutdown() {
	cConn.Close()
}

// -------------
// SYSTEM
// -------------

// Dispatch - Dispatch Command to Server
func Dispatch(target, command string) error {
	_, err := client.Dispatch(context.Background(), &systerapb.DispatchRequest{Target: target, Cmd: command})
	return err
}

// Announce - Announce to Server
func Announce(target, message string) error {
	_, err := client.Announce(context.Background(), &systerapb.AnnounceRequest{Target: target, Message: message})
	return err
}

// ----------------
// Player
// ----------------

// FetchPlayerProfileByName - Get Player's Profile By Name
func FetchPlayerProfileByName(playerName string) (*systerapb.PlayerEntry, error) {
	r, err := client.FetchPlayerProfileByName(
		context.Background(),
		&systerapb.FetchPlayerProfileByNameRequest{
			PlayerName: playerName,
		},
	)
	if r == nil {
		return nil, errors.New("player not found")
	}
	return r.Entry, err
}

// FetchPlayerProfile - Get Player's Profile
func FetchPlayerProfile(playerUUID string) (*systerapb.PlayerEntry, error) {
	r, err := client.FetchPlayerProfile(
		context.Background(),
		&systerapb.FetchPlayerProfileRequest{
			PlayerUUID: playerUUID,
		},
	)
	if r == nil {
		return nil, errors.New("player not found")
	}
	return r.Entry, err
}

// ----------------
// Group
// ----------------

// SetGroup - set players group
func SetGroup(playerUUID string, groups []string) error {
	_, err := client.SetPlayerGroups(context.Background(), &systerapb.SetPlayerGroupsRequest{
		PlayerUUID: playerUUID,
		Groups:     groups,
	})
	return err
}

// ------------
// PUNISH
// ------------

// LookupPunish - Lookup player's Punishments
func LookupPunish(playerUUID string) ([]*systerapb.PunishEntry, error) {
	r, err := client.GetPlayerPunish(context.Background(), &systerapb.GetPlayerPunishRequest{
		PlayerUUID:     playerUUID,
		FilterLevel:    0,
		IncludeExpired: true,
	})
	return r.Entry, err
}
