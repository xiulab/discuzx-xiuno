package lstr

// SubStr 截取字符串
func SubStr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	// 取末尾
	endLen := start + end

	// 开始和结束大于总长度
	if start > length || endLen > length {
		return ""
	}

	// 负数时，从后面起取值
	if start < 0 {
		startAbs := length + start
		if startAbs >= 0 {
			start = startAbs
			endLen = length
		}
	}

	if start < 0 || endLen > length {
		return ""
	}

	return string(rs[start:endLen])
}
