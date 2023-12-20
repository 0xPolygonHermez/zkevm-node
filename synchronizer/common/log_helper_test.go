package common

import "testing"

func TestLogComparedBytes(t *testing.T) {
	name1 := "file1.txt"
	name2 := "file2.txt"
	data1 := []byte{1, 2, 3, 4, 5}
	data2 := []byte{1, 2, 6, 4, 5}
	numBytesBefore := 2
	numBytesAfter := 2

	expected := "file1.txt(5): 0102*0304...(1)\nfile2.txt(5): 0102*0604...(1)"
	result := LogComparedBytes(name1, name2, data1, data2, numBytesBefore, numBytesAfter)
	if result != expected {
		t.Errorf("Unexpected result. Expected: %s, Got: %s", expected, result)
	}
}

func TestLogComparedBytes2(t *testing.T) {
	name1 := "file1.txt"
	name2 := "file2.txt"
	data1 := []byte{10, 20, 30, 1, 2, 3, 4, 5}
	data2 := []byte{10, 20, 30, 1, 2, 6, 4, 5}
	numBytesBefore := 2
	numBytesAfter := 2

	expected := "file1.txt(8): (3)...0102*0304...(1)\nfile2.txt(8): (3)...0102*0604...(1)"
	result := LogComparedBytes(name1, name2, data1, data2, numBytesBefore, numBytesAfter)
	if result != expected {
		t.Errorf("Unexpected result. Expected: %s, Got: %s", expected, result)
	}
}
