package sockpuppet

import (
	context "context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
EthClient is a wrapper around a client object that can
communicate with the ethereum blockchain.
*/
type EthClient struct {
	conn     *ethclient.Client
	data     bind.ContractBackend
	endpoint string
	token    string
	auth     *bind.TransactOpts
	err      errnie.Error
}

/*
NewEthClient instantiates a pointer to an Ethclient instance.
*/
func NewEthClient(endpoint, token string) *EthClient {
	return &EthClient{
		endpoint: endpoint,
		token:    token,
	}
}

/*
Dial to the ethereum blockchain.
*/
func (client *EthClient) Dial() *EthClient {
	var err error
	client.conn, err = ethclient.Dial(client.endpoint)
	client.err = errnie.Handles(err)
	client.auth, client.err = client.getAccountAuth(client.token)
	return client
}

/*
Conn returns the connection instance.
*/
func (client *EthClient) Conn() bind.ContractBackend {
	return client.conn
}

/*
Auth returns the authentication instance.
*/
func (client *EthClient) Auth() *bind.TransactOpts {
	return client.auth
}

/*
getAccountAuth retrieves an authentication instance from a private key.
*/
func (client *EthClient) getAccountAuth(
	privateKeyAddress string,
) (*bind.TransactOpts, errnie.Error) {
	privateKey, err := crypto.HexToECDSA(privateKeyAddress)
	errnie.Handles(err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		return &bind.TransactOpts{}, errnie.Handles(
			errnie.NewError(errors.New("invalid key")),
		)
	}

	var nonce uint64

	if nonce, err = client.conn.PendingNonceAt(
		context.Background(), crypto.PubkeyToAddress(*publicKeyECDSA),
	); err != nil {
		return &bind.TransactOpts{}, errnie.Handles(err)
	}

	fmt.Println("nounce=", nonce)
	chainID, err := client.conn.ChainID(context.Background())
	errnie.Handles(err)

	var auth *bind.TransactOpts

	if auth, err = bind.NewKeyedTransactorWithChainID(
		privateKey, chainID,
	); err != nil {
		return &bind.TransactOpts{}, errnie.Handles(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(1000000)

	return auth, errnie.Error{}
}
