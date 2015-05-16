package common

func NewCodes() map[string]int {
	Codes := map[string]int{
		"break":     1,
		"default":   2,
		"fn":        3,
		"interface": 4,
		"const":     5,
		"case":      6,
		"struct":    7,
		"if":        8,
		"else":      9,
		"package":   10,
		"switch":    11,
		"goto":      12,
		"range":     13,
		"type":      14,
		"continue":  15,
		"for":       16,
		"return":    17,
		"import":    18,
		"var":       19,
		"nil":       20,
		"int":       21,
		"double":    22,
		"string":    23,
		"byte":      24,
		"true":      25,
		"false":     26,
		"while":     27,
	}
	return Codes
}
