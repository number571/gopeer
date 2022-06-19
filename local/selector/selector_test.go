package selector

import (
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
)

var (
	tgPubKeys = []string{
		`Pub(go-peer/rsa){3082020A0282020100E26F6499A5D436FB8AF3A78AE3145E970FF6AECB330EBAF6B7716570990E9639BCB5D806F574D91621778FDAEE41F7820307662F61DABA5372973BE887EE6AB6BBBBD8C4EEAD8BA1C0DB2E21245EB8329D18A9B6C2C49B80B7D0A0089C9F8613DCAFA016F19D5396410096DA4C2C2F3868977CE4BB7D07BB8E7CF461D94EDFA85AB409D645177583507B739FE6EF4C21E06AB378935C2580CE159270897FBB0A32DF5029CBE215804845EC1AFF9038C72054EF38AFB5268A7047404D6FF627C3A8E413EEBC23F6037D9425AF078A51D71871A65ED2D1DEB7F6BF35787FB3906D6771B47DB1B70189700CE69090D7B8404D9F235AD759CD631E20AB98E894480D8D9B1E997C4BBF137095E8A9CDB6B0090594148D07B125ED2DD081A92249ED60720CCAF984C9881081EDEF772B86F0E10315CBB1E41A3429AEC78795ADC92AE383B62D54F8AAFA6410E56FB29B78382D5D570FA3B94CE58B4129A909F44C1E4194685C080062DDB44537744F08D52BF2D9E8C96B7F78150E31DF0AADE71137B48DE11598EDDBEDA816EE9487262686B5DDC63004DCC4692664A2E855DDA9A5CA86E8CAC6042F544C3ADBC47BF69E6A96224F567C36C9475D90C3129D2DD1C0F7D2EE5B8815C01D30A0A0BDF60BA1B1CE10CBC59E1DE4E18A3EAABBFF669B4B2BD6BEA0246B7227145961A75ED1DB28EFA8C36C77E5C0F8DF9A3342A0487BE3AD0203010001}`,
		`Pub(go-peer/rsa){3082020A028202010094292483DB9719B1CE64AABE186708653854C462E549D227C3A4E8C1C2A93A04E314A596907C261BE1AE4B29C697B59F53D3E1B67D660424F7B45371C806CE709911DD089BEC2FC7D4E0CBCDD2ABC01082ADE4F36443EE86858EA74CB61F687DD47743E6A0D95C4090022C2F9B489C114DEB91FC1D536A11A2F150B65E33DD3041C7FCECDDE5FDB660EE6C34F054E39303101BB8057CA878DA240196DBDA34222CC144932C336FCDC263CD14B4CF349B6EB7D1F9720EFE927810E03CB629B58BD10C03852A216A46443FAD86A8FE1050DBEFF6E8673EEC08B658395233096630AC2F631EBEFB7436EFD89603997533C3420A034FBC2334C2FE58E509906B8F8B94C286DA76B95E01F84CBA335787D04020C8CA82006D89393381BD9F26498B03C0E7389EB7D7CB1A1469C482C5205AFAAF3A8B1B5A9F45277A1B769C3E0041E16B231859BD2E676AB4E5E466CCDFEEBBA19AB5658DF81C78F0778813FE499CF2FF41D59DDC76367CDD598E9A74BA1266B604F7ED632E2315B0C93E2C2D539FDEE68E05DF9D1ADEFC76D20DE049417BA4DD2C3F06354D16C9F3A37CF75E4CFB948B9DB80152C790593F4EAA4255BE9CC77FF1953E3073D83C87B70EFDD2E1FC2882E5447C5E6AE38C66D39656C58D00DB34E6BFD15E97916ED112D3DB50F91D97D291D74DCD09A0863A53C16DE5A0D001A9E9295FA99581C9B5E2047E5FDB596D0203010001}`,
		`Pub(go-peer/rsa){3082020A0282020100DD5666A211B03F2B9E5C9BFFE8160F6C1CDBA1E9F332BC3B7D39E77D6E717A6B5361625B33D66E5C3C35092AED036FCE793CA15C59F9F015E3D87398337C005BE2DA7C08409309FB430A7EBB890C5E719AFC368DFDDAE32D359E11564043E1E1E33729ACA7F9CA35D835F1883F72DAFB43E75A11C555559533B0A898323EB3B8E44DE9D1F42593A42B20BDF01BB3245214EE4838B1D7BC8518DA371D2B7971D64225AF67E751F99D4EA50C9659B0AB673B2BD9EFA492504975F05A5F9754C72A30820980F9E99070A084F395C2FF86AF0A1EA363BB80E79289F7A0EC5C3161A924366488DBCC39285E2B5E81F2989E73A8696E49831650E1CC6323776B386971820F9D026FBE59E3F44937EE23A542CF9859D4DF85E1BA6E1AF54360C47ABBFF44CD4F562D7CD27E7B1B051AEC9E88B4BC48C697FF73BECB7A6B0A312662C932957B04AAAC947D94F505FF58C7A743A53806C06CF96DF6BC6A760B420F1E075213177D80128DC52223B05CF88A8024B501B8D870413AD7185DA74D3B99480B048C2DF25A9D8A5F70741CEEA7A2B27A808AA9A3565393FC5E183FE3CB4EEE9CC466FA2129EED70BC20F9725E9649083A5CCF95493BD7C30A4BC14912008B438A06E546B1F1FF568D34F62C51AB7032CB5DF859F29D923B9AF24EFAC593D2D04E17279A0EDBBD5F4DE0601C4BE625E2E7F6280396DD014617AC1944B362812F7C30203010001}`,
	}
)

func TestSelector(t *testing.T) {
	pubKeys := []asymmetric.IPubKey{}
	for _, sPubKey := range tgPubKeys {
		pubKeys = append(pubKeys, asymmetric.LoadRSAPubKey(sPubKey))
	}

	selector := NewSelector(pubKeys)
	for i := 0; i < 5; i++ {
		checkPubKeys := selector.Shuffle().Return(selector.Length())
		if !testAreUniq(checkPubKeys) {
			t.Error("selector's list's values not unique")
			return
		}
		for i := range pubKeys {
			if pubKeys[i].Address().String() != checkPubKeys[i].Address().String() {
				return
			}
		}
	}

	t.Error("selector's shuffle does not work")
}

func testAreUniq(pubKeys []asymmetric.IPubKey) bool {
	for i := 0; i < len(pubKeys); i++ {
		for j := i + 1; j < len(pubKeys); j++ {
			if pubKeys[i].Address().String() == pubKeys[j].Address().String() {
				return false
			}
		}
	}
	return true
}
