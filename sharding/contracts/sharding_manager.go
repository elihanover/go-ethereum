// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// SMCABI is the input ABI used to generate the binding from.
const SMCABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"collationTrees\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deregisterNotary\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"name\":\"periodHead\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"registerNotary\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"notaryRegistry\",\"outputs\":[{\"name\":\"deregisteredPeriod\",\"type\":\"uint256\"},{\"name\":\"poolIndex\",\"type\":\"uint256\"},{\"name\":\"deposited\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"releaseNotary\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"notaryPool\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"shardId\",\"type\":\"uint256\"},{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getNotaryInCommittee\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_shardId\",\"type\":\"uint256\"},{\"name\":\"period\",\"type\":\"uint256\"},{\"name\":\"chunkRoot\",\"type\":\"bytes32\"},{\"name\":\"proposerAddress\",\"type\":\"address\"}],\"name\":\"addHeader\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"shardId\",\"type\":\"uint256\"},{\"name\":\"parentHash\",\"type\":\"bytes32\"},{\"name\":\"chunkRoot\",\"type\":\"bytes32\"},{\"name\":\"period\",\"type\":\"uint256\"},{\"name\":\"proposerAddress\",\"type\":\"address\"}],\"name\":\"computeHeaderHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"int256\"},{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"collationHeaders\",\"outputs\":[{\"name\":\"shardId\",\"type\":\"uint256\"},{\"name\":\"chunkRoot\",\"type\":\"bytes32\"},{\"name\":\"period\",\"type\":\"uint256\"},{\"name\":\"proposerAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"notaryPoolLength\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"shardId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"chunkRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"period\",\"type\":\"int128\"},{\"indexed\":false,\"name\":\"proposerAddress\",\"type\":\"address\"}],\"name\":\"HeaderAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"notary\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"poolIndex\",\"type\":\"uint256\"}],\"name\":\"NotaryRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"notary\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"poolIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deregisteredPeriod\",\"type\":\"uint256\"}],\"name\":\"NotaryDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"notary\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"poolIndex\",\"type\":\"uint256\"}],\"name\":\"NotaryReleased\",\"type\":\"event\"}]"

// SMCBin is the compiled bytecode used for deploying new contracts.
const SMCBin = `0x608060405234801561001057600080fd5b506108ec806100206000396000f3006080604052600436106100b95763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166303dde97781146100be57806358377bd1146100eb578063584475db1461010257806368e9513e1461011a5780636bdd3271146101225780639910851d14610163578063a81f451014610178578063b09f427e146101ac578063b8bc055f146101c7578063b9505ea4146101f1578063b9d8ef961461021e578063f6f67d3614610268575b600080fd5b3480156100ca57600080fd5b506100d960043560243561027d565b60408051918252519081900360200190f35b3480156100f757600080fd5b5061010061029a565b005b34801561010e57600080fd5b506100d96004356103c5565b6101006103d7565b34801561012e57600080fd5b50610143600160a060020a036004351661058a565b604080519384526020840192909252151582820152519081900360600190f35b34801561016f57600080fd5b506101006105ad565b34801561018457600080fd5b506101906004356106d9565b60408051600160a060020a039092168252519081900360200190f35b3480156101b857600080fd5b50610190600435602435610701565b3480156101d357600080fd5b50610100600435602435604435600160a060020a0360643516610796565b3480156101fd57600080fd5b506100d9600435602435604435606435600160a060020a036084351661079c565b34801561022a57600080fd5b506102396004356024356107a7565b60408051948552602085019390935283830191909152600160a060020a03166060830152519081900360800190f35b34801561027457600080fd5b506100d96107e1565b600260209081526000928352604080842090915290825290205481565b33600160a060020a038116600090815260016020819052604082209081015460029091015490919060ff1615156102d057600080fd5b82600160a060020a03166000838154811015156102e957fe5b600091825260209091200154600160a060020a03161461030857600080fd5b6103106107e7565b50600160a060020a0382166000908152600160205260409020600543049081905561033a8261080a565b600080548390811061034857fe5b600091825260209182902001805473ffffffffffffffffffffffffffffffffffffffff191690556005805460001901905560408051600160a060020a0386168152918201849052818101839052517f90e5afdc8fd31453dcf6e37154fa117ddf3b0324c96c65015563df9d5e4b5a759181900360600190a1505050565b600b6020526000908152604090205481565b33600160a060020a03811660009081526001602052604081206002015460ff161561040157600080fd5b34683635c9adc5dea000001461041657600080fd5b61041e6107e7565b61042661087b565b156104895750600554600080546001810182559080527f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e56301805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0384161790556104d2565b610491610882565b9050816000828154811015156104a357fe5b9060005260206000200160006101000a815481600160a060020a030219169083600160a060020a031602179055505b6005805460019081019091556040805160608101825260008082526020808301868152838501868152600160a060020a0389168452918690529390912091518255915192810192909255516002909101805460ff1916911515919091179055600954811061054257600181016009555b60408051600160a060020a03841681526020810183905281517fa4fe15c53db34d35a5117acc26c27a2653dc68e2dadfc21ed211e38b7864d7a7929181900390910190a15050565b600160208190526000918252604090912080549181015460029091015460ff1683565b33600160a060020a038116600090815260016020819052604090912080820154600290910154909160ff9091161515146105e657600080fd5b600160a060020a038216600090815260016020526040902054151561060a57600080fd5b600160a060020a038216600090815260016020526040902054613f0001600543041161063557600080fd5b600160a060020a03821660008181526001602081905260408083208381559182018390556002909101805460ff1916905551683635c9adc5dea000009082818181858883f19350505050158015610690573d6000803e3d6000fd5b5060408051600160a060020a03841681526020810183905281517faee20171b64b7f3360a142659094ce929970d6963dcea8c34a9bf1ece8033680929181900390910190a15050565b60008054829081106106e757fe5b600091825260209091200154600160a060020a0316905081565b60008080808080600543049450600a54851115610722576009549350610728565b60085493505b6040805160001960058802018040808352602083018b90528284018c9052925191829003606001909120909450909250849081151561076357fe5b06905060008181548110151561077557fe5b600091825260209091200154600160a060020a031698975050505050505050565b50505050565b600095945050505050565b60036020818152600093845260408085209091529183529120805460018201546002830154929093015490929190600160a060020a031684565b60055481565b600a5460054304908110156107fb57610807565b600954600855600a8190555b50565b600754600654141561085057600680546001810182556000919091527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0181905561086f565b80600660075481548110151561086257fe5b6000918252602090912001555b50600780546001019055565b6007541590565b6000600160075411151561089557600080fd5b6007805460001901908190556006805490919081106108b057fe5b90600052602060002001549050905600a165627a7a723058206a18c219259d03462104b10265828f09d0a7edbae0d8cc91dad68df76b9ad2820029`

// DeploySMC deploys a new Ethereum contract, binding an instance of SMC to it.
func DeploySMC(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SMC, error) {
	parsed, err := abi.JSON(strings.NewReader(SMCABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SMCBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SMC{SMCCaller: SMCCaller{contract: contract}, SMCTransactor: SMCTransactor{contract: contract}, SMCFilterer: SMCFilterer{contract: contract}}, nil
}

// SMC is an auto generated Go binding around an Ethereum contract.
type SMC struct {
	SMCCaller     // Read-only binding to the contract
	SMCTransactor // Write-only binding to the contract
	SMCFilterer   // Log filterer for contract events
}

// SMCCaller is an auto generated read-only Go binding around an Ethereum contract.
type SMCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SMCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SMCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SMCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SMCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SMCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SMCSession struct {
	Contract     *SMC              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SMCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SMCCallerSession struct {
	Contract *SMCCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SMCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SMCTransactorSession struct {
	Contract     *SMCTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SMCRaw is an auto generated low-level Go binding around an Ethereum contract.
type SMCRaw struct {
	Contract *SMC // Generic contract binding to access the raw methods on
}

// SMCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SMCCallerRaw struct {
	Contract *SMCCaller // Generic read-only contract binding to access the raw methods on
}

// SMCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SMCTransactorRaw struct {
	Contract *SMCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSMC creates a new instance of SMC, bound to a specific deployed contract.
func NewSMC(address common.Address, backend bind.ContractBackend) (*SMC, error) {
	contract, err := bindSMC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SMC{SMCCaller: SMCCaller{contract: contract}, SMCTransactor: SMCTransactor{contract: contract}, SMCFilterer: SMCFilterer{contract: contract}}, nil
}

// NewSMCCaller creates a new read-only instance of SMC, bound to a specific deployed contract.
func NewSMCCaller(address common.Address, caller bind.ContractCaller) (*SMCCaller, error) {
	contract, err := bindSMC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SMCCaller{contract: contract}, nil
}

// NewSMCTransactor creates a new write-only instance of SMC, bound to a specific deployed contract.
func NewSMCTransactor(address common.Address, transactor bind.ContractTransactor) (*SMCTransactor, error) {
	contract, err := bindSMC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SMCTransactor{contract: contract}, nil
}

// NewSMCFilterer creates a new log filterer instance of SMC, bound to a specific deployed contract.
func NewSMCFilterer(address common.Address, filterer bind.ContractFilterer) (*SMCFilterer, error) {
	contract, err := bindSMC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SMCFilterer{contract: contract}, nil
}

// bindSMC binds a generic wrapper to an already deployed contract.
func bindSMC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SMCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SMC *SMCRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SMC.Contract.SMCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SMC *SMCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SMC.Contract.SMCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SMC *SMCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SMC.Contract.SMCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SMC *SMCCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SMC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SMC *SMCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SMC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SMC *SMCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SMC.Contract.contract.Transact(opts, method, params...)
}

// CollationHeaders is a free data retrieval call binding the contract method 0xb9d8ef96.
//
// Solidity: function collationHeaders( int256,  bytes32) constant returns(shardId uint256, chunkRoot bytes32, period uint256, proposerAddress address)
func (_SMC *SMCCaller) CollationHeaders(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte) (struct {
	ShardId         *big.Int
	ChunkRoot       [32]byte
	Period          *big.Int
	ProposerAddress common.Address
}, error) {
	ret := new(struct {
		ShardId         *big.Int
		ChunkRoot       [32]byte
		Period          *big.Int
		ProposerAddress common.Address
	})
	out := ret
	err := _SMC.contract.Call(opts, out, "collationHeaders", arg0, arg1)
	return *ret, err
}

// CollationHeaders is a free data retrieval call binding the contract method 0xb9d8ef96.
//
// Solidity: function collationHeaders( int256,  bytes32) constant returns(shardId uint256, chunkRoot bytes32, period uint256, proposerAddress address)
func (_SMC *SMCSession) CollationHeaders(arg0 *big.Int, arg1 [32]byte) (struct {
	ShardId         *big.Int
	ChunkRoot       [32]byte
	Period          *big.Int
	ProposerAddress common.Address
}, error) {
	return _SMC.Contract.CollationHeaders(&_SMC.CallOpts, arg0, arg1)
}

// CollationHeaders is a free data retrieval call binding the contract method 0xb9d8ef96.
//
// Solidity: function collationHeaders( int256,  bytes32) constant returns(shardId uint256, chunkRoot bytes32, period uint256, proposerAddress address)
func (_SMC *SMCCallerSession) CollationHeaders(arg0 *big.Int, arg1 [32]byte) (struct {
	ShardId         *big.Int
	ChunkRoot       [32]byte
	Period          *big.Int
	ProposerAddress common.Address
}, error) {
	return _SMC.Contract.CollationHeaders(&_SMC.CallOpts, arg0, arg1)
}

// CollationTrees is a free data retrieval call binding the contract method 0x03dde977.
//
// Solidity: function collationTrees( uint256,  bytes32) constant returns(bytes32)
func (_SMC *SMCCaller) CollationTrees(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _SMC.contract.Call(opts, out, "collationTrees", arg0, arg1)
	return *ret0, err
}

// CollationTrees is a free data retrieval call binding the contract method 0x03dde977.
//
// Solidity: function collationTrees( uint256,  bytes32) constant returns(bytes32)
func (_SMC *SMCSession) CollationTrees(arg0 *big.Int, arg1 [32]byte) ([32]byte, error) {
	return _SMC.Contract.CollationTrees(&_SMC.CallOpts, arg0, arg1)
}

// CollationTrees is a free data retrieval call binding the contract method 0x03dde977.
//
// Solidity: function collationTrees( uint256,  bytes32) constant returns(bytes32)
func (_SMC *SMCCallerSession) CollationTrees(arg0 *big.Int, arg1 [32]byte) ([32]byte, error) {
	return _SMC.Contract.CollationTrees(&_SMC.CallOpts, arg0, arg1)
}

// GetNotaryInCommittee is a free data retrieval call binding the contract method 0xb09f427e.
//
// Solidity: function getNotaryInCommittee(shardId uint256, _index uint256) constant returns(address)
func (_SMC *SMCCaller) GetNotaryInCommittee(opts *bind.CallOpts, shardId *big.Int, _index *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SMC.contract.Call(opts, out, "getNotaryInCommittee", shardId, _index)
	return *ret0, err
}

// GetNotaryInCommittee is a free data retrieval call binding the contract method 0xb09f427e.
//
// Solidity: function getNotaryInCommittee(shardId uint256, _index uint256) constant returns(address)
func (_SMC *SMCSession) GetNotaryInCommittee(shardId *big.Int, _index *big.Int) (common.Address, error) {
	return _SMC.Contract.GetNotaryInCommittee(&_SMC.CallOpts, shardId, _index)
}

// GetNotaryInCommittee is a free data retrieval call binding the contract method 0xb09f427e.
//
// Solidity: function getNotaryInCommittee(shardId uint256, _index uint256) constant returns(address)
func (_SMC *SMCCallerSession) GetNotaryInCommittee(shardId *big.Int, _index *big.Int) (common.Address, error) {
	return _SMC.Contract.GetNotaryInCommittee(&_SMC.CallOpts, shardId, _index)
}

// NotaryPool is a free data retrieval call binding the contract method 0xa81f4510.
//
// Solidity: function notaryPool( uint256) constant returns(address)
func (_SMC *SMCCaller) NotaryPool(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SMC.contract.Call(opts, out, "notaryPool", arg0)
	return *ret0, err
}

// NotaryPool is a free data retrieval call binding the contract method 0xa81f4510.
//
// Solidity: function notaryPool( uint256) constant returns(address)
func (_SMC *SMCSession) NotaryPool(arg0 *big.Int) (common.Address, error) {
	return _SMC.Contract.NotaryPool(&_SMC.CallOpts, arg0)
}

// NotaryPool is a free data retrieval call binding the contract method 0xa81f4510.
//
// Solidity: function notaryPool( uint256) constant returns(address)
func (_SMC *SMCCallerSession) NotaryPool(arg0 *big.Int) (common.Address, error) {
	return _SMC.Contract.NotaryPool(&_SMC.CallOpts, arg0)
}

// NotaryPoolLength is a free data retrieval call binding the contract method 0xf6f67d36.
//
// Solidity: function notaryPoolLength() constant returns(uint256)
func (_SMC *SMCCaller) NotaryPoolLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SMC.contract.Call(opts, out, "notaryPoolLength")
	return *ret0, err
}

// NotaryPoolLength is a free data retrieval call binding the contract method 0xf6f67d36.
//
// Solidity: function notaryPoolLength() constant returns(uint256)
func (_SMC *SMCSession) NotaryPoolLength() (*big.Int, error) {
	return _SMC.Contract.NotaryPoolLength(&_SMC.CallOpts)
}

// NotaryPoolLength is a free data retrieval call binding the contract method 0xf6f67d36.
//
// Solidity: function notaryPoolLength() constant returns(uint256)
func (_SMC *SMCCallerSession) NotaryPoolLength() (*big.Int, error) {
	return _SMC.Contract.NotaryPoolLength(&_SMC.CallOpts)
}

// NotaryRegistry is a free data retrieval call binding the contract method 0x6bdd3271.
//
// Solidity: function notaryRegistry( address) constant returns(deregisteredPeriod uint256, poolIndex uint256, deposited bool)
func (_SMC *SMCCaller) NotaryRegistry(opts *bind.CallOpts, arg0 common.Address) (struct {
	DeregisteredPeriod *big.Int
	PoolIndex          *big.Int
	Deposited          bool
}, error) {
	ret := new(struct {
		DeregisteredPeriod *big.Int
		PoolIndex          *big.Int
		Deposited          bool
	})
	out := ret
	err := _SMC.contract.Call(opts, out, "notaryRegistry", arg0)
	return *ret, err
}

// NotaryRegistry is a free data retrieval call binding the contract method 0x6bdd3271.
//
// Solidity: function notaryRegistry( address) constant returns(deregisteredPeriod uint256, poolIndex uint256, deposited bool)
func (_SMC *SMCSession) NotaryRegistry(arg0 common.Address) (struct {
	DeregisteredPeriod *big.Int
	PoolIndex          *big.Int
	Deposited          bool
}, error) {
	return _SMC.Contract.NotaryRegistry(&_SMC.CallOpts, arg0)
}

// NotaryRegistry is a free data retrieval call binding the contract method 0x6bdd3271.
//
// Solidity: function notaryRegistry( address) constant returns(deregisteredPeriod uint256, poolIndex uint256, deposited bool)
func (_SMC *SMCCallerSession) NotaryRegistry(arg0 common.Address) (struct {
	DeregisteredPeriod *big.Int
	PoolIndex          *big.Int
	Deposited          bool
}, error) {
	return _SMC.Contract.NotaryRegistry(&_SMC.CallOpts, arg0)
}

// PeriodHead is a free data retrieval call binding the contract method 0x584475db.
//
// Solidity: function periodHead( int256) constant returns(int256)
func (_SMC *SMCCaller) PeriodHead(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SMC.contract.Call(opts, out, "periodHead", arg0)
	return *ret0, err
}

// PeriodHead is a free data retrieval call binding the contract method 0x584475db.
//
// Solidity: function periodHead( int256) constant returns(int256)
func (_SMC *SMCSession) PeriodHead(arg0 *big.Int) (*big.Int, error) {
	return _SMC.Contract.PeriodHead(&_SMC.CallOpts, arg0)
}

// PeriodHead is a free data retrieval call binding the contract method 0x584475db.
//
// Solidity: function periodHead( int256) constant returns(int256)
func (_SMC *SMCCallerSession) PeriodHead(arg0 *big.Int) (*big.Int, error) {
	return _SMC.Contract.PeriodHead(&_SMC.CallOpts, arg0)
}

// AddHeader is a paid mutator transaction binding the contract method 0xb8bc055f.
//
// Solidity: function addHeader(_shardId uint256, period uint256, chunkRoot bytes32, proposerAddress address) returns()
func (_SMC *SMCTransactor) AddHeader(opts *bind.TransactOpts, _shardId *big.Int, period *big.Int, chunkRoot [32]byte, proposerAddress common.Address) (*types.Transaction, error) {
	return _SMC.contract.Transact(opts, "addHeader", _shardId, period, chunkRoot, proposerAddress)
}

// AddHeader is a paid mutator transaction binding the contract method 0xb8bc055f.
//
// Solidity: function addHeader(_shardId uint256, period uint256, chunkRoot bytes32, proposerAddress address) returns()
func (_SMC *SMCSession) AddHeader(_shardId *big.Int, period *big.Int, chunkRoot [32]byte, proposerAddress common.Address) (*types.Transaction, error) {
	return _SMC.Contract.AddHeader(&_SMC.TransactOpts, _shardId, period, chunkRoot, proposerAddress)
}

// AddHeader is a paid mutator transaction binding the contract method 0xb8bc055f.
//
// Solidity: function addHeader(_shardId uint256, period uint256, chunkRoot bytes32, proposerAddress address) returns()
func (_SMC *SMCTransactorSession) AddHeader(_shardId *big.Int, period *big.Int, chunkRoot [32]byte, proposerAddress common.Address) (*types.Transaction, error) {
	return _SMC.Contract.AddHeader(&_SMC.TransactOpts, _shardId, period, chunkRoot, proposerAddress)
}

// ComputeHeaderHash is a paid mutator transaction binding the contract method 0xb9505ea4.
//
// Solidity: function computeHeaderHash(shardId uint256, parentHash bytes32, chunkRoot bytes32, period uint256, proposerAddress address) returns(bytes32)
func (_SMC *SMCTransactor) ComputeHeaderHash(opts *bind.TransactOpts, shardId *big.Int, parentHash [32]byte, chunkRoot [32]byte, period *big.Int, proposerAddress common.Address) (*types.Transaction, error) {
	return _SMC.contract.Transact(opts, "computeHeaderHash", shardId, parentHash, chunkRoot, period, proposerAddress)
}

// ComputeHeaderHash is a paid mutator transaction binding the contract method 0xb9505ea4.
//
// Solidity: function computeHeaderHash(shardId uint256, parentHash bytes32, chunkRoot bytes32, period uint256, proposerAddress address) returns(bytes32)
func (_SMC *SMCSession) ComputeHeaderHash(shardId *big.Int, parentHash [32]byte, chunkRoot [32]byte, period *big.Int, proposerAddress common.Address) (*types.Transaction, error) {
	return _SMC.Contract.ComputeHeaderHash(&_SMC.TransactOpts, shardId, parentHash, chunkRoot, period, proposerAddress)
}

// ComputeHeaderHash is a paid mutator transaction binding the contract method 0xb9505ea4.
//
// Solidity: function computeHeaderHash(shardId uint256, parentHash bytes32, chunkRoot bytes32, period uint256, proposerAddress address) returns(bytes32)
func (_SMC *SMCTransactorSession) ComputeHeaderHash(shardId *big.Int, parentHash [32]byte, chunkRoot [32]byte, period *big.Int, proposerAddress common.Address) (*types.Transaction, error) {
	return _SMC.Contract.ComputeHeaderHash(&_SMC.TransactOpts, shardId, parentHash, chunkRoot, period, proposerAddress)
}

// DeregisterNotary is a paid mutator transaction binding the contract method 0x58377bd1.
//
// Solidity: function deregisterNotary() returns()
func (_SMC *SMCTransactor) DeregisterNotary(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SMC.contract.Transact(opts, "deregisterNotary")
}

// DeregisterNotary is a paid mutator transaction binding the contract method 0x58377bd1.
//
// Solidity: function deregisterNotary() returns()
func (_SMC *SMCSession) DeregisterNotary() (*types.Transaction, error) {
	return _SMC.Contract.DeregisterNotary(&_SMC.TransactOpts)
}

// DeregisterNotary is a paid mutator transaction binding the contract method 0x58377bd1.
//
// Solidity: function deregisterNotary() returns()
func (_SMC *SMCTransactorSession) DeregisterNotary() (*types.Transaction, error) {
	return _SMC.Contract.DeregisterNotary(&_SMC.TransactOpts)
}

// RegisterNotary is a paid mutator transaction binding the contract method 0x68e9513e.
//
// Solidity: function registerNotary() returns()
func (_SMC *SMCTransactor) RegisterNotary(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SMC.contract.Transact(opts, "registerNotary")
}

// RegisterNotary is a paid mutator transaction binding the contract method 0x68e9513e.
//
// Solidity: function registerNotary() returns()
func (_SMC *SMCSession) RegisterNotary() (*types.Transaction, error) {
	return _SMC.Contract.RegisterNotary(&_SMC.TransactOpts)
}

// RegisterNotary is a paid mutator transaction binding the contract method 0x68e9513e.
//
// Solidity: function registerNotary() returns()
func (_SMC *SMCTransactorSession) RegisterNotary() (*types.Transaction, error) {
	return _SMC.Contract.RegisterNotary(&_SMC.TransactOpts)
}

// ReleaseNotary is a paid mutator transaction binding the contract method 0x9910851d.
//
// Solidity: function releaseNotary() returns()
func (_SMC *SMCTransactor) ReleaseNotary(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SMC.contract.Transact(opts, "releaseNotary")
}

// ReleaseNotary is a paid mutator transaction binding the contract method 0x9910851d.
//
// Solidity: function releaseNotary() returns()
func (_SMC *SMCSession) ReleaseNotary() (*types.Transaction, error) {
	return _SMC.Contract.ReleaseNotary(&_SMC.TransactOpts)
}

// ReleaseNotary is a paid mutator transaction binding the contract method 0x9910851d.
//
// Solidity: function releaseNotary() returns()
func (_SMC *SMCTransactorSession) ReleaseNotary() (*types.Transaction, error) {
	return _SMC.Contract.ReleaseNotary(&_SMC.TransactOpts)
}

// SMCHeaderAddedIterator is returned from FilterHeaderAdded and is used to iterate over the raw logs and unpacked data for HeaderAdded events raised by the SMC contract.
type SMCHeaderAddedIterator struct {
	Event *SMCHeaderAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SMCHeaderAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SMCHeaderAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SMCHeaderAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SMCHeaderAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SMCHeaderAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SMCHeaderAdded represents a HeaderAdded event raised by the SMC contract.
type SMCHeaderAdded struct {
	ShardId         *big.Int
	ChunkRoot       [32]byte
	Period          *big.Int
	ProposerAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterHeaderAdded is a free log retrieval operation binding the contract event 0x5585108fa6e3c23be15f85641f8d891b72782aaeca2573e7360a4519e67a2bee.
//
// Solidity: event HeaderAdded(shardId indexed uint256, chunkRoot bytes32, period int128, proposerAddress address)
func (_SMC *SMCFilterer) FilterHeaderAdded(opts *bind.FilterOpts, shardId []*big.Int) (*SMCHeaderAddedIterator, error) {

	var shardIdRule []interface{}
	for _, shardIdItem := range shardId {
		shardIdRule = append(shardIdRule, shardIdItem)
	}

	logs, sub, err := _SMC.contract.FilterLogs(opts, "HeaderAdded", shardIdRule)
	if err != nil {
		return nil, err
	}
	return &SMCHeaderAddedIterator{contract: _SMC.contract, event: "HeaderAdded", logs: logs, sub: sub}, nil
}

// WatchHeaderAdded is a free log subscription operation binding the contract event 0x5585108fa6e3c23be15f85641f8d891b72782aaeca2573e7360a4519e67a2bee.
//
// Solidity: event HeaderAdded(shardId indexed uint256, chunkRoot bytes32, period int128, proposerAddress address)
func (_SMC *SMCFilterer) WatchHeaderAdded(opts *bind.WatchOpts, sink chan<- *SMCHeaderAdded, shardId []*big.Int) (event.Subscription, error) {

	var shardIdRule []interface{}
	for _, shardIdItem := range shardId {
		shardIdRule = append(shardIdRule, shardIdItem)
	}

	logs, sub, err := _SMC.contract.WatchLogs(opts, "HeaderAdded", shardIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SMCHeaderAdded)
				if err := _SMC.contract.UnpackLog(event, "HeaderAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SMCNotaryDeregisteredIterator is returned from FilterNotaryDeregistered and is used to iterate over the raw logs and unpacked data for NotaryDeregistered events raised by the SMC contract.
type SMCNotaryDeregisteredIterator struct {
	Event *SMCNotaryDeregistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SMCNotaryDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SMCNotaryDeregistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SMCNotaryDeregistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SMCNotaryDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SMCNotaryDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SMCNotaryDeregistered represents a NotaryDeregistered event raised by the SMC contract.
type SMCNotaryDeregistered struct {
	Notary             common.Address
	PoolIndex          *big.Int
	DeregisteredPeriod *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterNotaryDeregistered is a free log retrieval operation binding the contract event 0x90e5afdc8fd31453dcf6e37154fa117ddf3b0324c96c65015563df9d5e4b5a75.
//
// Solidity: event NotaryDeregistered(notary address, poolIndex uint256, deregisteredPeriod uint256)
func (_SMC *SMCFilterer) FilterNotaryDeregistered(opts *bind.FilterOpts) (*SMCNotaryDeregisteredIterator, error) {

	logs, sub, err := _SMC.contract.FilterLogs(opts, "NotaryDeregistered")
	if err != nil {
		return nil, err
	}
	return &SMCNotaryDeregisteredIterator{contract: _SMC.contract, event: "NotaryDeregistered", logs: logs, sub: sub}, nil
}

// WatchNotaryDeregistered is a free log subscription operation binding the contract event 0x90e5afdc8fd31453dcf6e37154fa117ddf3b0324c96c65015563df9d5e4b5a75.
//
// Solidity: event NotaryDeregistered(notary address, poolIndex uint256, deregisteredPeriod uint256)
func (_SMC *SMCFilterer) WatchNotaryDeregistered(opts *bind.WatchOpts, sink chan<- *SMCNotaryDeregistered) (event.Subscription, error) {

	logs, sub, err := _SMC.contract.WatchLogs(opts, "NotaryDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SMCNotaryDeregistered)
				if err := _SMC.contract.UnpackLog(event, "NotaryDeregistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SMCNotaryRegisteredIterator is returned from FilterNotaryRegistered and is used to iterate over the raw logs and unpacked data for NotaryRegistered events raised by the SMC contract.
type SMCNotaryRegisteredIterator struct {
	Event *SMCNotaryRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SMCNotaryRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SMCNotaryRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SMCNotaryRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SMCNotaryRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SMCNotaryRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SMCNotaryRegistered represents a NotaryRegistered event raised by the SMC contract.
type SMCNotaryRegistered struct {
	Notary    common.Address
	PoolIndex *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNotaryRegistered is a free log retrieval operation binding the contract event 0xa4fe15c53db34d35a5117acc26c27a2653dc68e2dadfc21ed211e38b7864d7a7.
//
// Solidity: event NotaryRegistered(notary address, poolIndex uint256)
func (_SMC *SMCFilterer) FilterNotaryRegistered(opts *bind.FilterOpts) (*SMCNotaryRegisteredIterator, error) {

	logs, sub, err := _SMC.contract.FilterLogs(opts, "NotaryRegistered")
	if err != nil {
		return nil, err
	}
	return &SMCNotaryRegisteredIterator{contract: _SMC.contract, event: "NotaryRegistered", logs: logs, sub: sub}, nil
}

// WatchNotaryRegistered is a free log subscription operation binding the contract event 0xa4fe15c53db34d35a5117acc26c27a2653dc68e2dadfc21ed211e38b7864d7a7.
//
// Solidity: event NotaryRegistered(notary address, poolIndex uint256)
func (_SMC *SMCFilterer) WatchNotaryRegistered(opts *bind.WatchOpts, sink chan<- *SMCNotaryRegistered) (event.Subscription, error) {

	logs, sub, err := _SMC.contract.WatchLogs(opts, "NotaryRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SMCNotaryRegistered)
				if err := _SMC.contract.UnpackLog(event, "NotaryRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SMCNotaryReleasedIterator is returned from FilterNotaryReleased and is used to iterate over the raw logs and unpacked data for NotaryReleased events raised by the SMC contract.
type SMCNotaryReleasedIterator struct {
	Event *SMCNotaryReleased // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SMCNotaryReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SMCNotaryReleased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SMCNotaryReleased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SMCNotaryReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SMCNotaryReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SMCNotaryReleased represents a NotaryReleased event raised by the SMC contract.
type SMCNotaryReleased struct {
	Notary    common.Address
	PoolIndex *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNotaryReleased is a free log retrieval operation binding the contract event 0xaee20171b64b7f3360a142659094ce929970d6963dcea8c34a9bf1ece8033680.
//
// Solidity: event NotaryReleased(notary address, poolIndex uint256)
func (_SMC *SMCFilterer) FilterNotaryReleased(opts *bind.FilterOpts) (*SMCNotaryReleasedIterator, error) {

	logs, sub, err := _SMC.contract.FilterLogs(opts, "NotaryReleased")
	if err != nil {
		return nil, err
	}
	return &SMCNotaryReleasedIterator{contract: _SMC.contract, event: "NotaryReleased", logs: logs, sub: sub}, nil
}

// WatchNotaryReleased is a free log subscription operation binding the contract event 0xaee20171b64b7f3360a142659094ce929970d6963dcea8c34a9bf1ece8033680.
//
// Solidity: event NotaryReleased(notary address, poolIndex uint256)
func (_SMC *SMCFilterer) WatchNotaryReleased(opts *bind.WatchOpts, sink chan<- *SMCNotaryReleased) (event.Subscription, error) {

	logs, sub, err := _SMC.contract.WatchLogs(opts, "NotaryReleased")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SMCNotaryReleased)
				if err := _SMC.contract.UnpackLog(event, "NotaryReleased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
