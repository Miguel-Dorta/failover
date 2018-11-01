package main

func boolToByte(bs [8]bool) byte {
	var result byte
	for i, bo := range bs {
		var by byte
		if bo {
			by = 1
		} else {
			by = 0
		}
		result |= by << uint(i)
	}
	return result
}

func byteToBool(b byte) [8]bool {
	var result [8]bool
	for i:=0; i<8; i++ {
		result[i] = (b & (1 << uint(i))) != 0
	}
	return result
}
