# OnlyLayer-faucet

## Get started

**Optional Flags**

| Flag        | Description                                      | Default Value |
| ----------- | ------------------------------------------------ | ------------- |
| -chainname  | Network name to display on the frontend          | testnet       |
| -httpport   | Listener port to serve HTTP connection           | 8080          |
| -interval   | Number of minutes to wait between funding rounds | 1440          |
| -payout     | Number of Ethers to transfer per user request    | 1             |
| -proxycount | Count of reverse proxies in front of the server  | 0             |
| -queuecap   | Maximum transactions waiting to be sent          | 100           |

## Development

### Prerequisites

- Go (1.16 or later)
- Node.js

### Build

1. Clone the repository and navigate to the app’s directory

```bash
git clone https://github.com/layeronly/faucet.git
cd onlyfaucet
```

2. Bundle Frontend web with Rollup

```bash
npm run build
```

3. Build binary application to run faucet

```bash
go build
export WEB3_PROVIDER=https://onlylayer.org
export PRIVATE_KEY=secret
./faucet
```

## License

This project is licensed under the MIT License
