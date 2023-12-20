package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// LogComparedBytes returns a string  the bytes of two []bytes, starting from the first byte that is different
func LogComparedBytes(name1 string, name2 string, data1 []byte, data2 []byte, numBytesBefore int, numBytesAfter int) string {
	findFirstByteDifferent := findFirstByteDifferent(data1, data2)
	if findFirstByteDifferent == -1 {
		return fmt.Sprintf("%s(%d) and %s(%d) are equal", name1, len(data1), name2, len(data2))
	}
	res := name1 + fmt.Sprintf("(%d)", len(data1)) + ": " + strSliceBytes(data1, findFirstByteDifferent, numBytesBefore, numBytesAfter) + "\n"
	res += name2 + fmt.Sprintf("(%d)", len(data1)) + ": " + strSliceBytes(data2, findFirstByteDifferent, numBytesBefore, numBytesAfter)
	return res
}

func strSliceBytes(data []byte, point int, before int, after int) string {
	res := ""
	startingPoint := max(0, point-before)
	if startingPoint > 0 {
		res += fmt.Sprintf("(%d)...", startingPoint)
	}
	endPoint := min(len(data), point+after)
	res += fmt.Sprintf("%s*%s", common.Bytes2Hex(data[startingPoint:point]), common.Bytes2Hex(data[point:endPoint]))

	if endPoint < len(data) {
		res += fmt.Sprintf("...(%d)", len(data)-endPoint)
	}
	return res
}

func findFirstByteDifferent(data1 []byte, data2 []byte) int {
	for i := 0; i < len(data1); i++ {
		if data1[i] != data2[i] {
			return i
		}
	}
	return -1
}
