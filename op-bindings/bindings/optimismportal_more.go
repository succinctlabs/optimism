// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/ethereum-optimism/optimism/op-bindings/solc"
)

const OptimismPortalStorageLayoutJSON = "{\"storage\":[{\"astId\":27194,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"_initialized\",\"offset\":0,\"slot\":\"0\",\"type\":\"t_uint8\"},{\"astId\":27197,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"_initializing\",\"offset\":1,\"slot\":\"0\",\"type\":\"t_bool\"},{\"astId\":1404,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"params\",\"offset\":0,\"slot\":\"1\",\"type\":\"t_struct(ResourceParams)1374_storage\"},{\"astId\":1409,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"__gap\",\"offset\":0,\"slot\":\"2\",\"type\":\"t_array(t_uint256)49_storage\"},{\"astId\":982,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"l2Sender\",\"offset\":0,\"slot\":\"51\",\"type\":\"t_address\"},{\"astId\":995,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"finalizedWithdrawals\",\"offset\":0,\"slot\":\"52\",\"type\":\"t_mapping(t_bytes32,t_bool)\"}],\"types\":{\"t_address\":{\"encoding\":\"inplace\",\"label\":\"address\",\"numberOfBytes\":\"20\"},\"t_array(t_uint256)49_storage\":{\"encoding\":\"inplace\",\"label\":\"uint256[49]\",\"numberOfBytes\":\"1568\"},\"t_bool\":{\"encoding\":\"inplace\",\"label\":\"bool\",\"numberOfBytes\":\"1\"},\"t_bytes32\":{\"encoding\":\"inplace\",\"label\":\"bytes32\",\"numberOfBytes\":\"32\"},\"t_mapping(t_bytes32,t_bool)\":{\"encoding\":\"mapping\",\"label\":\"mapping(bytes32 =\u003e bool)\",\"numberOfBytes\":\"32\",\"key\":\"t_bytes32\",\"value\":\"t_bool\"},\"t_struct(ResourceParams)1374_storage\":{\"encoding\":\"inplace\",\"label\":\"struct ResourceMetering.ResourceParams\",\"numberOfBytes\":\"32\"},\"t_uint128\":{\"encoding\":\"inplace\",\"label\":\"uint128\",\"numberOfBytes\":\"16\"},\"t_uint256\":{\"encoding\":\"inplace\",\"label\":\"uint256\",\"numberOfBytes\":\"32\"},\"t_uint64\":{\"encoding\":\"inplace\",\"label\":\"uint64\",\"numberOfBytes\":\"8\"},\"t_uint8\":{\"encoding\":\"inplace\",\"label\":\"uint8\",\"numberOfBytes\":\"1\"}}}"

var OptimismPortalStorageLayout = new(solc.StorageLayout)

var OptimismPortalDeployedBin = "0x6080604052600436106100f65760003560e01c8063a14238e71161008a578063cff0ab9611610059578063cff0ab96146102f7578063e9e05c4214610398578063f4daa291146103ab578063fdc9fe1d146103df57600080fd5b8063a14238e71461026d578063c4fc4798146102ad578063ca3e99ba146102cd578063cd7c9789146102e257600080fd5b80636bb0291e116100c65780636bb0291e146102005780638129fc1c14610215578063867ead131461022a5780639bf62d821461024057600080fd5b80621c2ff61461012257806313620abd1461018057806354fd4d50146101b957806364b79208146101db57600080fd5b3661011d5761011b3334620186a06000604051806020016040528060008152506103ff565b005b600080fd5b34801561012e57600080fd5b506101567f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b34801561018c57600080fd5b50610198633b9aca0081565b6040516fffffffffffffffffffffffffffffffff9091168152602001610177565b3480156101c557600080fd5b506101ce6108cc565b6040516101779190613307565b3480156101e757600080fd5b506101f2627a120081565b604051908152602001610177565b34801561020c57600080fd5b506101f2600481565b34801561022157600080fd5b5061011b61096f565b34801561023657600080fd5b506101f261271081565b34801561024c57600080fd5b506033546101569073ffffffffffffffffffffffffffffffffffffffff1681565b34801561027957600080fd5b5061029d61028836600461331a565b60346020526000908152604090205460ff1681565b6040519015158152602001610177565b3480156102b957600080fd5b5061029d6102c836600461331a565b610b2d565b3480156102d957600080fd5b506101f2610bf2565b3480156102ee57600080fd5b506101f2600881565b34801561030357600080fd5b5060015461035f906fffffffffffffffffffffffffffffffff81169067ffffffffffffffff7001000000000000000000000000000000008204811691780100000000000000000000000000000000000000000000000090041683565b604080516fffffffffffffffffffffffffffffffff909416845267ffffffffffffffff9283166020850152911690820152606001610177565b61011b6103a636600461345f565b6103ff565b3480156103b757600080fd5b506101f27f000000000000000000000000000000000000000000000000000000000000000081565b3480156103eb57600080fd5b5061011b6103fa36600461354d565b610c03565b8260005a905083156104b65773ffffffffffffffffffffffffffffffffffffffff8716156104b657604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f4f7074696d69736d506f7274616c3a206d7573742073656e6420746f2061646460448201527f72657373283029207768656e206372656174696e67206120636f6e747261637460648201526084015b60405180910390fd5b333281146104d7575033731111000000000000000000000000000000001111015b600034888888886040516020016104f2959493929190613641565b604051602081830303815290604052905060008973ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fb3813568d9991fc951961fcb4c784893574240a28925604d09fc577c55bb7c32846040516105629190613307565b60405180910390a450506001546000906105a2907801000000000000000000000000000000000000000000000000900467ffffffffffffffff16436136d5565b9050801561072b5760006105ba6004627a120061371b565b6001546105e59190700100000000000000000000000000000000900467ffffffffffffffff16613783565b9050600060086105f96004627a120061371b565b6001546106199085906fffffffffffffffffffffffffffffffff166137f7565b610623919061371b565b61062d919061371b565b600154909150600090610679906106639061065b9085906fffffffffffffffffffffffffffffffff166138b3565b6127106112a4565b6fffffffffffffffffffffffffffffffff6112bf565b905060018411156106ec576106e9610663670de0b6b3a76400006106d56106a160088361371b565b6106b390670de0b6b3a7640000613783565b6106be60018a6136d5565b6106d090670de0b6b3a7640000613927565b6112ce565b6106df90856137f7565b61065b919061371b565b90505b6fffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000067ffffffffffffffff4316021760015550505b6001805484919060109061075e908490700100000000000000000000000000000000900467ffffffffffffffff16613964565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550627a1200600160000160109054906101000a900467ffffffffffffffff1667ffffffffffffffff16131561083a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603e60248201527f5265736f757263654d65746572696e673a2063616e6e6f7420627579206d6f7260448201527f6520676173207468616e20617661696c61626c6520676173206c696d6974000060648201526084016104ad565b600154600090610866906fffffffffffffffffffffffffffffffff1667ffffffffffffffff8616613990565b6fffffffffffffffffffffffffffffffff169050600061088a48633b9aca006112ff565b61089490836139c8565b905060005a6108a390866136d5565b9050808211156108bf576108bf6108ba82846136d5565b61130f565b5050505050505050505050565b60606108f77f000000000000000000000000000000000000000000000000000000000000000061133d565b6109207f000000000000000000000000000000000000000000000000000000000000000061133d565b6109497f000000000000000000000000000000000000000000000000000000000000000061133d565b60405160200161095b939291906139dc565b604051602081830303815290604052905090565b600054610100900460ff161580801561098f5750600054600160ff909116105b806109a95750303b1580156109a9575060005460ff166001145b610a35576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016104ad565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790558015610a9357600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b603380547fffffffffffffffffffffffff00000000000000000000000000000000000000001661dead179055610ac761147a565b8015610b2a57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b6040517fa25ae55700000000000000000000000000000000000000000000000000000000815260048101829052600090819073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063a25ae557906024016040805180830381865afa158015610bbc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610be09190613a52565b9050610beb8161155d565b9392505050565b610c006004627a120061371b565b81565b60335473ffffffffffffffffffffffffffffffffffffffff1661dead14610cac576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603f60248201527f4f7074696d69736d506f7274616c3a2063616e206f6e6c79207472696767657260448201527f206f6e65207769746864726177616c20706572207472616e73616374696f6e0060648201526084016104ad565b3073ffffffffffffffffffffffffffffffffffffffff16856040015173ffffffffffffffffffffffffffffffffffffffff1603610d6b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603f60248201527f4f7074696d69736d506f7274616c3a20796f752063616e6e6f742073656e642060448201527f6d6573736167657320746f2074686520706f7274616c20636f6e74726163740060648201526084016104ad565b6040517fa25ae557000000000000000000000000000000000000000000000000000000008152600481018590526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a25ae557906024016040805180830381865afa158015610df8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e1c9190613a52565b9050610e278161155d565b610eb3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602d60248201527f4f7074696d69736d506f7274616c3a2070726f706f73616c206973206e6f742060448201527f7965742066696e616c697a65640000000000000000000000000000000000000060648201526084016104ad565b610eca610ec536869003860186613aa1565b611597565b815114610f59576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f4f7074696d69736d506f7274616c3a20696e76616c6964206f7574707574207260448201527f6f6f742070726f6f66000000000000000000000000000000000000000000000060648201526084016104ad565b6000610f64876115f3565b9050610fab81866040013586868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061162392505050565b611037576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603260248201527f4f7074696d69736d506f7274616c3a20696e76616c696420776974686472617760448201527f616c20696e636c7573696f6e2070726f6f66000000000000000000000000000060648201526084016104ad565b60008181526034602052604090205460ff16156110d6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603560248201527f4f7074696d69736d506f7274616c3a207769746864726177616c20686173206160448201527f6c7265616479206265656e2066696e616c697a6564000000000000000000000060648201526084016104ad565b600081815260346020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055608087015161111f90614e2090613b07565b5a10156111ae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603760248201527f4f7074696d69736d506f7274616c3a20696e73756666696369656e742067617360448201527f20746f2066696e616c697a65207769746864726177616c00000000000000000060648201526084016104ad565b8660200151603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000611211886040015189608001518a606001518b60a001516116ea565b603380547fffffffffffffffffffffffff00000000000000000000000000000000000000001661dead17905560405190915082907fdb5c7652857aa163daadd670e116628fb42e869d8ac4251ef8971d9e5727df1b9061127690841515815260200190565b60405180910390a25050505050505050565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b6000818312156112b457816112b6565b825b90505b92915050565b60008183126112b457816112b6565b60006112b6670de0b6b3a7640000836112e686611704565b6112f091906137f7565b6112fa919061371b565b611948565b6000818310156112b457816112b6565b6000805a90505b825a61132290836136d5565b10156113385761133182613b1f565b9150611316565b505050565b60608160000361138057505060408051808201909152600181527f3000000000000000000000000000000000000000000000000000000000000000602082015290565b8160005b81156113aa578061139481613b1f565b91506113a39050600a836139c8565b9150611384565b60008167ffffffffffffffff8111156113c5576113c561335c565b6040519080825280601f01601f1916602001820160405280156113ef576020820181803683370190505b5090505b8415611472576114046001836136d5565b9150611411600a86613b57565b61141c906030613b07565b60f81b81838151811061143157611431613b6b565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535061146b600a866139c8565b94506113f3565b949350505050565b600054610100900460ff16611511576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e6700000000000000000000000000000000000000000060648201526084016104ad565b60408051606081018252633b9aca00808252600060208301524367ffffffffffffffff169190920181905278010000000000000000000000000000000000000000000000000217600155565b60007f0000000000000000000000000000000000000000000000000000000000000000826020015161158f9190613b07565b421192915050565b600081600001518260200151836040015184606001516040516020016115d6949392919093845260208401929092526040830152606082015260800190565b604051602081830303815290604052805190602001209050919050565b80516020808301516040808501516060860151608087015160a088015193516000976115d6979096959101613b9a565b604080516020810185905260009181018290528190606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828252805160209182012090830181905292506116e19101604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828201909152600182527f01000000000000000000000000000000000000000000000000000000000000006020830152908587611b87565b95945050505050565b600080600080845160208601878a8af19695505050505050565b600080821361176f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f554e444546494e4544000000000000000000000000000000000000000000000060448201526064016104ad565b6000606061177c84611bab565b03609f8181039490941b90931c6c465772b2bbbb5f824b15207a3081018102606090811d6d0388eaa27412d5aca026815d636e018202811d6d0df99ac502031bf953eff472fdcc018202811d6d13cdffb29d51d99322bdff5f2211018202811d6d0a0f742023def783a307a986912e018202811d6d01920d8043ca89b5239253284e42018202811d6c0b7a86d7375468fac667a0a527016c29508e458543d8aa4df2abee7883018302821d6d0139601a2efabe717e604cbb4894018302821d6d02247f7a7b6594320649aa03aba1018302821d7fffffffffffffffffffffffffffffffffffffff73c0c716a594e00d54e3c4cbc9018302821d7ffffffffffffffffffffffffffffffffffffffdc7b88c420e53a9890533129f6f01830290911d7fffffffffffffffffffffffffffffffffffffff465fda27eb4d63ded474e5f832019091027ffffffffffffffff5f6af8f7b3396644f18e157960000000000000000000000000105711340daa0d5f769dba1915cef59f0815a5506027d0267a36c0c95b3975ab3ee5b203a7614a3f75373f047d803ae7b6687f2b393909302929092017d57115e47018c7177eebf7cd370a3356a1b7863008a5ae8028c72b88642840160ae1d92915050565b60007ffffffffffffffffffffffffffffffffffffffffffffffffdb731c958f34d94c1821361197957506000919050565b680755bf798b4a1bf1e582126119eb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f4558505f4f564552464c4f57000000000000000000000000000000000000000060448201526064016104ad565b6503782dace9d9604e83901b059150600060606bb17217f7d1cf79abc9e3b39884821b056b80000000000000000000000001901d6bb17217f7d1cf79abc9e3b39881029093037fffffffffffffffffffffffffffffffffffffffdbf3ccf1604d263450f02a550481018102606090811d6d0277594991cfc85f6e2461837cd9018202811d7fffffffffffffffffffffffffffffffffffffe5adedaa1cb095af9e4da10e363c018202811d6db1bbb201f443cf962f1a1d3db4a5018202811d7ffffffffffffffffffffffffffffffffffffd38dc772608b0ae56cce01296c0eb018202811d6e05180bb14799ab47a8a8cb2a527d57016d02d16720577bd19bf614176fe9ea6c10fe68e7fd37d0007b713f765084018402831d9081019084017ffffffffffffffffffffffffffffffffffffffe2c69812cf03b0763fd454a8f7e010290911d6e0587f503bb6ea29d25fcb7401964500190910279d835ebba824c98fb31b83b2ca45c000000000000000000000000010574029d9dc38563c32e5c2f6dc192ee70ef65f9978af30260c3939093039290921c92915050565b600080611b9386611c81565b9050611ba181868686611cb3565b9695505050505050565b6000808211611c16576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f554e444546494e4544000000000000000000000000000000000000000000000060448201526064016104ad565b5060016fffffffffffffffffffffffffffffffff821160071b82811c67ffffffffffffffff1060061b1782811c63ffffffff1060051b1782811c61ffff1060041b1782811c60ff10600390811b90911783811c600f1060021b1783811c909110821b1791821c111790565b60608180519060200120604051602001611c9d91815260200190565b6040516020818303038152906040529050919050565b6000806000611cc3878686611cf0565b91509150818015611ce557508051602080830191909120875191880191909120145b979650505050505050565b600060606000611cff85611e0b565b90506000806000611d11848a89611f06565b81519295509093509150158080611d255750815b611db1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f4d65726b6c65547269653a2070726f76696465642070726f6f6620697320696e60448201527f76616c696400000000000000000000000000000000000000000000000000000060648201526084016104ad565b600081611dcd5760405180602001604052806000815250611df9565b611df986611ddc6001886136d5565b81518110611dec57611dec613b6b565b602002602001015161248f565b919b919a509098505050505050505050565b60606000611e18836124b9565b90506000815167ffffffffffffffff811115611e3657611e3661335c565b604051908082528060200260200182016040528015611e7b57816020015b6040805180820190915260608082526020820152815260200190600190039081611e545790505b50905060005b8251811015611efe576000611eae848381518110611ea157611ea1613b6b565b60200260200101516124ec565b90506040518060400160405280828152602001611eca836124b9565b815250838381518110611edf57611edf613b6b565b6020026020010181905250508080611ef690613b1f565b915050611e81565b509392505050565b60006060818080611f16876125b3565b90506000869050600080611f3d604051806040016040528060608152602001606081525090565b60005b8c5181101561244b578c8181518110611f5b57611f5b613b6b565b602002602001015191508284611f719190613b07565b9350611f7e600188613b07565b965083600003611fff57815180516020909101208514611ffa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f4d65726b6c65547269653a20696e76616c696420726f6f74206861736800000060448201526064016104ad565b61213b565b8151516020116120a157815180516020909101208514611ffa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602760248201527f4d65726b6c65547269653a20696e76616c6964206c6172676520696e7465726e60448201527f616c20686173680000000000000000000000000000000000000000000000000060648201526084016104ad565b815185906120ae90613bf1565b1461213b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4d65726b6c65547269653a20696e76616c696420696e7465726e616c206e6f6460448201527f652068617368000000000000000000000000000000000000000000000000000060648201526084016104ad565b61214760106001613b07565b826020015151036121b9578551841461244b57600086858151811061216e5761216e613b6b565b602001015160f81c60f81b60f81c9050600083602001518260ff168151811061219957612199613b6b565b602002602001015190506121ac81612736565b9650600194505050612439565b6002826020015151036123b15760006121d18361276c565b90506000816000815181106121e8576121e8613b6b565b016020015160f81c905060006121ff600283613c33565b61220a906002613c55565b9050600061221b848360ff16612790565b905060006122298b8a612790565b9050600061223783836127c6565b905060ff85166002148061224e575060ff85166003145b156122a4578083511480156122635750808251145b1561227557612272818b613b07565b99505b507f8000000000000000000000000000000000000000000000000000000000000000995061244b945050505050565b60ff851615806122b7575060ff85166001145b1561232957825181146122f357507f8000000000000000000000000000000000000000000000000000000000000000995061244b945050505050565b61231a886020015160018151811061230d5761230d613b6b565b6020026020010151612736565b9a509750612439945050505050565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603260248201527f4d65726b6c65547269653a2072656365697665642061206e6f6465207769746860448201527f20616e20756e6b6e6f776e20707265666978000000000000000000000000000060648201526084016104ad565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f4d65726b6c65547269653a20726563656976656420616e20756e70617273656160448201527f626c65206e6f646500000000000000000000000000000000000000000000000060648201526084016104ad565b8061244381613b1f565b915050611f40565b507f800000000000000000000000000000000000000000000000000000000000000084148661247a8786612790565b909e909d50909b509950505050505050505050565b602081015180516060916112b9916124a9906001906136d5565b81518110611ea157611ea1613b6b565b6040805180820182526000808252602091820152815180830190925282518252808301908201526060906112b990612872565b606060008060006124fc85612acb565b91945092509050600081600181111561251757612517613c78565b146125a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f524c505265616465723a20696e76616c696420524c502062797465732076616c60448201527f756500000000000000000000000000000000000000000000000000000000000060648201526084016104ad565b6116e185602001518484612fb6565b60606000825160026125c59190613927565b67ffffffffffffffff8111156125dd576125dd61335c565b6040519080825280601f01601f191660200182016040528015612607576020820181803683370190505b50905060005b835181101561272f57600484828151811061262a5761262a613b6b565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c8261265f836002613927565b8151811061266f5761266f613b6b565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060108482815181106126b2576126b2613b6b565b01602001516126c4919060f81c613c33565b60f81b826126d3836002613927565b6126de906001613b07565b815181106126ee576126ee613b6b565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508061272781613b1f565b91505061260d565b5092915050565b600060606020836000015110156127575761275083613094565b9050612763565b612760836124ec565b90505b610beb81613bf1565b60606112b961278b8360200151600081518110611ea157611ea1613b6b565b6125b3565b6060825182106127af57506040805160208101909152600081526112b9565b6112b683838486516127c191906136d5565b61309f565b6000805b8084511180156127da5750808351115b801561285b57508281815181106127f3576127f3613b6b565b602001015160f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191684828151811061283257612832613b6b565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016145b156112b6578061286a81613b1f565b9150506127ca565b606060008061288084612acb565b9193509091506001905081600181111561289c5761289c613c78565b14612929576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f524c505265616465723a20696e76616c696420524c50206c6973742076616c7560448201527f650000000000000000000000000000000000000000000000000000000000000060648201526084016104ad565b6040805160208082526104208201909252600091816020015b60408051808201909152600080825260208201528152602001906001900390816129425790505090506000835b8651811015612ac05760208210612a08576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603460248201527f524c505265616465723a2070726f766964656420524c50206c6973742065786360448201527f65656473206d6178206c697374206c656e67746800000000000000000000000060648201526084016104ad565b600080612a456040518060400160405280858c60000151612a2991906136d5565b8152602001858c60200151612a3e9190613b07565b9052612acb565b509150915060405180604001604052808383612a619190613b07565b8152602001848b60200151612a769190613b07565b815250858581518110612a8b57612a8b613b6b565b6020908102919091010152612aa1600185613b07565b9350612aad8183613b07565b612ab79084613b07565b9250505061296f565b508152949350505050565b600080600080846000015111612b63576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f524c505265616465723a20524c50206974656d2063616e6e6f74206265206e7560448201527f6c6c00000000000000000000000000000000000000000000000000000000000060648201526084016104ad565b6020840151805160001a607f8111612b88576000600160009450945094505050612faf565b60b78111612c44576000612b9d6080836136d5565b905080876000015111612c32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f524c505265616465723a20696e76616c696420524c502073686f72742073747260448201527f696e67000000000000000000000000000000000000000000000000000000000060648201526084016104ad565b60019550935060009250612faf915050565b60bf8111612db3576000612c5960b7836136d5565b905080876000015111612cee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f524c505265616465723a20696e76616c696420524c50206c6f6e67207374726960448201527f6e67206c656e677468000000000000000000000000000000000000000000000060648201526084016104ad565b600183015160208290036101000a9004612d088183613b07565b885111612d97576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f524c505265616465723a20696e76616c696420524c50206c6f6e67207374726960448201527f6e6700000000000000000000000000000000000000000000000000000000000060648201526084016104ad565b612da2826001613b07565b9650945060009350612faf92505050565b60f78111612e6e576000612dc860c0836136d5565b905080876000015111612e5d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f524c505265616465723a20696e76616c696420524c502073686f7274206c697360448201527f740000000000000000000000000000000000000000000000000000000000000060648201526084016104ad565b600195509350849250612faf915050565b6000612e7b60f7836136d5565b905080876000015111612f10576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602760248201527f524c505265616465723a20696e76616c696420524c50206c6f6e67206c69737460448201527f206c656e6774680000000000000000000000000000000000000000000000000060648201526084016104ad565b600183015160208290036101000a9004612f2a8183613b07565b885111612f93576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f524c505265616465723a20696e76616c696420524c50206c6f6e67206c69737460448201526064016104ad565b612f9e826001613b07565b9650945060019350612faf92505050565b9193909250565b606060008267ffffffffffffffff811115612fd357612fd361335c565b6040519080825280601f01601f191660200182016040528015612ffd576020820181803683370190505b5090508051600003613010579050610beb565b600061301c8587613b07565b90506020820160005b6130306020876139c8565b8110156130675782518252613046602084613b07565b9250613053602083613b07565b91508061305f81613b1f565b915050613025565b5060006001602087066020036101000a039050808251168119845116178252839450505050509392505050565b60606112b982613277565b60608182601f01101561310e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f736c6963655f6f766572666c6f7700000000000000000000000000000000000060448201526064016104ad565b82828401101561317a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f736c6963655f6f766572666c6f7700000000000000000000000000000000000060448201526064016104ad565b818301845110156131e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f736c6963655f6f75744f66426f756e647300000000000000000000000000000060448201526064016104ad565b606082158015613206576040519150600082526020820160405261326e565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101561323f578051835260209283019201613227565b5050858452601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016604052505b50949350505050565b60606112b9826020015160008460000151612fb6565b60005b838110156132a8578181015183820152602001613290565b838111156132b7576000848401525b50505050565b600081518084526132d581602086016020860161328d565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006112b660208301846132bd565b60006020828403121561332c57600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461335757600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160c0810167ffffffffffffffff811182821017156133ae576133ae61335c565b60405290565b600082601f8301126133c557600080fd5b813567ffffffffffffffff808211156133e0576133e061335c565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156134265761342661335c565b8160405283815286602085880101111561343f57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080600060a0868803121561347757600080fd5b61348086613333565b945060208601359350604086013567ffffffffffffffff80821682146134a557600080fd5b90935060608701359081151582146134bc57600080fd5b909250608087013590808211156134d257600080fd5b506134df888289016133b4565b9150509295509295909350565b6000608082840312156134fe57600080fd5b50919050565b60008083601f84011261351657600080fd5b50813567ffffffffffffffff81111561352e57600080fd5b60208301915083602082850101111561354657600080fd5b9250929050565b600080600080600060e0868803121561356557600080fd5b853567ffffffffffffffff8082111561357d57600080fd5b9087019060c0828a03121561359157600080fd5b61359961338b565b823581526135a960208401613333565b60208201526135ba60408401613333565b6040820152606083013560608201526080830135608082015260a0830135828111156135e557600080fd5b6135f18b8286016133b4565b60a08301525096506020880135955061360d8960408a016134ec565b945060c088013591508082111561362357600080fd5b5061363088828901613504565b969995985093965092949392505050565b8581528460208201527fffffffffffffffff0000000000000000000000000000000000000000000000008460c01b16604082015282151560f81b60488201526000825161369581604985016020870161328d565b919091016049019695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000828210156136e7576136e76136a6565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261372a5761372a6136ec565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561377e5761377e6136a6565b500590565b6000808312837f8000000000000000000000000000000000000000000000000000000000000000018312811516156137bd576137bd6136a6565b837f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0183138116156137f1576137f16136a6565b50500390565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600084136000841385830485118282161615613838576138386136a6565b7f80000000000000000000000000000000000000000000000000000000000000006000871286820588128184161615613873576138736136a6565b6000871292508782058712848416161561388f5761388f6136a6565b878505871281841616156138a5576138a56136a6565b505050929093029392505050565b6000808212827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038413811516156138ed576138ed6136a6565b827f8000000000000000000000000000000000000000000000000000000000000000038412811615613921576139216136a6565b50500190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561395f5761395f6136a6565b500290565b600067ffffffffffffffff808316818516808303821115613987576139876136a6565b01949350505050565b60006fffffffffffffffffffffffffffffffff808316818516818304811182151516156139bf576139bf6136a6565b02949350505050565b6000826139d7576139d76136ec565b500490565b600084516139ee81846020890161328d565b80830190507f2e000000000000000000000000000000000000000000000000000000000000008082528551613a2a816001850160208a0161328d565b60019201918201528351613a4581600284016020880161328d565b0160020195945050505050565b600060408284031215613a6457600080fd5b6040516040810181811067ffffffffffffffff82111715613a8757613a8761335c565b604052825181526020928301519281019290925250919050565b600060808284031215613ab357600080fd5b6040516080810181811067ffffffffffffffff82111715613ad657613ad661335c565b8060405250823581526020830135602082015260408301356040820152606083013560608201528091505092915050565b60008219821115613b1a57613b1a6136a6565b500190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203613b5057613b506136a6565b5060010190565b600082613b6657613b666136ec565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b868152600073ffffffffffffffffffffffffffffffffffffffff808816602084015280871660408401525084606083015283608083015260c060a0830152613be560c08301846132bd565b98975050505050505050565b805160208083015191908110156134fe577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b600060ff831680613c4657613c466136ec565b8060ff84160691505092915050565b600060ff821660ff841680821015613c6f57613c6f6136a6565b90039392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fdfea164736f6c634300080f000a"

func init() {
	if err := json.Unmarshal([]byte(OptimismPortalStorageLayoutJSON), OptimismPortalStorageLayout); err != nil {
		panic(err)
	}

	layouts["OptimismPortal"] = OptimismPortalStorageLayout
	deployedBytecodes["OptimismPortal"] = OptimismPortalDeployedBin
}
