package systera

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"gitlab.com/Startail/Systera-API/systerapb"
	"google.golang.org/grpc"
)

var cConn *grpc.ClientConn
var client systerapb.SysteraClient

func NewClient() {
	//conn, err := grpc.Dial("argon.synchthia.net:17300", grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:17300", grpc.WithInsecure())
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

// Execute - Dispatch Command to Server
func Execute(target, command string) {
	client.Dispatch(context.Background(), &systerapb.DispatchRequest{Target: target, Cmd: command})
}

// Announce - Announce to Server
func Announce(target, message string) {
	client.Announce(context.Background(), &systerapb.AnnounceRequest{Target: target, Message: message})
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
