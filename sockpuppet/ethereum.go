package sockpuppet

import (
	context "context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/theapemachine/wrkspc/errnie"
)

type EthClient struct {
	conn     *ethclient.Client
	data     bind.ContractBackend
	endpoint string
	token    string
	auth     *bind.TransactOpts
	err      error
}

func NewEthClient(endpoint, token string) *EthClient {
	return &EthClient{
		endpoint: endpoint,
		token:    token,
	}
}

func (client *EthClient) Dial() *EthClient {
	client.conn, client.err = ethclient.Dial(client.endpoint)
	errnie.Handles(client.err)
	client.auth = getAccountAuth(client.conn, client.token)
	return client
}

func (client *EthClient) Conn() bind.ContractBackend {
	return client.conn
}

func (client *EthClient) Auth() *bind.TransactOpts {
	return client.auth
}

// function to create auth for any account from its private key
func getAccountAuth(
	client *ethclient.Client, privateKeyAddress string,
) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(privateKeyAddress)
	errnie.Handles(err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	errnie.Handles(err)

	fmt.Println("nounce=", nonce)
	chainID, err := client.ChainID(context.Background())
	errnie.Handles(err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	errnie.Handles(err)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(1000000)

	return auth
}
