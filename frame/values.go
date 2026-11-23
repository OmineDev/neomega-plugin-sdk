package frame

import (
	"encoding/binary"
	"errors"

	"github.com/google/uuid"
)

// Values is an ordered list of byte fragments used in node messages and API calls. Scalar accessors read
// the first fragment without consuming it.
type Values [][]byte

// ErrInsufficientData indicates that a required fragment is missing or has the wrong length.
var ErrInsufficientData = errors.New("insufficient data")

// IsEmpty reports whether there are no fragments. A single empty fragment is not an empty Values.
func (r Values) IsEmpty() bool {
	return len(r) == 0
}

// EqualString compares val with the first fragment, treating a missing fragment as an empty string.
func (r Values) EqualString(val string) bool {
	if empty := r.IsEmpty(); empty {
		return val == ""
	}
	return string(r[0]) == val
}

// ToString returns the first fragment as a string, or ErrInsufficientData if there are no fragments.
func (r Values) ToString() (string, error) {
	if empty := r.IsEmpty(); empty {
		return "", ErrInsufficientData
	}
	return string(r[0]), nil
}

// ToBool reports whether the first byte of the first fragment is nonzero. A missing or empty first
// fragment returns ErrInsufficientData.
func (r Values) ToBool() (bool, error) {
	if empty := r.IsEmpty(); empty {
		return false, ErrInsufficientData
	}
	if len(r[0]) == 0 {
		return false, ErrInsufficientData
	}
	return r[0][0] != 0, nil
}

// ToByte returns the first byte of the first fragment, or ErrInsufficientData if it is missing.
func (r Values) ToByte() (byte, error) {
	if empty := r.IsEmpty(); empty {
		return 0, ErrInsufficientData
	}
	if len(r[0]) == 0 {
		return 0, ErrInsufficientData
	}
	return r[0][0], nil
}

// ConsumeHead returns the remaining fragments after the first, or the empty input unchanged.
func (r Values) ConsumeHead() Values {
	if len(r) == 0 {
		return r
	}
	return r[1:]
}

// Extend appends the fragments from each argument and returns the resulting slice.
func (r Values) Extend(vals ...Values) Values {
	newR := r
	for _, val := range vals {
		newR = append(newR, val...)
	}
	return newR
}

// ExtendFrags appends the given byte slices as separate fragments.
func (r Values) ExtendFrags(frags ...[]byte) Values {
	return append(r, frags...)
}

// ToStrings returns one string per fragment, in order.
func (r Values) ToStrings() []string {
	if empty := r.IsEmpty(); empty {
		return []string{}
	}
	ss := make([]string, len(r))
	for i, v := range r {
		ss[i] = string(v)
	}
	return ss
}

// ToInt64 decodes the first fragment as an eight-byte, little-endian signed integer.
func (r Values) ToInt64() (int64, error) {
	if empty := r.IsEmpty(); empty || len(r[0]) != 8 {
		return 0, ErrInsufficientData
	}
	return int64(binary.LittleEndian.Uint64(r[0])), nil
}

// ToInt32 decodes the first fragment as a four-byte, little-endian signed integer.
func (r Values) ToInt32() (int32, error) {
	if empty := r.IsEmpty(); empty || len(r[0]) != 4 {
		return 0, ErrInsufficientData
	}
	return int32(binary.LittleEndian.Uint32(r[0])), nil
}

// ToUint32 decodes the first fragment as a four-byte, little-endian unsigned integer.
func (r Values) ToUint32() (uint32, error) {
	if empty := r.IsEmpty(); empty || len(r[0]) != 4 {
		return 0, ErrInsufficientData
	}
	return uint32(binary.LittleEndian.Uint32(r[0])), nil
}

// ToBytes returns the first fragment without copying it, or ErrInsufficientData if it is missing.
func (r Values) ToBytes() ([]byte, error) {
	if empty := r.IsEmpty(); empty {
		return nil, ErrInsufficientData
	}
	return r[0], nil
}

// ToUUID parses the first fragment as the 16-byte binary form of a UUID.
func (r Values) ToUUID() (ud uuid.UUID, err error) {
	if empty := r.IsEmpty(); empty {
		return ud, ErrInsufficientData
	}
	return uuid.FromBytes(r[0])
}
