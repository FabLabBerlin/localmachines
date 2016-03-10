/*
Commandline tool for automatic configuration of mfi switches.

TODO: use github.com/ziutek/telnet
      (golang.org/x/crypto/ssh doesn't support cbc ciphers because they're
      unsafe)
      use parser for hwaddr to check for equality

*/
package main

import (
	"flag"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/mfi"
	"log"
)

func usage(ethernetConfig bool) {
	fmt.Printf("Please hold the power switch's reset button until it\n")
	fmt.Printf("blinks blue-orange.  ")
	if ethernetConfig {
		fmt.Printf("Be sure to be connected with\n")
		fmt.Printf("the device with an Ethernet cable.\n")
		fmt.Printf("Be on 192.168.1.x (The device is now on 192.168.1.20)\n")
	} else {
		fmt.Printf("Afterwards wait until a Wifi\n")
		fmt.Printf("named 'mfi...' comes up and connect to it.\n")
		fmt.Printf("(The device is now on 192.168.2.20)\n")
	}
}

func main() {
	ethernetConfig := flag.Bool("ethConfig", false, "Connect Switch to computer via Ethernet Cable (be on 192.168.1.x)")
	netmask := flag.String("network", "172.26.0.128/24", "Wifi network's netmask on which we auto discover the switch")
	ssid := flag.String("ssid", "FabLab-M", "Wifi SSID")
	wifiPwManually := flag.Bool("wifiPwManual", false, "Enter Wifi password manually")
	flag.Parse()

	c := &mfi.Config{
		EthernetConfig: *ethernetConfig,
		WifiSSID:       *ssid,
	}

	if *wifiPwManually {
		fmt.Printf("What is the Wifi password?\n")
		if _, err := fmt.Scanln(&c.WifiPassword); err != nil {
			log.Fatalf("scan ln: %v", err)
		}
	}

	if err := c.Run(); err == nil {
		fmt.Printf("Your switch is properly configured and its hardware")
		fmt.Printf(" address is: '%v'\n", c.HwAddr)
		fmt.Printf("Wait until the LED starts blinking blue, connect to the")
		fmt.Printf(" normal Wifi and then press enter...\n")
		var tmp string
		fmt.Scanln(&tmp)
		if ip, err := c.FindDeviceOn(*netmask); err == nil {
			fmt.Printf("Successfully discovered device: %v\n", ip.String())
		} else {
			fmt.Printf("Unable to find device on %v\n.", *netmask)
		}
	} else if err == mfi.ErrWifiPasswordNotPresent {
		fmt.Printf("Wifi password not present, either configure via cmd line\n")
		fmt.Printf("or use the automatic Ubiquitiy web configuration.\n")
		flag.Usage()
	} else {
		log.Printf("config: %v", err)
		fmt.Printf("Make sure the switch is properly connected:\n\n")
		usage(c.EthernetConfig)
	}
}
