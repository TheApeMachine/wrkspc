package sockpuppet

import (
	context "context"
	"crypto/ecdsa"
	"errors"
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
	err      error
}

/*
NewEthClient instantiates a pointer to an Ethclient instance.
*/
func NewEthClient(endpoint, token string) *EthClient {
	errnie.Traces()
	errnie.Debugs(endpoint, token)

	return &EthClient{
		endpoint: endpoint,
		token:    token,
	}
}

/*
Dial to the ethereum blockchain.
*/
func (client *EthClient) Dial() *EthClient {
	errnie.Traces()

	if client.conn, client.err = ethclient.Dial(client.endpoint); client.err != nil {
		errnie.Handles(client.err)
		return nil
	}

	if client.auth, client.err = client.getAccountAuth(client.token); client.err != nil {
		errnie.Handles(client.err)
		return nil
	}

	return client
}

/*
Error implements the error interface and returns the latest error.
*/
func (client *EthClient) Error() error {
	errnie.Traces()
	return client.err
}

/*
Conn returns the connection instance.
*/
func (client *EthClient) Conn() bind.ContractBackend {
	errnie.Traces()
	return client.conn
}

/*
Auth returns the authentication instance.
*/
func (client *EthClient) Auth() *bind.TransactOpts {
	errnie.Traces()
	return client.auth
}

/*
getAccountAuth retrieves an authentication instance from a private key.
*/
func (client *EthClient) getAccountAuth(
	privateKeyAddress string,
) (*bind.TransactOpts, errnie.Error) {
	errnie.Traces()
	privateKey, err := crypto.HexToECDSA(privateKeyAddress)

	if e := errnie.Handles(err); e.Type != errnie.NIL {
		return &bind.TransactOpts{}, e
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	errnie.Debugs("publicKeyECDSA", publicKeyECDSA)

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

	errnie.Debugs("nonce", nonce)

	chainID, err := client.conn.ChainID(context.Background())
	if e := errnie.Handles(err); e.Type != errnie.NIL {
		return &bind.TransactOpts{}, e
	}

	errnie.Debugs("chainID", nonce)

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

	errnie.Debugs("auth", nonce)
	return auth, errnie.Error{}
}
