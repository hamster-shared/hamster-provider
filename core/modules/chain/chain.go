package chain

import (
	"errors"
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/hamster-shared/hamster-provider/core/modules/config"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
)

// ChainClient blockchain chain connection
type ChainClient struct {
	cm  *config.ConfigManager
	api *gsrpc.SubstrateAPI
}

func NewChainClient(cm *config.ConfigManager, api *gsrpc.SubstrateAPI) (*ChainClient, error) {
	return &ChainClient{
		cm:  cm,
		api: api,
	}, nil
}

func (cc *ChainClient) getPeerId() string {
	cf, err := cc.cm.GetConfig()
	if err != nil {
		return ""
	}
	return cf.Identity.PeerID
}

//func (cc *ChainClient) call(c types.Call, meta *types.Metadata) error {
//
//	cf, err := cc.cm.GetConfig()
//
//	// Create the extrinsic
//	ext := types.NewExtrinsic(c)
//	genesisHash, err := cc.api.RPC.Chain.GetBlockHash(0)
//	if err != nil {
//		return err
//	}
//
//	rv, err := cc.api.RPC.State.GetRuntimeVersionLatest()
//	if err != nil {
//		return err
//	}
//
//	keypair, err := signature.KeyringPairFromSecret(cf.SeedOrPhrase, 42)
//	if err != nil {
//		return err
//	}
//
//	// Get the nonce for Account
//	key, err := types.CreateStorageKey(meta, "System", "Account", keypair.PublicKey)
//	if err != nil {
//		return err
//	}
//
//	var accountInfo types.AccountInfo
//	ok, err := cc.api.RPC.State.GetStorageLatest(key, &accountInfo)
//	if err != nil {
//		return err
//	}
//	if !ok {
//		return errors.New("GetStorageLatest fail")
//	}
//
//	nonce := uint32(accountInfo.Nonce)
//	o := types.SignatureOptions{
//		BlockHash:          genesisHash,
//		Era:                types.ExtrinsicEra{IsMortalEra: false},
//		GenesisHash:        genesisHash,
//		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
//		SpecVersion:        rv.SpecVersion,
//		Tip:                types.NewUCompactFromUInt(0),
//		TransactionVersion: rv.TransactionVersion,
//	}
//
//	// Sign the transaction using User's default account
//	err = ext.Sign(keypair, o)
//	if err != nil {
//		return err
//	}
//
//	res, err := cc.api.RPC.Author.SubmitExtrinsic(ext)
//	if err != nil {
//		logrus.Errorf("extrinsic submit failed: %v", err)
//		return err
//	}
//
//	hex, err := types.Hex(res)
//	if err != nil {
//		return err
//	}
//	if hex == "" {
//		return errors.New("hex is empty")
//	}
//	return nil
//}

func (cc *ChainClient) callAndWatch(c types.Call, meta *types.Metadata, hook func(header *types.Header) error) error {

	cf, err := cc.cm.GetConfig()

	// Create the extrinsic
	ext := types.NewExtrinsic(c)
	genesisHash, err := cc.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return err
	}

	rv, err := cc.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return err
	}

	keypair, err := signature.KeyringPairFromSecret(cf.SeedOrPhrase, 42)
	address := keypair.Address
	println(address)
	if err != nil {
		return err
	}

	// Get the nonce for Account
	key, err := types.CreateStorageKey(meta, "System", "Account", keypair.PublicKey)
	if err != nil {
		return err
	}

	var accountInfo types.AccountInfo
	ok, err := cc.api.RPC.State.GetStorageLatest(key, &accountInfo)
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

	sub, err := cc.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer sub.Unsubscribe()

	for {
		status := <-sub.Chan()
		fmt.Printf("Transaction status: %#v\n", status)

		if status.IsInBlock {
			fmt.Printf("Completed at block hash: %#x\n", status.AsInBlock)

			if hook != nil {
				hh, err := cc.api.RPC.Chain.GetHeader(status.AsInBlock)
				if err != nil {
					return err
				}
				return hook(hh)
			}

			return nil
		}

		if status.IsDropped || status.IsInvalid {
			return errors.New("submit tx fail")
		}
	}
}

func (cc *ChainClient) getBlock(blockNumber uint64) {
	cf, _ := cc.cm.GetConfig()
	kp, _ := signature.KeyringPairFromSecret(cf.SeedOrPhrase, 42)
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	hash, err := cc.api.RPC.Chain.GetBlockHash(uint64(blockNumber))
	if err != nil {
		err = fmt.Errorf("get block hash error: %s", err)
		return
	}
	block, err := cc.api.RPC.Chain.GetBlock(hash)
	if err != nil {
		err = fmt.Errorf("get block error: %s", err)
		return
	}
	fmt.Println("blocknumber", block.Block.Header.Number)
	for _, ext := range block.Block.Extrinsics {
		//callIndex, err := meta.FindCallIndex("ResourceOrder.order_exec")
		callIndex, err := meta.FindCallIndex("ResourceOrder.staking_amount")
		if err != nil {
			fmt.Println(err)
			continue
		}
		if ext.Method.CallIndex != callIndex {
			continue
		}

		if string(ext.Signature.Signer.AsID[:]) != string(kp.PublicKey) {
			continue
		}
		fmt.Println("callIndex:", ext.Method.CallIndex.SectionIndex)
		fmt.Println("args:", ext.Method.Args)
		fmt.Println()
	}
}

func (cc *ChainClient) GetEvent(blockNumber uint64) (*MyEventRecords, error) {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	bh, err := cc.api.RPC.Chain.GetBlockHash(blockNumber)
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return nil, err
	}
	raw, err := cc.api.RPC.State.GetStorageRaw(key, bh)
	if err != nil {
		return nil, err
	}
	// Decode the event records
	events := MyEventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(meta, &events)

	return &events, err
}

// Register chain
func (cc *ChainClient) RegisterResource(r ResourceInfo) error {

	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	cf, err := cc.cm.GetConfig()

	if cf.ChainRegInfo.ResourceIndex > 0 {
		resource, err := cc.GetResource(cf.ChainRegInfo.ResourceIndex)
		if err != nil {
			fmt.Println(err)
		}
		if resource != nil {
			return nil
		}
	}

	peerId := cf.Identity.PeerID
	cpu := types.NewU64(r.Cpu)
	memory := types.NewU64(r.Memory)
	system := r.System
	cpuModel := r.CpuModel
	price := types.NewU128(*big.NewInt(int64(r.Price)))
	hours := r.ExpireTime.Sub(time.Now()).Hours()
	rentDurationHour := types.NewU32(uint32(hours))

	c, err := types.NewCall(meta, "Provider.register_resource", peerId, cpu, memory, system, cpuModel, price, rentDurationHour)

	if err != nil {
		return err
	}

	hook := func(header *types.Header) error {
		events, err := cc.GetEvent(uint64(header.Number))
		if err != nil {
			return err
		}
		if len(events.Provider_RegisterResourceSuccess) > 0 {
			for _, e := range events.Provider_RegisterResourceSuccess {
				if e.PeerId == cc.getPeerId() {
					cf.ChainRegInfo.ResourceIndex = uint64(e.Index)
					return cc.cm.Save(cf)
				}
			}
		}

		return errors.New("cannot get Order Index")
	}

	return cc.callAndWatch(c, meta, hook)
}

func (cc *ChainClient) RemoveResource(index uint64) error {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "Provider.remove_resource", types.NewU64(index))

	if err != nil {
		return err
	}

	return cc.callAndWatch(c, meta, nil)
}

func (cc *ChainClient) ChangeResourceStatus(index uint64) error {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}
	c, err := types.NewCall(meta, "Provider.change_resource_status", types.NewU64(index))

	if err != nil {
		return err
	}

	return cc.callAndWatch(c, meta, nil)
}

func (cc *ChainClient) ModifyResourcePrice(index uint64, unitPrice int64) error {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "Provider.modify_resource_price", types.NewU64(index), types.NewU128(*big.NewInt(unitPrice)))

	if err != nil {
		return err
	}

	return cc.callAndWatch(c, meta, nil)
}

func (cc *ChainClient) AddResourceDuration(index uint64, duration int) error {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "Provider.add_resource_duration", types.NewU64(index), types.NewU32(uint32(duration)))

	if err != nil {
		return err
	}

	return cc.callAndWatch(c, meta, nil)
}

func (cc *ChainClient) Heartbeat(agreementindex uint64) error {

	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "ResourceOrder.heartbeat", types.NewU64(agreementindex))

	if err != nil {
		return err
	}

	return cc.callAndWatch(c, meta, nil)
}

// LoadKeyFromChain Get the public Yue of the current node from the chain
func (cc *ChainClient) LoadKeyFromChain() ([]string, error) {
	return []string{}, nil
}

// ReportStatus Timing the virtual machine to report the service status
func (cc *ChainClient) ReportStatus() error {
	return nil
}

// OrderExec order execution
func (cc *ChainClient) OrderExec(orderIndex uint64) error {

	fmt.Printf("orderExec : %d\n", orderIndex)

	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	callName := "ResourceOrder.order_exec"

	c, err := types.NewCall(meta, callName, types.NewU64(orderIndex))

	if err != nil {
		return err
	}

	hook := func(header *types.Header) error {
		// Determine whether the transaction is successfully executed
		err := cc.CheckExtrinsicSuccess(header, callName)
		if err != nil {
			return err
		}

		// get protocol id
		events, err := cc.GetEvent(uint64(header.Number))
		if err != nil {
			return err
		}
		if len(events.ResourceOrder_OrderExecSuccess) > 0 {
			for _, e := range events.ResourceOrder_OrderExecSuccess {
				if uint64(e.OrderIndex) == orderIndex {
					cfg, err := cc.cm.GetConfig()
					if err != nil {
						return err
					}
					cfg.ChainRegInfo.AgreementIndex = uint64(e.AgreementIndex)
					return cc.cm.Save(cfg)
				}
			}
			return errors.New("cannot get agreementIndex")
		} else {
			return errors.New("cannot get agreementIndex")
		}
	}

	return cc.callAndWatch(c, meta, hook)
}

// CheckExtrinsicSuccess verify that the transaction is successful
func (cc *ChainClient) CheckExtrinsicSuccess(header *types.Header, call string) error {

	fmt.Printf("check tx exec Success, blockNumber : %d\n", header.Number)

	cf, _ := cc.cm.GetConfig()
	kp, _ := signature.KeyringPairFromSecret(cf.SeedOrPhrase, 42)
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		logrus.Errorf("get block hash error: %s", err)
		return err
	}
	bh, err := cc.api.RPC.Chain.GetBlockHash(uint64(header.Number))
	if err != nil {
		return err
	}
	block, err := cc.api.RPC.Chain.GetBlock(bh)
	if err != nil {
		logrus.Errorf("get block error: %s", err)
		return err
	}
	extrinsics := block.Block.Extrinsics
	// get the event corresponding to the block
	events, err := cc.GetEvent(uint64(header.Number))

	callIndex, err := meta.FindCallIndex(call)

	for _, e := range events.System_ExtrinsicFailed {
		extrinsicIndex := e.Phase.AsApplyExtrinsic

		extrinsic := extrinsics[extrinsicIndex]
		//who := extrinsic.Signature.Signer.AsID
		if extrinsic.Method.CallIndex == callIndex && string(extrinsic.Signature.Signer.AsID[:]) == string(kp.PublicKey) {
			return err
		}
	}

	return nil
}

// calculate instance expiration time
func (cc *ChainClient) CalculateProviderOverdue(agreementIndex uint64) time.Duration {
	// get current block
	header, err := cc.api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return time.Second
	}
	currentNumber := int64(header.Number)
	agreement, err := cc.getRentalAgreement(agreementIndex)
	if err != nil {
		return time.Second
	}
	overdueNumber := int64(agreement.RentalInfo.EndOfRent)

	duration := overdueNumber - currentNumber

	return time.Duration(int64(time.Second) * duration * 6)
}

func (cc *ChainClient) CalculateInstanceOverdue(orderIndex uint64) time.Duration {
	header, err := cc.api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return time.Second
	}
	currentNumber := int64(header.Number)
	order, err := cc.GetOrder(orderIndex)
	if err != nil {
		return time.Second
	}
	overdueNumber := int64(order.RentDuration) + int64(order.Create)

	duration := overdueNumber - currentNumber

	return time.Duration(int64(time.Second) * duration * 6)
}

func (cc *ChainClient) GetResource(resourceIndex uint64) (*ComputingResource, error) {

	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	bytes, err := types.EncodeToBytes(types.NewU64(resourceIndex))
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(meta, "Provider", "Resources", bytes)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(key.Hex())

	rows, err := cc.api.RPC.State.GetStorageRawLatest(key)
	fmt.Println("rows", len(*rows))
	fmt.Println("err:", err)
	fmt.Println("row:", rows)

	var computingResource ComputingResource

	ok, err := cc.api.RPC.State.GetStorageLatest(key, &computingResource)
	if !ok {
		fmt.Println(err)
		return nil, errors.New("cannot get state with computingResource")
	}

	return &computingResource, err
}

func (cc *ChainClient) CalculateResourceOverdue(expireBlock uint64) (time.Duration, error) {
	header, err := cc.api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return time.Second, err
	}
	currentNumber := int64(header.Number)
	duration := int64(expireBlock) - currentNumber
	return time.Duration(int64(time.Microsecond) * duration * 6), nil
}

func (cc *ChainClient) getRentalAgreement(agreementIndex uint64) (*RentalAgreement, error) {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	param, _ := types.EncodeToBytes(types.NewU64(agreementIndex))
	key, err := types.CreateStorageKey(meta, "ResourceOrder", "RentalAgreements", param)
	var data RentalAgreement
	ok, err := cc.api.RPC.State.GetStorageLatest(key, &data)
	if !ok {
		return nil, errors.New("get rentalAgreement state fail")
	}
	return &data, err
}

func (cc *ChainClient) GetGatewayNodes() ([]string, error) {
	var nodes []string
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nodes, err
	}
	key, err := types.CreateStorageKey(meta, "Gateway", "Gateways")
	var data []string
	ok, err := cc.api.RPC.State.GetStorageLatest(key, &data)
	if !ok {
		return data, errors.New("gateway nodes is empty")
	}
	return data, err
}

func (cc *ChainClient) GetOrder(orderIndex uint64) (*ComputingOrder, error) {

	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	bytes, err := types.EncodeToBytes(types.NewU64(orderIndex))
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(meta, "ResourceOrder", "ResourceOrders", bytes)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(key.Hex())

	var order ComputingOrder
	ok, err := cc.api.RPC.State.GetStorageLatest(key, &order)
	if !ok {
		return nil, errors.New("cannot get state with computingResource")
	}

	return &order, err
}

func (cc *ChainClient) GetAgreementIndex(orderIndex uint64) (uint64, error) {
	order, err := cc.GetOrder(orderIndex)
	if err != nil {
		return 0, err
	}
	if order.AgreementIndex.IsSome() {
		ok, value := order.AgreementIndex.Unwrap()
		if ok {
			return uint64(value), nil
		}
	}
	return 0, errors.New("no agreementIndex")
}

func (cc *ChainClient) ReceiveIncome() error {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}
	config, err := cc.cm.GetConfig()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "ResourceOrder.withdraw_rental_amount", types.NewU64(config.ChainRegInfo.AgreementIndex))

	if err != nil {
		return err
	}

	return cc.callAndWatch(c, meta, nil)
}

func (cc *ChainClient) GetAccountInfo() (*AccountInfo, error) {
	cf, err := cc.cm.GetConfig()
	if err != nil {
		return nil, err
	}
	keypair, err := signature.KeyringPairFromSecret(cf.SeedOrPhrase, 42)
	if err != nil {
		return nil, err
	}
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	// Get the nonce for Account
	key, err := types.CreateStorageKey(meta, "System", "Account", keypair.PublicKey)
	if err != nil {
		return nil, err
	}

	var account AccountInfoCustom
	ok, err := cc.api.RPC.State.GetStorageLatest(key, &account)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("failed to get account information")
	}
	var accountInfo AccountInfo
	accountInfo.Address = keypair.Address
	accountInfo.Amount = account.Data.Free
	return &accountInfo, nil
}

func (cc *ChainClient) GetStakingInfo() (*StakingAmount, error) {
	cf, err := cc.cm.GetConfig()
	if err != nil {
		return nil, err
	}
	keypair, err := signature.KeyringPairFromSecret(cf.SeedOrPhrase, 42)
	if err != nil {
		return nil, err
	}
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	key, err := types.CreateStorageKey(meta, "ResourceOrder", "Staking", keypair.PublicKey)
	if err != nil {
		return nil, err
	}
	var stakingInfo StakingAmount
	ok, err := cc.api.RPC.State.GetStorageLatest(key, &stakingInfo)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("failed to get account information")
	}
	return &stakingInfo, nil
}

func (cc *ChainClient) StakingAmount(unitPrice int64) error {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "ResourceOrder.staking_amount", types.NewU128(*big.NewInt(unitPrice)))

	if err != nil {
		return err
	}

	return cc.callAndWatch(c, meta, nil)
}

func (cc *ChainClient) WithdrawStakingAmount(unitPrice int64) error {
	meta, err := cc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}

	c, err := types.NewCall(meta, "ResourceOrder.withdraw_amount", types.NewU128(*big.NewInt(unitPrice)))

	if err != nil {
		return err
	}

	return cc.callAndWatch(c, meta, nil)
}

func (cc *ChainClient) ReceiveIncomeJudge() bool {
	config, err := cc.cm.GetConfig()
	if err != nil {
		return false
	}
	agreement, err := cc.getRentalAgreement(config.ChainRegInfo.AgreementIndex)
	if err != nil {
		return false
	}
	price := agreement.ReceiveAmount.Int64()
	if price > 0 {
		return true
	}
	return false
}
