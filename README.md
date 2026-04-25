# vanity-onion

[![GitHub Repo stars](https://img.shields.io/github/stars/feloex/vanity-onion?style=flat&color=yellow&link=https%3A%2F%2Fgithub.com%2Ffeloex%2Fvanity-onion)](https://github.com/feloex/vanity-onion)
[![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/feloex/vanity-onion/total?style=flat&color=dark-green)](https://github.com/feloex/vanity-onion/releases)
[![GitHub Downloads (all assets, latest release)](https://img.shields.io/github/downloads/feloex/vanity-onion/latest/total?style=flat&color=dark-green)](https://github.com/feloex/vanity-onion/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/feloex/vanity-onion)](https://goreportcard.com/report/github.com/feloex/vanity-onion)

This tool easily generates a set amount of .onion addresses with a set prefix (vanity addresses).

## Usage
1. Download the latest release from the [releases page](https://github.com/feloex/vanity-onion/releases)
2. On Linux/MacOS, you need to give the binary execute permissions:
```bash
chmod +x vanity-onion-linux-amd64
```
3. Run the binary with the desired prefix and amount of addresses to generate:
```bashbash
./vanity-onion-linux-amd64 hello 5
```
This will generate 5 .onion addresses that start with "hello"

The generated addresses and their corresponding keys will be saved in the `keys` directory.

## Sources
Great video to understand how .onion addresses are generated:
https://www.youtube.com/watch?v=kRQvE5x36t4

Tor specification:
https://spec.torproject.org/

---
## License
This project is licensed under the [GNU General Public License v3.0](LICENSE). 