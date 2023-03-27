package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	textdraw "github.com/rockwell-uk/go-draw/draw"

	"github.com/rockwell-uk/go-text/fonts"
	"github.com/rockwell-uk/go-text/fonts/ttf"
)

var (
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	pink  = color.RGBA{0xEC, 0x74, 0xB4, 0xFF}
)

func TestGetCharMetrics(t *testing.T) {
	tests := map[string]struct {
		wkt      string
		label    string
		fontSize float64
		spacing  float64
		expected []float64
	}{
		"Pilsworth Road": {
			"MULTILINESTRING((384342 409455.9999997657,384476.99999999994 409567.9999997657,384504.99999999994 409597.9999997654,384563.00000000006 409669.9999997661))",
			"Pilsworth Road",
			float64(34),
			float64(0),
			[]float64{20.77, 9.45, 9.45, 17, 30.23, 20.77, 13.23, 13.23, 20.77, 9.45, 22.67, 20.77, 18.91, 20.77},
		},
		"Mellor Street": {
			"MULTILINESTRING ((388874 413258.9999997683,388844.99999999994 413290.99999976775,388740.99999999994 413427.9999997701,388659.00000000006 413499.99999976833,388648 413512.9999997685,388583.00000000006 413820.9999997686))",
			"Mellor Street",
			float64(34),
			float64(0),
			[]float64{32.09, 18.91, 9.45, 9.45, 20.77, 13.23, 9.45, 22.67, 13.23, 13.23, 18.91, 18.91, 13.23},
		},
	}

	for name, tt := range tests {
		// Font
		f, err := truetype.Parse(ttf.UniversBold)
		if err != nil {
			t.Fatal(err)
		}

		// Truetype stuff
		opts := truetype.Options{
			Size: tt.fontSize,
		}
		face := truetype.NewFace(f, &opts)

		// strokestyle
		strokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: 1.0,
		}

		backgroundStrokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: 1.0,
		}

		// font
		fontData := draw2d.FontData{
			Name:   "bold",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleNormal,
		}

		typeFace := fonts.TypeFace{
			StrokeStyle:           strokeStyle,
			Color:                 pink,
			Size:                  tt.fontSize,
			FontData:              fontData,
			Face:                  face,
			BackgroundColor:       pink,
			BackgroundStrokeStyle: backgroundStrokeStyle,
			Spacing:               tt.spacing,
		}

		metrics := getCharMetrics(tt.label, typeFace)
		actual := []float64{}

		for _, m := range metrics {
			actual = append(actual, m.Width)
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("@%v@ %v: Expected [%+v]\nActual [%+v]", name, tt.label, tt.expected, actual)
		}
	}
}

//nolint:dupl
func TestGetLetterPositions(t *testing.T) {
	tests := map[string]struct {
		points   [][]float64
		label    string
		fontSize float64
		spacing  float64
		zoom     float64
		expected []LetterPosition
	}{
		"Mellor Street": {
			[][]float64{{388874, 413258.9999997683}, {388844.99999999994, 413290.99999976775}, {388740.99999999994, 413427.9999997701}, {388659.00000000006, 413499.99999976833}, {388648, 413512.9999997685}, {388583.00000000006, 413820.9999997686}},
			"Mellor Street",
			float64(60),
			float64(0),
			float64(1),
			[]LetterPosition{
				{Char: "M", X: 388859.1802609714, Y: 413245.56961127336, Angle: 132.18444331631272},
				{Char: "e", X: 388820.9410296373, Y: 413289.6155843894, Angle: 127.20300888929643},
				{Char: "l", X: 388800.7581149327, Y: 413316.20269318344, Angle: 127.20300888929643},
				{Char: "l", X: 388790.67875039927, Y: 413329.4803176171, Angle: 127.20300888929643},
				{Char: "o", X: 388780.59938586585, Y: 413342.7579420508, Angle: 127.20300888929643},
				{Char: "r", X: 388758.43324886553, Y: 413371.9575648306, Angle: 127.20300888929643},
				{Char: " ", X: 388744.3209292368, Y: 413390.5478320341, Angle: 127.20300888929643},
				{Char: "S", X: 388739.2022372033, Y: 413402.96296639036, Angle: 138.7152891060774},
				{Char: "t", X: 388709.1295993965, Y: 413429.36820934206, Angle: 138.7152891060774},
				{Char: "r", X: 388691.5909845436, Y: 413444.7679687247, Angle: 138.7152891060774},
				{Char: "e", X: 388674.05236969073, Y: 413460.16772810736, Angle: 138.7152891060774},
				{Char: "e", X: 388646.4531998187, Y: 413483.865515681, Angle: 130.23635830904382},
				{Char: "t", X: 388625.9245866997, Y: 413520.7468649566, Angle: 101.91678096578067},
			},
		},
		"Pilsworth Road 1": {
			[][]float64{{384342, 409455.9999997657}, {384476.99999999994, 409567.9999997657}, {384504.99999999994, 409597.9999997654}, {384563.00000000006, 409669.9999997661}},
			"Pilsworth Road",
			float64(34),
			float64(0),
			float64(1),
			[]LetterPosition{
				{Char: "P", X: 384334.76365949906, Y: 409464.72237447667, Angle: 39.68010608220529},
				{Char: "i", X: 384350.7486938591, Y: 409477.9840326124, Angle: 39.68010608220529},
				{Char: "l", X: 384358.0216151255, Y: 409484.01786358893, Angle: 39.68010608220529},
				{Char: "s", X: 384365.29453639186, Y: 409490.05169456545, Angle: 39.68010608220529},
				{Char: "w", X: 384378.37809845834, Y: 409500.9062053169, Angle: 39.68010608220529},
				{Char: "o", X: 384401.6437502977, Y: 409520.20807943546, Angle: 39.68010608220529},
				{Char: "r", X: 384417.62878465775, Y: 409533.4697375712, Angle: 39.68010608220529},
				{Char: "t", X: 384427.81087443064, Y: 409541.9171009383, Angle: 39.68010608220529},
				{Char: "h", X: 384437.9929642035, Y: 409550.36446430546, Angle: 39.68010608220529},
				{Char: " ", X: 384453.97799856355, Y: 409563.6261224412, Angle: 39.68010608220529},
				{Char: "R", X: 384461.16762159235, Y: 409567.6467770428, Angle: 46.97493401060472},
				{Char: "o", X: 384476.6357763281, Y: 409584.2197999738, Angle: 46.97493401060472},
				{Char: "a", X: 384490.74300740194, Y: 409598.36766078236, Angle: 51.14662565986203},
				{Char: "d", X: 384502.60580896307, Y: 409613.09389720316, Angle: 51.14662565986203},
			},
		},
		"Pilsworth Road 2": {
			[][]float64{{300, 200}, {298.078528040323, 180.49096779838717}, {292.3879532511287, 161.73165676349103}, {283.14696123025453, 144.44297669803979}, {270.71067811865476, 129.28932188134524}, {255.55702330196021, 116.85303876974548}, {238.268343236509, 107.61204674887132}, {219.50903220161283, 101.92147195967695}, {200, 100}, {180.49096779838717, 101.92147195967695}, {161.73165676349103, 107.61204674887132}, {144.4429766980398, 116.85303876974545}, {129.28932188134524, 129.28932188134524}, {116.85303876974547, 144.44297669803979}, {107.61204674887132, 161.73165676349106}, {101.92147195967695, 180.49096779838723}, {100, 200.00000000000009}, {101.92147195967698, 219.5090322016129}, {107.61204674887136, 238.26834323650908}, {116.85303876974555, 255.55702330196033}, {129.28932188134536, 270.71067811865487}, {144.44297669803993, 283.14696123025465}, {161.7316567634912, 292.38795325112875}, {180.4909677983874, 298.0785280403231}, {200.00000000000026, 300}, {219.50903220161308, 298.078528040323}, {238.26834323650925, 292.3879532511286}, {255.5570233019605, 283.14696123025436}, {270.710678118655, 270.7106781186545}, {283.1469612302547, 255.55702330195993}, {292.3879532511288, 238.26834323650863}, {298.07852804032314, 219.50903220161243}, {300, 200}},
			"Pilsworth Road",
			float64(34),
			float64(0),
			float64(1),
			[]LetterPosition{
				{Char: "P", X: 311.2787602356182, Y: 198.88913907626497, Angle: -95.62500000000007},
				{Char: "i", X: 308.5852138879073, Y: 176.08473505124144, Angle: -106.87499999999989},
				{Char: "l", X: 305.84202368785265, Y: 167.04164887857206, Angle: -106.87499999999989},
				{Char: "s", X: 302.1647364056044, Y: 155.98070398699963, Angle: -118.125},
				{Char: "w", X: 293.2655282896111, Y: 138.90764820220244, Angle: -129.375},
				{Char: "o", X: 271.3404795420657, Y: 115.14489279581252, Angle: -140.62500000000003},
				{Char: "r", X: 252.38645744752372, Y: 102.30760455744638, Angle: -151.87499999999994},
				{Char: "t", X: 238.42001586474572, Y: 95.814754199008, Angle: -163.125},
				{Char: "h", X: 223.69899114809584, Y: 90.94597640939344, Angle: -174.37500000000003},
				{Char: " ", X: 200.8072825276443, Y: 88.53231912401506, Angle: 174.37500000000003},
				{Char: "R", X: 191.40278686059204, Y: 89.45858110012941, Angle: 174.37500000000003},
				{Char: "o", X: 156.3891604127964, Y: 97.61693908625664, Angle: 151.87500000000003},
				{Char: "a", X: 136.35141384453564, Y: 108.8323190321136, Angle: 140.62500000000003},
				{Char: "d", X: 120.22837745082033, Y: 122.46527648353494, Angle: 129.375},
			},
		},
		"Turf Hill Road": {
			[][]float64{{390902, 411492.9999997673}, {390951.00000000006, 411523.99999976787}, {391010, 411571.999999767}, {391052, 411608.9999997665}, {391092.99999999994, 411656.99999976583}},
			"Turf Hill Road",
			float64(34),
			float64(0),
			float64(1),
			[]LetterPosition{
				{Char: "T", X: 390895.94072725705, Y: 411502.57755990914, Angle: 32.319616508635505},
				{Char: "u", X: 390913.4930146818, Y: 411513.68206828006, Angle: 32.319616508635505},
				{Char: "r", X: 390931.04530210653, Y: 411524.786576651, Angle: 32.319616508635505},
				{Char: "f", X: 390941.35550297465, Y: 411530.7638687093, Angle: 39.13039955650276},
				{Char: " ", X: 390950.14431629435, Y: 411537.9140897151, Angle: 39.13039955650276},
				{Char: "H", X: 390957.4747916581, Y: 411543.8778662821, Angle: 39.13039955650276},
				{Char: "i", X: 390976.5185133703, Y: 411559.37106360705, Angle: 39.13039955650276},
				{Char: "l", X: 390983.84898873407, Y: 411565.33484017407, Angle: 39.13039955650276},
				{Char: "l", X: 390991.17946409783, Y: 411571.2986167411, Angle: 39.13039955650276},
				{Char: " ", X: 390998.5099394616, Y: 411577.2623933081, Angle: 39.13039955650276},
				{Char: "R", X: 391005.4032478821, Y: 411583.05436153454, Angle: 41.37851529552499},
				{Char: "o", X: 391022.4138862951, Y: 411598.0399239458, Angle: 41.37851529552499},
				{Char: "a", X: 391037.748043557, Y: 411609.76448261284, Angle: 49.49715161429619},
				{Char: "d", X: 391050.0298209829, Y: 411624.1431488674, Angle: 49.49715161429619},
			},
		},
	}

	for name, tt := range tests {
		// Font
		f, err := truetype.Parse(ttf.UniversBold)
		if err != nil {
			t.Fatal(err)
		}

		// Truetype stuff
		opts := truetype.Options{
			Size: tt.fontSize,
		}
		face := truetype.NewFace(f, &opts)

		// strokestyle
		strokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: 1.0,
		}

		backgroundStrokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: 1.0,
		}

		// font
		fontData := draw2d.FontData{
			Name:   "bold",
			Family: draw2d.FontFamilySans,
			Style:  draw2d.FontStyleNormal,
		}

		typeFace := fonts.TypeFace{
			StrokeStyle:           strokeStyle,
			Color:                 pink,
			Size:                  tt.fontSize,
			FontData:              fontData,
			Face:                  face,
			BackgroundColor:       pink,
			BackgroundStrokeStyle: backgroundStrokeStyle,
			Spacing:               tt.spacing,
		}

		actual, err := GetLetterPositions(tt.label, tt.points, typeFace)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%v: Expected [%+v]\nActual [%+v]", name, tt.expected, actual)
		}
	}
}

func TestLoveHeart(t *testing.T) {
	tests := map[string]struct {
		dim            int
		fontData       draw2d.FontData
		fontSize       float64
		fontStroke     float64
		fontSpacing    float64
		zoom           float64
		text           string
		heartCoords    [][]float64
		heartLineWidth float64
	}{
		"4000x4000": {
			dim: 4000,
			fontData: draw2d.FontData{
				Name:   "bold",
				Family: draw2d.FontFamilySans,
				Style:  draw2d.FontStyleBold,
			},
			fontSize:       60,
			zoom:           float64(1),
			fontStroke:     float64(10),
			fontSpacing:    float64(0.5),
			text:           "Why are there two equation expressions? Because the equation expression corresponding to the horizontal and vertical direction of the cardioid is different, and the cardioid drawn with the same equation expression will change the direction by exchanging the X coordinate and Y coordinate of each point, so there will be two equation expressions.",
			heartCoords:    [][]float64{{0.0000000000000000, -500.0000000000000000}, {9.1810914497437768, -564.3260795987359870}, {30.3628542785403788, -639.0153015860223604}, {59.3413017892321690, -711.1976489719144183}, {94.4863763834901960, -779.4743388649623057}, {139.9963162096800602, -850.6263897627829920}, {186.1455181508988233, -909.9454982794251237}, {239.8061566258299990, -967.4239328283274517}, {288.0319157227647224, -1010.9372126369872831}, {354.6398008817002960, -1060.8845654086310333}, {412.6589934667308057, -1096.3519836999989820}, {474.5072050341198064, -1127.0740705274497486}, {539.7601751290404763, -1152.4384452814203996}, {607.9283086152656779, -1171.9206179716484257}, {678.4629992117854727, -1185.0908631032300491}, {750.7638291000669142, -1191.6190290204676785}, {824.1865421571777688, -1191.2772275481331690}, {898.0516771976353994, -1183.9403976669395888}, {971.6537380669594768, -1169.5847853482273422}, {1044.2707696596808091, -1148.2844283875117526}, {1115.1742030459379293, -1120.2057790128342276}, {1183.6388289684473420, -1085.6006371738271810}, {1248.9527570721511438, -1044.7976028147827492}, {1310.4272183795553701, -998.1922852975735623}, {1367.4060707236424150, -946.2365318089043740}, {1419.2748710657358515, -889.4269535604267958}, {1465.4693847972221192, -828.2930385390453694}, {1505.4834101629887755, -763.3851423344608520}, {1538.8758057352622473, -695.2626441873478598}, {1565.2766202684863401, -624.4825440678642963}, {1584.3922371150879371, -551.5887586848328965}, {1596.0094594937568218, -477.1023503758869992}, {1599.9984780732888794, -401.5128935188396895}, {1596.3146783480426620, -325.2711492428353495}, {1584.9992619054307852, -248.7831817189328092}, {1566.1786726825725964, -172.4060091703453281}, {1540.0628364339604559, -96.4448410156565075}, {1506.9422386384680976, -21.1519103261171395}, {1467.1838827178753490, 53.2731308757702280}, {1421.2261864808631344, 126.6813259673862433}, {1369.5728899156413263, 198.9711475959726670}, {1312.7860616123386990, 270.0839947484936943}, {1251.4783039993926650, 339.9986261234109293}, {1100.4497404732073846, 493.0137280682280334}, {937.8617066122181996, 640.5171071108630940}, {493.0886711217313518, 1025.5627992503702899}, {371.0385330344309978, 1140.7286811174888044}, {265.2651900839405243, 1250.5243825337227008}, {159.0014426165467682, 1377.6398026807380575}, {117.1222630306786243, 1435.8008397015571518}, {76.5946050503770124, 1499.7909834420215702}, {46.1031007431970892, 1556.4892781877260859}, {21.9115030239104378, 1612.1050070224321189}, {6.6636117958010166, 1660.0500988255234915}, {0.0000064637329671, 1699.9960683595550108}, {-7.5342274678077148, 1656.6577906824850288}, {-23.7786003038882079, 1607.2271300794532181}, {-49.0884420731063003, 1550.4484926466343495}, {-80.6779974477544783, 1492.8711987654280620}, {-122.3902534904586190, 1428.1171995959880405}, {-165.2911658251633753, 1369.4075719613606452}, {-273.5807864484402216, 1241.4190048736857079}, {-380.8394509165758564, 1131.1078586028577320}, {-504.1408056785006693, 1015.5451374495208938}, {-931.9981933553258386, 645.6448070911789046}, {-1094.8442893343776632, 498.3245234793522513}, {-1246.4211189698589806, 345.5147136081608323}, {-1308.0612425637555134, 275.6954437192920295}, {-1365.2311437224013844, 204.6776858869530429}, {-1417.3145523031307675, 132.4799729060658251}, {-1463.7450800980323038, 59.1577140119696097}, {-1504.0140747393381844, -15.1911619897985304}, {-1537.6776809096982106, -90.4216123900500151}, {-1564.3630077051723219, -166.3381102731082990}, {-1583.7733137937682386, -242.6926444395007820}, {-1595.6921360776077563, -319.1842022580980256}, {-1599.9863027001322280, -395.4598110247487739}, {-1596.6077872235689483, -471.1171735605520894}, {-1585.5943774064344325, -545.7088922210486999}, {-1567.0691489973883108, -618.7482332619154022}, {-1541.2387520860538643, -689.7163416997453851}, {-1508.3905345673522334, -758.0707765015177984}, {-1468.8885439387465794, -823.2551981897320275}, {-1423.1684647200645486, -884.7100067640983525}, {-1371.7315640321705814, -941.8836981370134254}, {-1315.1377320745939414, -994.2446828764179827}, {-1253.9977171989628459, -1041.2932926376165597}, {-1188.9646667986241937, -1082.5736877923754946}, {-1120.7250951594221533, -1117.6853748167693539}, {-1049.9894075996237461, -1146.2940441960088265}, {-977.4821165508660670, -1168.1414489848527865}, {-903.9318896058210839, -1183.0540605866603983}, {-830.0615719192908273, -1190.9502614658047150}, {-756.5783256637257637, -1191.8458638996619356}, {-684.1640275033022363, -1185.8577788606214654}, {-613.4660612877770518, -1173.2056989041150246}, {-545.0886374325199313, -1154.2117026083226392}, {-479.5847628270122982, -1129.2977346442539783}, {-417.4489757104216778, -1098.9809638516742325}, {-359.1109489049012495, -1063.8670706051707384}, {-292.0694322704976571, -1014.2798468845521711}, {-232.0547691706165665, -959.7612061912803938}, {-179.4120091470793739, -901.9303522616603459}, {-134.2891928624968898, -842.4951464222806408}, {-89.9844167635740888, -771.5193617812279854}, {-55.9733917526247211, -703.7553436248874732}, {-28.1727337804661246, -632.5666983210579701}, {-8.1871672360787056, -559.6898432325492649}, {0.0005063475946655, -500.0951977201792147}},
			heartLineWidth: 10,
		},
		"400x400": {
			dim: 400,
			fontData: draw2d.FontData{
				Name:   "bold",
				Family: draw2d.FontFamilySans,
				Style:  draw2d.FontStyleNormal,
			},
			fontSize:       6,
			zoom:           float64(1),
			fontStroke:     float64(1),
			fontSpacing:    float64(0.5),
			text:           "Why are there two equation expressions? Because the equation expression corresponding to the horizontal and vertical direction of the cardioid is different, and the cardioid drawn with the same equation expression will change the direction by exchanging the X coordinate and Y coordinate of each",
			heartCoords:    [][]float64{{0.0000000000000000, -50.0000000000000000}, {9.4486376383490196, -77.9474338864962277}, {28.8031915722764680, -101.0937212636987397}, {53.9760175129040505, -115.2438445281420343}, {82.4186542157177655, -119.1277227548133197}, {111.5174203045937986, -112.0205779012834313}, {136.7406070723642415, -94.6236531808904289}, {153.8875805735262077, -69.5262644187347831}, {159.9998478073288766, -40.1512893518839675}, {154.0062836433960456, -9.6444841015656504}, {136.9572889915641554, 19.8971147595972688}, {110.0449740473207498, 49.3013728068228048}, {26.5265190083940539, 125.0524382533722871}, {7.6594605050377016, 149.9790983442021570}, {0.0000006463732967, 169.9996068359554897}, {-8.0677997447754475, 149.2871198765427891}, {-27.3580786448440207, 124.1419004873685736}, {-109.4844289334377549, 49.8324523479352237}, {-136.5231143722401441, 20.4677685886953071}, {-153.7677680909698381, -9.0421612390050026}, {-159.9986302700132228, -39.5459811024748760}, {-154.1238752086053978, -68.9716341699745499}, {-137.1731564032170638, -94.1883698137013425}, {-112.0725095159422011, -111.7685374816769439}, {-83.0061571919290770, -119.0950261465804658}, {-54.5088637432519931, -115.4211702608322554}, {-29.2069432270497664, -101.4279846884552114}, {-8.9984416763574089, -77.1519361781227957}, {0.0000506347594666, -50.0095197720179243}},
			heartLineWidth: 1,
		},
	}

	for name, tt := range tests {
		tileWidth := float64(tt.dim)
		tileHeight := float64(tt.dim)
		centerX := tileWidth / 2
		centerY := tileHeight / 2

		bounds := [][]float64{
			{0, 0},
			{0, tileWidth},
			{tileWidth, tileHeight},
			{tileHeight, 0},
			{0, 0},
		}

		envelope, err := textdraw.ToEnvelope(bounds)
		if err != nil {
			t.Fatal(err)
		}

		scale := func(x, y float64) (float64, float64) {
			nx := envelope.Px(x) * tileWidth
			ny := tileHeight - (envelope.Py(y) * tileHeight)
			return nx, ny
		}

		m := image.NewRGBA(image.Rect(0, 0, tt.dim, tt.dim))
		draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
		gc := draw2dimg.NewGraphicContext(m)

		gc.SetDPI(72)

		// draw the line
		fillColour := white
		strokeColour := pink
		gc.Translate(centerX, centerY)
		err = textdraw.DrawCoordLine(gc, tt.heartCoords, tt.heartLineWidth, fillColour, tt.heartLineWidth, strokeColour, scale)
		if err != nil {
			t.Fatal(err)
		}

		fontSize := tt.fontSize * tt.zoom
		fontSpacing := tt.fontSpacing * tt.zoom
		strokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: tt.fontStroke,
		}

		// font options
		face := fonts.GetFace(gc, tt.fontData, fontSize)

		typeFace := fonts.TypeFace{
			StrokeStyle: strokeStyle,
			Color:       pink,
			Size:        tt.fontSize,
			FontData:    tt.fontData,
			Face:        face,
			Spacing:     fontSpacing,
		}
		fonts.SetFont(gc, typeFace)

		// text along line
		gc.Translate(0, 0)
		glyphs, err := TextAlongLine(gc, tt.text, tt.heartCoords, typeFace)
		if err != nil {
			t.Fatalf("%v: %v", name, err)
		}
		for _, glyph := range glyphs {
			err = textdraw.DrawRune(gc, glyph.Pos, face, glyph.Rotation, glyph.Char)
			if err != nil {
				t.Fatal(err)
			}
		}

		err = savePNG(fmt.Sprintf("test-output/heart-test/%v.png", tt.dim), m)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTestAlongLineOutlines(t *testing.T) {
	tests := map[string]struct {
		dim         int
		fontData    draw2d.FontData
		fontSize    float64
		fontStroke  float64
		fontSpacing float64
		zoom        float64
		text        string
		lineWidth   float64
	}{
		"400x400": {
			dim: 400,
			fontData: draw2d.FontData{
				Name:   "bold",
				Family: draw2d.FontFamilySans,
				Style:  draw2d.FontStyleNormal,
			},
			fontSize:    40,
			zoom:        float64(1),
			fontStroke:  float64(1),
			fontSpacing: float64(0.5),
			text:        "Why are there two equation expressions?",
			lineWidth:   1,
		},
	}

	for _, tt := range tests {
		tileWidth := float64(tt.dim)
		tileHeight := float64(tt.dim)
		centerX := tileWidth / 2
		centerY := tileHeight / 2

		bounds := [][]float64{
			{0, 0},
			{0, tileWidth},
			{tileWidth, tileHeight},
			{tileHeight, 0},
			{0, 0},
		}

		envelope, err := textdraw.ToEnvelope(bounds)
		if err != nil {
			t.Fatal(err)
		}

		scale := func(x, y float64) (float64, float64) {
			nx := envelope.Px(x) * tileWidth
			ny := tileHeight - (envelope.Py(y) * tileHeight)
			return nx, ny
		}

		m := image.NewRGBA(image.Rect(0, 0, tt.dim, tt.dim))
		draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.Point{0, 0}, draw.Src)
		gc := draw2dimg.NewGraphicContext(m)

		gc.SetDPI(72)
		gc.Translate(-200, -200)

		// generate a circular line to use for the test
		radius := 140.0
		numPoints := 40 // IMPORTANT: the text will never fit on the line unless the distance between each point is greater than the charWidth/2
		origin := []float64{
			200.00,
			200.00,
		}
		circleCoords, err := textdraw.Circle(
			origin,
			radius,
			numPoints,
		)
		if err != nil {
			t.Fatal(err)
		}

		// draw the line
		fillColour := black
		strokeColour := white
		gc.Translate(centerX, centerY)
		err = textdraw.DrawCoordLine(gc, circleCoords, tt.lineWidth, fillColour, tt.lineWidth, strokeColour, scale)
		if err != nil {
			t.Fatal(err)
		}

		fontSize := tt.fontSize * tt.zoom
		fontSpacing := tt.fontSpacing * tt.zoom
		strokeStyle := draw2d.StrokeStyle{
			Color: white,
			Width: tt.fontStroke,
		}

		// font options
		face := fonts.GetFace(gc, tt.fontData, fontSize)

		typeFace := fonts.TypeFace{
			StrokeStyle: strokeStyle,
			Color:       pink,
			Size:        tt.fontSize,
			FontData:    tt.fontData,
			Face:        face,
			Spacing:     fontSpacing,
		}
		fonts.SetFont(gc, typeFace)

		// text along line
		glyphs, err := TextAlongLine(gc, tt.text, circleCoords, typeFace)
		if err != nil {
			t.Fatal(err)
		}
		for _, glyph := range glyphs {
			err = textdraw.DrawRune(gc, glyph.Pos, face, glyph.Rotation, glyph.Char)
			if err != nil {
				t.Fatal(err)
			}
		}
		err = DrawGlyphOutlines(gc, tt.text, circleCoords, typeFace)
		if err != nil {
			t.Fatal(err)
		}

		err = savePNG(fmt.Sprintf("test-output/outline-test/%v.png", tt.dim), m)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func savePNG(fname string, m image.Image) error {
	dir, _ := path.Split(fname)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	return draw2dimg.SaveToPngFile(fname, m)
}
