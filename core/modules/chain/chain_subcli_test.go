package chain

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	"github.com/minio/blake2b-simd"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestSampleCall(t *testing.T) {
	api, err := gsrpc.NewSubstrateAPI("wss://chain.stage-ttchain.tntlinking.com")
	if err != nil {
		panic(err)
	}

	chain, err := api.RPC.System.Chain()
	if err != nil {
		panic(err)
	}
	nodeName, err := api.RPC.System.Name()
	if err != nil {
		panic(err)
	}
	nodeVersion, err := api.RPC.System.Version()
	if err != nil {
		panic(err)
	}

	fmt.Printf("You are connected to chain %v using %v v%v\n", chain, nodeName, nodeVersion)
}

func TestTransform(t *testing.T) {
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
	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	//bob, err := types.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	amount := types.NewUCompactFromUInt(12345)

	c, err := types.NewCall(meta, "Balances.transfer", bob, amount)
	if err != nil {
		panic(err)
	}

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

func TestListenToNewBlocks(t *testing.T) {
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	sub, err := api.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	count := 0

	for {
		head := <-sub.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		count++

		if count == 10 {
			sub.Unsubscribe()
			break
		}
	}
}

func TestAuthor_SubmitExtrinsic(t *testing.T) {
	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	assert.NoError(t, err)

	meta, err := api.RPC.State.GetMetadataLatest()
	assert.NoError(t, err)

	// Create a call, transferring 12345 units to Bob
	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	assert.NoError(t, err)

	amount := types.NewUCompactFromUInt(12345)
	c, err := types.NewCall(meta, "Balances.transfer", bob, amount)
	assert.NoError(t, err)

	for {
		// Create the extrinsic
		ext := types.NewExtrinsic(c)
		genesisHash, err := api.RPC.Chain.GetBlockHash(0)
		assert.NoError(t, err)

		rv, err := api.RPC.State.GetRuntimeVersionLatest()
		assert.NoError(t, err)

		// Get the nonce for Alice
		key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey)
		assert.NoError(t, err)

		var accountInfo types.AccountInfo
		ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
		assert.NoError(t, err)
		assert.True(t, ok)
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

		fmt.Printf("Sending %v from %#x to %#x with nonce %v\n", amount, signature.TestKeyringPairAlice.PublicKey,
			bob.AsID, nonce)

		// Sign the transaction using Alice's default account
		err = ext.Sign(signature.TestKeyringPairAlice, o)
		assert.NoError(t, err)

		res, err := api.RPC.Author.SubmitExtrinsic(ext)
		if err != nil {
			t.Logf("extrinsic submit failed: %v", err)
			continue
		}

		hex, err := types.Hex(res)
		assert.NoError(t, err)
		assert.NotEmpty(t, hex)
		break
	}
}

func TestStorageOrder(t *testing.T) {
	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	assert.NoError(t, err)

	meta, err := api.RPC.State.GetMetadataLatest()
	assert.NoError(t, err)

	cid := "abcd"
	fileName := "ac"
	tips := types.NewU128(*big.NewInt(123))
	duration := types.NewU32(32)
	size := types.NewU64(5000)
	c, err := types.NewCall(meta, "StorageOrder.create_order", cid, fileName, tips, duration, size)

	//c,err := types.NewCall(meta,"StorageOrder.check_param",size)

	assert.NoError(t, err)

	//for {
	// Create the extrinsic
	ext := types.NewExtrinsic(c)
	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	assert.NoError(t, err)

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	assert.NoError(t, err)

	// Get the nonce for Alice
	key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey)
	assert.NoError(t, err)

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	assert.NoError(t, err)
	assert.True(t, ok)
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

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	res, err := api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		t.Logf("extrinsic submit failed: %v", err)
	}

	hex, err := types.Hex(res)
	assert.NoError(t, err)
	assert.NotEmpty(t, hex)
	//}
}

type EventSomethingStored struct {
	Phase     types.Phase
	Something types.U32
	Who       types.AccountID
	Topics    []types.Hash
}

type EventBalancesWithdraw struct {
	Phase   types.Phase
	Who     types.AccountID
	Balance types.U128
	Topics  []types.Hash
}

func TestDisplaySystemEvents(t *testing.T) {
	// Create our API with a default connection to the local node
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		panic(err)
	}

	sub, err := api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		set := <-sub.Chan()
		// inner loop for the changes within one of those notifications
		for _, chng := range set.Changes {
			if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
				// skip, we are only interested in events with content
				continue
			}

			type MyEventRecords struct {
				types.EventRecords
				TemplateModule_SomethingStored []EventSomethingStored //nolint:stylecheck,golint
				Balances_Withdraw              []EventBalancesWithdraw
			}

			// Decode the event records
			events := MyEventRecords{}
			storageData := chng.StorageData
			err = types.EventRecordsRaw(storageData).DecodeEventRecords(meta, &events)
			if err != nil {
				panic(err)
			}

			// Show what we are busy with
			for _, e := range events.Balances_Endowed {
				fmt.Printf("\tBalances:Endowed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Balances_DustLost {
				fmt.Printf("\tBalances:DustLost:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Balances_Transfer {
				fmt.Printf("\tBalances:Transfer:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v, %v\n", e.From, e.To, e.Value)
			}
			for _, e := range events.Balances_BalanceSet {
				fmt.Printf("\tBalances:BalanceSet:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v, %v\n", e.Who, e.Free, e.Reserved)
			}
			for _, e := range events.Balances_Deposit {
				fmt.Printf("\tBalances:Deposit:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v, %v\n", e.Who, e.Balance)
			}
			for _, e := range events.Grandpa_NewAuthorities {
				fmt.Printf("\tGrandpa:NewAuthorities:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.NewAuthorities)
			}
			for _, e := range events.Grandpa_Paused {
				fmt.Printf("\tGrandpa:Paused:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.Grandpa_Resumed {
				fmt.Printf("\tGrandpa:Resumed:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.ImOnline_HeartbeatReceived {
				fmt.Printf("\tImOnline:HeartbeatReceived:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x\n", e.AuthorityID)
			}
			for _, e := range events.ImOnline_AllGood {
				fmt.Printf("\tImOnline:AllGood:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.ImOnline_SomeOffline {
				fmt.Printf("\tImOnline:SomeOffline:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.IdentificationTuples)
			}
			for _, e := range events.Indices_IndexAssigned {
				fmt.Printf("\tIndices:IndexAssigned:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x%v\n", e.AccountID, e.AccountIndex)
			}
			for _, e := range events.Indices_IndexFreed {
				fmt.Printf("\tIndices:IndexFreed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.AccountIndex)
			}
			for _, e := range events.Offences_Offence {
				fmt.Printf("\tOffences:Offence:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v%v\n", e.Kind, e.OpaqueTimeSlot)
			}
			for _, e := range events.Session_NewSession {
				fmt.Printf("\tSession:NewSession:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.SessionIndex)
			}
			for _, e := range events.Staking_Reward {
				fmt.Printf("\tStaking:Reward:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.Amount)
			}
			for _, e := range events.Staking_Slash {
				fmt.Printf("\tStaking:Slash:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x%v\n", e.AccountID, e.Balance)
			}
			for _, e := range events.Staking_OldSlashingReportDiscarded {
				fmt.Printf("\tStaking:OldSlashingReportDiscarded:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.SessionIndex)
			}
			for _, e := range events.System_ExtrinsicSuccess {
				fmt.Printf("\tSystem:ExtrinsicSuccess:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.System_ExtrinsicFailed {
				fmt.Printf("\tSystem:ExtrinsicFailed:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%v\n", e.DispatchError)
			}
			for _, e := range events.System_CodeUpdated {
				fmt.Printf("\tSystem:CodeUpdated:: (phase=%#v)\n", e.Phase)
			}
			for _, e := range events.System_NewAccount {
				fmt.Printf("\tSystem:NewAccount:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#x\n", e.Who)
			}
			for _, e := range events.System_KilledAccount {
				fmt.Printf("\tSystem:KilledAccount:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%#X\n", e.Who)
			}
			for _, e := range events.TemplateModule_SomethingStored {
				fmt.Printf("\tTemplateModule:SomethingStored:: (phase=%#v)\n", e.Phase)
				fmt.Printf("\t\t%d\n", e.Who)
			}
			for _, e := range events.Balances_Withdraw {
				fmt.Printf("\t\t%d\n", e.Who)
			}
		}
	}
}

//func TestListenToBalanceChange(t *testing.T) {
//	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
//	if err != nil {
//		panic(err)
//	}
//
//	meta, err := api.RPC.State.GetMetadataLatest()
//	if err != nil {
//		panic(err)
//	}
//
//	alice := signature.TestKeyringPairAlice.PublicKey
//	key, err := types.CreateStorageKey(meta, "System", "Account", alice)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(key.Hex())
//
//	var accountInfo AccountInfo
//	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
//	if err != nil || !ok {
//		panic(err)
//	}
//
//	previous := accountInfo.Data.Free
//	fmt.Printf("%#x has a balance of %v\n", alice, previous)
//	fmt.Println(previous.Bytes())
//
//	latest, err := api.RPC.State.GetStorageRawLatest(key)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(latest.Hex())
//	fmt.Printf("You may leave this example running and transfer any value to %#x\n", alice)
//
//	//Here we subscribe to any balance changes
//	sub, err := api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
//	if err != nil {
//		panic(err)
//	}
//	defer sub.Unsubscribe()
//
//	//outer for loop for subscription notifications
//	for {
//		// inner loop for the changes within one of those notifications
//		for _, chng := range (<-sub.Chan()).Changes {
//			if !chng.HasStorageData {
//				continue
//			}
//
//			var acc AccountInfo
//			if err = types.DecodeFromBytes(chng.StorageData, &acc); err != nil {
//				panic(err)
//			}
//
//			// Calculate the delta
//			current := acc.Data.Free
//			var change = types.U128{Int: big.NewInt(0).Sub(current.Int, previous.Int)}
//
//			// Only display positive value changes (Since we are pulling `previous` above already,
//			// the initial balance change will also be zero)
//			if change.Cmp(big.NewInt(0)) != 0 {
//				fmt.Printf("New balance change of: %v %v %v\n", change, previous, current)
//				previous = current
//				return
//			}
//		}
//	}
//}

func TestResourceList(t *testing.T) {
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	bytes, err := types.EncodeToBytes(types.NewU64(0))
	if err != nil {
		return
	}
	key, err := types.CreateStorageKey(meta, "Provider", "Resources", bytes)

	if err != nil {
		panic(err)
	}
	fmt.Println(key.Hex())

	fmt.Println("0x1ef7b0947b9fbafef1c12486bf8512c22111e0df19de9563b58301e5f7e00743 bb1bdbcacd6ac9340000000000000000")
	var computingResource ComputingResource
	ok, err := api.RPC.State.GetStorageLatest(key, &computingResource)
	if err != nil || !ok {
		panic(err)
	}

	fmt.Printf("computingResource: %+v\n", computingResource)
}

func TestGetRentalAgreement(t *testing.T) {
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}
	param, _ := types.EncodeToBytes(types.NewU64(0))
	key, err := types.CreateStorageKey(meta, "ResourceOrder", "RentalAgreements", param)
	var data RentalAgreement
	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	assert.True(t, ok)
	assert.NoError(t, err)
	fmt.Printf("price: %+v\n ", data)
}

func TestSomeEvent(t *testing.T) {
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	assert.NoError(t, err)
	meta, err := api.RPC.State.GetMetadataLatest()
	assert.NoError(t, err)

	bh, err := api.RPC.Chain.GetBlockHash(493888)
	assert.NoError(t, err)
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	assert.NoError(t, err)
	raw, err := api.RPC.State.GetStorageRaw(key, bh)
	/// [accountId, index, peerId, cpu, memory, system, cpu_model, price_hour, rent_duration_hour]
	//RegisterResourceSuccess(T::AccountId, u64, Vec<u8>, u64, u64, Vec<u8>, Vec<u8>, Balance, u32),

	// Decode the event records
	events := MyEventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(meta, &events)

	if err != nil {
		fmt.Println(err)
		return

	}
	fmt.Println(len(events.Balances_Transfer))

}

// 连接2个比特数组
func concatBytes(a, b []byte) []byte {
	s := make([]byte, 0)
	for _, n := range a {
		s = append(s, n)
	}
	for _, n := range b {
		s = append(s, n)
	}
	return s
}

// 将types.AccountId　转换为 string 类型的地址
func AccountIdToAddress(id types.AccountID) string {
	s := concatBytes([]byte{42}, id[:])
	hash := blake2b.Sum512(concatBytes([]byte("SS58PRE"), s))
	address := base58.Encode(concatBytes(s, hash[0:2]))
	return address
}

func TestAccountAmount(t *testing.T) {

	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	keypair, _ := signature.KeyringPairFromSecret("betray extend distance category chimney globe employ scrap armor success kiss forum", 42)

	key, err := types.CreateStorageKey(meta, "System", "Account", keypair.PublicKey)
	if err != nil {
		panic(err)
	}

	var accountInfo AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

}

func TestGetQueryMarket(t *testing.T) {
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}
	param, err := types.EncodeToBytes(Provider_MarketUserStatus)
	pk, err := utils.AddressToPublicKey("5Ck8UKvwPkx6ALYib5gZCQu95se6VgDMEvohfQS6gvQ4LtqQ")
	assert.NoError(t, err)
	key, err := types.CreateStorageKey(meta, "Market", "StakerInfo", param, pk)
	assert.NoError(t, err)

	var data MarketUser
	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	assert.True(t, ok)
	assert.NoError(t, err)
	fmt.Printf("data: %+v\n ", data)
}

func TestQueryReward(t *testing.T) {
	api, err := gsrpc.NewSubstrateAPI("ws://183.66.65.207:49944")
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}
	pk, err := utils.AddressToPublicKey("5Ck8UKvwPkx6ALYib5gZCQu95se6VgDMEvohfQS6gvQ4LtqQ")
	assert.NoError(t, err)
	key, err := types.CreateStorageKey(meta, "Market", "ProviderReward", pk)
	fmt.Println(key.Hex())
	assert.NoError(t, err)

	var data MarketIncome
	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	assert.True(t, ok)
	assert.NoError(t, err)
	fmt.Printf("data: %+v\n ", data)
}
