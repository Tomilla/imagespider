package util

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestEqualSliceGeneric(t *testing.T) {
    type TestInfo struct {
        lhs      interface{}
        rhs      interface{}
        excepted bool
    }
    var (
        sourceAndExcepted = []TestInfo{
            // test: string slice
            {[]string{"hello", "world"}, []string{"hello", "world"}, true,},
            {[]string{"hello", "world"}, []string{"hello", "tomi"}, false,},
            {[]string{"hello", "world"}, []string{"hello",}, false,},
            {[]bool{true, true}, []string{"hello", "world"}, false,},
            {[]bool{true, true}, []interface{}{"hello", true}, false,},
            // test: boolean slice
            {[]bool{true, true}, []bool{true, true}, true,},
            {[]bool{false, false}, []bool{false, true}, false,},
            // test: int slice
            {[]int{1, 3, 5, 7, 9}, []int{1, 3, 5, 7, 9}, true,},
            {[]int{1, 3, 5, 7, 9}, []int{1, 3, 5, 7, 10}, false,},
            // test negative int slice
            {[]int{-1, -3, -5, -7, -9}, []int{-1, -3, -5, -7, -9}, true,},
            // test: int slice v.s. int16 slice
            {[]int{1, 3, 5, 7, 9}, []int16{1, 3, 5, 7, 9}, false,},
            {[]int16{1, 3, 5, 7, 9}, []int16{1, 3, 5, 7, 9}, true,},
            // test: uint slice
            {[]uint{1, 3, 5, 7, 9}, []uint{1, 3, 5, 7, 9}, true,},
            {[]uint{1, 3, 5, 7, 9}, []uint{1, 3, 5, 7, 10}, false,},
            // test: uint64 slice
            {[]uint64{1, 3, 5, 7, 9}, []uint64{1, 3, 5, 7, 9}, true,},
            {[]uint64{1, 3, 5, 7, 9}, []uint64{1, 3, 5, 7, 10}, false,},
            // test: float32 slice
            {[]float32{1.2, 3.4, 5.6, 7.8, 9.10}, []float32{1.2, 3.4, 5.6, 7.8, 9.10}, true,},
            {[]float32{1.2, 3.4, 5.6, 7.8, 9.10}, []float32{1.2, 3.5, 5.7, 7.8, 9.10}, false,},
            // test: float64 slice
            {[]float64{1.2, 3.4, 5.6, 7.8, 9.10}, []float64{1.2, 3.4, 5.6, 7.8, 9.10}, true,},
            {[]float64{1.2, 3.3, 5.6, 7.8, 9.10}, []float64{1.2, 3.3, 5.6, 7.8, 9.10}, true,},
            // test: complex64 slice
            {[]complex64{1.2 + 3i, 2.3 + 4i, 3.4 + 5i, 4.5 + 6i, 5.6 + 7i}, []complex64{1.2 + 3i, 2.3 + 4i, 3.4 + 5i, 4.5 + 6i, 5.6 + 7i},
                true,},
            {[]complex64{1.2 + 3i, 2.3 + 4i, 3.4 + 5i, 4.5 + 6i, 5.6 + 7i}, []complex64{1.3 + 3i, 2.3 + 4i, 3.4 + 5i,
                4.5 + 6i, 5.6 + 7i},
                false,},
            // test complex128 slice
            {[]complex128{1.2 + 3i, 2.3 + 4i, 3.4 + 5i, 4.5 + 6i, 5.6 + 7i}, []complex128{1.2 + 3i, 2.3 + 4i, 3.4 + 5i,
                4.5 + 6i, 5.6 + 7i},
                true,},
            {[]complex128{1.2 + 3i, 2.3 + 4i, 3.4 + 5i, 4.5 + 6i, 5.6 + 7i}, []complex128{1.2 + 3i, 2.3 + 4i, 3.4 + 5i,
                4.5 + 6i, 5.6 + 8i},
                false,},
        }
    )
    for _, info := range sourceAndExcepted {
        assert.Equal(t, EqualSliceGeneric(info.lhs, info.rhs), info.excepted)
    }
}
