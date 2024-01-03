// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/ethereum-optimism/optimism/op-bindings/solc"
)

const FaultDisputeGameStorageLayoutJSON = "{\"storage\":[{\"astId\":1000,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"createdAt\",\"offset\":0,\"slot\":\"0\",\"type\":\"t_userDefinedValueType(Timestamp)1016\"},{\"astId\":1001,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"resolvedAt\",\"offset\":8,\"slot\":\"0\",\"type\":\"t_userDefinedValueType(Timestamp)1016\"},{\"astId\":1002,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"status\",\"offset\":16,\"slot\":\"0\",\"type\":\"t_enum(GameStatus)1009\"},{\"astId\":1003,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"l1Head\",\"offset\":0,\"slot\":\"1\",\"type\":\"t_userDefinedValueType(Hash)1014\"},{\"astId\":1004,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"claimData\",\"offset\":0,\"slot\":\"2\",\"type\":\"t_array(t_struct(ClaimData)1010_storage)dyn_storage\"},{\"astId\":1005,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"claims\",\"offset\":0,\"slot\":\"3\",\"type\":\"t_mapping(t_userDefinedValueType(ClaimHash)1012,t_bool)\"},{\"astId\":1006,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"subgames\",\"offset\":0,\"slot\":\"4\",\"type\":\"t_mapping(t_uint256,t_array(t_uint256)dyn_storage)\"},{\"astId\":1007,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"subgameAtRootResolved\",\"offset\":0,\"slot\":\"5\",\"type\":\"t_bool\"},{\"astId\":1008,\"contract\":\"src/dispute/FaultDisputeGame.sol:FaultDisputeGame\",\"label\":\"initialized\",\"offset\":1,\"slot\":\"5\",\"type\":\"t_bool\"}],\"types\":{\"t_address\":{\"encoding\":\"inplace\",\"label\":\"address\",\"numberOfBytes\":\"20\"},\"t_array(t_struct(ClaimData)1010_storage)dyn_storage\":{\"encoding\":\"dynamic_array\",\"label\":\"struct IFaultDisputeGame.ClaimData[]\",\"numberOfBytes\":\"32\",\"base\":\"t_struct(ClaimData)1010_storage\"},\"t_array(t_uint256)dyn_storage\":{\"encoding\":\"dynamic_array\",\"label\":\"uint256[]\",\"numberOfBytes\":\"32\",\"base\":\"t_uint256\"},\"t_bool\":{\"encoding\":\"inplace\",\"label\":\"bool\",\"numberOfBytes\":\"1\"},\"t_enum(GameStatus)1009\":{\"encoding\":\"inplace\",\"label\":\"enum GameStatus\",\"numberOfBytes\":\"1\"},\"t_mapping(t_uint256,t_array(t_uint256)dyn_storage)\":{\"encoding\":\"mapping\",\"label\":\"mapping(uint256 =\u003e uint256[])\",\"numberOfBytes\":\"32\",\"key\":\"t_uint256\",\"value\":\"t_array(t_uint256)dyn_storage\"},\"t_mapping(t_userDefinedValueType(ClaimHash)1012,t_bool)\":{\"encoding\":\"mapping\",\"label\":\"mapping(ClaimHash =\u003e bool)\",\"numberOfBytes\":\"32\",\"key\":\"t_userDefinedValueType(ClaimHash)1012\",\"value\":\"t_bool\"},\"t_struct(ClaimData)1010_storage\":{\"encoding\":\"inplace\",\"label\":\"struct IFaultDisputeGame.ClaimData\",\"numberOfBytes\":\"160\"},\"t_uint128\":{\"encoding\":\"inplace\",\"label\":\"uint128\",\"numberOfBytes\":\"16\"},\"t_uint256\":{\"encoding\":\"inplace\",\"label\":\"uint256\",\"numberOfBytes\":\"32\"},\"t_uint32\":{\"encoding\":\"inplace\",\"label\":\"uint32\",\"numberOfBytes\":\"4\"},\"t_userDefinedValueType(Claim)1011\":{\"encoding\":\"inplace\",\"label\":\"Claim\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(ClaimHash)1012\":{\"encoding\":\"inplace\",\"label\":\"ClaimHash\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(Clock)1013\":{\"encoding\":\"inplace\",\"label\":\"Clock\",\"numberOfBytes\":\"16\"},\"t_userDefinedValueType(Hash)1014\":{\"encoding\":\"inplace\",\"label\":\"Hash\",\"numberOfBytes\":\"32\"},\"t_userDefinedValueType(Position)1015\":{\"encoding\":\"inplace\",\"label\":\"Position\",\"numberOfBytes\":\"16\"},\"t_userDefinedValueType(Timestamp)1016\":{\"encoding\":\"inplace\",\"label\":\"Timestamp\",\"numberOfBytes\":\"8\"}}}"

var FaultDisputeGameStorageLayout = new(solc.StorageLayout)

var FaultDisputeGameDeployedBin = "0x6080604052600436106101b75760003560e01c80638d450a95116100ec578063d8cc1a3c1161008a578063f8f43ff611610064578063f8f43ff614610624578063fa24f74314610644578063fa315aa914610668578063fdffbb281461069b57600080fd5b8063d8cc1a3c146105ab578063e1f0c376146105be578063ec5e6308146105f157600080fd5b8063c395e1ca116100c6578063c395e1ca146104cc578063c55cd0c7146104ed578063c6f0308c14610500578063cf09e0d01461058a57600080fd5b80638d450a951461041e578063bbdc02db14610451578063bcef3b551461048f57600080fd5b8063609d33341161015957806368800abf1161013357806368800abf1461038e5780638129fc1c146103c15780638980e0cc146103c95780638b85902b146103de57600080fd5b8063609d333414610350578063632247ea146103655780636361506d1461037857600080fd5b80632810e1d6116101955780632810e1d61461027f57806335fef567146102945780633a768463146102a957806354fd4d50146102fa57600080fd5b80630356fe3a146101bc57806319effeb4146101fe578063200d2ed214610244575b600080fd5b3480156101c857600080fd5b507f00000000000000000000000000000000000000000000000000000000000000005b6040519081526020015b60405180910390f35b34801561020a57600080fd5b5060005461022b9068010000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016101f5565b34801561025057600080fd5b5060005461027290700100000000000000000000000000000000900460ff1681565b6040516101f591906130a6565b34801561028b57600080fd5b506102726106ae565b6102a76102a23660046130e7565b6108ab565b005b3480156102b557600080fd5b5060405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101f5565b34801561030657600080fd5b506103436040518060400160405280600681526020017f302e302e3232000000000000000000000000000000000000000000000000000081525081565b6040516101f59190613174565b34801561035c57600080fd5b506103436108bb565b6102a761037336600461319c565b6108cd565b34801561038457600080fd5b506101eb60015481565b34801561039a57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006101eb565b6102a76110cc565b3480156103d557600080fd5b506002546101eb565b3480156103ea57600080fd5b50367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c9003602001356101eb565b34801561042a57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006101eb565b34801561045d57600080fd5b5060405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101f5565b34801561049b57600080fd5b50367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c9003356101eb565b3480156104d857600080fd5b506101eb6104e73660046131d1565b50600090565b6102a76104fb3660046130e7565b6113f2565b34801561050c57600080fd5b5061052061051b36600461320a565b6113fe565b6040805163ffffffff909816885273ffffffffffffffffffffffffffffffffffffffff968716602089015295909416948601949094526fffffffffffffffffffffffffffffffff9182166060860152608085015291821660a08401521660c082015260e0016101f5565b34801561059657600080fd5b5060005461022b9067ffffffffffffffff1681565b6102a76105b936600461326c565b611495565b3480156105ca57600080fd5b507f000000000000000000000000000000000000000000000000000000000000000061022b565b3480156105fd57600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006101eb565b34801561063057600080fd5b506102a761063f3660046132f6565b611a73565b34801561065057600080fd5b50610659611edc565b6040516101f593929190613322565b34801561067457600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006101eb565b6102a76106a936600461320a565b611f39565b600080600054700100000000000000000000000000000000900460ff1660028111156106dc576106dc613077565b14610713576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60055460ff1661074f576040517f9a07664600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16600260008154811061077b5761077b61334d565b6000918252602090912060059091020154640100000000900473ffffffffffffffffffffffffffffffffffffffff16146107b65760016107b9565b60025b6000805467ffffffffffffffff421668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff82168117835592935083927fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff167fffffffffffffffffffffffffffffff000000000000000000ffffffffffffffff9091161770010000000000000000000000000000000083600281111561086a5761086a613077565b02179055600281111561087f5761087f613077565b6040517f5e186f09b9c93491f14e277eea7faa5de6a2d4bda75a79af7a3684fbfb42da6090600090a290565b6108b7828260006108cd565b5050565b60606108c8602080612369565b905090565b60008054700100000000000000000000000000000000900460ff1660028111156108f9576108f9613077565b14610930576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600284815481106109455761094561334d565b600091825260208083206040805160e0810182526005909402909101805463ffffffff808216865273ffffffffffffffffffffffffffffffffffffffff6401000000009092048216948601949094526001820154169184019190915260028101546fffffffffffffffffffffffffffffffff90811660608501526003820154608085015260049091015480821660a0850181905270010000000000000000000000000000000090910490911660c0840152919350909190610a0a908390869061240016565b90506000610aaa826fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff169050861580610aec5750610ae97f000000000000000000000000000000000000000000000000000000000000000060026133ab565b81145b8015610af6575084155b15610b2d576040517fa42637bc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f0000000000000000000000000000000000000000000000000000000000000000811115610b87576040517f56f57b2b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610bb27f000000000000000000000000000000000000000000000000000000000000000060016133ab565b8103610bc457610bc486888588612408565b835160009063ffffffff90811614610c24576002856000015163ffffffff1681548110610bf357610bf361334d565b906000526020600020906005020160040160109054906101000a90046fffffffffffffffffffffffffffffffff1690505b60c0850151600090610c489067ffffffffffffffff165b67ffffffffffffffff1690565b67ffffffffffffffff1642610c72610c3b856fffffffffffffffffffffffffffffffff1660401c90565b67ffffffffffffffff16610c8691906133ab565b610c9091906133c3565b90507f000000000000000000000000000000000000000000000000000000000000000060011c677fffffffffffffff1667ffffffffffffffff82161115610d03576040517f3381d11400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000604082901b421760008a8152608087901b6fffffffffffffffffffffffffffffffff8d1617602052604081209192509060008181526003602052604090205490915060ff1615610d81576040517f80497e3b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016003600083815260200190815260200160002060006101000a81548160ff02191690831515021790555060026040518060e001604052808d63ffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff168152602001346fffffffffffffffffffffffffffffffff1681526020018c8152602001886fffffffffffffffffffffffffffffffff168152602001846fffffffffffffffffffffffffffffffff16815250908060018154018082558091505060019003906000526020600020906005020160009091909190915060008201518160000160006101000a81548163ffffffff021916908363ffffffff16021790555060208201518160000160046101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060608201518160020160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff1602179055506080820151816003015560a08201518160040160006101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555060c08201518160040160106101000a8154816fffffffffffffffffffffffffffffffff02191690836fffffffffffffffffffffffffffffffff16021790555050503360028c815481106110065761100661334d565b600091825260208083206005909202909101805473ffffffffffffffffffffffffffffffffffffffff94909416640100000000027fffffffffffffffff0000000000000000000000000000000000000000ffffffff909416939093179092558c8152600490915260409020600254611080906001906133c3565b8154600181018355600092835260208320015560405133918c918e917f9b3245740ec3b155098a55be84957a4da13eaf7f14a8bc6f53126c0b9350f2be91a45050505050505050505050565b600554610100900460ff161561110e576040517f0dc149f000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f0000000000000000000000000000000000000000000000000000000000000000367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c900360200135116111c5576040517ff40239db000000000000000000000000000000000000000000000000000000008152367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90033560048201526024015b60405180910390fd5b60463611156111dc5763c407e0256000526004601cfd5b6040805160e08101825263ffffffff8152600060208201523291810191909152346fffffffffffffffffffffffffffffffff16606082015260029060808101367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c900335815260016020820152604001426fffffffffffffffffffffffffffffffff90811690915282546001808201855560009485526020808620855160059094020180549186015163ffffffff9094167fffffffffffffffff0000000000000000000000000000000000000000000000009092169190911764010000000073ffffffffffffffffffffffffffffffffffffffff94851602178155604085015181830180547fffffffffffffffffffffffff000000000000000000000000000000000000000016919094161790925560608401516002830180547fffffffffffffffffffffffffffffffff00000000000000000000000000000000169185169190911790556080840151600383015560a084015160c09094015193831670010000000000000000000000000000000094909316939093029190911760049091015581547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000164267ffffffffffffffff16179091556113c090436133c3565b40600155600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16610100179055565b6108b7828260016108cd565b6002818154811061140e57600080fd5b60009182526020909120600590910201805460018201546002830154600384015460049094015463ffffffff8416955064010000000090930473ffffffffffffffffffffffffffffffffffffffff908116949216926fffffffffffffffffffffffffffffffff91821692918082169170010000000000000000000000000000000090041687565b60008054700100000000000000000000000000000000900460ff1660028111156114c1576114c1613077565b146114f8576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006002878154811061150d5761150d61334d565b6000918252602082206005919091020160048101549092506fffffffffffffffffffffffffffffffff16908715821760011b905061156c7f000000000000000000000000000000000000000000000000000000000000000060016133ab565b611608826fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff1614611649576040517f5f53dd9800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008089156117385761169c7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000006133c3565b6001901b6116bb846fffffffffffffffffffffffffffffffff166125c0565b67ffffffffffffffff166116cf9190613409565b1561170c576117036116f460016fffffffffffffffffffffffffffffffff871661341d565b865463ffffffff166000612666565b6003015461172e565b7f00000000000000000000000000000000000000000000000000000000000000005b9150849050611762565b6003850154915061175f6116f46fffffffffffffffffffffffffffffffff8616600161344e565b90505b600882901b60088a8a604051611779929190613482565b6040518091039020901b146117ba576040517f696550ff00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006117c58c61274a565b905060006117d4836003015490565b6040517fe14ced320000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063e14ced329061184e908f908f908f908f908a906004016134db565b6020604051808303816000875af115801561186d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906118919190613515565b60048501549114915060009060029061193c906fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b6119d8896fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b6119e2919061352e565b6119ec919061354f565b67ffffffffffffffff161590508115158103611a34576040517ffb4e40dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505085547fffffffffffffffff0000000000000000000000000000000000000000ffffffff163364010000000002179095555050505050505050505050565b60008054700100000000000000000000000000000000900460ff166002811115611a9f57611a9f613077565b14611ad6576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600080611ae586612779565b93509350935093506000611afb85858585612ba6565b905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16637dc0d1d06040518163ffffffff1660e01b8152600401602060405180830381865afa158015611b6a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b8e9190613576565b905060018903611c565773ffffffffffffffffffffffffffffffffffffffff81166352f0f3ad8a846001545b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e086901b16815260048101939093526024830191909152604482015260206064820152608481018a905260a4015b6020604051808303816000875af1158015611c2c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c509190613515565b50611ed1565b60028903611c825773ffffffffffffffffffffffffffffffffffffffff81166352f0f3ad8a8489611bba565b60038903611cae5773ffffffffffffffffffffffffffffffffffffffff81166352f0f3ad8a8487611bba565b60048903611e265760006fffffffffffffffffffffffffffffffff861615611d4657611d0c6fffffffffffffffffffffffffffffffff87167f0000000000000000000000000000000000000000000000000000000000000000612c65565b611d36907f00000000000000000000000000000000000000000000000000000000000000006133ab565b611d419060016133ab565b611d68565b7f00000000000000000000000000000000000000000000000000000000000000005b905073ffffffffffffffffffffffffffffffffffffffff82166352f0f3ad8b8560405160e084901b7fffffffff000000000000000000000000000000000000000000000000000000001681526004810192909252602482015260c084901b604482015260086064820152608481018b905260a4016020604051808303816000875af1158015611dfb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611e1f9190613515565b5050611ed1565b60058903611e9f576040517f52f0f3ad000000000000000000000000000000000000000000000000000000008152600481018a9052602481018390524660c01b6044820152600860648201526084810188905273ffffffffffffffffffffffffffffffffffffffff8216906352f0f3ad9060a401611c0d565b6040517fff137e6500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050505050505050565b7f0000000000000000000000000000000000000000000000000000000000000000367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c9003356060611f326108bb565b9050909192565b60008054700100000000000000000000000000000000900460ff166002811115611f6557611f65613077565b14611f9c576040517f67fe195000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060028281548110611fb157611fb161334d565b600091825260208220600591909102016004810154909250611ff390700100000000000000000000000000000000900460401c67ffffffffffffffff16610c3b565b600483015490915060009061202590700100000000000000000000000000000000900467ffffffffffffffff16610c3b565b61202f904261352e565b9050677fffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000060011c1661206982846135ac565b67ffffffffffffffff16116120aa576040517ff2440b5300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008481526004602052604090208054851580156120ca575060055460ff165b806120dd5750801580156120dd57508515155b15612114576040517ff1a9458100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000805b828110156121f35760008482815481106121345761213461334d565b600091825260208083209091015480835260049091526040909120549091501561218a576040517f9a07664600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006002828154811061219f5761219f61334d565b600091825260209091206005909102018054909150640100000000900473ffffffffffffffffffffffffffffffffffffffff166121e05733935050506121f3565b5050806121ec906135cf565b9050612118565b5073ffffffffffffffffffffffffffffffffffffffff81161561226e57600286015460405173ffffffffffffffffffffffffffffffffffffffff8316916fffffffffffffffffffffffffffffffff1680156108fc02916000818181858888f19350505050158015612268573d6000803e3d6000fd5b506122d0565b6001860154600287015460405173ffffffffffffffffffffffffffffffffffffffff909216916fffffffffffffffffffffffffffffffff90911680156108fc02916000818181858888f193505050501580156122ce573d6000803e3d6000fd5b505b85547fffffffffffffffff0000000000000000000000000000000000000000ffffffff1664010000000073ffffffffffffffffffffffffffffffffffffffff831602178655600087815260046020526040812061232c9161303d565b8660000361236057600580547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790555b50505050505050565b606060006123a084367ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe81013560f01c90036133ab565b90508267ffffffffffffffff1667ffffffffffffffff8111156123c5576123c5613607565b6040519080825280601f01601f1916602001820160405280156123ef576020820181803683370190505b509150828160208401375092915050565b151760011b90565b60006124276fffffffffffffffffffffffffffffffff8416600161344e565b9050600061243782866001612666565b9050600086901a838061252a575061247060027f0000000000000000000000000000000000000000000000000000000000000000613409565b6004830154600290612514906fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b61251e919061354f565b67ffffffffffffffff16145b156125825760ff811660011480612544575060ff81166002145b61257d576040517ff40239db000000000000000000000000000000000000000000000000000000008152600481018890526024016111bc565b612360565b60ff811615612360576040517ff40239db000000000000000000000000000000000000000000000000000000008152600481018890526024016111bc565b60008061264d837e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b600167ffffffffffffffff919091161b90920392915050565b600080826126af576126aa6fffffffffffffffffffffffffffffffff86167f0000000000000000000000000000000000000000000000000000000000000000612d1a565b6126ca565b6126ca856fffffffffffffffffffffffffffffffff16612ee1565b9050600284815481106126df576126df61334d565b906000526020600020906005020191505b60048201546fffffffffffffffffffffffffffffffff82811691161461274257815460028054909163ffffffff1690811061272d5761272d61334d565b906000526020600020906005020191506126f0565b509392505050565b600080600080600061275b86612779565b935093509350935061276f84848484612ba6565b9695505050505050565b60008060008060008590506000600282815481106127995761279961334d565b600091825260209091206004600590920201908101549091507f000000000000000000000000000000000000000000000000000000000000000090612870906fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff16116128b1576040517fb34b5c2200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000815b60048301547f000000000000000000000000000000000000000000000000000000000000000090612978906fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff1692508211156129f457825463ffffffff166129be7f000000000000000000000000000000000000000000000000000000000000000060016133ab565b83036129c8578391505b600281815481106129db576129db61334d565b90600052602060002090600502019350809450506128b5565b600481810154908401546fffffffffffffffffffffffffffffffff91821691166000816fffffffffffffffffffffffffffffffff16612a5d612a48856fffffffffffffffffffffffffffffffff1660011c90565b6fffffffffffffffffffffffffffffffff1690565b6fffffffffffffffffffffffffffffffff161490508015612b42576000612a95836fffffffffffffffffffffffffffffffff166125c0565b67ffffffffffffffff161115612af8576000612acf612ac760016fffffffffffffffffffffffffffffffff861661341d565b896001612666565b6003810154600490910154909c506fffffffffffffffffffffffffffffffff169a50612b1c9050565b7f00000000000000000000000000000000000000000000000000000000000000009a505b600386015460048701549099506fffffffffffffffffffffffffffffffff169750612b98565b6000612b64612ac76fffffffffffffffffffffffffffffffff8516600161344e565b6003808901546004808b015492840154930154909e506fffffffffffffffffffffffffffffffff9182169d50919b50169850505b505050505050509193509193565b60006fffffffffffffffffffffffffffffffff84168103612c0c578282604051602001612bef9291909182526fffffffffffffffffffffffffffffffff16602082015260400190565b604051602081830303815290604052805190602001209050612c5d565b60408051602081018790526fffffffffffffffffffffffffffffffff8087169282019290925260608101859052908316608082015260a0016040516020818303038152906040528051906020012090505b949350505050565b600080612cf2847e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff1690508083036001841b600180831b0386831b17039250505092915050565b600081612db9846fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff1611612dfa576040517fb34b5c2200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612e0383612ee1565b905081612ea2826fffffffffffffffffffffffffffffffff167e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff1611612edb57612ed8612ebf8360016133ab565b6fffffffffffffffffffffffffffffffff831690612f8d565b90505b92915050565b60008119600183011681612f75827e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff169390931c8015179392505050565b60008061301a847e09010a0d15021d0b0e10121619031e080c141c0f111807131b17061a05041f7f07c4acdd0000000000000000000000000000000000000000000000000000000067ffffffffffffffff831160061b83811c63ffffffff1060051b1792831c600181901c17600281901c17600481901c17600881901c17601081901c170260fb1c1a1790565b67ffffffffffffffff169050808303600180821b0385821b179250505092915050565b508054600082559060005260206000209081019061305b919061305e565b50565b5b80821115613073576000815560010161305f565b5090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60208101600383106130e1577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b600080604083850312156130fa57600080fd5b50508035926020909101359150565b6000815180845260005b8181101561312f57602081850181015186830182015201613113565b81811115613141576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612ed86020830184613109565b8035801515811461319757600080fd5b919050565b6000806000606084860312156131b157600080fd5b83359250602084013591506131c860408501613187565b90509250925092565b6000602082840312156131e357600080fd5b81356fffffffffffffffffffffffffffffffff8116811461320357600080fd5b9392505050565b60006020828403121561321c57600080fd5b5035919050565b60008083601f84011261323557600080fd5b50813567ffffffffffffffff81111561324d57600080fd5b60208301915083602082850101111561326557600080fd5b9250929050565b6000806000806000806080878903121561328557600080fd5b8635955061329560208801613187565b9450604087013567ffffffffffffffff808211156132b257600080fd5b6132be8a838b01613223565b909650945060608901359150808211156132d757600080fd5b506132e489828a01613223565b979a9699509497509295939492505050565b60008060006060848603121561330b57600080fd5b505081359360208301359350604090920135919050565b60ff841681528260208201526060604082015260006133446060830184613109565b95945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600082198211156133be576133be61337c565b500190565b6000828210156133d5576133d561337c565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082613418576134186133da565b500690565b60006fffffffffffffffffffffffffffffffff838116908316818110156134465761344661337c565b039392505050565b60006fffffffffffffffffffffffffffffffff8083168185168083038211156134795761347961337c565b01949350505050565b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6060815260006134ef606083018789613492565b8281036020840152613502818688613492565b9150508260408301529695505050505050565b60006020828403121561352757600080fd5b5051919050565b600067ffffffffffffffff838116908316818110156134465761344661337c565b600067ffffffffffffffff8084168061356a5761356a6133da565b92169190910692915050565b60006020828403121561358857600080fd5b815173ffffffffffffffffffffffffffffffffffffffff8116811461320357600080fd5b600067ffffffffffffffff8083168185168083038211156134795761347961337c565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036136005761360061337c565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c634300080f000a"


func init() {
	if err := json.Unmarshal([]byte(FaultDisputeGameStorageLayoutJSON), FaultDisputeGameStorageLayout); err != nil {
		panic(err)
	}

	layouts["FaultDisputeGame"] = FaultDisputeGameStorageLayout
	deployedBytecodes["FaultDisputeGame"] = FaultDisputeGameDeployedBin
	immutableReferences["FaultDisputeGame"] = true
}
