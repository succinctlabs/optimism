// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/ethereum-optimism/optimism/op-bindings/solc"
)

const FaultDisputeGameStorageLayoutJSON = "{\"storage\":[{\"astId\":1000,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"gameStart\",\"offset\":0,\"slot\":\"0\",\"type\":\"t_userDefinedValueType(Timestamp)1014\"},{\"astId\":1001,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"status\",\"offset\":8,\"slot\":\"0\",\"type\":\"t_enum(GameStatus)1007\"},{\"astId\":1002,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"bondManager\",\"offset\":9,\"slot\":\"0\",\"type\":\"t_contract(IBondManager)1006\"},{\"astId\":1003,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"l1Head\",\"offset\":0,\"slot\":\"1\",\"type\":\"t_userDefinedValueType(Hash)1012\"},{\"astId\":1004,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"claimData\",\"offset\":0,\"slot\":\"2\",\"type\":\"t_array(t_struct(ClaimData)1008_storage)dyn_storage\"},{\"astId\":1005,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"claims\",\"offset\":0,\"slot\":\"3\",\"type\":\"t_mapping(t_userDefinedValueType(ClaimHash)1010,t_bool)\"}],\"types\":{\"t_array(t_struct(ClaimData)1008_storage)dyn_storage\":{\"encoding\":\"dynamic_array\",\"label\":\"struct IFaultDisputeGame.ClaimData[]\",\"numberOfBytes\":\"32\",\"base\":\"t_struct(ClaimData)1008_storage\"},\"t_bool\":{\"encoding\":\"inplace\",\"label\":\"bool\",\"numberOfBytes\":\"1\"},\"t_contract(IBondManager)1006\":{\"encoding\":\"inplace\",\"label\":\"contract IBondManager\",\"numberOfBytes\":\"20\"},\"t_enum(GameStatus)1007\":{\"encoding\":\"inplace\",\"label\":\"enum GameStatus\",\"numberOfBytes\":\"1\"},\"t_mapping(t_userDefinedValueType(ClaimHash)1010,t_bool)\":{\"encoding\":\"mapping\",\"label\":\"mapping(ClaimHash =\u003e bool)\",\"numberOfBytes\":\"32\",\"key\":\"t_userDefinedValueType(ClaimHash)1010\",\"value\":\"t_bool\"},\"t_struct(ClaimData)1008_storage\":{\"encoding\":\"inplace\",\"label\":\"struct IFaultDisputeGame.ClaimData\",\"numberOfBytes\":\"96\"},\"t_uint32\":{\"encoding\":\"inplace\",\"label\":\"uint32\",\"numberOfBytes\":\"4\"},\"t_userDefinedValueType(Claim)1009\":{\"encoding\":\"inplace\",\"label\":\"Claim\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(ClaimHash)1010\":{\"encoding\":\"inplace\",\"label\":\"ClaimHash\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(Clock)1011\":{\"encoding\":\"inplace\",\"label\":\"Clock\",\"numberOfBytes\":\"16\"},\"t_userDefinedValueType(Hash)1012\":{\"encoding\":\"inplace\",\"label\":\"Hash\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(Position)1013\":{\"encoding\":\"inplace\",\"label\":\"Position\",\"numberOfBytes\":\"16\"},\"t_userDefinedValueType(Timestamp)1014\":{\"encoding\":\"inplace\",\"label\":\"Timestamp\",\"numberOfBytes\":\"8\"}}}"

var FaultDisputeGameStorageLayout = new(solc.StorageLayout)

var FaultDisputeGameDeployedBin = "0x6080604052600436106101965760003560e01c80638129fc1c116100e1578063c0c3a0921161008a578063c6f0308c11610064578063c6f0308c14610502578063cf09e0d014610566578063d8cc1a3c14610585578063fa24f743146105a557600080fd5b8063c0c3a09214610487578063c31b29ce146104bb578063c55cd0c7146104ef57600080fd5b806392931298116100bb57806392931298146103fa578063bbdc02db1461042e578063bcef3b551461044a57600080fd5b80638129fc1c146103905780638980e0cc146103a55780638b85902b146103ba57600080fd5b8063363cc42711610143578063609d33341161011d578063609d333414610352578063632247ea146103675780636361506d1461037a57600080fd5b8063363cc4271461029d5780634778efe8146102fc57806354fd4d501461033057600080fd5b80632810e1d6116101745780632810e1d61461023b5780633218b99d1461025057806335fef5671461028a57600080fd5b80631e27052a1461019b578063200d2ed2146101bd578063266198f9146101f9575b600080fd5b3480156101a757600080fd5b506101bb6101b636600461217c565b6105c9565b005b3480156101c957600080fd5b506000546101e39068010000000000000000900460ff1681565b6040516101f091906121cd565b60405180910390f35b34801561020557600080fd5b5061022d7f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020016101f0565b34801561024757600080fd5b506101e3610ad2565b34801561025c57600080fd5b506000546102719067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016101f0565b6101bb61029836600461217c565b610ef3565b3480156102a957600080fd5b506000546102d7906901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f0565b34801561030857600080fd5b5061022d7f000000000000000000000000000000000000000000000000000000000000000081565b34801561033c57600080fd5b50610345610f03565b6040516101f09190612284565b34801561035e57600080fd5b50610345610fa6565b6101bb6103753660046122b3565b610fb8565b34801561038657600080fd5b5061022d60015481565b34801561039c57600080fd5b506101bb6115d4565b3480156103b157600080fd5b5060025461022d565b3480156103c657600080fd5b50367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90036020013561022d565b34801561040657600080fd5b506102d77f000000000000000000000000000000000000000000000000000000000000000081565b34801561043a57600080fd5b50604051600081526020016101f0565b34801561045657600080fd5b50367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90033561022d565b34801561049357600080fd5b506102d77f000000000000000000000000000000000000000000000000000000000000000081565b3480156104c757600080fd5b506102717f000000000000000000000000000000000000000000000000000000000000000081565b6101bb6104fd36600461217c565b611725565b34801561050e57600080fd5b5061052261051d3660046122e8565b611731565b6040805163ffffffff90961686529315156020860152928401919091526fffffffffffffffffffffffffffffffff908116606084015216608082015260a0016101f0565b34801561057257600080fd5b5060005467ffffffffffffffff16610271565b34801561059157600080fd5b506101bb6105a036600461234a565b6117a2565b3480156105b157600080fd5b506105ba611cc6565b6040516101f0939291906123d4565b6000805468010000000000000000900460ff1660028111156105ed576105ed61219e565b14610624576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16637dc0d1d06040518163ffffffff1660e01b8152600401602060405180830381865afa158015610691573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106b591906123ff565b905082600103610770576001546040517f9a1f5e7f000000000000000000000000000000000000000000000000000000008152600481018590526024810191909152602060448201526064810183905273ffffffffffffffffffffffffffffffffffffffff821690639a1f5e7f906084015b6020604051808303816000875af1158015610746573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061076a9190612435565b50505050565b82600203610909576040517fcf8e5cf0000000000000000000000000000000000000000000000000000000008152367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c900360200135600482015260009073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063cf8e5cf090602401606060405180830381865afa158015610832573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610856919061249d565b80516040517f9a1f5e7f000000000000000000000000000000000000000000000000000000008152600481018790526024810191909152602060448201526064810185905290915073ffffffffffffffffffffffffffffffffffffffff831690639a1f5e7f906084016020604051808303816000875af11580156108de573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109029190612435565b5050505050565b826003036109a2576040517f9a1f5e7f00000000000000000000000000000000000000000000000000000000815260048101849052367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c9003356024820152602060448201526064810183905273ffffffffffffffffffffffffffffffffffffffff821690639a1f5e7f90608401610727565b82600403610a41576040517f9a1f5e7f00000000000000000000000000000000000000000000000000000000815260048101849052367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90036020013560c01b6024820152600860448201526064810183905273ffffffffffffffffffffffffffffffffffffffff821690639a1f5e7f90608401610727565b82600503610acd576040517f9a1f5e7f000000000000000000000000000000000000000000000000000000008152600481018490524660c01b6024820152600860448201526064810183905273ffffffffffffffffffffffffffffffffffffffff821690639a1f5e7f906084016020604051808303816000875af1158015610746573d6000803e3d6000fd5b505050565b60008060005468010000000000000000900460ff166002811115610af857610af861219e565b14610b2f576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600254600090610b4190600190612558565b90506fffffffffffffffffffffffffffffffff815b67ffffffffffffffff811015610c2b57600060028281548110610b7b57610b7b61256f565b6000918252602090912060039091020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9093019290915060ff6401000000009091041615610bcc5750610b56565b6002810154600090610c10906fffffffffffffffffffffffffffffffff167f0000000000000000000000000000000000000000000000000000000000000000611d04565b905083811015610c24578093508260010194505b5050610b56565b50600060028381548110610c4157610c4161256f565b600091825260208220600390910201805490925063ffffffff90811691908214610cab5760028281548110610c7857610c7861256f565b906000526020600020906003020160020160109054906101000a90046fffffffffffffffffffffffffffffffff16610cd7565b600283015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff165b9050677fffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000060011c16610d1b67ffffffffffffffff831642612558565b610d37836fffffffffffffffffffffffffffffffff1660401c90565b67ffffffffffffffff16610d4b919061259e565b11610d82576040517ff2440b5300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600283810154610e24906fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b610e2e91906125e5565b67ffffffffffffffff16158015610e5557506fffffffffffffffffffffffffffffffff8414155b15610e635760029550610e68565b600195505b600080548791907fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff1668010000000000000000836002811115610ead57610ead61219e565b021790556002811115610ec257610ec261219e565b6040517f5e186f09b9c93491f14e277eea7faa5de6a2d4bda75a79af7a3684fbfb42da6090600090a2505050505090565b610eff82826000610fb8565b5050565b6060610f2e7f0000000000000000000000000000000000000000000000000000000000000000611db9565b610f577f0000000000000000000000000000000000000000000000000000000000000000611db9565b610f807f0000000000000000000000000000000000000000000000000000000000000000611db9565b604051602001610f929392919061260c565b604051602081830303815290604052905090565b6060610fb3602080611ef6565b905090565b6000805468010000000000000000900460ff166002811115610fdc57610fdc61219e565b14611013576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8215801561101f575080155b15611056576040517fa42637bc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006002848154811061106b5761106b61256f565b60009182526020918290206040805160a0810182526003909302909101805463ffffffff8116845260ff64010000000090910416151593830193909352600180840154918301919091526002928301546fffffffffffffffffffffffffffffffff808216606085015270010000000000000000000000000000000090910416608083015282549193509190869081106111065761110661256f565b6000918252602082206003909102018054921515640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff909316929092179091556060820151611170906fffffffffffffffffffffffffffffffff1684151760011b90565b90507f000000000000000000000000000000000000000000000000000000000000000061122f826fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff161115611271576040517f56f57b2b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815160009063ffffffff908116146112d1576002836000015163ffffffff16815481106112a0576112a061256f565b906000526020600020906003020160020160109054906101000a90046fffffffffffffffffffffffffffffffff1690505b608083015160009067ffffffffffffffff1667ffffffffffffffff164261130a846fffffffffffffffffffffffffffffffff1660401c90565b67ffffffffffffffff1661131e919061259e565b6113289190612558565b9050677fffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000060011c1667ffffffffffffffff8216111561139b576040517f3381d11400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000604082901b4217905060006113bc888660009182526020526040902090565b60008181526003602052604090205490915060ff1615611408576040517f80497e3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016003600083815260200190815260200160002060006101000a81548160ff02191690831515021790555060026040518060a001604052808b63ffffffff1681526020016000151581526020018a8152602001876fffffffffffffffffffffffffffffffff168152602001846fffffffffffffffffffffffffffffffff16815250908060018154018082558091505060019003906000526020600020906003020160009091909190915060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548160ff0219169083151502179055506040820151816001015560608201518160020160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555060808201518160020160106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555050503373ffffffffffffffffffffffffffffffffffffffff16888a7f9b3245740ec3b155098a55be84957a4da13eaf7f14a8bc6f53126c0b9350f2be60405160405180910390a4505050505050505050565b600080547fffffffffffffffffffffffffffffffffffffffffffffff000000000000000000164267ffffffffffffffff161781556040805160a08101825263ffffffff815260208101929092526002919081016116597ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe369081013560f01c90033590565b815260016020820152604001426fffffffffffffffffffffffffffffffff908116909152825460018181018555600094855260209485902084516003909302018054958501511515640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000090961663ffffffff909316929092179490941781556040830151818501556060830151608090930151821670010000000000000000000000000000000002929091169190911760029091015561171f9043612558565b40600155565b610eff82826001610fb8565b6002818154811061174157600080fd5b600091825260209091206003909102018054600182015460029092015463ffffffff8216935064010000000090910460ff1691906fffffffffffffffffffffffffffffffff8082169170010000000000000000000000000000000090041685565b6000805468010000000000000000900460ff1660028111156117c6576117c661219e565b146117fd576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600287815481106118125761181261256f565b6000918252602082206003919091020160028101549092506fffffffffffffffffffffffffffffffff16908715821760011b90506118717f0000000000000000000000000000000000000000000000000000000000000000600161259e565b61190d826fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff161461194e576040517f5f53dd9800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008089156119d457611972836fffffffffffffffffffffffffffffffff16611f8d565b67ffffffffffffffff166000036119ab577f000000000000000000000000000000000000000000000000000000000000000091506119cd565b6119c66119b9600186612682565b865463ffffffff16612033565b6001015491505b50836119ee565b846001015491506119eb8460016119b991906126b3565b90505b8189896040516119ff9291906126e7565b604051809103902014611a3e576040517f696550ff00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081600101547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f8e0cb968c8c8c8c6040518563ffffffff1660e01b8152600401611aa49493929190612740565b6020604051808303816000875af1158015611ac3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611ae79190612435565b600284810154929091149250600091611b92906fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b611c2e886fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b611c389190612772565b611c4291906125e5565b67ffffffffffffffff161590508115158103611c8a576040517ffb4e40dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505084547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff166401000000001790945550505050505050505050565b6000367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c9003356060611cfd610fa6565b9050909192565b600080611d91847e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff1690508083036001841b600180831b0386831b17039250505092915050565b606081600003611dfc57505060408051808201909152600181527f3000000000000000000000000000000000000000000000000000000000000000602082015290565b8160005b8115611e265780611e1081612793565b9150611e1f9050600a836127cb565b9150611e00565b60008167ffffffffffffffff811115611e4157611e4161244e565b6040519080825280601f01601f191660200182016040528015611e6b576020820181803683370190505b5090505b8415611eee57611e80600183612558565b9150611e8d600a866127df565b611e9890603061259e565b60f81b818381518110611ead57611ead61256f565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350611ee7600a866127cb565b9450611e6f565b949350505050565b60606000611f2d84367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c900361259e565b90508267ffffffffffffffff1667ffffffffffffffff811115611f5257611f5261244e565b6040519080825280601f01601f191660200182016040528015611f7c576020820181803683370190505b509150828160208401375092915050565b60008061201a837e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b600167ffffffffffffffff919091161b90920392915050565b600080612051846fffffffffffffffffffffffffffffffff166120d0565b9050600283815481106120665761206661256f565b906000526020600020906003020191505b60028201546fffffffffffffffffffffffffffffffff8281169116146120c957815460028054909163ffffffff169081106120b4576120b461256f565b90600052602060002090600302019150612077565b5092915050565b60008119600183011681612164827e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff169390931c8015179392505050565b6000806040838503121561218f57600080fd5b50508035926020909101359150565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6020810160038310612208577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b60005b83811015612229578181015183820152602001612211565b8381111561076a5750506000910152565b6000815180845261225281602086016020860161220e565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612297602083018461223a565b9392505050565b803580151581146122ae57600080fd5b919050565b6000806000606084860312156122c857600080fd5b83359250602084013591506122df6040850161229e565b90509250925092565b6000602082840312156122fa57600080fd5b5035919050565b60008083601f84011261231357600080fd5b50813567ffffffffffffffff81111561232b57600080fd5b60208301915083602082850101111561234357600080fd5b9250929050565b6000806000806000806080878903121561236357600080fd5b863595506123736020880161229e565b9450604087013567ffffffffffffffff8082111561239057600080fd5b61239c8a838b01612301565b909650945060608901359150808211156123b557600080fd5b506123c289828a01612301565b979a9699509497509295939492505050565b60ff841681528260208201526060604082015260006123f6606083018461223a565b95945050505050565b60006020828403121561241157600080fd5b815173ffffffffffffffffffffffffffffffffffffffff8116811461229757600080fd5b60006020828403121561244757600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b80516fffffffffffffffffffffffffffffffff811681146122ae57600080fd5b6000606082840312156124af57600080fd5b6040516060810181811067ffffffffffffffff821117156124f9577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040528251815261250c6020840161247d565b602082015261251d6040840161247d565b60408201529392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008282101561256a5761256a612529565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082198211156125b1576125b1612529565b500190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600067ffffffffffffffff80841680612600576126006125b6565b92169190910692915050565b6000845161261e81846020890161220e565b80830190507f2e00000000000000000000000000000000000000000000000000000000000000808252855161265a816001850160208a0161220e565b6001920191820152835161267581600284016020880161220e565b0160020195945050505050565b60006fffffffffffffffffffffffffffffffff838116908316818110156126ab576126ab612529565b039392505050565b60006fffffffffffffffffffffffffffffffff8083168185168083038211156126de576126de612529565b01949350505050565b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6040815260006127546040830186886126f7565b82810360208401526127678185876126f7565b979650505050505050565b600067ffffffffffffffff838116908316818110156126ab576126ab612529565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036127c4576127c4612529565b5060010190565b6000826127da576127da6125b6565b500490565b6000826127ee576127ee6125b6565b50069056fea164736f6c634300080f000a"

func init() {
	if err := json.Unmarshal([]byte(FaultDisputeGameStorageLayoutJSON), FaultDisputeGameStorageLayout); err != nil {
		panic(err)
	}

	layouts["FaultDisputeGame"] = FaultDisputeGameStorageLayout
	deployedBytecodes["FaultDisputeGame"] = FaultDisputeGameDeployedBin
}
