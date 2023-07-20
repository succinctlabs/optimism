// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/ethereum-optimism/optimism/op-bindings/solc"
)

const FaultDisputeGameStorageLayoutJSON = "{\"storage\":[{\"astId\":1000,\"contract\":\"contracts/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"gameStart\",\"offset\":0,\"slot\":\"0\",\"type\":\"t_userDefinedValueType(Timestamp)1014\"},{\"astId\":1001,\"contract\":\"contracts/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"status\",\"offset\":8,\"slot\":\"0\",\"type\":\"t_enum(GameStatus)1007\"},{\"astId\":1002,\"contract\":\"contracts/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"bondManager\",\"offset\":9,\"slot\":\"0\",\"type\":\"t_contract(IBondManager)1006\"},{\"astId\":1003,\"contract\":\"contracts/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"l1Head\",\"offset\":0,\"slot\":\"1\",\"type\":\"t_userDefinedValueType(Hash)1012\"},{\"astId\":1004,\"contract\":\"contracts/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"claimData\",\"offset\":0,\"slot\":\"2\",\"type\":\"t_array(t_struct(ClaimData)1008_storage)dyn_storage\"},{\"astId\":1005,\"contract\":\"contracts/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"claims\",\"offset\":0,\"slot\":\"3\",\"type\":\"t_mapping(t_userDefinedValueType(ClaimHash)1010,t_bool)\"}],\"types\":{\"t_array(t_struct(ClaimData)1008_storage)dyn_storage\":{\"encoding\":\"dynamic_array\",\"label\":\"struct IFaultDisputeGame.ClaimData[]\",\"numberOfBytes\":\"32\",\"base\":\"t_struct(ClaimData)1008_storage\"},\"t_bool\":{\"encoding\":\"inplace\",\"label\":\"bool\",\"numberOfBytes\":\"1\"},\"t_contract(IBondManager)1006\":{\"encoding\":\"inplace\",\"label\":\"contract IBondManager\",\"numberOfBytes\":\"20\"},\"t_enum(GameStatus)1007\":{\"encoding\":\"inplace\",\"label\":\"enum GameStatus\",\"numberOfBytes\":\"1\"},\"t_mapping(t_userDefinedValueType(ClaimHash)1010,t_bool)\":{\"encoding\":\"mapping\",\"label\":\"mapping(ClaimHash =\u003e bool)\",\"numberOfBytes\":\"32\",\"key\":\"t_userDefinedValueType(ClaimHash)1010\",\"value\":\"t_bool\"},\"t_struct(ClaimData)1008_storage\":{\"encoding\":\"inplace\",\"label\":\"struct IFaultDisputeGame.ClaimData\",\"numberOfBytes\":\"96\"},\"t_uint32\":{\"encoding\":\"inplace\",\"label\":\"uint32\",\"numberOfBytes\":\"4\"},\"t_userDefinedValueType(Claim)1009\":{\"encoding\":\"inplace\",\"label\":\"Claim\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(ClaimHash)1010\":{\"encoding\":\"inplace\",\"label\":\"ClaimHash\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(Clock)1011\":{\"encoding\":\"inplace\",\"label\":\"Clock\",\"numberOfBytes\":\"16\"},\"t_userDefinedValueType(Hash)1012\":{\"encoding\":\"inplace\",\"label\":\"Hash\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(Position)1013\":{\"encoding\":\"inplace\",\"label\":\"Position\",\"numberOfBytes\":\"16\"},\"t_userDefinedValueType(Timestamp)1014\":{\"encoding\":\"inplace\",\"label\":\"Timestamp\",\"numberOfBytes\":\"8\"}}}"

var FaultDisputeGameStorageLayout = new(solc.StorageLayout)

var FaultDisputeGameDeployedBin = "0x6080604052600436106101965760003560e01c80638980e0cc116100e1578063c31b29ce1161008a578063cf09e0d011610064578063cf09e0d014610546578063d8cc1a3c14610565578063f05a6c3914610585578063fa24f743146105a557600080fd5b8063c31b29ce1461049b578063c55cd0c7146104cf578063c6f0308c146104e257600080fd5b8063bbdc02db116100bb578063bbdc02db1461040e578063bcef3b551461042a578063c0c3a0921461046757600080fd5b80638980e0cc146103855780638b85902b1461039a57806392931298146103da57600080fd5b80634778efe811610143578063632247ea1161011d578063632247ea146103475780636361506d1461035a5780638129fc1c1461037057600080fd5b80634778efe8146102dc57806354fd4d5014610310578063609d33341461033257600080fd5b80633218b99d116101745780633218b99d1461022e57806335fef56714610268578063363cc4271461027d57600080fd5b8063200d2ed21461019b578063266198f9146101d75780632810e1d614610219575b600080fd5b3480156101a757600080fd5b506000546101c19068010000000000000000900460ff1681565b6040516101ce9190612143565b60405180910390f35b3480156101e357600080fd5b5061020b7f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020016101ce565b34801561022557600080fd5b506101c16105c9565b34801561023a57600080fd5b5060005461024f9067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016101ce565b61027b610276366004612184565b6109ea565b005b34801561028957600080fd5b506000546102b7906901000000000000000000900473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101ce565b3480156102e857600080fd5b5061020b7f000000000000000000000000000000000000000000000000000000000000000081565b34801561031c57600080fd5b506103256109fa565b6040516101ce9190612220565b34801561033e57600080fd5b50610325610a9d565b61027b61035536600461224f565b610aaf565b34801561036657600080fd5b5061020b60015481565b34801561037c57600080fd5b5061027b6110cb565b34801561039157600080fd5b5060025461020b565b3480156103a657600080fd5b50367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90036020013561020b565b3480156103e657600080fd5b506102b77f000000000000000000000000000000000000000000000000000000000000000081565b34801561041a57600080fd5b50604051600081526020016101ce565b34801561043657600080fd5b50367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90033561020b565b34801561047357600080fd5b506102b77f000000000000000000000000000000000000000000000000000000000000000081565b3480156104a757600080fd5b5061024f7f000000000000000000000000000000000000000000000000000000000000000081565b61027b6104dd366004612184565b61121c565b3480156104ee57600080fd5b506105026104fd366004612284565b611228565b6040805163ffffffff90961686529315156020860152928401919091526fffffffffffffffffffffffffffffffff908116606084015216608082015260a0016101ce565b34801561055257600080fd5b5060005467ffffffffffffffff1661024f565b34801561057157600080fd5b5061027b6105803660046122e6565b611299565b34801561059157600080fd5b5061027b6105a0366004612284565b6117bd565b3480156105b157600080fd5b506105ba611c5e565b6040516101ce93929190612370565b60008060005468010000000000000000900460ff1660028111156105ef576105ef612114565b14610626576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600254600090610638906001906123ca565b90506fffffffffffffffffffffffffffffffff815b67ffffffffffffffff81101561072257600060028281548110610672576106726123e1565b6000918252602090912060039091020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9093019290915060ff64010000000090910416156106c3575061064d565b6002810154600090610707906fffffffffffffffffffffffffffffffff167f0000000000000000000000000000000000000000000000000000000000000000611c9c565b90508381101561071b578093508260010194505b505061064d565b50600060028381548110610738576107386123e1565b600091825260208220600390910201805490925063ffffffff908116919082146107a2576002828154811061076f5761076f6123e1565b906000526020600020906003020160020160109054906101000a90046fffffffffffffffffffffffffffffffff166107ce565b600283015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff165b9050677fffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000060011c1661081267ffffffffffffffff8316426123ca565b61082e836fffffffffffffffffffffffffffffffff1660401c90565b67ffffffffffffffff166108429190612410565b11610879576040517ff2440b5300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60028381015461091b906fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b6109259190612457565b67ffffffffffffffff1615801561094c57506fffffffffffffffffffffffffffffffff8414155b1561095a576002955061095f565b600195505b600080548791907fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000008360028111156109a4576109a4612114565b0217905560028111156109b9576109b9612114565b6040517f5e186f09b9c93491f14e277eea7faa5de6a2d4bda75a79af7a3684fbfb42da6090600090a2505050505090565b6109f682826000610aaf565b5050565b6060610a257f0000000000000000000000000000000000000000000000000000000000000000611d51565b610a4e7f0000000000000000000000000000000000000000000000000000000000000000611d51565b610a777f0000000000000000000000000000000000000000000000000000000000000000611d51565b604051602001610a899392919061247e565b604051602081830303815290604052905090565b6060610aaa602080611e8e565b905090565b6000805468010000000000000000900460ff166002811115610ad357610ad3612114565b14610b0a576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b82158015610b16575080155b15610b4d576040517fa42637bc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060028481548110610b6257610b626123e1565b60009182526020918290206040805160a0810182526003909302909101805463ffffffff8116845260ff64010000000090910416151593830193909352600180840154918301919091526002928301546fffffffffffffffffffffffffffffffff80821660608501527001000000000000000000000000000000009091041660808301528254919350919086908110610bfd57610bfd6123e1565b6000918252602082206003909102018054921515640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff909316929092179091556060820151610c67906fffffffffffffffffffffffffffffffff1684151760011b90565b90507f0000000000000000000000000000000000000000000000000000000000000000610d26826fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff161115610d68576040517f56f57b2b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815160009063ffffffff90811614610dc8576002836000015163ffffffff1681548110610d9757610d976123e1565b906000526020600020906003020160020160109054906101000a90046fffffffffffffffffffffffffffffffff1690505b608083015160009067ffffffffffffffff1667ffffffffffffffff1642610e01846fffffffffffffffffffffffffffffffff1660401c90565b67ffffffffffffffff16610e159190612410565b610e1f91906123ca565b9050677fffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000060011c1667ffffffffffffffff82161115610e92576040517f3381d11400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000604082901b421790506000610eb3888660009182526020526040902090565b60008181526003602052604090205490915060ff1615610eff576040517f80497e3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016003600083815260200190815260200160002060006101000a81548160ff02191690831515021790555060026040518060a001604052808b63ffffffff1681526020016000151581526020018a8152602001876fffffffffffffffffffffffffffffffff168152602001846fffffffffffffffffffffffffffffffff16815250908060018154018082558091505060019003906000526020600020906003020160009091909190915060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548160ff0219169083151502179055506040820151816001015560608201518160020160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555060808201518160020160106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555050503373ffffffffffffffffffffffffffffffffffffffff16888a7f9b3245740ec3b155098a55be84957a4da13eaf7f14a8bc6f53126c0b9350f2be60405160405180910390a4505050505050505050565b600080547fffffffffffffffffffffffffffffffffffffffffffffff000000000000000000164267ffffffffffffffff161781556040805160a08101825263ffffffff815260208101929092526002919081016111507ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe369081013560f01c90033590565b815260016020820152604001426fffffffffffffffffffffffffffffffff908116909152825460018181018555600094855260209485902084516003909302018054958501511515640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000090961663ffffffff909316929092179490941781556040830151818501556060830151608090930151821670010000000000000000000000000000000002929091169190911760029091015561121690436123ca565b40600155565b6109f682826001610aaf565b6002818154811061123857600080fd5b600091825260209091206003909102018054600182015460029092015463ffffffff8216935064010000000090910460ff1691906fffffffffffffffffffffffffffffffff8082169170010000000000000000000000000000000090041685565b6000805468010000000000000000900460ff1660028111156112bd576112bd612114565b146112f4576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060028781548110611309576113096123e1565b6000918252602082206003919091020160028101549092506fffffffffffffffffffffffffffffffff16908715821760011b90506113687f00000000000000000000000000000000000000000000000000000000000000006001612410565b611404826fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff1614611445576040517f5f53dd9800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008089156114cb57611469836fffffffffffffffffffffffffffffffff16611f25565b67ffffffffffffffff166000036114a2577f000000000000000000000000000000000000000000000000000000000000000091506114c4565b6114bd6114b06001866124f4565b865463ffffffff16611fcb565b6001015491505b50836114e5565b846001015491506114e28460016114b09190612525565b90505b8189896040516114f6929190612559565b604051809103902014611535576040517f696550ff00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081600101547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f8e0cb968c8c8c8c6040518563ffffffff1660e01b815260040161159b94939291906125b2565b6020604051808303816000875af11580156115ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115de91906125e4565b600284810154929091149250600091611689906fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b611725886fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b61172f91906125fd565b6117399190612457565b67ffffffffffffffff161590508115158103611781576040517ffb4e40dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505084547fffffffffffffffffffffffffffffffffffffffffffffffffffffff00ffffffff166401000000001790945550505050505050505050565b6000805468010000000000000000900460ff1660028111156117e1576117e1612114565b14611818576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16637dc0d1d06040518163ffffffff1660e01b8152600401602060405180830381865afa158015611885573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118a9919061261e565b90508160010361194a576001546040517fe52f09370000000000000000000000000000000000000000000000000000000081526004810184905260248101919091526020604482015273ffffffffffffffffffffffffffffffffffffffff82169063e52f0937906064015b600060405180830381600087803b15801561192e57600080fd5b505af1158015611942573d6000803e3d6000fd5b505050505050565b81600203611ac9576040517fcf8e5cf0000000000000000000000000000000000000000000000000000000008152367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c900360200135600482015260009073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063cf8e5cf090602401606060405180830381865afa158015611a0c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a3091906126a3565b80516040517fe52f09370000000000000000000000000000000000000000000000000000000081526004810186905260248101919091526020604482015290915073ffffffffffffffffffffffffffffffffffffffff83169063e52f093790606401600060405180830381600087803b158015611aac57600080fd5b505af1158015611ac0573d6000803e3d6000fd5b50505050505050565b81600303611b5b576040517fe52f093700000000000000000000000000000000000000000000000000000000815260048101839052367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90033560248201526020604482015273ffffffffffffffffffffffffffffffffffffffff82169063e52f093790606401611914565b81600403611bf3576040517fe52f093700000000000000000000000000000000000000000000000000000000815260048101839052367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90036020013560c01b60248201526008604482015273ffffffffffffffffffffffffffffffffffffffff82169063e52f093790606401611914565b816005036109f6576040517fe52f0937000000000000000000000000000000000000000000000000000000008152600481018390524660c01b60248201526008604482015273ffffffffffffffffffffffffffffffffffffffff82169063e52f093790606401611914565b6000367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c9003356060611c95610a9d565b9050909192565b600080611d29847e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff1690508083036001841b600180831b0386831b17039250505092915050565b606081600003611d9457505060408051808201909152600181527f3000000000000000000000000000000000000000000000000000000000000000602082015290565b8160005b8115611dbe5780611da88161272f565b9150611db79050600a83612767565b9150611d98565b60008167ffffffffffffffff811115611dd957611dd9612654565b6040519080825280601f01601f191660200182016040528015611e03576020820181803683370190505b5090505b8415611e8657611e186001836123ca565b9150611e25600a8661277b565b611e30906030612410565b60f81b818381518110611e4557611e456123e1565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350611e7f600a86612767565b9450611e07565b949350505050565b60606000611ec584367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c9003612410565b90508267ffffffffffffffff1667ffffffffffffffff811115611eea57611eea612654565b6040519080825280601f01601f191660200182016040528015611f14576020820181803683370190505b509150828160208401375092915050565b600080611fb2837e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b600167ffffffffffffffff919091161b90920392915050565b600080611fe9846fffffffffffffffffffffffffffffffff16612068565b905060028381548110611ffe57611ffe6123e1565b906000526020600020906003020191505b60028201546fffffffffffffffffffffffffffffffff82811691161461206157815460028054909163ffffffff1690811061204c5761204c6123e1565b9060005260206000209060030201915061200f565b5092915050565b600081196001830116816120fc827e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff169390931c8015179392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b602081016003831061217e577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b6000806040838503121561219757600080fd5b50508035926020909101359150565b60005b838110156121c15781810151838201526020016121a9565b838111156121d0576000848401525b50505050565b600081518084526121ee8160208601602086016121a6565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061223360208301846121d6565b9392505050565b8035801515811461224a57600080fd5b919050565b60008060006060848603121561226457600080fd5b833592506020840135915061227b6040850161223a565b90509250925092565b60006020828403121561229657600080fd5b5035919050565b60008083601f8401126122af57600080fd5b50813567ffffffffffffffff8111156122c757600080fd5b6020830191508360208285010111156122df57600080fd5b9250929050565b600080600080600080608087890312156122ff57600080fd5b8635955061230f6020880161223a565b9450604087013567ffffffffffffffff8082111561232c57600080fd5b6123388a838b0161229d565b9096509450606089013591508082111561235157600080fd5b5061235e89828a0161229d565b979a9699509497509295939492505050565b60ff8416815282602082015260606040820152600061239260608301846121d6565b95945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000828210156123dc576123dc61239b565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082198211156124235761242361239b565b500190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600067ffffffffffffffff8084168061247257612472612428565b92169190910692915050565b600084516124908184602089016121a6565b80830190507f2e0000000000000000000000000000000000000000000000000000000000000080825285516124cc816001850160208a016121a6565b600192019182015283516124e78160028401602088016121a6565b0160020195945050505050565b60006fffffffffffffffffffffffffffffffff8381169083168181101561251d5761251d61239b565b039392505050565b60006fffffffffffffffffffffffffffffffff8083168185168083038211156125505761255061239b565b01949350505050565b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6040815260006125c6604083018688612569565b82810360208401526125d9818587612569565b979650505050505050565b6000602082840312156125f657600080fd5b5051919050565b600067ffffffffffffffff8381169083168181101561251d5761251d61239b565b60006020828403121561263057600080fd5b815173ffffffffffffffffffffffffffffffffffffffff8116811461223357600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b80516fffffffffffffffffffffffffffffffff8116811461224a57600080fd5b6000606082840312156126b557600080fd5b6040516060810181811067ffffffffffffffff821117156126ff577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040528251815261271260208401612683565b602082015261272360408401612683565b60408201529392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036127605761276061239b565b5060010190565b60008261277657612776612428565b500490565b60008261278a5761278a612428565b50069056fea164736f6c634300080f000a"

func init() {
	if err := json.Unmarshal([]byte(FaultDisputeGameStorageLayoutJSON), FaultDisputeGameStorageLayout); err != nil {
		panic(err)
	}

	layouts["FaultDisputeGame"] = FaultDisputeGameStorageLayout
	deployedBytecodes["FaultDisputeGame"] = FaultDisputeGameDeployedBin
}
