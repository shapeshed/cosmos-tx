package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func main() {
	// Create a new background context
	_ = context.Background()

	// Set prefix for config
	sdk.GetConfig().SetBech32PrefixForAccount("osmo", "osmo")

	// Get simapp encoding configuration
	encodingConfig := simapp.MakeTestEncodingConfig()

	// Access keyring
	kr, err := keyring.New("test", keyring.BackendPass, "/home/go/", nil)
	if err != nil {
		fmt.Printf("failed to create keyring: %s", err.Error())
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("failed to get info: %s", err.Error())
		os.Exit(1)
	}

	// Create  node client
	nodeClient, err := client.NewClientFromNode("https://rpc.testnet.osmosis.zone:443")
	if err != nil {
		fmt.Printf("failed to create node client: %s", err.Error())
		os.Exit(1)
	}

	// Create a new client context
	clientCtx := client.Context{}.
		WithAccountRetriever(types.AccountRetriever{}).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithCodec(encodingConfig.Marshaler).
		WithKeyring(kr).
		WithClient(nodeClient).
		WithChainID("osmo-test-5").
		WithSignModeStr("SIGN_MODE_UNSPECIFIED").
		WithBroadcastMode("block")

	// Create a new transaction builder
	txBuilder := clientCtx.TxConfig.NewTxBuilder()
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("uosmo", 100000)))
	txBuilder.SetGasLimit(200000)
	txBuilder.SetMemo("test transaction")

	// Fetch the sender's account key
	senderKey, err := clientCtx.Keyring.Key("bot-1")

	clientCtx.WithFeeGranterAddress(senderKey.GetAddress())
	clientCtx = clientCtx.WithFromAddress(senderKey.GetAddress())

	if err != nil {
		fmt.Printf("failed to get sender key: %s", err.Error())
		os.Exit(1)
	}

	// Fetch sender & recipient addresses
	fromAddr := senderKey.GetAddress()
	toAddr, err := sdk.AccAddressFromBech32("osmo1ys5lj28zrvpm9cjj0958awqmjt8aa8vpdflqc2")
	if err != nil {
		fmt.Printf("failed to parse recipient address: %s", err.Error())
		os.Exit(1)
	}

	// Create a new message to send funds to the recipient
	msg := &bankTypes.MsgSend{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uosmo", 50000)),
	}
	if err := msg.ValidateBasic(); err != nil {
		fmt.Printf("failed to validate message: %s", err.Error())
		os.Exit(1)
	}

	// Add the message to the transaction
	err = txBuilder.SetMsgs(msg)
	if err != nil {
		fmt.Printf("failed to set message: %s", err.Error())
		os.Exit(1)
	}

	// Sign the transaction
	txf := tx.Factory{}.
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithKeybase(clientCtx.Keyring).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithChainID(clientCtx.ChainID).
		WithTxConfig(clientCtx.TxConfig)

	// This checks the from address exists and sets the sequence (nonce)
	txf, err = txf.Prepare(clientCtx)
	if err != nil {
		fmt.Printf("failed to prepare txf: %s", err.Error())
		os.Exit(1)
	}

	err = tx.Sign(txf, senderKey.GetName(), txBuilder, true)
	if err != nil {
		fmt.Printf("failed to sign transaction: %s", err.Error())
		os.Exit(1)
	}

	// Encode the transaction and broadcast it to the network
	txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		fmt.Printf("failed to encode transaction: %s", err.Error())
		os.Exit(1)
	}

	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		fmt.Printf("failed to broadcast transaction: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("transaction sent: %s", res.TxHash)
}
