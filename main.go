package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type phoneNumber string

type config struct {
	Mac          string        `json:"mac"`
	Interface    string        `json:"interface"`
	PhoneNumbers []phoneNumber `json:"phoneNumbers"`
	TwilioSid    string        `json:"twilioSid"`
	TwilioToken  string        `json:"twilioToken"`
}

type smsSender interface {
	SendSMS(to phoneNumber, from phoneNumber, body string) error
}

func main() {
	// read into config struct
	cfgFile, err := os.Open("config.json")

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

	listenAndSendSMS(packets, cfg.Mac)
}

func listenAndSendSMS(packets <-chan gopacket.Packet, mac string) {
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
		if net.HardwareAddr(arp.SourceHwAddress).String() != mac {
			continue
		}

		// send SMS - we got an arp from the button (stubbed for now)
		log.Printf("Got a button click!")
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
