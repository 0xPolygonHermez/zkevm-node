package evm

// import (
// 	"github.com/holiman/uint256"
// )

// type Stack struct {
// 	data []uint256.Int
// }

// // NewStack is the constructor
// func NewStack() *Stack {
// 	return &Stack{}
// }

// func (st *Stack) push(d *uint256.Int) {
// 	// TODO: check pos <= 1024
// 	st.data = append(st.data, *d)
// }

// func (st *Stack) pop() (ret uint256.Int) {
// 	ret = st.data[len(st.data)-1]
// 	st.data = st.data[:len(st.data)-1]
// 	return ret
// }

// func (st *Stack) len() int {
// 	return len(st.data)
// }

// func (st *Stack) swap(n int) {
// 	st.data[st.len()-n], st.data[st.len()-1] = st.data[st.len()-1], st.data[st.len()-n]
// }

// func (st *Stack) dup(n int) {
// 	st.push(&st.data[st.len()-n])
// }

// func (st *Stack) peek() *uint256.Int {
// 	return &st.data[st.len()-1]
// }
