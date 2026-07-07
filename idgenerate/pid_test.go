package idgenerate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"unicode"
)

func TestGenerateRandomID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		uid       uint64
		secretKey string
		length    int
	}{
		{
			name:      "six digits",
			uid:       12345,
			secretKey: "secret",
			length:    6,
		},
		{
			name:      "eighteen digits fast path",
			uid:       1<<63 + 99,
			secretKey: "another-secret",
			length:    18,
		},
		{
			name:      "twenty digits big integer fallback",
			uid:       987654321,
			secretKey: "secret",
			length:    20,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GenerateRandomID(tt.uid, tt.secretKey, tt.length)
			if len(got) != tt.length {
				t.Fatalf("GenerateRandomID() length = %d, want %d; id = %q", len(got), tt.length, got)
			}

			for _, r := range got {
				if !unicode.IsDigit(r) {
					t.Fatalf("GenerateRandomID() = %q, contains non-digit rune %q", got, r)
				}
			}
		})
	}
}

func TestGenerateRandomIDDeterministic(t *testing.T) {
	t.Parallel()

	got := GenerateRandomID(12345, "secret", 8)
	for i := 0; i < 100; i++ {
		if next := GenerateRandomID(12345, "secret", 8); next != got {
			t.Fatalf("GenerateRandomID() = %q, want deterministic value %q", next, got)
		}
	}
}

func TestGenerateRandomIDMatchesBigIntAlgorithm(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		uid       uint64
		secretKey string
		length    int
	}{
		{
			name:      "one digit",
			uid:       1,
			secretKey: "secret",
			length:    1,
		},
		{
			name:      "six digits",
			uid:       12345,
			secretKey: "secret",
			length:    6,
		},
		{
			name:      "eighteen digits",
			uid:       1<<63 + 99,
			secretKey: "another-secret",
			length:    18,
		},
		{
			name:      "twenty digits",
			uid:       987654321,
			secretKey: "secret",
			length:    20,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GenerateRandomID(tt.uid, tt.secretKey, tt.length)
			want := generateRandomIDWithBigInt(tt.uid, tt.secretKey, tt.length)
			if got != want {
				t.Fatalf("GenerateRandomID() = %q, want %q", got, want)
			}
		})
	}
}

func TestGenerateRandomIDDifferentInputs(t *testing.T) {
	t.Parallel()

	base := GenerateRandomID(12345, "secret", 18)
	if got := GenerateRandomID(12346, "secret", 18); got == base {
		t.Fatalf("GenerateRandomID() with different uid = %q, want different from %q", got, base)
	}

	if got := GenerateRandomID(12345, "other-secret", 18); got == base {
		t.Fatalf("GenerateRandomID() with different secret = %q, want different from %q", got, base)
	}
}

func generateRandomIDWithBigInt(uid uint64, secretKey string, length int) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	uidBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(uidBytes, uid)
	_, _ = mac.Write(uidBytes)
	hash := mac.Sum(nil)

	bigInt := new(big.Int).SetBytes(hash)
	modulus := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	resultInt := new(big.Int).Mod(bigInt, modulus)

	formatStr := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(formatStr, resultInt)
}

func TestGenerateRandomIDInvalidLength(t *testing.T) {
	t.Parallel()

	if got := GenerateRandomID(12345, "secret", 0); got != "" {
		t.Fatalf("GenerateRandomID() = %q, want empty string for zero length", got)
	}

	if got := GenerateRandomID(12345, "secret", -1); got != "" {
		t.Fatalf("GenerateRandomID() = %q, want empty string for negative length", got)
	}
}

func TestFormatFixedLengthUint64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		id     uint64
		length int
		want   string
	}{
		{
			name:   "pads leading zeroes",
			id:     42,
			length: 6,
			want:   "000042",
		},
		{
			name:   "keeps exact length",
			id:     123456,
			length: 6,
			want:   "123456",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := formatFixedLengthUint64(tt.id, tt.length); got != tt.want {
				t.Fatalf("formatFixedLengthUint64() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomIDPadsLeadingZeroes(t *testing.T) {
	t.Parallel()

	got := GenerateRandomID(1, "secret", 18)
	if len(got) != 18 {
		t.Fatalf("GenerateRandomID() length = %d, want 18", len(got))
	}

	if strings.TrimLeft(got, "0") == "" {
		t.Fatalf("GenerateRandomID() = %q, want at least one non-zero digit", got)
	}
}
