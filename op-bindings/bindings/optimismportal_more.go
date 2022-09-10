// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"github.com/ethereum-optimism/optimism/op-bindings/solc"
)

const OptimismPortalStorageLayoutJSON = "{\"storage\":[{\"astId\":26128,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"_initialized\",\"offset\":0,\"slot\":\"0\",\"type\":\"t_uint8\"},{\"astId\":26131,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"_initializing\",\"offset\":1,\"slot\":\"0\",\"type\":\"t_bool\"},{\"astId\":1374,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"params\",\"offset\":0,\"slot\":\"1\",\"type\":\"t_struct(ResourceParams)1344_storage\"},{\"astId\":1379,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"__gap\",\"offset\":0,\"slot\":\"2\",\"type\":\"t_array(t_uint256)49_storage\"},{\"astId\":947,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"l2Sender\",\"offset\":0,\"slot\":\"51\",\"type\":\"t_address\"},{\"astId\":960,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"finalizedWithdrawals\",\"offset\":0,\"slot\":\"52\",\"type\":\"t_mapping(t_bytes32,t_bool)\"},{\"astId\":965,\"contract\":\"contracts/L1/OptimismPortal.sol:OptimismPortal\",\"label\":\"__gap\",\"offset\":0,\"slot\":\"53\",\"type\":\"t_array(t_uint256)48_storage\"}],\"types\":{\"t_address\":{\"encoding\":\"inplace\",\"label\":\"address\",\"numberOfBytes\":\"20\"},\"t_array(t_uint256)48_storage\":{\"encoding\":\"inplace\",\"label\":\"uint256[48]\",\"numberOfBytes\":\"1536\"},\"t_array(t_uint256)49_storage\":{\"encoding\":\"inplace\",\"label\":\"uint256[49]\",\"numberOfBytes\":\"1568\"},\"t_bool\":{\"encoding\":\"inplace\",\"label\":\"bool\",\"numberOfBytes\":\"1\"},\"t_bytes32\":{\"encoding\":\"inplace\",\"label\":\"bytes32\",\"numberOfBytes\":\"32\"},\"t_mapping(t_bytes32,t_bool)\":{\"encoding\":\"mapping\",\"label\":\"mapping(bytes32 =\u003e bool)\",\"numberOfBytes\":\"32\",\"key\":\"t_bytes32\",\"value\":\"t_bool\"},\"t_struct(ResourceParams)1344_storage\":{\"encoding\":\"inplace\",\"label\":\"struct ResourceMetering.ResourceParams\",\"numberOfBytes\":\"32\"},\"t_uint128\":{\"encoding\":\"inplace\",\"label\":\"uint128\",\"numberOfBytes\":\"16\"},\"t_uint256\":{\"encoding\":\"inplace\",\"label\":\"uint256\",\"numberOfBytes\":\"32\"},\"t_uint64\":{\"encoding\":\"inplace\",\"label\":\"uint64\",\"numberOfBytes\":\"8\"},\"t_uint8\":{\"encoding\":\"inplace\",\"label\":\"uint8\",\"numberOfBytes\":\"1\"}}}"

var OptimismPortalStorageLayout = new(solc.StorageLayout)

var OptimismPortalDeployedBin = "0x6080604052600436106100f65760003560e01c8063a14238e71161008a578063cff0ab9611610059578063cff0ab96146102f7578063e9e05c4214610398578063f4daa291146103ab578063fdc9fe1d146103df57600080fd5b8063a14238e71461026d578063c4fc4798146102ad578063ca3e99ba146102cd578063cd7c9789146102e257600080fd5b80636bb0291e116100c65780636bb0291e146102005780638129fc1c14610215578063867ead131461022a5780639bf62d821461024057600080fd5b80621c2ff61461012257806313620abd1461018057806354fd4d50146101b957806364b79208146101db57600080fd5b3661011d5761011b3334620186a06000604051806020016040528060008152506103f2565b005b600080fd5b34801561012e57600080fd5b506101567f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b34801561018c57600080fd5b50610198633b9aca0081565b6040516fffffffffffffffffffffffffffffffff9091168152602001610177565b3480156101c557600080fd5b506101ce6108bf565b60405161017791906132fa565b3480156101e757600080fd5b506101f2627a120081565b604051908152602001610177565b34801561020c57600080fd5b506101f2600481565b34801561022157600080fd5b5061011b610962565b34801561023657600080fd5b506101f261271081565b34801561024c57600080fd5b506033546101569073ffffffffffffffffffffffffffffffffffffffff1681565b34801561027957600080fd5b5061029d61028836600461330d565b60346020526000908152604090205460ff1681565b6040519015158152602001610177565b3480156102b957600080fd5b5061029d6102c836600461330d565b610b20565b3480156102d957600080fd5b506101f2610be5565b3480156102ee57600080fd5b506101f2600881565b34801561030357600080fd5b5060015461035f906fffffffffffffffffffffffffffffffff81169067ffffffffffffffff7001000000000000000000000000000000008204811691780100000000000000000000000000000000000000000000000090041683565b604080516fffffffffffffffffffffffffffffffff909416845267ffffffffffffffff9283166020850152911690820152606001610177565b61011b6103a6366004613452565b6103f2565b3480156103b757600080fd5b506101f27f000000000000000000000000000000000000000000000000000000000000000081565b61011b6103ed366004613540565b610bf6565b8260005a905083156104a95773ffffffffffffffffffffffffffffffffffffffff8716156104a957604080517f08c379a00000000000000000000000000000000000000000000000000000000081526020600482015260248101919091527f4f7074696d69736d506f7274616c3a206d7573742073656e6420746f2061646460448201527f72657373283029207768656e206372656174696e67206120636f6e747261637460648201526084015b60405180910390fd5b333281146104ca575033731111000000000000000000000000000000001111015b600034888888886040516020016104e5959493929190613634565b604051602081830303815290604052905060008973ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fb3813568d9991fc951961fcb4c784893574240a28925604d09fc577c55bb7c328460405161055591906132fa565b60405180910390a45050600154600090610595907801000000000000000000000000000000000000000000000000900467ffffffffffffffff16436136c8565b9050801561071e5760006105ad6004627a120061370e565b6001546105d89190700100000000000000000000000000000000900467ffffffffffffffff16613776565b9050600060086105ec6004627a120061370e565b60015461060c9085906fffffffffffffffffffffffffffffffff166137ea565b610616919061370e565b610620919061370e565b60015490915060009061066c906106569061064e9085906fffffffffffffffffffffffffffffffff166138a6565b612710611297565b6fffffffffffffffffffffffffffffffff6112b2565b905060018411156106df576106dc610656670de0b6b3a76400006106c861069460088361370e565b6106a690670de0b6b3a7640000613776565b6106b160018a6136c8565b6106c390670de0b6b3a764000061391a565b6112c1565b6106d290856137ea565b61064e919061370e565b90505b6fffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000067ffffffffffffffff4316021760015550505b60018054849190601090610751908490700100000000000000000000000000000000900467ffffffffffffffff16613957565b92506101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550627a1200600160000160109054906101000a900467ffffffffffffffff1667ffffffffffffffff16131561082d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603e60248201527f5265736f757263654d65746572696e673a2063616e6e6f7420627579206d6f7260448201527f6520676173207468616e20617661696c61626c6520676173206c696d6974000060648201526084016104a0565b600154600090610859906fffffffffffffffffffffffffffffffff1667ffffffffffffffff8616613983565b6fffffffffffffffffffffffffffffffff169050600061087d48633b9aca006112f2565b61088790836139bb565b905060005a61089690866136c8565b9050808211156108b2576108b26108ad82846136c8565b611302565b5050505050505050505050565b60606108ea7f0000000000000000000000000000000000000000000000000000000000000000611330565b6109137f0000000000000000000000000000000000000000000000000000000000000000611330565b61093c7f0000000000000000000000000000000000000000000000000000000000000000611330565b60405160200161094e939291906139cf565b604051602081830303815290604052905090565b600054610100900460ff16158080156109825750600054600160ff909116105b8061099c5750303b15801561099c575060005460ff166001145b610a28576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016104a0565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790558015610a8657600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b603380547fffffffffffffffffffffffff00000000000000000000000000000000000000001661dead179055610aba61146d565b8015610b1d57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b6040517fa25ae55700000000000000000000000000000000000000000000000000000000815260048101829052600090819073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063a25ae557906024016040805180830381865afa158015610baf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bd39190613a45565b9050610bde81611550565b9392505050565b610bf36004627a120061370e565b81565b60335473ffffffffffffffffffffffffffffffffffffffff1661dead14610c9f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603f60248201527f4f7074696d69736d506f7274616c3a2063616e206f6e6c79207472696767657260448201527f206f6e65207769746864726177616c20706572207472616e73616374696f6e0060648201526084016104a0565b3073ffffffffffffffffffffffffffffffffffffffff16856040015173ffffffffffffffffffffffffffffffffffffffff1603610d5e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603f60248201527f4f7074696d69736d506f7274616c3a20796f752063616e6e6f742073656e642060448201527f6d6573736167657320746f2074686520706f7274616c20636f6e74726163740060648201526084016104a0565b6040517fa25ae557000000000000000000000000000000000000000000000000000000008152600481018590526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063a25ae557906024016040805180830381865afa158015610deb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e0f9190613a45565b9050610e1a81611550565b610ea6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602d60248201527f4f7074696d69736d506f7274616c3a2070726f706f73616c206973206e6f742060448201527f7965742066696e616c697a65640000000000000000000000000000000000000060648201526084016104a0565b610ebd610eb836869003860186613a94565b61158a565b815114610f4c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f4f7074696d69736d506f7274616c3a20696e76616c6964206f7574707574207260448201527f6f6f742070726f6f66000000000000000000000000000000000000000000000060648201526084016104a0565b6000610f57876115e6565b9050610f9e81866040013586868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061161692505050565b61102a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603260248201527f4f7074696d69736d506f7274616c3a20696e76616c696420776974686472617760448201527f616c20696e636c7573696f6e2070726f6f66000000000000000000000000000060648201526084016104a0565b60008181526034602052604090205460ff16156110c9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603560248201527f4f7074696d69736d506f7274616c3a207769746864726177616c20686173206160448201527f6c7265616479206265656e2066696e616c697a6564000000000000000000000060648201526084016104a0565b600081815260346020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055608087015161111290614e2090613afa565b5a10156111a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603760248201527f4f7074696d69736d506f7274616c3a20696e73756666696369656e742067617360448201527f20746f2066696e616c697a65207769746864726177616c00000000000000000060648201526084016104a0565b8660200151603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000611204886040015189608001518a606001518b60a001516116dd565b603380547fffffffffffffffffffffffff00000000000000000000000000000000000000001661dead17905560405190915082907fdb5c7652857aa163daadd670e116628fb42e869d8ac4251ef8971d9e5727df1b9061126990841515815260200190565b60405180910390a25050505050505050565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b6000818312156112a757816112a9565b825b90505b92915050565b60008183126112a757816112a9565b60006112a9670de0b6b3a7640000836112d9866116f7565b6112e391906137ea565b6112ed919061370e565b61193b565b6000818310156112a757816112a9565b6000805a90505b825a61131590836136c8565b101561132b5761132482613b12565b9150611309565b505050565b60608160000361137357505060408051808201909152600181527f3000000000000000000000000000000000000000000000000000000000000000602082015290565b8160005b811561139d578061138781613b12565b91506113969050600a836139bb565b9150611377565b60008167ffffffffffffffff8111156113b8576113b861334f565b6040519080825280601f01601f1916602001820160405280156113e2576020820181803683370190505b5090505b8415611465576113f76001836136c8565b9150611404600a86613b4a565b61140f906030613afa565b60f81b81838151811061142457611424613b5e565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535061145e600a866139bb565b94506113e6565b949350505050565b600054610100900460ff16611504576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e6700000000000000000000000000000000000000000060648201526084016104a0565b60408051606081018252633b9aca00808252600060208301524367ffffffffffffffff169190920181905278010000000000000000000000000000000000000000000000000217600155565b60007f000000000000000000000000000000000000000000000000000000000000000082602001516115829190613afa565b421192915050565b600081600001518260200151836040015184606001516040516020016115c9949392919093845260208401929092526040830152606082015260800190565b604051602081830303815290604052805190602001209050919050565b80516020808301516040808501516060860151608087015160a088015193516000976115c9979096959101613b8d565b604080516020810185905260009181018290528190606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828252805160209182012090830181905292506116d49101604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828201909152600182527f01000000000000000000000000000000000000000000000000000000000000006020830152908587611b7a565b95945050505050565b600080600080845160208601878a8af19695505050505050565b6000808213611762576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f554e444546494e4544000000000000000000000000000000000000000000000060448201526064016104a0565b6000606061176f84611b9e565b03609f8181039490941b90931c6c465772b2bbbb5f824b15207a3081018102606090811d6d0388eaa27412d5aca026815d636e018202811d6d0df99ac502031bf953eff472fdcc018202811d6d13cdffb29d51d99322bdff5f2211018202811d6d0a0f742023def783a307a986912e018202811d6d01920d8043ca89b5239253284e42018202811d6c0b7a86d7375468fac667a0a527016c29508e458543d8aa4df2abee7883018302821d6d0139601a2efabe717e604cbb4894018302821d6d02247f7a7b6594320649aa03aba1018302821d7fffffffffffffffffffffffffffffffffffffff73c0c716a594e00d54e3c4cbc9018302821d7ffffffffffffffffffffffffffffffffffffffdc7b88c420e53a9890533129f6f01830290911d7fffffffffffffffffffffffffffffffffffffff465fda27eb4d63ded474e5f832019091027ffffffffffffffff5f6af8f7b3396644f18e157960000000000000000000000000105711340daa0d5f769dba1915cef59f0815a5506027d0267a36c0c95b3975ab3ee5b203a7614a3f75373f047d803ae7b6687f2b393909302929092017d57115e47018c7177eebf7cd370a3356a1b7863008a5ae8028c72b88642840160ae1d92915050565b60007ffffffffffffffffffffffffffffffffffffffffffffffffdb731c958f34d94c1821361196c57506000919050565b680755bf798b4a1bf1e582126119de576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f4558505f4f564552464c4f57000000000000000000000000000000000000000060448201526064016104a0565b6503782dace9d9604e83901b059150600060606bb17217f7d1cf79abc9e3b39884821b056b80000000000000000000000001901d6bb17217f7d1cf79abc9e3b39881029093037fffffffffffffffffffffffffffffffffffffffdbf3ccf1604d263450f02a550481018102606090811d6d0277594991cfc85f6e2461837cd9018202811d7fffffffffffffffffffffffffffffffffffffe5adedaa1cb095af9e4da10e363c018202811d6db1bbb201f443cf962f1a1d3db4a5018202811d7ffffffffffffffffffffffffffffffffffffd38dc772608b0ae56cce01296c0eb018202811d6e05180bb14799ab47a8a8cb2a527d57016d02d16720577bd19bf614176fe9ea6c10fe68e7fd37d0007b713f765084018402831d9081019084017ffffffffffffffffffffffffffffffffffffffe2c69812cf03b0763fd454a8f7e010290911d6e0587f503bb6ea29d25fcb7401964500190910279d835ebba824c98fb31b83b2ca45c000000000000000000000000010574029d9dc38563c32e5c2f6dc192ee70ef65f9978af30260c3939093039290921c92915050565b600080611b8686611c74565b9050611b9481868686611ca6565b9695505050505050565b6000808211611c09576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f554e444546494e4544000000000000000000000000000000000000000000000060448201526064016104a0565b5060016fffffffffffffffffffffffffffffffff821160071b82811c67ffffffffffffffff1060061b1782811c63ffffffff1060051b1782811c61ffff1060041b1782811c60ff10600390811b90911783811c600f1060021b1783811c909110821b1791821c111790565b60608180519060200120604051602001611c9091815260200190565b6040516020818303038152906040529050919050565b6000806000611cb6878686611ce3565b91509150818015611cd857508051602080830191909120875191880191909120145b979650505050505050565b600060606000611cf285611dfe565b90506000806000611d04848a89611ef9565b81519295509093509150158080611d185750815b611da4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602560248201527f4d65726b6c65547269653a2070726f76696465642070726f6f6620697320696e60448201527f76616c696400000000000000000000000000000000000000000000000000000060648201526084016104a0565b600081611dc05760405180602001604052806000815250611dec565b611dec86611dcf6001886136c8565b81518110611ddf57611ddf613b5e565b6020026020010151612482565b919b919a509098505050505050505050565b60606000611e0b836124ac565b90506000815167ffffffffffffffff811115611e2957611e2961334f565b604051908082528060200260200182016040528015611e6e57816020015b6040805180820190915260608082526020820152815260200190600190039081611e475790505b50905060005b8251811015611ef1576000611ea1848381518110611e9457611e94613b5e565b60200260200101516124df565b90506040518060400160405280828152602001611ebd836124ac565b815250838381518110611ed257611ed2613b5e565b6020026020010181905250508080611ee990613b12565b915050611e74565b509392505050565b60006060818080611f09876125a6565b90506000869050600080611f30604051806040016040528060608152602001606081525090565b60005b8c5181101561243e578c8181518110611f4e57611f4e613b5e565b602002602001015191508284611f649190613afa565b9350611f71600188613afa565b965083600003611ff257815180516020909101208514611fed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f4d65726b6c65547269653a20696e76616c696420726f6f74206861736800000060448201526064016104a0565b61212e565b81515160201161209457815180516020909101208514611fed576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602760248201527f4d65726b6c65547269653a20696e76616c6964206c6172676520696e7465726e60448201527f616c20686173680000000000000000000000000000000000000000000000000060648201526084016104a0565b815185906120a190613be4565b1461212e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4d65726b6c65547269653a20696e76616c696420696e7465726e616c206e6f6460448201527f652068617368000000000000000000000000000000000000000000000000000060648201526084016104a0565b61213a60106001613afa565b826020015151036121ac578551841461243e57600086858151811061216157612161613b5e565b602001015160f81c60f81b60f81c9050600083602001518260ff168151811061218c5761218c613b5e565b6020026020010151905061219f81612729565b965060019450505061242c565b6002826020015151036123a45760006121c48361275f565b90506000816000815181106121db576121db613b5e565b016020015160f81c905060006121f2600283613c26565b6121fd906002613c48565b9050600061220e848360ff16612783565b9050600061221c8b8a612783565b9050600061222a83836127b9565b905060ff851660021480612241575060ff85166003145b15612297578083511480156122565750808251145b1561226857612265818b613afa565b99505b507f8000000000000000000000000000000000000000000000000000000000000000995061243e945050505050565b60ff851615806122aa575060ff85166001145b1561231c57825181146122e657507f8000000000000000000000000000000000000000000000000000000000000000995061243e945050505050565b61230d886020015160018151811061230057612300613b5e565b6020026020010151612729565b9a50975061242c945050505050565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603260248201527f4d65726b6c65547269653a2072656365697665642061206e6f6465207769746860448201527f20616e20756e6b6e6f776e20707265666978000000000000000000000000000060648201526084016104a0565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f4d65726b6c65547269653a20726563656976656420616e20756e70617273656160448201527f626c65206e6f646500000000000000000000000000000000000000000000000060648201526084016104a0565b8061243681613b12565b915050611f33565b507f800000000000000000000000000000000000000000000000000000000000000084148661246d8786612783565b909e909d50909b509950505050505050505050565b602081015180516060916112ac9161249c906001906136c8565b81518110611e9457611e94613b5e565b6040805180820182526000808252602091820152815180830190925282518252808301908201526060906112ac90612865565b606060008060006124ef85612abe565b91945092509050600081600181111561250a5761250a613c6b565b14612597576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f524c505265616465723a20696e76616c696420524c502062797465732076616c60448201527f756500000000000000000000000000000000000000000000000000000000000060648201526084016104a0565b6116d485602001518484612fa9565b60606000825160026125b8919061391a565b67ffffffffffffffff8111156125d0576125d061334f565b6040519080825280601f01601f1916602001820160405280156125fa576020820181803683370190505b50905060005b835181101561272257600484828151811061261d5761261d613b5e565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c8261265283600261391a565b8151811061266257612662613b5e565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060108482815181106126a5576126a5613b5e565b01602001516126b7919060f81c613c26565b60f81b826126c683600261391a565b6126d1906001613afa565b815181106126e1576126e1613b5e565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508061271a81613b12565b915050612600565b5092915050565b6000606060208360000151101561274a5761274383613087565b9050612756565b612753836124df565b90505b610bde81613be4565b60606112ac61277e8360200151600081518110611e9457611e94613b5e565b6125a6565b6060825182106127a257506040805160208101909152600081526112ac565b6112a983838486516127b491906136c8565b613092565b6000805b8084511180156127cd5750808351115b801561284e57508281815181106127e6576127e6613b5e565b602001015160f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191684828151811061282557612825613b5e565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016145b156112a9578061285d81613b12565b9150506127bd565b606060008061287384612abe565b9193509091506001905081600181111561288f5761288f613c6b565b1461291c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f524c505265616465723a20696e76616c696420524c50206c6973742076616c7560448201527f650000000000000000000000000000000000000000000000000000000000000060648201526084016104a0565b6040805160208082526104208201909252600091816020015b60408051808201909152600080825260208201528152602001906001900390816129355790505090506000835b8651811015612ab357602082106129fb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603460248201527f524c505265616465723a2070726f766964656420524c50206c6973742065786360448201527f65656473206d6178206c697374206c656e67746800000000000000000000000060648201526084016104a0565b600080612a386040518060400160405280858c60000151612a1c91906136c8565b8152602001858c60200151612a319190613afa565b9052612abe565b509150915060405180604001604052808383612a549190613afa565b8152602001848b60200151612a699190613afa565b815250858581518110612a7e57612a7e613b5e565b6020908102919091010152612a94600185613afa565b9350612aa08183613afa565b612aaa9084613afa565b92505050612962565b508152949350505050565b600080600080846000015111612b56576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f524c505265616465723a20524c50206974656d2063616e6e6f74206265206e7560448201527f6c6c00000000000000000000000000000000000000000000000000000000000060648201526084016104a0565b6020840151805160001a607f8111612b7b576000600160009450945094505050612fa2565b60b78111612c37576000612b906080836136c8565b905080876000015111612c25576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f524c505265616465723a20696e76616c696420524c502073686f72742073747260448201527f696e67000000000000000000000000000000000000000000000000000000000060648201526084016104a0565b60019550935060009250612fa2915050565b60bf8111612da6576000612c4c60b7836136c8565b905080876000015111612ce1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f524c505265616465723a20696e76616c696420524c50206c6f6e67207374726960448201527f6e67206c656e677468000000000000000000000000000000000000000000000060648201526084016104a0565b600183015160208290036101000a9004612cfb8183613afa565b885111612d8a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f524c505265616465723a20696e76616c696420524c50206c6f6e67207374726960448201527f6e6700000000000000000000000000000000000000000000000000000000000060648201526084016104a0565b612d95826001613afa565b9650945060009350612fa292505050565b60f78111612e61576000612dbb60c0836136c8565b905080876000015111612e50576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f524c505265616465723a20696e76616c696420524c502073686f7274206c697360448201527f740000000000000000000000000000000000000000000000000000000000000060648201526084016104a0565b600195509350849250612fa2915050565b6000612e6e60f7836136c8565b905080876000015111612f03576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602760248201527f524c505265616465723a20696e76616c696420524c50206c6f6e67206c69737460448201527f206c656e6774680000000000000000000000000000000000000000000000000060648201526084016104a0565b600183015160208290036101000a9004612f1d8183613afa565b885111612f86576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f524c505265616465723a20696e76616c696420524c50206c6f6e67206c69737460448201526064016104a0565b612f91826001613afa565b9650945060019350612fa292505050565b9193909250565b606060008267ffffffffffffffff811115612fc657612fc661334f565b6040519080825280601f01601f191660200182016040528015612ff0576020820181803683370190505b5090508051600003613003579050610bde565b600061300f8587613afa565b90506020820160005b6130236020876139bb565b81101561305a5782518252613039602084613afa565b9250613046602083613afa565b91508061305281613b12565b915050613018565b5060006001602087066020036101000a039050808251168119845116178252839450505050509392505050565b60606112ac8261326a565b60608182601f011015613101576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f736c6963655f6f766572666c6f7700000000000000000000000000000000000060448201526064016104a0565b82828401101561316d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f736c6963655f6f766572666c6f7700000000000000000000000000000000000060448201526064016104a0565b818301845110156131da576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f736c6963655f6f75744f66426f756e647300000000000000000000000000000060448201526064016104a0565b6060821580156131f95760405191506000825260208201604052613261565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101561323257805183526020928301920161321a565b5050858452601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016604052505b50949350505050565b60606112ac826020015160008460000151612fa9565b60005b8381101561329b578181015183820152602001613283565b838111156132aa576000848401525b50505050565b600081518084526132c8816020860160208601613280565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006112a960208301846132b0565b60006020828403121561331f57600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461334a57600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160c0810167ffffffffffffffff811182821017156133a1576133a161334f565b60405290565b600082601f8301126133b857600080fd5b813567ffffffffffffffff808211156133d3576133d361334f565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019082821181831017156134195761341961334f565b8160405283815286602085880101111561343257600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080600060a0868803121561346a57600080fd5b61347386613326565b945060208601359350604086013567ffffffffffffffff808216821461349857600080fd5b90935060608701359081151582146134af57600080fd5b909250608087013590808211156134c557600080fd5b506134d2888289016133a7565b9150509295509295909350565b6000608082840312156134f157600080fd5b50919050565b60008083601f84011261350957600080fd5b50813567ffffffffffffffff81111561352157600080fd5b60208301915083602082850101111561353957600080fd5b9250929050565b600080600080600060e0868803121561355857600080fd5b853567ffffffffffffffff8082111561357057600080fd5b9087019060c0828a03121561358457600080fd5b61358c61337e565b8235815261359c60208401613326565b60208201526135ad60408401613326565b6040820152606083013560608201526080830135608082015260a0830135828111156135d857600080fd5b6135e48b8286016133a7565b60a0830152509650602088013595506136008960408a016134df565b945060c088013591508082111561361657600080fd5b50613623888289016134f7565b969995985093965092949392505050565b8581528460208201527fffffffffffffffff0000000000000000000000000000000000000000000000008460c01b16604082015282151560f81b604882015260008251613688816049850160208701613280565b919091016049019695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000828210156136da576136da613699565b500390565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261371d5761371d6136df565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561377157613771613699565b500590565b6000808312837f8000000000000000000000000000000000000000000000000000000000000000018312811516156137b0576137b0613699565b837f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0183138116156137e4576137e4613699565b50500390565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60008413600084138583048511828216161561382b5761382b613699565b7f8000000000000000000000000000000000000000000000000000000000000000600087128682058812818416161561386657613866613699565b6000871292508782058712848416161561388257613882613699565b8785058712818416161561389857613898613699565b505050929093029392505050565b6000808212827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038413811516156138e0576138e0613699565b827f800000000000000000000000000000000000000000000000000000000000000003841281161561391457613914613699565b50500190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561395257613952613699565b500290565b600067ffffffffffffffff80831681851680830382111561397a5761397a613699565b01949350505050565b60006fffffffffffffffffffffffffffffffff808316818516818304811182151516156139b2576139b2613699565b02949350505050565b6000826139ca576139ca6136df565b500490565b600084516139e1818460208901613280565b80830190507f2e000000000000000000000000000000000000000000000000000000000000008082528551613a1d816001850160208a01613280565b60019201918201528351613a38816002840160208801613280565b0160020195945050505050565b600060408284031215613a5757600080fd5b6040516040810181811067ffffffffffffffff82111715613a7a57613a7a61334f565b604052825181526020928301519281019290925250919050565b600060808284031215613aa657600080fd5b6040516080810181811067ffffffffffffffff82111715613ac957613ac961334f565b8060405250823581526020830135602082015260408301356040820152606083013560608201528091505092915050565b60008219821115613b0d57613b0d613699565b500190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203613b4357613b43613699565b5060010190565b600082613b5957613b596136df565b500690565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b868152600073ffffffffffffffffffffffffffffffffffffffff808816602084015280871660408401525084606083015283608083015260c060a0830152613bd860c08301846132b0565b98975050505050505050565b805160208083015191908110156134f1577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60209190910360031b1b16919050565b600060ff831680613c3957613c396136df565b8060ff84160691505092915050565b600060ff821660ff841680821015613c6257613c62613699565b90039392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fdfea164736f6c634300080f000a"

func init() {
	if err := json.Unmarshal([]byte(OptimismPortalStorageLayoutJSON), OptimismPortalStorageLayout); err != nil {
		panic(err)
	}

	layouts["OptimismPortal"] = OptimismPortalStorageLayout
	deployedBytecodes["OptimismPortal"] = OptimismPortalDeployedBin
}
