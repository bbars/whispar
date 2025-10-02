package vp

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"

	"github.com/bbars/whispar/pkg/vpencoding"
)

type ID vpencoding.ID

func (id ID) GetId() ID {
	return id
}

type Path vpencoding.Path

func (p Path) GetId() ID {
	if l := len(p); l > 0 {
		return ID(p[l-1])
	} else {
		return ""
	}
}

func (p Path) String() string {
	sb := strings.Builder{}
	for i, id := range p {
		if i > 0 {
			sb.WriteRune(':')
		}
		sb.WriteString(string(id))
	}
	return sb.String()
}

type PathSub vpencoding.PathSub

func (p PathSub) GetId() ID {
	for i := len(p) - 1; i >= 0; i-- {
		if l := len(p[i]); l > 0 {
			return ID(p[i][l-1])
		}
	}
	return ""
}

func (p PathSub) String() string {
	sb := strings.Builder{}
	for i, ids := range p {
		if i > 0 {
			sb.WriteRune('$')
		}
		for j, id := range ids {
			if j > 0 {
				sb.WriteRune(':')
			}
			sb.WriteString(string(id))
		}
	}
	return sb.String()
}

const (
	idChars     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	idTimeStart = 0 // 2015-01-01T00:00:00Z
	idLen       = 16
)

func TimeID() ID {
	bb := make([]byte, idLen)
	t := time.Now().Unix() - idTimeStart
	base64.RawURLEncoding.Encode(bb, []byte{
		byte((t >> (0 * 8)) & 0xff),
		byte((t >> (1 * 8)) & 0xff),
		byte((t >> (2 * 8)) & 0xff),
		byte((t >> (3 * 8)) & 0xff),
		byte((t >> (4 * 8)) & 0xff),
		byte((t >> (5 * 8)) & 0xff),
		byte((t >> (6 * 8)) & 0xff),
		byte((t >> (7 * 8)) & 0xff),
	})

	// idLen MUST be greater than 8
	bb[8] = '.'
	for i := 9; i < idLen; i++ {
		bb[i] = idChars[rand.Intn(len(idChars))]
	}

	return ID(bb)
}

func RandID() ID {
	bb := make([]byte, idLen)
	for i := 0; i < idLen; i++ {
		if (i+1)%9 == 0 {
			bb[i] = '.'
		} else {
			bb[i] = idChars[rand.Intn(len(idChars))]
		}
	}

	return ID(bb)
}

func MakeSeededID(seed uint64) func() ID {
	static := make([]byte, idLen)
	for i := 0; i < 8; i++ {
		ci := (seed >> (i * 8)) & 0xff
		static[i] = idChars[int(ci)%len(idChars)]
	}

	return func() ID {
		bb := make([]byte, idLen)
		copy(bb, static)

		// idLen MUST be greater than 8
		bb[8] = '.'
		for i := 9; i < idLen; i++ {
			bb[i] = idChars[rand.Intn(len(idChars))]
		}

		return ID(bb)
	}
}
