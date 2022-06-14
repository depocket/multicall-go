package call

type Chain string

type ChainConfig struct {
	MultiCallAddress string
	Url              string
}

const (
	Arbitrum  Chain = "Arbitrum"
	Aurora          = "aurora"
	Avalanche       = "avalanche"
	Bsc             = "bsc"
	Ethereum        = "ethereum"
	Fantom          = "fantom"
	Moonbeam        = "moonbeam"
	Moonriver       = "moonriver"
)

var DefaultChainConfigs = map[Chain]ChainConfig{
	Arbitrum: {
		MultiCallAddress: "0x7a7443f8c577d537f1d8cd4a629d40a3148dd7ee",
		Url:              "https://arb1.arbitrum.io/rpc",
	},
	Aurora: {
		MultiCallAddress: "0x88b373B83166E72FD55648Ce114712633f1782E2",
		Url:              "https://mainnet.aurora.dev",
	},
	Avalanche: {
		MultiCallAddress: "0xa00FB557AA68d2e98A830642DBbFA534E8512E5f",
		Url:              "https://api.avax.network/ext/bc/C/rpc",
	},
	Bsc: {
		MultiCallAddress: "0x41263cBA59EB80dC200F3E2544eda4ed6A90E76C",
		Url:              "https://bsc-dataseed1.ninicoin.io",
	},
	Ethereum: {
		MultiCallAddress: "0xeefba1e63905ef1d7acba5a8513c70307c1ce441",
		Url:              "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
	},
	Fantom: {
		MultiCallAddress: "0x7F4e475462A0fA0F1e2C69d50866D54505F99D72",
		Url:              "https://rpcapi.fantom.network",
	},
	Moonbeam: {
		MultiCallAddress: "0x6477204E12A7236b9619385ea453F370aD897bb2",
		Url:              "https://moonbeam.public.blastapi.io",
	},
	Moonriver: {
		MultiCallAddress: "0xaef00a0cf402d9dedd54092d9ca179be6f9e5ce3",
		Url:              "https://moonriver.public.blastapi.io",
	},
}
