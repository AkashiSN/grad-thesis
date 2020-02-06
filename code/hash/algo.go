// hash5 Hash#5のアルゴリズム
func hash5(heads, tails [][]byte) string {
	byte1 := heads[0]
	byte2 := tails[len(tails)-1]
	byte3 := tails[len(tails)-2]
	bytes1 := (byte1[0] ^ byte2[0]) ^ byte3[0]
	bytes2 := (byte1[1] ^ byte2[1]) ^ byte3[1]
	bytes3 := (byte1[2] ^ byte2[2]) ^ byte3[2]

	hash := fmt.Sprintf("%02x%02x%02x%02x", byte1[0], bytes1, bytes2, bytes3)
	return hash
}

// hash6 Hash#6のアルゴリズム
func hash6(heads, tails [][]byte) string {
	byte1 := heads[0]
	byte2 := tails[len(tails)-1]
	bytes1 := (byte1[0] ^ byte1[1]) ^ byte1[2]
	bytes2 := (byte2[0] ^ byte2[1]) ^ byte2[2]

	hash := fmt.Sprintf("%02x%02x%02x", byte1[0], bytes1, bytes2)
	return hash
}

// hash7 Hash#7のアルゴリズム
func hash7(heads, tails [][]byte) string {
	byte1 := heads[0]
	byte2 := tails[len(tails)-1]
	byte3 := tails[len(tails)-2]
	bytes1 := (byte1[0] ^ byte1[1]) ^ byte1[2]
	bytes2 := (byte2[0] ^ byte2[1]) ^ byte2[2]
	bytes3 := (byte3[0] ^ byte3[1]) ^ byte3[2]

	hash := fmt.Sprintf("%02x%02x%02x%02x", byte1[0], bytes1, bytes2, bytes3)
	return hash
}
