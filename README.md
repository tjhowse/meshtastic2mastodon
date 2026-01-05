# Meshastic2mastodon

This talks to an MQTT broker that a meshtastic node is also talking to. It collects a map of node ID to username, and relays any messages published to channel 0 to a Mastodon instance. Set the following environment variables, run the binary and you're away. The message format is hardcoded, there are serious warts all over this. Use at your own risk.

    MASTODON_SERVER
    MASTODON_CLIENT_ID
    MASTODON_CLIENT_SECRET
    MASTODON_ACCESS_TOKEN
    MQTT_BROKER_HOSTNAME
    MQTT_BROKER_PORT
    MQTT_USERNAME
    MQTT_PASSWORD