package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"

	_ "embed"
)

var (
	//go:embed test_binary.msg
	tgBinaryMessage []byte

	//go:embed test_string.msg
	tgStringMessage string
)

const (
	tcEnck = "61713777ce5bd64ec2a43c299c0d125d8187a55158f0894482c36ea36d5c9463f8ffb68080c5b4fc23abc83975c75504256b356bdb2a54567fe552ab41d2cf561772406a366f5b17e213f19629170fe9e9e533afaff63c486c5fdd98deaa95739a4e9e2f7608e01aa71b11de92ffd07d967a52e6c14a9c1c9122cab9f171c9fbc63f89e493c4d551808be6b57c45a096b5c13785d818f3baa0b33cb4fa10b9c97e44252d92c0b5438dc332781fd3d2b4ea27d0a15ef430f9ef61264dab21059604f351b940d2321bc8c82c2c0e523432c890866c0bd561642769e75d6b8696823c4a53a875bbeb6f788f8dfee665a32e23518c73f289900d1029616bc25c84f44346175bbe88663b5d335ca118f1ed7007c4173d2f009b7b83136e024f78592513bc07321a1679b89f6544279032160d6f301cc91b487efce4feac06dd30443228c8b9b542f5a4d0c429f27f977038f1a38375106981f77a487f99d4d9831a7a57485aa4fa99b19ada7f0a45c29a65dc483279f2485f03c1d346749ddabfe3c95ff04770ecfd57d54c4b4eb56755af2ff36e110b0d63fcee7d44602bb391909a99526222cb8ab552c902f728eff898a5cb8a906545480e2dfaafef921d2e6338cbf0bc420078526c34c2980f316458c26c8d11114f413d4a26a58f670be521e5430be7279a42757fa0c742782e9874ce0e914f83ebe3491d0ae16483d44e0b05486ff8c2663be2eae67696264fab8308e0051176b37fe4892b18562d041346f3c8811dbda228822229ed19d585ceb066900455945c877adcf40978580c5581c849a548aa239936f84e4e6a0dafa9b9c23ae1602b80d90d4e548a6ebc40c8f3662377706bc95e861f7db8e555b238a5934fdec4ed5ec9f8d6b73745a69e914683da64cd3fc53e2c9c2b39bade5e0e3d2e7ce144031fb50e9824361a2d021a10ec1c37a73bbe55d5a7681787614ab321d6e8a3d0359dcb40496120000949a94e2a48fb2fc9581c694ec715b875ce08d035121bf0c32e522cb92d9cbcaaf20f8efa5b470db10aa9a7e3794a1d92aa2c48e50c53571e3eefddd2cbfc4aa1477eec6d8d77cf1c55523987eced6011e51c34b40826a12869870e18d9f4c853a18881b33c6239bef36b324a9df6c31e3420ad408f9423b218068d45892693e054f463e54b1e0bfcf988af716c40dd2ed2d03df3baf4370b795fadd988e0cda26dd549b4798e19909ca6bab96efd3405b96fe57269ba10762b9b2677d0c9080d0c2ca0bc15a99322257d84666df045459a3c9c2de381f964b371924a82c8e3043788439eabdc2c3018a0349933ceb20481ebc19366fda201b1c3c03a8321313467367d43a88972f0f25ffc79abf78b7c416f4db9c2e86ee5d397a9760b51851cff6a082997c09e3a8c7af66f4afe7c1858ecd943522e40db780c3a20c37a0eacb779a829c0dbc42a7e5ccf93d0073b6d71068aa60afe261a14dc8264f7a851cf660cac47cfd67d1c66bc5658b9a902030ded710e09fdbdfe5ab72dcea679e9ace4dd5181"
	tcEncd = "a60b5a3cf8480d09f96d7378c72c47a186e4952a9cf5ff58bd475466211340cbff1f8b1c5fb62d73d84eeab0eb3889b74f082cafe78cdd16ed575d054e68521a1babc416db8883a3d285e960158d54408acbcc86df66e134ee8a5862a9b344a4734561d3fc778c765b6d5072e1492943478e86ecba59bce78f3338145dfa1bbf721460a545b74c86fad5f61786e87040b4f21cf003720bda88312de2356b71a856bc1c9c39e565205e8b3c5f3e13d898e3b9f5be6c0327f0cb894ae03bf17babf38ad2726738e7870c31e9b1a29d1ee17247d8d6b0683c94fbea0248fa50858300ab66b84487764cd24de6ddd441358ddaf664bfeed881769f1b7efa89922e04a6dac400f78bd8949b1150066fe7a0dad14e45ac5aefe4fd524618494af32f32a1feeccb4ebe06e8d7325dfbe920d3470e506a247a53329d90b9a0778f439dc787aa60960010b6796eac18896e0c3c5b608ad443d63fb353d7aa76638720be73b07e861226cc382b7330ea55d93ac877bdd2a4e3d9dfeec285a2c0992261546284d035f211070e3fd90797d4583331e965e64858a9160882a1b41d0e440060a2447059bbe15c32edfb11a185c9983f32d6b26e3984b8a9525e87ca41e4f3b5cec3815333518040ab81c2873f2a61eca648af5647c7920750bd4124c9603361289939eaab36ea39b6d5c7bbb293d85a44fa012c543aa30f78458fc1f84f7c65ea1f4cf0d6a1c615fde7202f9fe6b7e080466adec0611536f14d83c4175e470f23144ea13b6e7b3359dffe41343efe00b7fa046f951ad8b24b76dbb5e51775aa30e971fcee8988027a3e2135528156d02806310060b304c8c56656fbe00e008aa14e0b7cf3181ce5ae61f22da1bd85d09ec50c379653f29fe50aaa294ca0dc41fbaaee89aa5da530e94ddc6f1389e2605192a5fe057dfafa60cbd92d7547dca2c9f0b76187f4b615ad191e748aaa8fb17226df249ee77475ae3a5b225f9d0b1c424c6903d8e33e42730e4b4f66bc6e0ad00281ed4b4a5429bc3ee9d1dc9dce3ab35f36cb1d17b2e1d190056ef57fe8a851d7cef451db1499dad1b0de62264b0060b1fe7b5f358c90b07050ea97313b6f10b7f169a73e38421b6ec3a0fe65d9f3add095fcb0d41c3fe80a58b2c0f1364d6369b9cfb7c48430714229499d7fdd0102e9f28655db49324a258115a3fc3687403c46612002dcbcf749c06dd3bc1a7e446eed0f90b87717abe7107de37b6ff7bae9bce6e9e8dc330c9e828821fe696fbd2049906cc1621d1016c3a96a39ed073ca0be07220b3836f29c7e38e4fa86ad8b07f1741b08e0c78730d3b4593869638ee9c1be61e237507e26086066fd1eef31bece38926ebdf7efd631431ab1e09cd3ad5e16ae411cb43a3a7ee3cd35944602df7cd2bcf6888402b8597b39e3ac29fee2ca161373c5cdbeca01724bb169e33422aec0efcce7e0d817eb077c6d54f5f3bcecc90ec7f211e8f30ede4598a1ef6d1ce3c7ea18d6798134f7997fe2713fd2c8ceef6f00d37e346a55d943b300b9e521e809c9847954506588ab2953640f0b1debe49f537895dc7c460b9e810e21f9a016b20179d6952081fd1e2081127941b6b16b8cca576cd0d0666316158785b4ab4b2cc535943ef0df9144dba098a09c10703b82a43eabaa1403fe37983cbb5d9fab9e05c3b563cb5123eb2348984a96dbcf3924c63f0e274d37c0759cf7726040e55c0643413642662671115f9430c816f6fa776caa3a6e61d3e847b2c7613052e89d13fcf75401a66f69829af0bd1880ea3cd8648713acadb5e24aeaaf6d0158093827629bca1b8f8515e9fc8f9a6436b966d8e6bd69cfd60e7b02e377e756489fcfebda141d2f38a9e40e92e21f7c362447984df4cfabdc328a9371236958ecbbd34d2e249469b2c72c11e95477bbfa8411419fb8d9a84c273da259c71c1d92860d1465113fae85c05ee723a9238c5a09a2b2c708177f5be242abf0dadeab60cdf917dbe369a04103a6821a36374692dd9f041cf4ea4a7a0e4550195acb8ebb3d0c9872437576f2eda22c4f4d8b2a90441b554aaefba268bd2a77a592fc10b476b38a6e71e15bdd49f55839c3292d59726df33a52b1a2343aa9f52e956bd6367b38897ed028d4969f5bb592b63c79376e5326a7ea172e3bdae661e4ed9789287a3657f854bef04518cc87a8ed9a165fde489d22c0289eee5f288953bcc71110c8eece894f05eccf7dd8eeed823e3e9b16aa6efdfb2af2ebb8861d43b8012ad4ca7406e2de041321d75194e4eb55d2ec0e4ea0937c5b3ca7a424f74991f644b850b511b039dfeae8ce1e8e564b945dbb6ea4a8c1bf645c558107912cc1207542bff0bbccc56780ee7fcaf5623ca527981ecfacab40ed53913cfc88bd84479e6957105cb8c411ae5f566cf1e2ef794d45e29b2dfee9b0ec58fa2e7abda78ed085b1c4cd37bc0f5dbe5fe04d1b39fcf186a61a84955057edef706b9c48e9900250218c04e5a5712da45c6bac31dc2e175f2efabc77b49a71f8777c97fd70e69be4a84803336b9d4284bec1eb3c3b72dd1305a2d86ca59132776548a0456b1fdffac21ab823b2bb45310a327764fe361a2fd353633addc875b8f0f5f9598c94b5f20df7a1b218adb61d2a62b62e2979b79964cf23f98b6fddb7a6e1dac86c3b701f05fa422198e682f35408ebb28a538a4c3b2e246aeae0274480f522b2c1a3c9676a654138b3215389abb7afdd4523cf0f8eb7c3da802e392748ca6adfe6ee69af23d7ff59cdff5cb27b82de4ff753cc00ce316fabbcb4bcafe5dab1793f5857a82260a6b63733caee8d5aa464cb1616db7adde5d4d4c550f4c81e4d941e45d447bdba55c5f99dcc62268d0038f2591eb939920a3ce6b5e80299b6a7e9f7d9488429d6dd0a92831ea6b57d3114776fe7a1e7b39a125ef5a12ede41a7e28e8d71bd43058ff0d20dbcfa45a3cfbb826f8d4586caf387d46bba4173d26c60b8b9a06e89c34c9bec1b68505e1ab8af3ae934337475d38ef0c00422081487a5db546bf613cb40844dda459d0bc4b7b1937afb7b9c57c486a0359e52d13628aba30173b6aaf6de3920b031c927ddfd22419b945ae6678770e8975fab3260e1db857dddb4697a127b85516c04af3dbd738ce37a9af3af390c6cea90db5626faa4ef510259d3e5e0f927af76efee23225b38f65a342df4155bde77c43cd06b0b4edcd612981dab54acb0f4a5abee51cdaf07eddca5b0bb5b3579316f836221eedd517f36bddbe0f50a79400f5522c19a1c48732fd5315a7ae22b530c3d877992286a02b4e4567d27937e4cf94829450c57111fec05a781326ece47cfb17ccb52337733213edea1dfebcd22cd4132f1548c69307633ced89ee8a3425c50d545acba139f491a09dd42aec948f95dd1c616ee3d60b43024d9bdb936f84efbca4d5600212aeb2d524d5b60b4382db7134d818bc32999c8ed4cbd02406d5b979b7d949124fc7316631c7f85177a536f9052fda702e58f0f453aa61406d9174b4a154d5e972ac37ae1688dda4dcfd4e75bb109b7b960f124efa9682db6819d20ffd06236d1dba879c1322406ab309541054dcb9cd0fa655fc8c9749dc08e4a8103989a253a87c616318f2e82e9169768aee446747867f00543c64a176d650e7fc3bd8818fd5c98c44990b597aa6fd929e62ccfc353670900120ed9ef7a307ffaed0246fca0c5e40021b743de364df7f2f028fe2c5a07638c8df832c8dea253e784a455ec799e2b90e83ef37feed9e5056353cd93080ca7af668d39e989921b7978e6d1c987f9f128cdcf00b32fe562f5049a1e1b7be7c6cbb93a1e1e5972aac9cd415d29aab7c19dbe119aa015752f57a87309454e62b8884811afcdec8bede1893f4d4e68727d7b16fdb95594d87aab6d66ae04b17ee1b6a8dc2e4b43fb0a4952c4f0837edc07fb1b328d57363ac52369ab3e7855825718ca0feef814826e72e74022e044fbf59e70be75f4328af5f00f87e268c5d642ba8c031270070920d44aabf84f75e6a642cacdb1673b9b1c8462b50aa09cf31e61d512918af4b07483f1878abe2be9958074349f6f84ea86616e19d4aaf829436286416bfa33722b9b9b056744a77077c306fb2bb7710e6c403df386b0588256c8d24192ced6c8fc6b12736aa518aaff811edd277c251312b43aa7a264f1aa18317a55f64a96fbd4554d0e7a7c545400c54a14f3ec8bfa32a8a0150e8e5de910704b309e77986e8997a87bced0c66eb963fbbde4c697e19ecac0895832d344cf7df04186bc7438802e9464824b2469f05654613fa301cfaaf80b2dd40a26e7d347f4ce406c03e317eb840c1d2b73c6422431c4955b53d05994e6ab5c536e1c233b630fb91421d926b2ddcc26af11d17871bfc06458379d25728a125bed0b7415beced210f8c13f52fe92e6ec667112df18995262f65f2a058a17cb768c659c7759526a7824f4f6ff9c85688f2466bcddec93e2307b5a1a5604b226ec0359e24d27d5ef7d8a21f648698ee6f87b72ef2127b05ab9b2ca5c520e3b6f2d007da343b92ce19cf8e60371623e3d12b1efc563860de726175ab05f5659a1d605e87ef13767883d91813db5b69851983d11bd1205deba099c1d8df6f89aecda902fc8b3e264821c253fe1c59804e5b04d0e864333f7696bf34012b7fa4d1cff7c4943bb1c73438615a2d7703774c087bb035e498361b904992450143a9eeea62cf3b3b3e6422f1e3b17a530fe1219c6593c1d63661ea33083550e8cd00b7b9c3efadf74a3ed2405fbd2bf51a03ebb303603dd1a46a75ced0af6ea13fc8d9663beb7e50df5b08acc1b994c99a42d66d2836718a1175f9fbdacd7a2214d507fae91eb996b3cc06c339e077aa4beefc2da72e9184f8eb8c6ef8a2febf9e796adac6aeb951c171dbbc73041ff3146a50cc15245e1458afc5cf5762b4a7f73b4198bdd945251944ab1e6b39b9063ad93fb70dd4d648d5a5047a84653fd3b6abc7a604c31537dad4778f9d057b0344b7b4251adda60f62686231e240a45d95c80c09f5f9c0693d5ec9070eeaf754d56fcbd8801bb8c8b49a7202a9313004f413f85c0f439345edd90b4d26aa6fe693ec48b6cdbc987fcf0f62e0e3c48155b7145b799d8ce9596e6f2d290f57dae540e4530d9623519587742630b393eb7ccc781390b83165b1c3c42e6f89c01dc5a713dcedaf456831c22cbb64d85159932ade3416ffb2e32a74a19672cc5e2e9f877b1e2af11fffcb0057e3ff6dc2546677070c15fd690485c0caac6e529e21bf58cd64508871e4bcd1b33c8da16e3426c5c2133b98b20f657e107a1b4493f0cb7ebf20cfa87a66736a6b429c3b466a307612e50ea7aa622d427bb37296a5468d9205a04cba602debcbe8c4229089d1e46ad38b5285b8eec9115e75622bf4012ab1ee552cfb21cad3b5bd334a2c1222aa66c132f063e5e5caf33795d2457905b965d744b75756a817f11d3097a393efce38513137b21e6ec245aca0b6d66b66eb809809900445bcc3a7d13b9134a3333042c0724ab61dae82fb9e3cbb5ad52f31e7844ddb0f5ceaab6ce0ad4fee6887961a1f5a1707e9f4592130472646e3300319fb1016bbc8e82ce4dc21bc0e8817ee99b4d8b5892576eca2845724b9e6ec5adc6227ebad2514e3325ce95a66756871f0fdbfbabeff1c8ab5e9f35a2878d8ff5a09ed306dc354c755d2e9ac27585b5c96055fc7961cca6424f06b23f310aaefd35a244037ce34014aba9bd236afcb06e905a5e5407090aa26daa2f0e2e896ea612f908f9bc6d95b37c90bbf02dc1a6c86e12a161644976bec4c0bbda27a1de49c1415d23cfb3fbe76cf980a5f31314267be406ea68f25fa58ed510349d4236bee5b6e1814609f1eef86637bd7d79b1974b932b62af32ecf6f9243611daae7aeeddd014e403ffd0d85acc6ecc2b6e1a5fe4ac13e514c097f103ab4abac56b711c654f547480e8a0e365841e16f1c75136d266277a454c54387eb3459a0b688d0767467905b6f4502c72c4461b888c76237ac87596962ab7d4f3691d80cad8e5072217634dd6bc852376677d48f265be785d0909b818cfc89f84fafc00606a48873dc2a55b3f761ccf7f692db23c1c920061c105369da530a3bef57be4455cbbbdec6eebed28a2b8b211f84262c7306b8358881ad414de274f903c2159c2a1e9d85e1b05135c142863305e7909005de60e5bd447379c67f32dbf338082d85002aa2bb5c92035229c44e87b4bff2a9c5cc3dc1c181e9fe97005d2efec52ae1b2bf491d8727a438c5b6f549af09f290040db78ff8ebf78f5218901c41d5addbf66c42c7998cf81dc914ff3c13f5808e87df4c6ee752e0ff2e24740787785a16083c43adfff1af27a87daa63e05ec16e7d44d51d7e7837baeb590efe65e3614418ffa6a38c55d692bbde4764cdb05cd048f4feeba3e3c0fd4fcee1c2bb570e7809af3fbdc615e621b509e90dbee59c9f252a1530b4adb4ad7a789f48e81dcc7d479a01bf2b00a62012cb3174e6ef533444ed2b631edef2c0f77cf0399b0880925b5d50beb04d1ad2b41dd98ca1225efeaab6f068c9eb9ac992cd3f303ffe187b4c717190a29237c26d0d03b09bad70eb5668737df548154e639593f130c50f00592a0a91a756f0aabb39eb76806990ba29e1f3f32fbab1145b70637cc4047bbb1f109396e2a42f7c9147d94a82ca2179a340f5a3fab5bf83fc3a9273c166d1fa22108a732034145cd3f6b0fcfc5f9440840e4c4d8e85e04e3f35a4fa738a8326ae53b75f155d8d704e19d55c5d2c3c07c9f1eac55a0066c59280dd5de8b2937e9204735e15ceb8f45ef52c546a8a189b96eecbfc146675cc77dfd965c3a0943b682d3c22d501dc1d9f23b2c9d4f8e0803401006b7c979f86ff4a47f5c29c72e970e75e98d211eb716f761614d5c84d9c384e94b767394837d15cf4cd8f7ee72eae9fad41ce2cacdfdb7ecf9a81b47be1c7c6dc9fd28e9e2eda54ff51b0675df15eae2d291f7d850406daea5efc6701aa535c4b79deb1df427b97fcdacea0e1fd9456bcf215eda5b8bc837d373fab166e137e0cd08aa01284da9f6f26f78d0a0b1d13ad366f5043e420b1253e232a66018a71bf373ef0957a285677cf52d54c9fa6904e0ee814cbd4fa007307e7987db51088825f378f7144efcb43f7a2499a793bbb943891e13b724429a4ce757636126150111684a16e3da559f2be75acfc07328def5ff158da6e35d1f35cfa55b9f95104cfffad88447ab1373880e94a992b42c023bc3303744743d61bee531439dab9301ccdbed9656c3d51a8e817708f929df6054dadab34dea0c070895341accc97f0f2e1bd2ac054d32856dbb2fc0f2e954313c9455ac04ddf961c70f891f1bd33e4b447942130e3d06b38f135ecb2aff072c4fef5a1acb9c27fd31863d4d984815f2a2853d22cd2627871fd6ae966e285a13ddeec11f991b8d5482d4e7e1f0575467f874111dfb129d1dc600026e01600efc0fc4cb45993d78ea1dd3c1d4407c5e245d59880e74044974bdcf566722b0e8879540a9b629e9f37658e939b2d6ea9e3ab4366f4d4e5f855715b4bdc5a4c3eaa0403ac7656d5a0aba10f8e7513d89654a24d4723452c398be4edb45d57e67d2cb446db99a7eb8f41aad56da2838fe0f2a38efa4dedefd119cd6f3b48e1407efc15f460aaee2ecf9fcf7fc0963ee2aacbb1a8652264e32e3d9def1cd8154ea8241812ace789de9ef7495f496ccf2f7500266ed87f4943873f78462e712aee6c48122a00e5ec1d5e8d245da4c5cef3d47a2499a1bea9508f117154093036b7a2a11df0b0783b2175487d80a9d8132bf27938d560e2aa7983b83d873313df89ed10309a07b28c5e2609537818215d2fc1f7180da099feb123f1118da80b26e54c4a95000f5a0871f9d21aaa8395238fdae7ff1351d4116a4e9f8866994c6348f3b0b37d0cc878c030a84b7bcdb3bfaa255d592edb8cfd84b47abb0c1e2ba664ac5935a26e94a452172ebd4a97d12a450718e058b921e977d1dd96a84e33da5251f13f145abcb91171e5aa4c85b7e9ebcf2f753f0d24da65a743e6b18f1177444891bffaceebf26c04ba224ff906502eba2ab938967a0e587214e632c03cd67d46bd9c1775e8601771ad0aa619b6121bac25a2c42f5b24bc3b677fb25a0c7bf85f144f0088a826ec58a7f33f3c70be85cc43822e65a19b4c1ded8f171f9f4b53163da980d3232ed5efa253a8230baecf3d00b69fbadc11a15e808127e413ecb98b5ffdef7e6a8cf2b7964354250f289b82fe6a3cb5bc04f1fd69b51a99943bc8dafc2c1a467c496dc404df8073f6f4babeddda7c40e91eac387ca2d418b65f50ec3c86b4cf5baec09ff16cd26c9f621448913a3b42eb5ea211a65bcdf5befb7b7b278801b83e020174c8ff6a1320ff6af736b5be70b3a39c8137dff516dc451bc13bcdecfe39ce3143a09056ea319a08443c000e343ce14a7d3ddcf0cf2f9861030e6b906f9bceade28a0c963f62dbba28eac6d6d2b477245150223bba74a95835f6bf9387eaf9e39cfab3e7eb785fb4ae83c9930c28f5a57bfe59331b068af48dab59ef7ac3b23e63b6c8f464ebcda0bb78aa4c3ed7342be3103377018897a59a9db70d3b7b62ea94c94d4e02750af3bbb57d88c58ece2ad3a12505562daa957c596bf2341d7683863c3a0d2106fef6ab7a1e262d1e4a1bf8c6ac6412cd36e3d3e67b7d1ce0261c2d26003055926a0f3c50ece1431a696332733673af31f854e01c3fa48205a7b4e2c4254f9b4886a4bd75571c84fd2aa1b58161b67c9c3b4e8f6268c14f0e12f0b0b5e47ab369abc5d65d44c16a1962a52d193c1565bda73c42d552b6aa4a1c0525c0ee4001e45b331b76e5751be568ab237cdc0ccccb30a33d2c443dd6c852d57c8cfc72665124058fde729775871b8bb22cf3623095f179f47046eaf87bb4fb965f60a3559588ded46455a9ab5a32cebc5213b99f4cf355343ee710dba9aa89b51c1f36485701e8cff3a2fefa663bedf70f7b964a687d654c3492a38623d48d1fbe8f02c00a271fa9de3c9956567378c9940b5a9fdfb242c8f05bb46ca80a72de9ccb0212ab7cb7453cdf5e7919738a5a71f01384511962b6d27201e49eccfc15bc9653ba3fdb4e976aa0f407a4fc4009d2374df83da50a9a76bef5c90dd9c4c68845728a0c6e0a999b0b2de0a6b192f4aaa0581280f02177bfacf500ac5e9b18314b223f7e87777b6189de6a3a6f5ae785f295d68d51d734457764bb449cb8d5d4df28d2df7ed9fc56bb39d25049ca9a6361376e2cd9c50753272ca5e7dbce787bbdfb4da2e5c85ce507e9c154ba6a7279081630a31cfe885f62f279142d985b2b42ad2d109b5d4438efa03e827e563c64780d226df9e75f4f8f80c512c79204753ef76a83cf46f180819c5bcc11fea3d873250d532d8ca6ee8504e959c790e9a5826f6b234187f470acdff2030cd6d8c1aa74cadba1afafe432f6e293163a5d97ebe9c3559da217bb878c0664677e64ef793610081227e5199aa1e699c4aa953222fe4748ef1fa5c8dc6dd8987d5b8137a0aaa169c5d1a8b298929b5a5a56714dce34b0551d6bfffab03d33464bbc038d4137bec8abfed5bf75f70909846f39a49cf41415ce0683e6d1bb3f7b873ab376ead831ed9ea799cbaca426f0ce5f23755127d8912f3c2e80eb6dc6bada43fcb4d4d7e3801a4af74423f04e67c2cef975c746fbb611cdfbe172322fc35d646bdd5ff07f298cecea59e843cec052697727928bcda6d12790f5ad63f23a5861702dded1097f22a0c9e53326e2027117e584e311ee7395552a0f237971d6862e4b52a85d789888f953d3eadca3f18a6476776788f5589364d74c790e28eec2d8d5b91eeba7e458790d54693ea9eac88f3e13157f87fa8f376029768a394a8900560460861b66258c5bf33d76d7c1e222aedd15af3e07a0837b1f008331a5e9fe8efc38aec87f0168c998cb8027e87e96e56b0bc12f5ddb80df4f9a424fd6bed5cdc22d20bc7b0f6664ff10356dca97ebebb10"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SMessageError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestPanicNewMessage(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewMessage([]byte{}, []byte{})
}

func TestInvalidMessage(t *testing.T) {
	t.Parallel()

	msgSize := uint64(2 << 10)

	if _, err := LoadMessage(msgSize, struct{}{}); err == nil {
		t.Error("success load message with unknown type")
		return
	}

	if _, err := LoadMessage(msgSize, []byte{123}); err == nil {
		t.Error("success load invalid message")
		return
	}

	msgBytes := joiner.NewBytesJoiner32([][]byte{[]byte("aaa"), []byte("bbb")})
	if _, err := LoadMessage(msgSize, msgBytes); err == nil {
		t.Error("success load invalid message")
		return
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()

	msgSize := uint64(8 << 10)

	msg1, err := LoadMessage(msgSize, tgBinaryMessage)
	if err != nil {
		t.Error(err)
		return
	}
	testMessage(t, msgSize, msg1)

	msg2, err := LoadMessage(msgSize, tgStringMessage)
	if err != nil {
		t.Error(err)
		return
	}
	testMessage(t, msgSize, msg2)
}

func testMessage(t *testing.T, msgSize uint64, msg IMessage) {
	if !bytes.Equal(msg.ToBytes(), tgBinaryMessage) {
		t.Error("invalid convert to bytes")
		return
	}

	if msg.ToString() != tgStringMessage {
		t.Error("invalid convert to string")
		return
	}

	if !bytes.Equal(msg.GetEnck(), encoding.HexDecode(tcEnck)) {
		t.Error("incorrect enck")
		return
	}

	if !bytes.Equal(msg.GetEncd(), encoding.HexDecode(tcEncd)) {
		t.Error("incorrect encd")
		return
	}

	msgBytes := bytes.Join([][]byte{msg.GetEnck(), msg.GetEncd()}, []byte{})
	if _, err := LoadMessage(msgSize, msgBytes); err != nil {
		t.Error("new message is invalid")
		return
	}
}
