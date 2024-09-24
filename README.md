# MITMagic: Network Packet Sniffer & ARP Spoofing Tool

MITMagic is a Go-based tool designed to perform network packet sniffing and ARP spoofing. It captures network traffic and identifies common protocols such as HTTP, HTTPS, DNS, and MDNS. The tool can be used for network security auditing, traffic monitoring, and understanding how devices communicate over a network.

## Features

- **ARP Spoofing**: Perform ARP spoofing on the target device and gateway to capture traffic.
- **Packet Sniffing**: Capture and analyze packets in real-time.
- **Protocol Identification**: Classify captured traffic into common protocols such as:
  - HTTP
  - HTTPS
  - DNS
  - MDNS
  - TCP/UDP with port information
- **Color-Coded Output**: Display captured packets with a clean, color-coded format for easy readability.
- **Hostname Resolution**: Resolve IP addresses to their corresponding domain names (if available).

## Example Output

```
[2024-09-24T03:09:07-04:00] » 192.168.1.91 > https://unn-149-102-226-193.datapacket.com. (149.102.226.193) [HTTPS]
[2024-09-24T03:09:07-04:00] » 192.168.1.251 > WEBDADDY.attlocal.net. (192.168.1.91) [DNS]
```


## Prerequisites

- **Go**: Make sure you have Go installed. You can download it [here](https://golang.org/doc/install).
- **ARP-Scan**: Used for network device discovery.
- **Arpspoof**: Tool for ARP spoofing, which is part of the `dsniff` package.
  
  On Debian-based systems (e.g., Ubuntu), you can install these tools with:

```bash
sudo apt-get install arp-scan dsniff
```
- libpcap: Go's gopacket package requires libpcap for packet capturing. Install it using:
```bash
sudo apt-get install libpcap-dev
```

# Installation

1. Clone the repository:
```bash
git clone https://github.com/blackmagic2023/MITMagic.git
```
2. Change into the project directory:
```bash
cd MITMagic
```
3. Install Go dependencies:
```bash
go get -u github.com/google/gopacket
go get -u github.com/google/gopacket/pcap
```
4. Build the program:
```bash
go build -o MITMagic MITMagic.go
```

# Usage

1. Run the tool as root: To capture traffic and perform ARP spoofing, root privileges are required. You can run the program using `sudo:`
```bash
sudo ./MITMagic
```
2. Select the network interface: When prompted, enter the interface (e.g., `eth0`, `wlan0`) you wish to monitor.
```
Enter the network interface (e.g., eth0, wlan0): eth0
```
3. Probe the network: The tool will scan the network and display a list of connected devices. Choose a target from this list and provide the gateway IP.
```bash
Enter the target IP address from the list above: 192.168.1.100
Enter the gateway IP address: 192.168.1.1
```
4. Packet capturing and ARP spoofing: MITMagic will start capturing packets, displaying them with protocol labels and color-coded for easy readability.

# Protocols Identified

- HTTP: Detected on TCP port 80.
- HTTPS: Detected on TCP port 443 (shown with https:// prefix).
- DNS: Detected on UDP port 53.
- MDNS: Detected on UDP port 5353.
- Other TCP/UDP protocols: Detected by their port numbers.

# Troubleshooting

- ARP Spoofing Hanging: If the ARP spoofing process hangs, ensure you have the correct target and gateway IPs, and that the `arpspoof` tool is installed.
- Permission Denied: If you encounter permission issues, try running the program with `sudo`.
- Packet Loss: If you're missing packets or not seeing all traffic, try adjusting your capture filter or increasing the buffer size by modifying the program.

# Contribution

Feel free to contribute by submitting issues, suggesting features, or making pull requests. Here’s how you can contribute:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.

# License

MITMagic is licensed under the MIT License.
