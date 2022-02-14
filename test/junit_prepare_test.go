package test

import (
	"errors"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/hamster-shared/hamster-provider/cmd"
	"github.com/hamster-shared/hamster-provider/core/context"
	"github.com/sirupsen/logrus"
	"math/big"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := cmd.NewContext()

	// 转账
	cf := ctx.GetConfig()
	kp, _ := signature.KeyringPairFromSecret(cf.SeedOrPhrase, 42)
	address := types.NewMultiAddressFromAccountID(kp.PublicKey)
	transformAmountToTestAccount(address, 100_000000000000)
	// 质押
	//err := staking(ctx, 99_000000000000)
	//assert.NoError(t, err)
}

func transformAmountToTestAccount(target types.MultiAddress, amount uint64) {
	// Display the events that occur during a transfer by sending a value to bob

	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Bob
	//bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	//bob, err := types.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}
	a := types.NewUCompactFromUInt(amount)

	call, err := types.NewCall(meta, "Balances.transfer", target, a)
	if err != nil {
		panic(err)
	}

	c, err := types.NewCall(meta, "Utility.batch", []types.Call{call})

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	keypair := signature.TestKeyringPairAlice

	// Get the nonce for Alice
	key, err := types.CreateStorageKey(meta, "System", "Account", keypair.PublicKey)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	//	fmt.Printf("Sending %v from %#x to %#x with nonce %v", amount, signature.TestKeyringPairAlice.PublicKey, bob.AsAccountID, nonce)

	// Sign the transaction using Alice's default account
	err = ext.Sign(keypair, o)
	if err != nil {
		panic(err)
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		status := <-sub.Chan()
		fmt.Printf("Transaction status: %#v\n", status)

		if status.IsInBlock {
			fmt.Printf("Completed at block hash: %#x\n", status.AsInBlock)
			return
		}
	}
}

func staking(context context.CoreContext, amount uint64) error {
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		return err
	}
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}
	c, err := types.NewCall(meta, "ResourceOrder.staking_amount", types.NewU128(*big.NewInt(int64(amount))))
	if err != nil {
		return err
	}
	err = call(context, api, c, meta)
	return err
}

func call(ctx context.CoreContext, api *gsrpc.SubstrateAPI, c types.Call, meta *types.Metadata) error {

	cf := ctx.GetConfig()

	// Create the extrinsic
	ext := types.NewExtrinsic(c)
	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return err
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return err
	}

	keypair, err := signature.KeyringPairFromSecret(cf.SeedOrPhrase, 42)
	if err != nil {
		return err
	}

	// Get the nonce for Account
	key, err := types.CreateStorageKey(meta, "System", "Account", keypair.PublicKey)
	if err != nil {
		return err
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("GetStorageLatest fail")
	}

	nonce := uint32(accountInfo.Nonce)
	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using User's default account
	err = ext.Sign(keypair, o)
	if err != nil {
		return err
	}

	res, err := api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		logrus.Errorf("extrinsic submit failed: %v", err)
		return err
	}

	hex, err := types.Hex(res)
	if err != nil {
		return err
	}
	if hex == "" {
		return errors.New("hex is empty")
	}
	return nil
}
