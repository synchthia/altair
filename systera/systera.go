package systera

import (
	"context"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/synchthia/systera-api/systerapb"
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
	logrus.WithFields(logrus.Fields{
		"address": address,
	}).Debugf("[Systera] Connecting to Systera-API")

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

// GetProfile - GetProfile UUID or String
func GetProfile(uuidOrString string) (systerapb.PlayerEntry, error) {
	if len(uuidOrString) == 32 {
		r, err := FetchPlayerProfile(uuidOrString)
		if err != nil {
			logrus.WithError(err).Errorf("Failed Lookup Player's Profile(from UUID)")
			return systerapb.PlayerEntry{}, err
		}
		return *r, nil
	}
	r, err := FetchPlayerProfileByName(uuidOrString)
	if err != nil {
		logrus.WithError(err).Errorf("Failed Lookup Player's Profile(from Name)")
		return systerapb.PlayerEntry{}, err
	}
	return *r, nil
}

// FetchPlayerProfileByName - Get Player's Profile By Name
func FetchPlayerProfileByName(playerName string) (*systerapb.PlayerEntry, error) {
	r, err := client.FetchPlayerProfileByName(
		context.Background(),
		&systerapb.FetchPlayerProfileByNameRequest{
			Name: playerName,
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
			Uuid: playerUUID,
		},
	)
	if r == nil {
		return nil, errors.New("player not found")
	}
	return r.Entry, err
}

// ----------------
// AltLookup
// ----------------

// AltLookup - Lookup Player's Alternative Accounts
func AltLookup(playerUUID string) ([]*systerapb.AltLookupEntry, error) {
	r, err := client.AltLookup(
		context.Background(),
		&systerapb.AltLookupRequest{
			PlayerUuid: playerUUID,
		},
	)
	if r == nil {
		return nil, err
	}
	return r.Entries, err
}

// ----------------
// Group
// ----------------

// GetGroups - get groups
func GetGroups(server string) (*systerapb.FetchGroupsResponse, error) {
	r, err := client.FetchGroups(context.Background(), &systerapb.FetchGroupsRequest{
		ServerName: server,
	})
	return r, err
}

// SetGroup - set players group
func SetGroup(playerUUID string, groups []string) error {
	_, err := client.SetPlayerGroups(context.Background(), &systerapb.SetPlayerGroupsRequest{
		Uuid:   playerUUID,
		Groups: groups,
	})
	return err
}

// CreateGroup - create group
func CreateGroup(name, prefix string) error {
	_, err := client.CreateGroup(context.Background(), &systerapb.CreateGroupRequest{
		GroupName:   name,
		GroupPrefix: prefix,
	})
	return err
}

// RemoveGroup - remove group
func RemoveGroup(name string) error {
	_, err := client.RemoveGroup(context.Background(), &systerapb.RemoveGroupRequest{
		GroupName: name,
	})
	return err
}

// ------------
// Permissions
// ------------
// AddPermission - add permission
func AddPermission(name, target string, permissions []string) error {
	_, err := client.AddPermission(context.Background(), &systerapb.AddPermissionRequest{
		GroupName:   name,
		Target:      target,
		Permissions: permissions,
	})
	return err
}

// RemovePermission - remove permission
func RemovePermission(name, target string, permissions []string) error {
	_, err := client.RemovePermission(context.Background(), &systerapb.RemovePermissionRequest{
		GroupName:   name,
		Target:      target,
		Permissions: permissions,
	})
	return err
}

// ------------
// PUNISH
// ------------

// LookupPunish - Lookup player's Punishments
func LookupPunish(playerUUID string) ([]*systerapb.PunishEntry, error) {
	r, err := client.GetPlayerPunish(context.Background(), &systerapb.GetPlayerPunishRequest{
		Uuid:           playerUUID,
		FilterLevel:    0,
		IncludeExpired: true,
	})
	return r.Entry, err
}

// Punish - Punish the Player
func Punish(playerUUID, reason string) error {
	_, err := client.SetPlayerPunish(context.Background(), &systerapb.SetPlayerPunishRequest{
		Force: true,
		Entry: &systerapb.PunishEntry{
			Available: true,
			Level:     systerapb.PunishLevel_WARN,
			Reason:    reason,
			PunishedFrom: &systerapb.PlayerIdentity{
				Uuid: "",
				Name: "Console",
			},
			PunishedTo: &systerapb.PlayerIdentity{
				Uuid: playerUUID,
				Name: "",
			},
		},
	})
	return err
}
