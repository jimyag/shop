package paseto

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasetoPublicMaker(t *testing.T) {
	b, _ := hex.DecodeString("b4cbfb43df4ce210727d953e4a713307fa19bb7d9f85041438d9e11b942a37741eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	privateKey := ed25519.PrivateKey(b)

	b, _ = hex.DecodeString("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	publicKey := ed25519.PublicKey(b)
	maker, err := NewPasetoMaker(privateKey, publicKey, 6)
	require.NoError(t, err)

	payload := &Payload{
		UID:  11,
		Role: 1,
	}

	createToken, err := maker.CreateToken(payload)
	require.NoError(t, err)
	require.NotEmpty(t, createToken)

	payloads, err := maker.VerifyToken(createToken)
	require.NoError(t, err)
	require.NotEmpty(t, createToken)

	require.NotZero(t, payloads.UID)
	require.Equal(t, int32(11), payloads.UID)
}
