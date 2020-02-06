// splitToArrayFillByte ICNを各セクションの前後3Byteの配列に分解する
func splitToArrayFillByte(icn string) ([][]byte, [][]byte, bool) {
	idx, dx := 5, 0    // 現在の位置(最初のセクション(icn:/)は無視する),現在の位置から"/"までの距離
	var heads [][]byte // 各セクションの最初の3Byte
	var tails [][]byte // 各セクションの最後の3Byte

	for {
		if idx >= len(icn) { // 終端に達したとき
			break
		}

		dx = strings.Index(icn[idx:], "/") // 現在の位置から"/"を検索
		if dx == -1 {                      // 最後のセクション
			dx = len(icn[idx:]) // 末尾までの距離
		}
		var head []byte
		var tail []byte

		switch dx {
		case 0:
			idx += dx + 1 // 次の"/"に進む
			continue
		case 1: // セクションが1Byteの場合は2ByteでICNの長さとスラッシュの数を掛けたものを連結する
			c := len(icn) * strings.Count(icn, "/")
			buf := make([]byte, 2)
			binary.LittleEndian.PutUint16(buf, uint16(c))
			head = []byte(icn[idx : idx+1])
			head = append(head, buf...)
			tail = append(tail, buf...)
			tail = append(tail, []byte(icn[idx:idx+1])...)
		case 2: // セクションが2Byteの場合は1Byteでスラッシュの数を連結する
			c := strings.Count(icn, "/")
			head = []byte(icn[idx : idx+2])
			head = append(head, byte(c))
			tail = append(tail, byte(c))
			tail = append(tail, []byte(icn[idx:idx+2])...)
		default:
			head = []byte(icn[idx : idx+3])
			tail = []byte(icn[idx+dx-3 : idx+dx])
		}

		heads = append(heads, head)
		tails = append(tails, tail)

		idx += dx + 1 // 次の"/"に進む
	}

	if len(heads) < 3 { // カウントしたセクションが３列以下のとき
		return [][]byte{}, [][]byte{}, false
	}

	return heads, tails, true
}
