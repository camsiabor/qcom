package qio

func BytesEqual(a []byte, b []byte) bool {
	if a == nil || b == nil {
		panic("one of the input is nil")
	}
	var alen = len(a)
	var blen = len(b)
	if alen != blen {
		return false
	}

	for i := 0; i < alen; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true

}
