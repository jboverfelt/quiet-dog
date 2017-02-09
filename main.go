package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type phoneNumber string

type config struct {
	Mac          string        `json:"mac"`
	Interface    string        `json:"interface"`
	PhoneNumbers []phoneNumber `json:"phoneNumbers"`
	TwilioNumber phoneNumber   `json:"twilioNumber"`
	TwilioSid    string        `json:"twilioSid"`
	TwilioToken  string        `json:"twilioToken"`
}

const baseURL = "https://api.twilio.com/2010-04-01/Accounts/"

func main() {
	// read into config struct
	cfgFile, err := os.Open("config.json")
	defer cfgFile.Close()

	if err != nil {
		log.Fatalf("Could not open config.json file: %v", err)
	}

	var cfg config
	err = json.NewDecoder(cfgFile).Decode(&cfg)

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	packets, err := setupPacketSource(cfg.Interface)

	if err != nil {
		log.Fatalf("error on packet source setup: %v", err)
	}

	sender := twilioSmsSender{
		baseURL: baseURL,
		sid:     cfg.TwilioSid,
		token:   cfg.TwilioToken,
	}

	listenAndSendSMS(packets, cfg, sender)
}

func listenAndSendSMS(packets <-chan gopacket.Packet, cfg config, sender smsSender) {
	var packet gopacket.Packet
	for {
		// block until a packet is received
		packet = <-packets

		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer == nil {
			continue
		}

		arp := arpLayer.(*layers.ARP)

		// discard the packet unless it comes from the dash button
		if net.HardwareAddr(arp.SourceHwAddress).String() != cfg.Mac {
			continue
		}

		log.Println("ring! sending sms")

		err := sender.SendSMS(cfg.PhoneNumbers, cfg.TwilioNumber, strings.NewReader("Someone is at the door!"))

		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func setupPacketSource(ifaceName string) (<-chan gopacket.Packet, error) {
	handle, err := pcap.OpenLive(ifaceName, 65536, true, pcap.BlockForever)

	if err != nil {
		return nil, err
	}

	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)

	return src.Packets(), nil
}
