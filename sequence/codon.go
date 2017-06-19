package sequence

var (
	// 氨基酸对应的标准密码子
	ToCodon = map[byte][]string{
		'L': {"CUG", "UUA", "CUU", "UUG", "CUC", "CUA"},
		'R': {"CGU", "CGC", "CGG", "AGA", "CGA", "AGG"},
		'S': {"AGC", "AGU", "UCA", "UCU", "UCC", "UCG"},
		'A': {"GCA", "GCC", "GCG", "GCU"},
		'G': {"GGU", "GGC", "GGA", "GGG"},
		'P': {"CCG", "CCU", "CCA", "CCC"},
		'T': {"ACC", "ACA", "ACG", "ACU"},
		'V': {"GUU", "GUG", "GUC", "GUA"},
		'I': {"AUU", "AUC", "AUA"},
		'*': {"UAA", "UGA", "UAG"},
		'C': {"UGU", "UGC"},
		'D': {"GAU", "GAC"},
		'E': {"GAA", "GAG"},
		'F': {"UUU", "UUC"},
		'H': {"CAU", "CAC"},
		'K': {"AAA", "AAG"},
		'N': {"AAU", "AAC"},
		'Q': {"CAG", "CAA"},
		'Y': {"UAU", "UAC"},
		'M': {"AUG"},
		'W': {"UGG"}}

	// 标准密码子对应的氨基酸，添加了简并碱基和gap的Codon
	FromCodon = map[string]byte{
		"AGU": 'S', "GAA": 'E', "GGC": 'G', "GAY": 'D', "GCD": 'A', "CCN": 'P', "UCD": 'S', "GAG": 'E', "AGR": 'R',
		"GUH": 'V', "CUB": 'L', "ACR": 'T', "GUN": 'V', "GCB": 'A', "UAG": '*', "GGN": 'G', "UGU": 'W', "ACV": 'T',
		"CCD": 'P', "ACB": 'T', "CCH": 'P', "AAC": 'N', "CCG": 'P', "CUU": 'L', "GCA": 'A', "CAG": 'Q', "AUU": 'I',
		"UAC": 'Y', "CUV": 'L', "UCW": 'S', "UGY": 'C', "GCR": 'A', "CGN": 'R', "GCY": 'A', "GCU": 'A', "CGG": 'R',
		"UCR": 'S', "UUY": 'F', "AGY": 'S', "CGW": 'R', "GUB": 'V', "UAR": '*', "UCG": 'S', "UGG": 'W', "UCY": 'S',
		"AAU": 'N', "AGC": 'S', "AAA": 'K', "UAA": '*', "ACW": 'T', "CGV": 'R', "CCS": 'P', "CGR": 'R', "GCN": 'A',
		"UCN": 'S', "CCC": 'P', "GGU": 'G', "UCC": 'S', "CAC": 'H', "GUR": 'V', "GUW": 'V', "UCS": 'S', "AUY": 'I',
		"CGD": 'R', "GUD": 'V', "UCV": 'S', "AUW": 'I', "CUW": 'L', "GAR": 'E', "AUA": 'I', "CUY": 'L', "AUC": 'I',
		"GGW": 'G', "CUN": 'L', "CAU": 'H', "GAU": 'D', "ACC": 'T', "CGC": 'R', "CCR": 'P', "GUC": 'V', "CGH": 'R',
		"AGA": 'R', "ACD": 'T', "AUM": 'I', "GGR": 'G', "CAA": 'Q', "UUU": 'F', "GCG": 'A', "GAC": 'D', "CCU": 'P',
		"GCV": 'A', "UAU": 'Y', "CGY": 'R', "UUA": 'L', "AAR": 'K', "CGB": 'R', "UAY": 'Y', "UGA": '*', "AUH": 'I',
		"CUR": 'L', "CAR": 'Q', "CCY": 'P', "UUC": 'F', "CCV": 'P', "GCS": 'A', "UUR": 'L', "CCW": 'P', "AAY": 'N',
		"GGH": 'G', "CAY": 'H', "CUH": 'L', "CGS": 'R', "GGD": 'G', "GGA": 'G', "GUG": 'V', "CUC": 'L', "ACA": 'T',
		"GUU": 'V', "CCA": 'P', "ACN": 'T', "GGG": 'G', "GGB": 'G', "CUA": 'L', "CUD": 'L', "GUS": 'V', "GCC": 'A',
		"GUY": 'V', "UCB": 'S', "CCB": 'P', "UCU": 'S', "CGU": 'R', "CGA": 'R', "CUS": 'L', "GGS": 'G', "UCA": 'S',
		"AGG": 'R', "GUV": 'V', "GCH": 'A', "GGV": 'G', "CUG": 'L', "AAG": 'K', "UCH": 'S', "GUA": 'V', "ACG": 'T',
		"UGC": 'W', "ACY": 'T', "UUG": 'L', "---": '-', "ACH": 'T', "ACS": 'T', "GGY": 'G', "ACU": 'T', "GCW": 'A'}

	// 氨基酸单字母表示对应的三字母表示和中文名
	AminoAcid = map[byte][2]string{
		'A': {"Ala", "丙氨酸"},
		'R': {"Arg", "精氨酸"},
		'N': {"Asn", "天冬酰胺"},
		'D': {"Asp", "天冬氨酸"},
		'C': {"Cys", "半胱氨酸"},
		'E': {"Glu", "谷氨酸"},
		'Q': {"Gln", "谷氨酰胺"},
		'G': {"Gly", "甘氨酸"},
		'H': {"His", "组氨酸"},
		'I': {"Ile", "异亮氨酸"},
		'L': {"Leu", "亮氨酸"},
		'K': {"Lys", "赖氨酸"},
		'M': {"Met", "半胱氨酸"},
		'F': {"Phe", "苯丙氨酸"},
		'P': {"Pro", "脯氨酸"},
		'S': {"Ser", "丝氨酸"},
		'T': {"Thr", "苏氨酸"},
		'W': {"Trp", "色氨酸"},
		'Y': {"Tyr", "酪氨酸"},
		'V': {"Val", "缬氨酸"},
		'*': {"End", "终止密码子"}}

	// 大肠杆菌中某个密码子相对于所有密码子的频率
	CodonProbInEcoli = map[string]float64{
		"UUU": 24.4, "UCU": 13.1, "UAU": 21.6, "UGU": 5.90,
		"UUC": 13.9, "UCC": 9.70, "UAC": 11.7, "UGC": 5.50,
		"UUA": 17.4, "UCA": 13.1, "UAA": 2.00, "UGA": 1.10,
		"UUG": 12.9, "UCG": 8.20, "UAG": 0.30, "UGG": 13.4,
		"CUU": 14.5, "CCU": 9.50, "CAU": 12.4, "CGU": 15.9,
		"CUC": 9.50, "CCC": 6.20, "CAC": 7.30, "CGC": 14.0,
		"CUA": 5.60, "CCA": 9.10, "CAA": 14.4, "CGA": 4.80,
		"CUG": 37.4, "CCG": 14.5, "CAG": 26.7, "CGG": 7.90,
		"AUU": 29.6, "ACU": 13.1, "AAU": 29.3, "AGU": 13.2,
		"AUC": 19.4, "ACC": 18.9, "AAC": 20.3, "AGC": 14.3,
		"AUA": 13.3, "ACA": 15.1, "AAA": 37.2, "AGA": 7.10,
		"AUG": 23.7, "ACG": 13.6, "AAG": 15.3, "AGG": 4.00,
		"GUU": 21.6, "GCU": 18.9, "GAU": 33.7, "GGU": 23.7,
		"GUC": 13.1, "GCC": 21.6, "GAC": 17.9, "GGC": 20.6,
		"GUA": 13.1, "GCA": 23.0, "GAA": 35.1, "GGA": 13.6,
		"GUG": 19.9, "GCG": 21.1, "GAG": 19.4, "GGG": 12.3}
)

/*
// 标准密码子对应的氨基酸，无简并碱基版本
FromCodon = map[string]byte{
	"UUU": 'F', "UCU": 'S', "UAU": 'Y', "UGU": 'C',
	"UUC": 'F', "UCC": 'S', "UAC": 'Y', "UGC": 'C',
	"UUA": 'L', "UCA": 'S', "UAA": '*', "UGA": '*',
	"UUG": 'L', "UCG": 'S', "UAG": '*', "UGG": 'W',
	"CUU": 'L', "CCU": 'P', "CAU": 'H', "CGU": 'R',
	"CUC": 'L', "CCC": 'P', "CAC": 'H', "CGC": 'R',
	"CUA": 'L', "CCA": 'P', "CAA": 'Q', "CGA": 'R',
	"CUG": 'L', "CCG": 'P', "CAG": 'Q', "CGG": 'R',
	"AUU": 'I', "ACU": 'T', "AAU": 'N', "AGU": 'S',
	"AUC": 'I', "ACC": 'T', "AAC": 'N', "AGC": 'S',
	"AUA": 'I', "ACA": 'T', "AAA": 'K', "AGA": 'R',
	"AUG": 'M', "ACG": 'T', "AAG": 'K', "AGG": 'R',
	"GUU": 'V', "GCU": 'A', "GAU": 'D', "GGU": 'G',
	"GUC": 'V', "GCC": 'A', "GAC": 'D', "GGC": 'G',
	"GUA": 'V', "GCA": 'A', "GAA": 'E', "GGA": 'G',
	"GUG": 'V', "GCG": 'A', "GAG": 'E', "GGG": 'G'}
*/
