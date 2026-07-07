package idgenerate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
)

var pow10Uint64 = [...]uint64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
}

// GenerateRandomID 从 UID 生成固定长度的随机数字 ID
func GenerateRandomID(uid uint64, secretKey string, length int) string {
	if length <= 0 {
		return ""
	}

	// 1. 使用 HMAC-SHA256 进行签名防篡改
	mac := hmac.New(sha256.New, []byte(secretKey))
	var uidBytes [8]byte
	binary.BigEndian.PutUint64(uidBytes[:], uid)
	_, _ = mac.Write(uidBytes[:])
	hash := mac.Sum(nil)

	if length < len(pow10Uint64) {
		return formatFixedLengthUint64(hashModUint64(hash, pow10Uint64[length]), length)
	}

	// 2. 将 hash 转换为大整数并取模
	// big.Int 可以安全处理巨大的 hash 字节并转换为纯数字
	bigInt := new(big.Int).SetBytes(hash)

	// 计算所需的模数 (例如长度为 6 则模数为 1000000)
	modulus := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)

	// 取模得到纯数字
	resultInt := new(big.Int).Mod(bigInt, modulus)

	// 3. 格式化为固定长度（不足前补 0）
	formatStr := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(formatStr, resultInt)
}

func formatFixedLengthUint64(id uint64, length int) string {
	raw := strconv.FormatUint(id, 10)
	if len(raw) >= length {
		return raw
	}

	return strings.Repeat("0", length-len(raw)) + raw
}

func hashModUint64(hash []byte, modulus uint64) uint64 {
	var remainder uint64
	for _, b := range hash {
		hi, lo := bits.Mul64(remainder, 256)
		lo, carry := bits.Add64(lo, uint64(b), 0)
		hi += carry
		_, remainder = bits.Div64(hi, lo, modulus)
	}

	return remainder
}
