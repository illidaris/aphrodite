package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCount_EmptySlice 测试空切片场景
func TestCount_EmptySlice(t *testing.T) {
	srcs := []int{} // 空切片
	opts := []Option{WithBatch(5)}
	assert.Equal(t, 0, Count(srcs, opts...)) // 预期 0
}

// TestCount_ExactDivision 测试整除场景
func TestCount_ExactDivision(t *testing.T) {
	srcs := make([]int, 10) // 长度 10
	opts := []Option{WithBatch(5)}
	assert.Equal(t, 2, Count(srcs, opts...)) // 10/5=2
}

// TestCount_NonExactDivision 测试非整除场景
func TestCount_NonExactDivision(t *testing.T) {
	srcs := make([]int, 10) // 长度 10
	opts := []Option{WithBatch(3)}
	assert.Equal(t, 4, Count(srcs, opts...)) // 10/3=3.333 → 4
}

// TestCount_BatchLargerThanTotal 测试批次大小超过总长度
func TestCount_BatchLargerThanTotal(t *testing.T) {
	// 防止验证负提升
	srcs := make([]int, 5) // 长度 5
	opts := []Option{WithBatch(10)}
	assert.Equal(t, 5, Count(srcs, opts...)) // 5/10=0.5
}

// TestCount_BatchAutoCorrection 测试批次自动修正逻辑
func TestCount_BatchAutoCorrection(t *testing.T) {
	srcs := make([]int, 15)
	// 假设 WithBatchSize(0) 会被修正为默认值 1
	opts := []Option{WithBatch(0)}
	assert.Equal(t, 15, Count(srcs, opts...)) // 15/1=15
}

// TestCount_ResultAdjustment 测试结果调整逻辑
func TestCount_ResultAdjustment(t *testing.T) {
	srcs := make([]int, 500)
	// 设置最小批次数为 3
	opts := []Option{WithBatch(5), WithParallelismMax(3)}
	// 实际计算值 1 → 被调整为 3
	assert.Equal(t, 100, Count(srcs, opts...))
}
