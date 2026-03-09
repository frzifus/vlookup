# vlookup

[![Build Status](https://github.com/frzifus/vlookup/workflows/Go/badge.svg)](https://github.com/frzifus/vlookup/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/frzifus/vlookup)](https://github.com/frzifus/vlookup/blob/master/go.mod)
[![License](https://img.shields.io/github/license/frzifus/vlookup)](https://github.com/frzifus/vlookup/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/frzifus/vlookup.svg)](https://github.com/frzifus/vlookup/releases/)

**Network device discovery and MAC address vendor lookup via ARP**

## Overview

vlookup is a command-line tool for discovering devices on local networks and identifying their manufacturers. It performs ARP (Address Resolution Protocol) scans to find active devices, then matches their MAC addresses against IEEE vendor databases to show which organization manufactured each device.

Designed for network administrators, security professionals, and IT support staff, vlookup helps you:
- Inventory all devices on your network segments
- Identify unauthorized or unknown devices
- Audit network security posture
- Troubleshoot connectivity issues
- Monitor network topology changes

The tool includes embedded IEEE OUI (Organizationally Unique Identifier) databases containing over 2.7 million vendor entries, with utilities to download the latest database updates.

## Features

- **Active ARP Scanning**: Discover all devices on your network by sending ARP probes
- **Passive Mode**: Query system ARP cache without active scanning (no root required)
- **MAC Vendor Lookup**: Match MAC addresses to manufacturers using IEEE databases
- **Multi-Interface Support**: Scan all network interfaces or filter to specific ones
- **Embedded Databases**: Includes small, medium, and large IEEE vendor databases
- **Database Updates**: Download latest IEEE databases with included crawler utility
- **Custom Databases**: Support for organization-specific MAC address mappings
- **Flexible Output**: Display results in terminal and/or save to files
- **Configurable Timeouts**: Adjust scanning duration for different network sizes
- **Cross-Architecture**: Builds for AMD64 and ARM Linux systems

## Requirements

**System Requirements**:
- Linux kernel 2.6 or later
- IPv4 network connectivity
- AMD64 or ARM architecture

**Privileges**:
- **Active scanning**: Root privileges or `CAP_NET_RAW` capability
- **Passive mode**: No special privileges required

**Grant capability to avoid sudo**:
```bash
sudo setcap cap_net_raw+ep /usr/local/bin/vlookup
vlookup  # Now works without sudo
```

**Limitations**:
- Linux only (no Windows or macOS support)
- IPv4 only (IPv6 not supported)
- Ethernet networks only

## Installation

### Option 1: Download Binary (Recommended)

Download pre-built binaries from the [Releases page](https://github.com/frzifus/vlookup/releases):

```bash
# Download latest release (replace vX.Y.Z with actual version)
wget https://github.com/frzifus/vlookup/releases/download/vX.Y.Z/vlookup-linux-amd64
chmod +x vlookup-linux-amd64
sudo mv vlookup-linux-amd64 /usr/local/bin/vlookup

# Verify installation
vlookup --version
```

### Option 2: Build from Source (Primary Method)

**Prerequisites**: Go 1.18 or later

```bash
# Clone repository
git clone https://github.com/frzifus/vlookup.git
cd vlookup

# Build for your architecture
make amd64  # For AMD64
# or
make arm    # For ARM

# Install binary
sudo cp build/bin/vlookup-linux-amd64 /usr/local/bin/vlookup
sudo cp build/bin/crawler-linux-amd64 /usr/local/bin/crawler

# Verify installation
vlookup --version
crawler --version
```

### Option 3: Using Go (Developers)

```bash
# Install main tool
go install github.com/frzifus/vlookup/cmd/vlookup@latest

# Install database crawler
go install github.com/frzifus/vlookup/cmd/crawler@latest

# Ensure $GOPATH/bin is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Verification

```bash
# Check vlookup is installed
which vlookup

# Check version
vlookup --version

# Test passive mode (no root required)
vlookup --arp.scan=false
```

## Command Reference

### vlookup

Network scanning and MAC address lookup tool.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-src` | string | `"embd-l"` | Database source: `ieee-s`, `ieee-m`, `ieee-l`, `embd-s`, `embd-m`, `embd-l` |
| `--src.local-file` | string | `""` | Path to local CSV database file |
| `--trim.address` | int | `40` | Maximum length of vendor address field in output |
| `--arp.scan` | bool | `true` | Enable active ARP scanning (requires root) |
| `--arp.timeout` | duration | `10s` | Maximum time to wait for ARP responses |
| `-i` | string | `""` | Filter to specific network interface (empty = all) |
| `-o` | string | `""` | Output file path (empty = stdout only) |
| `--version` | bool | `false` | Print version information and exit |

**Database Source Options**:
- `ieee-l`, `ieee-m`, `ieee-s`: Download from IEEE (large, medium, small)
- `embd-l`, `embd-m`, `embd-s`: Use embedded databases (large, medium, small)

### crawler

Database download utility for updating IEEE vendor databases.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--src.fetch-l` | bool | `false` | Download large (OUI) database (~2.7MB) |
| `--src.fetch-m` | bool | `false` | Download medium (MAM) database (~380KB) |
| `--src.fetch-s` | bool | `false` | Download small (OUI36) database (~390KB) |
| `--src.fetch-all` | bool | `false` | Download all three database sizes |
| `--src.fetch-custom` | string | `""` | Download from custom URL |
| `--timeout` | duration | `30s` | HTTP request timeout |
| `-o` | string | `""` | Output filename prefix for downloaded files |
| `--version` | bool | `false` | Print version information and exit |

**Database Sizes**:
- **Large (OUI)**: ~2.7M entries, most comprehensive
- **Medium (MAM)**: ~380K entries, moderate coverage
- **Small (OUI36)**: ~390K entries, minimal size

## Quick Start

Already installed? Get started with these essential commands:

```bash
# Scan all network interfaces
sudo vlookup

# Scan specific interface only
sudo vlookup -i eth0

# View ARP cache without active scanning (no root required)
vlookup --arp.scan=false
```

See [Installation](#installation) to set up vlookup, or [Examples](#examples) for more usage scenarios.

## Usage

### Basic Scanning

Scan all active network interfaces (requires root):

```bash
sudo vlookup
```

Output format:
```
idx   interface  IP                   MAC                  Name                 Address
---   ---------  --                   ---                  ----                 -------
0     eth0       192.168.1.1          aa:bb:cc:dd:ee:ff    Cisco Systems        San Jose, CA US
1     eth0       192.168.1.100        11:22:33:44:55:66    Apple Inc            Cupertino, CA US
```

### Passive Mode

Query ARP cache without active scanning (no root required):

```bash
vlookup --arp.scan=false
```

Use this mode for:
- Forensic analysis
- Environments where active scanning is prohibited
- Quick MAC address lookups
- Running as non-root user

### Interface Filtering

Scan only a specific network interface:

```bash
sudo vlookup -i eth0
```

List available interfaces:
```bash
ip addr show
```

### Custom Timeout

Adjust scanning timeout for large or slow networks:

```bash
# Wait 30 seconds for responses (default: 10s)
sudo vlookup --arp.timeout=30s
```

### Save Results

Write results to a file:

```bash
# Output to both terminal and file
sudo vlookup -o network-scan.txt

# Suppress terminal output
sudo vlookup -o network-scan.txt > /dev/null
```

## Examples

### Example 1: Basic Network Inventory

Discover all devices on all network interfaces:

```bash
sudo vlookup
```

**Expected Output**:
```
idx   interface  IP                   MAC                  Name                 Address
---   ---------  --                   ---                  ----                 -------
0     eth0       192.168.1.1          aa:bb:cc:dd:ee:ff    Cisco Systems        San Jose, CA US
1     eth0       192.168.1.100        11:22:33:44:55:66    Apple Inc            Cupertino, CA US
2     eth0       192.168.1.150        ff:ee:dd:cc:bb:aa    Dell Inc             Round Rock, TX US
```

**Use Case**: Quick inventory of all active devices for security audit

---

### Example 2: Scan Wireless Network Only

Scan only the wireless interface on a multi-homed system:

```bash
sudo vlookup -i wlan0
```

**Use Case**: Focus troubleshooting on specific network segment without disrupting others

---

### Example 3: Offline MAC Lookup

Look up vendors for devices already in ARP cache:

```bash
vlookup --arp.scan=false
```

**Use Case**: Forensic analysis, read-only mode, no network traffic generated

**Note**: Results depend on ARP cache contents; use `ip neigh show` to view cache

---

### Example 4: Large Network Scan

Scan a large or slow network with extended timeout:

```bash
sudo vlookup --arp.timeout=30s
```

**Use Case**: Corporate networks, congested links, or networks with slow-responding devices

**Performance**: /24 network typically completes in 10-15 seconds

---

### Example 5: Periodic Audit

Save scan results for later comparison:

```bash
# Initial scan
sudo vlookup -o baseline-$(date +%Y%m%d).txt

# Weekly scan
sudo vlookup -o audit-$(date +%Y%m%d).txt

# Compare results
diff baseline-20251201.txt audit-20251208.txt
```

**Use Case**: Track network changes, detect unauthorized devices, compliance auditing

---

### Example 6: Custom Vendor Database

Use organization-specific MAC address database:

```bash
# Create custom database (IEEE CSV format)
cat > my-vendors.csv << EOF
Registry,Assignment,Organization Name,Organization Address
MA-L,001122,Internal Device Pool,"Corporate IT, Building 1"
MA-L,AABBCC,Lab Equipment,"R&D Lab, Building 2"
EOF

# Use custom database
vlookup --arp.scan=false --src.local-file=my-vendors.csv
```

**Use Case**: Corporate environments with private MAC ranges, custom device naming

---

### Example 7: Download Database Updates

Download latest IEEE vendor databases:

```bash
# Download all databases
crawler --src.fetch-all -o updated

# Files created with timestamps:
# - 2025-12-03_0_updated.csv (large)
# - 2025-12-03_1_updated.csv (medium)
# - 2025-12-03_2_updated.csv (small)

# Use updated database
vlookup --arp.scan=false --src.local-file=2025-12-03_0_updated.csv
```

**Use Case**: Quarterly database refresh to identify newly allocated MAC ranges

**Recommendation**: Update databases quarterly or when encountering "not found" vendors

---

### Example 8: Fast Scan with Small Database

Use smaller database for faster startup:

```bash
vlookup -src embd-s
```

**Use Case**: Quick scans where comprehensive vendor coverage isn't critical

**Trade-off**: Faster load time (~0.5s vs ~2s) but less vendor coverage

---

### Example 9: Combined Filters

Combine multiple options for specific use case:

```bash
sudo vlookup -i eth0 --arp.timeout=5s -o quick-scan.txt
```

**Use Case**: Scripted monitoring, cron jobs, automated audits

---

### Example 10: Troubleshooting Connection

Diagnose why a specific device isn't accessible:

```bash
# Check if device responds to ARP
sudo vlookup -i eth0

# Verify device's MAC is in cache
cat /proc/net/arp | grep 192.168.1.100

# Check vendor to confirm device type
vlookup --arp.scan=false | grep 192.168.1.100
```

**Use Case**: Network troubleshooting, device identification, connectivity issues

## Database Management

### Embedded Databases

vlookup includes three pre-loaded IEEE vendor databases:

| Size | Entries | File Size | Load Time | Use Case |
|------|---------|-----------|-----------|----------|
| Large (OUI) | ~2.7M | ~2.7MB | ~2s | Maximum coverage |
| Medium (MAM) | ~380K | ~380KB | ~0.8s | Balance coverage/speed |
| Small (OUI36) | ~390K | ~390KB | ~0.5s | Minimal size, major vendors |

**Default**: Large database (`embd-l`)

### Updating Databases

IEEE allocates new MAC address ranges monthly. Update databases quarterly:

```bash
# Download latest databases
crawler --src.fetch-all -o updated

# Verify download
ls -lh *_updated.csv

# Use updated database
vlookup --src.local-file=2025-12-03_0_updated.csv --arp.scan=false
```

### Database Sources

**Embedded** (`embd-*`):
- Compiled into binary
- No network access required
- Updated with each vlookup release

**Remote** (`ieee-*`):
- Downloaded from IEEE servers
- Always current
- Requires internet connectivity
- May be slow on first run

**Local** (`--src.local-file`):
- Custom or downloaded databases
- Organization-specific mappings
- Complete control over content

### CSV Format

Custom databases must follow IEEE format:

```csv
Registry,Assignment,Organization Name,Organization Address
MA-L,AABBCC,Example Corporation,"123 Main Street, City, State, Country ZIP"
MA-M,AABBCCD,Medium Block Corp,"456 Oak Avenue, City, Country"
MA-S,AABBCCDEF,Small Block Inc,"789 Elm Street, City, ZIP"
```

**Requirements**:
- Header row required
- Exactly 4 columns
- Registry: MA-L, MA-M, or MA-S
- Assignment: Hex MAC prefix (6, 7, or 9 characters, no colons)
- Organization Name: UTF-8 string
- Organization Address: Quoted if contains commas

## Troubleshooting

### Error: "user has insufficient permissions"

**Cause**: Active ARP scanning requires root or CAP_NET_RAW capability

**Solutions**:

1. **Run with sudo** (quick fix):
   ```bash
   sudo vlookup
   ```

2. **Grant capability** (persistent):
   ```bash
   sudo setcap cap_net_raw+ep /usr/local/bin/vlookup
   vlookup  # Now works without sudo
   ```

3. **Use passive mode** (no privileges):
   ```bash
   vlookup --arp.scan=false
   ```

---

### No Devices Found

**Symptoms**: Scan completes but results table is empty

**Causes**:
- Network isolation (no other devices)
- Wrong interface selected
- Timeout too short
- Firewall blocking ARP

**Diagnostics**:
```bash
# Check network connectivity
ping 192.168.1.1

# List interfaces
ip addr show

# Check ARP cache
cat /proc/net/arp
```

**Solutions**:
```bash
# Scan specific interface
sudo vlookup -i eth0

# Increase timeout
sudo vlookup --arp.timeout=30s

# Try passive mode
vlookup --arp.scan=false
```

---

### Vendor Shows "not found"

**Cause**: MAC address not in database

**Reasons**:
- Newly allocated MAC range (update database)
- Locally administered MAC (won't have IEEE vendor)
- Private/custom MAC range

**Solutions**:
```bash
# Update database
crawler --src.fetch-all -o fresh

# Use updated database
vlookup --src.local-file=*_fresh.csv --arp.scan=false

# Create custom database for private MACs
# (see Database Management section)
```

**Note**: Locally administered MACs have second hex digit of 2, 6, A, or E (e.g., `02:`, `06:`, `0A:`, `0E:`)

---

### Tool Hangs or Very Slow

**Cause**: Large network or network issues

**Solutions**:
```bash
# Reduce timeout
sudo vlookup --arp.timeout=5s

# Scan specific interface
sudo vlookup -i eth0

# Use passive mode
vlookup --arp.scan=false
```

**Expected Performance**:
- /24 network (254 hosts): 10-15 seconds
- /16 network (65K hosts): Several minutes (not recommended)

---

### CSV Parsing Error

**Error**: `failed to parse CSV at line X`

**Cause**: Corrupted or invalid database file

**Solutions**:
```bash
# Re-download database
crawler --src.fetch-all -o fresh

# Verify CSV format
head -5 *_fresh.csv

# Expected format:
# Registry,Assignment,Organization Name,Organization Address
# MA-L,AABBCC,Example Corp,"123 Main St, City, Country"
```

---

### Interface Not Found

**Symptoms**: Specified interface produces no results

**Diagnostics**:
```bash
# List all interfaces (case-sensitive)
ip link show

# Bring interface up if needed
sudo ip link set eth0 up
```

**Solution**:
```bash
# Use exact interface name
sudo vlookup -i eth0  # Not ETH0 or Eth0
```

---

### Still Having Issues?

1. Check [GitHub Issues](https://github.com/frzifus/vlookup/issues) for known problems
2. Search for your specific error message
3. Report new issues with:
   - Output of `vlookup --version`
   - Full error text
   - Network environment details (Linux distro, kernel version)
   - Steps to reproduce

## Performance

### Expected Scan Times

| Network Size | Active Scan | Passive Mode |
|--------------|-------------|--------------|
| /24 (254 hosts) | 10-15 seconds | < 1 second |
| /23 (510 hosts) | 15-25 seconds | < 1 second |
| Single device | 10 seconds (timeout) | Instant |

**Note**: Times assume default 10-second timeout and responsive devices

### Memory Usage

- Embedded databases: ~10MB total
- Scan results: ~100 bytes per device
- 1000 devices: ~10.1MB total memory

### Optimization Tips

**For Large Networks**:
```bash
# Reduce timeout
sudo vlookup --arp.timeout=5s

# Scan interfaces individually
sudo vlookup -i eth0
sudo vlookup -i eth1
```

**For Speed**:
```bash
# Use smaller database
vlookup -src embd-s

# Skip active scanning
vlookup --arp.scan=false
```

### Limitations

- **Network Size**: Practical limit around 1000-2000 active devices
- **Timeout**: Minimum ~10 seconds for network round-trip
- **Concurrency**: Single-threaded; one scan at a time per interface
- **Protocols**: ARP only; no IPv6, no ICMP, no TCP/UDP

## Security

### Privilege Requirements

**Active Scanning** requires:
- Root privileges (`sudo vlookup`), OR
- `CAP_NET_RAW` capability

**Why**: Raw socket access needed to send/receive ARP packets

**Grant Capability**:
```bash
sudo setcap cap_net_raw+ep /usr/local/bin/vlookup
```

**Risks**:
- Grants raw packet access to the binary
- Binary can craft arbitrary network packets
- Only grant to trusted binaries

---

### Network Impact

**Active Scanning**:
- Sends one ARP request per IP in subnet
- /24 network = 254 ARP packets
- May trigger IDS/IPS alerts
- Generates visible network traffic

**Passive Mode**:
- Read-only access to ARP cache
- No network traffic generated
- Safe for production environments

---

### Data Privacy

**Collected Data**:
- IP addresses
- MAC addresses
- Network interface names
- Vendor information (public IEEE data)

**Storage**:
- In-memory only (unless `-o` used)
- No telemetry or external transmission
- No persistent storage by default

**Saved Files** (`-o` option):
- Contain network topology information
- Should be protected as sensitive data
- Recommend secure file permissions:
  ```bash
  sudo vlookup -o scan.txt
  chmod 600 scan.txt  # Owner read/write only
  ```

---

### Best Practices

1. **Authorization**: Only scan networks you own or have permission to scan
2. **Notification**: Inform security teams before scanning production networks
3. **Scheduling**: Avoid scanning during peak hours
4. **Storage**: Protect scan results as sensitive data
5. **Updates**: Keep vlookup and databases current
6. **Auditing**: Log scan activities for security compliance
7. **Least Privilege**: Use passive mode when active scanning isn't needed

---

### Compliance

**Legal Considerations**:
- Scanning networks you don't own may be illegal
- Corporate policies may restrict network scanning
- Compliance frameworks (PCI DSS, HIPAA) may have requirements

**Your Responsibilities**:
- Obtain proper authorization
- Follow organizational policies
- Secure scan results appropriately
- Use tool ethically and legally

## Contributing

### Reporting Issues

Found a bug? [Open an issue](https://github.com/frzifus/vlookup/issues) with:

- **vlookup version**: Output of `vlookup --version`
- **System info**: Linux distro, kernel version (`uname -a`)
- **Network details**: Interface type (eth, wlan, etc.), network size
- **Error message**: Full error text
- **Steps to reproduce**: Exact commands run
- **Expected vs actual**: What should happen vs what does happen

### Feature Requests

Want a new feature? [Open an issue](https://github.com/frzifus/vlookup/issues) describing:

- **Use case**: What problem are you trying to solve?
- **Proposed solution**: How should it work?
- **Alternatives considered**: Other ways to achieve the goal
- **Impact**: Who benefits? How common is this need?

### Code Contributions

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/my-feature`
3. **Make** your changes
4. **Test**: Run `make test` and `make lint`
5. **Commit**: Use clear, descriptive commit messages
6. **Push**: `git push origin feature/my-feature`
7. **Open** a Pull Request

**Code Style**:
- Follow existing Go code conventions
- Run `gofmt` before committing
- Add tests for new functionality
- Update documentation as needed

### Database Contributions

Custom vendor databases or IEEE URL updates:
- Verify CSV format
- Test with sample data
- Document any special requirements

## License

vlookup is licensed under the **GNU General Public License v3.0**.

See [LICENSE](LICENSE) for the full license text.

**In Summary**:
- ✅ Use for any purpose
- ✅ Study and modify source code
- ✅ Distribute copies
- ✅ Distribute modified versions
- ⚠️ Must disclose source
- ⚠️ Must license under GPL-3.0
- ⚠️ Must state changes
- ⚠️ No warranty provided

**Commercial Use**: Permitted under GPL-3.0 terms

## Acknowledgments

**Dependencies**:
- [mdlayher/arp](https://github.com/mdlayher/arp) - ARP protocol implementation
- [mdlayher/ethernet](https://github.com/mdlayher/ethernet) - Ethernet frame handling
- [google/go-cmp](https://github.com/google/go-cmp) - Testing comparisons

**Data Sources**:
- [IEEE Registration Authority](https://standards.ieee.org/products-programs/regauth/) - MAC address vendor databases

**Inspiration**:
- [arp-scan](https://github.com/royhills/arp-scan) - Original ARP scanning utility
- [nmap](https://nmap.org/) - Network exploration and security auditing

**Contributors**:
- See [GitHub Contributors](https://github.com/frzifus/vlookup/graphs/contributors)

**Maintained By**: [frzifus](https://github.com/frzifus)
