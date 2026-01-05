package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tjhowse/meshtastic2mastodon/protobufs/generated"
	mgc "github.com/tjhowse/mqttgochannels"
	"google.golang.org/protobuf/proto"
)

func main() {
	fullPath := fmt.Sprintf("tcp://%s:%s", os.Getenv("MQTT_BROKER_HOSTNAME"), os.Getenv("MQTT_BROKER_PORT"))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fullPath)
	opts.SetClientID("testing")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetAutoReconnect(true)
	opts.SetResumeSubs(true)
	opts.SetCleanSession(false)

	masto, err := NewMastodon(os.Getenv("MASTODON_SERVER"), os.Getenv("MASTODON_CLIENT_ID"), os.Getenv("MASTODON_CLIENT_SECRET"), os.Getenv("MASTODON_ACCESS_TOKEN"))
	if err != nil {
		log.Fatal("Error creating Mastodon client:", err)
	}

	client := mgc.NewMQTTgoChannels(mqtt.NewClient(opts))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Error connecting to MQTT broker:", token.Error())
	}

	msgChannel, err := client.SubscribeGetChannel("msh/ANZ/2/e/#", 0)
	if err != nil {
		fmt.Println("Error subscribing to topic:", err)
		return
	}

	userData := loadUserInfo()

	go func() {
		for msg := range msgChannel {
			packet := &generated.ServiceEnvelope{}
			if err := proto.Unmarshal(msg.Payload(), packet); err != nil {
				fmt.Printf("Error decoding Meshtastic packet: %v\n", err)
				continue
			}
			if packet.Packet.PkiEncrypted {
				// Encrypted packet, skip
				continue
			}
			decoded, ok := packet.Packet.PayloadVariant.(*generated.MeshPacket_Decoded)
			if !ok {
				continue
			}
			switch decoded.Decoded.Portnum {
			case generated.PortNum_TEXT_MESSAGE_APP:
				textMsg := string(decoded.Decoded.Payload)
				longName, exists := userData[packet.Packet.From]
				if exists {
					textMsg = fmt.Sprintf("%s: %s", longName.LongName, textMsg)
				} else {
					textMsg = fmt.Sprintf("!%x: %s", packet.Packet.From, textMsg)
				}
				fmt.Printf("Channel %d: %s\n", packet.Packet.Channel, textMsg)
				if packet.Packet.Channel != 0 {
					// Not the main channel, don't post it to mastodon.
					continue
				}
				// Prepend "MediumFast: " to the message
				textMsg = "MediumFast: " + textMsg
				// Post to Mastodon
				if err := masto.PostStatus(textMsg); err != nil {
					fmt.Printf("Error posting to Mastodon: %v\n", err)
				} else {
					fmt.Println("Posted to Mastodon successfully.")
				}
			case generated.PortNum_NODEINFO_APP:
				userInfo := &generated.User{}
				if err := proto.Unmarshal(decoded.Decoded.Payload, userInfo); err != nil {
					fmt.Printf("Error decoding NodeInfo payload: %v\n", err)
					continue
				}

				fmt.Printf("Updated user info: %x (%s)is %s\n", packet.Packet.From, userInfo.Id, userInfo.LongName)
				userData[userIdToNodeId(userInfo.Id)] = userInfoFromProto(userInfo)
				saveUserInfo(userData)
			default:
				continue
			}
			// if decoded.Decoded.Portnum != generated.PortNum_TEXT_MESSAGE_APP {
			// 	// Not a text message, skip
			// 	continue
			// }
		}
	}()

	select {}

}
