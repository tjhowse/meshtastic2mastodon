package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tjhowse/meshtastic2mastodon/protobufs/generated"
)

// This code handles storing disk-persistent data to a sqlite table.

// // This stores the detail of a single find event.
// type MeshtasticUserInfo struct {
// 	gorm.Model
// 	Name      string
// 	FindTime  time.Time
// 	FindType  string
// 	CacheCode string
// 	LogString string
// }

type userInfo struct {
	Id        string `json:"id,omitempty"`
	LongName  string `json:"long_name,omitempty"`
	ShortName string `json:"short_name,omitempty"`
}

func userInfoFromProto(user *generated.User) userInfo {
	var u userInfo
	u.Id = user.Id
	u.LongName = user.LongName
	u.ShortName = user.ShortName
	return u
}

type userInfoCache map[uint32]userInfo

func userIdToNodeId(userId string) uint32 {
	// Convert user ID string to uint32 node ID.
	// Strip the "!" prefix off and convert from hex to uint32.
	if len(userId) > 0 && userId[0] == '!' {
		userId = userId[1:]
	}
	var nodeId uint32
	fmt.Sscanf(userId, "%x", &nodeId)
	return nodeId
}

func loadUserInfo() userInfoCache {
	// This loads user info from user.json file into a map.
	userMap := make(userInfoCache)
	u, err := os.ReadFile("user.json")
	if err != nil {
		fmt.Printf("Error reading user.json: %v\n", err)
		return userMap
	}
	var users []userInfo
	if err := json.Unmarshal(u, &users); err != nil {
		fmt.Printf("Error unmarshaling user.json: %v\n", err)
		return userMap
	}
	for _, user := range users {
		userMap[userIdToNodeId(user.Id)] = user
	}
	return userMap
}

func saveUserInfo(userMap userInfoCache) {
	// This saves user info from a map into user.json file.
	var users []userInfo
	for _, user := range userMap {
		users = append(users, user)
	}
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling user info: %v\n", err)
		return
	}
	if err := os.WriteFile("user.json", data, 0644); err != nil {
		fmt.Printf("Error writing user.json: %v\n", err)
		return
	}
}
