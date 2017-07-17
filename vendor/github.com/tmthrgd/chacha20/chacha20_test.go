// Copyright 2014 Coda Hale. All rights reserved.
// Use of this source code is governed by an MIT
// License that can be found in the LICENSE file.
//
// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package chacha20

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	codahale "github.com/codahale/chacha20"
	"github.com/tmthrgd/chacha20/internal/ref"
)

func mustHexDecode(v string) []byte {
	b, err := hex.DecodeString(v)
	if err != nil {
		panic(err)
	}

	return b
}

type testVector struct {
	key       []byte
	nonce     []byte
	keyStream []byte
	counter   uint64
}

// stolen from https://tools.ietf.org/html/rfc7539
var rfcTestVectors = []testVector{
	testVector{
		mustHexDecode("0000000000000000000000000000000000000000000000000000000000000000"),
		mustHexDecode("000000000000000000000000"),
		mustHexDecode("76b8e0ada0f13d90405d6ae55386bd28bdd219b8a08ded1aa836efcc8b770dc7" +
			"da41597c5157488d7724e03fb8d84a376a43b8f41518a11cc387b669b2ee6586"),
		0,
	},
	testVector{
		mustHexDecode("0000000000000000000000000000000000000000000000000000000000000000"),
		mustHexDecode("000000000000000000000000"),
		mustHexDecode("9f07e7be5551387a98ba977c732d080dcb0f29a048e3656912c6533e32ee7aed" +
			"29b721769ce64e43d57133b074d839d531ed1f28510afb45ace10a1f4b794d6f"),
		1,
	},
	testVector{
		mustHexDecode("0000000000000000000000000000000000000000000000000000000000000001"),
		mustHexDecode("000000000000000000000000"),
		mustHexDecode("3aeb5224ecf849929b9d828db1ced4dd832025e8018b8160b82284f3c949aa5a" +
			"8eca00bbb4a73bdad192b5c42f73f2fd4e273644c8b36125a64addeb006c13a0"),
		1,
	},
	testVector{
		mustHexDecode("00ff000000000000000000000000000000000000000000000000000000000000"),
		mustHexDecode("000000000000000000000000"),
		mustHexDecode("72d54dfbf12ec44b362692df94137f328fea8da73990265ec1bbbea1ae9af0ca" +
			"13b25aa26cb4a648cb9b9d1be65b2c0924a66c54d545ec1b7374f4872e99f096"),
		2,
	},
	testVector{
		mustHexDecode("0000000000000000000000000000000000000000000000000000000000000000"),
		mustHexDecode("000000000000000000000002"),
		mustHexDecode("c2c64d378cd536374ae204b9ef933fcd1a8b2288b3dfa49672ab765b54ee27c7" +
			"8a970e0e955c14f3a88e741b97c286f75f8fc299e8148362fa198a39531bed6d"),
		0,
	},
	testVector{
		mustHexDecode("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
		mustHexDecode("000000000000004a00000000"),
		mustHexDecode("224f51f3401bd9e12fde276fb8631ded8c131f823d2c06" +
			"e27e4fcaec9ef3cf788a3b0aa372600a92b57974cded2b" +
			"9334794cba40c63e34cdea212c4cf07d41b769a6749f3f" +
			"630f4122cafe28ec4dc47e26d4346d70b98c73f3e9c53a" +
			"c40c5945398b6eda1a832c89c167eacd901d7e2bf363"),
		1,
	},
}

// stolen from https://tools.ietf.org/html/draft-strombergson-chacha-test-vectors-01
var draftTestVectors = []testVector{
	testVector{
		mustHexDecode("0000000000000000000000000000000000000000000000000000000000000000"),
		mustHexDecode("0000000000000000"),
		mustHexDecode("76b8e0ada0f13d90405d6ae55386bd28bdd219b8a08ded1aa836efcc8b770dc7" +
			"da41597c5157488d7724e03fb8d84a376a43b8f41518a11cc387b669b2ee6586" +
			"9f07e7be5551387a98ba977c732d080dcb0f29a048e3656912c6533e32ee7aed" +
			"29b721769ce64e43d57133b074d839d531ed1f28510afb45ace10a1f4b794d6f"),
		0,
	},
	testVector{
		mustHexDecode("0100000000000000000000000000000000000000000000000000000000000000"),
		mustHexDecode("0000000000000000"),
		mustHexDecode("c5d30a7ce1ec119378c84f487d775a8542f13ece238a9455e8229e888de85bbd" +
			"29eb63d0a17a5b999b52da22be4023eb07620a54f6fa6ad8737b71eb0464dac0" +
			"10f656e6d1fd55053e50c4875c9930a33f6d0263bd14dfd6ab8c70521c19338b" +
			"2308b95cf8d0bb7d202d2102780ea3528f1cb48560f76b20f382b942500fceac"),
		0,
	},
	testVector{
		mustHexDecode("0000000000000000000000000000000000000000000000000000000000000000"),
		mustHexDecode("0100000000000000"),
		mustHexDecode("ef3fdfd6c61578fbf5cf35bd3dd33b8009631634d21e42ac33960bd138e50d32" +
			"111e4caf237ee53ca8ad6426194a88545ddc497a0b466e7d6bbdb0041b2f586b" +
			"5305e5e44aff19b235936144675efbe4409eb7e8e5f1430f5f5836aeb49bb532" +
			"8b017c4b9dc11f8a03863fa803dc71d5726b2b6b31aa32708afe5af1d6b69058"),
		0,
	},
	testVector{
		mustHexDecode("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		mustHexDecode("ffffffffffffffff"),
		mustHexDecode("d9bf3f6bce6ed0b54254557767fb57443dd4778911b606055c39cc25e674b836" +
			"3feabc57fde54f790c52c8ae43240b79d49042b777bfd6cb80e931270b7f50eb" +
			"5bac2acd86a836c5dc98c116c1217ec31d3a63a9451319f097f3b4d6dab07787" +
			"19477d24d24b403a12241d7cca064f790f1d51ccaff6b1667d4bbca1958c4306"),
		0,
	},
	testVector{
		mustHexDecode("5555555555555555555555555555555555555555555555555555555555555555"),
		mustHexDecode("5555555555555555"),
		mustHexDecode("bea9411aa453c5434a5ae8c92862f564396855a9ea6e22d6d3b50ae1b3663311" +
			"a4a3606c671d605ce16c3aece8e61ea145c59775017bee2fa6f88afc758069f7" +
			"e0b8f676e644216f4d2a3422d7fa36c6c4931aca950e9da42788e6d0b6d1cd83" +
			"8ef652e97b145b14871eae6c6804c7004db5ac2fce4c68c726d004b10fcaba86"),
		0,
	},
	testVector{
		mustHexDecode("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		mustHexDecode("aaaaaaaaaaaaaaaa"),
		mustHexDecode("9aa2a9f656efde5aa7591c5fed4b35aea2895dec7cb4543b9e9f21f5e7bcbcf3" +
			"c43c748a970888f8248393a09d43e0b7e164bc4d0b0fb240a2d72115c4808906" +
			"72184489440545d021d97ef6b693dfe5b2c132d47e6f041c9063651f96b623e6" +
			"2a11999a23b6f7c461b2153026ad5e866a2e597ed07b8401dec63a0934c6b2a9"),
		0,
	},
	testVector{
		mustHexDecode("00112233445566778899aabbccddeeffffeeddccbbaa99887766554433221100"),
		mustHexDecode("0f1e2d3c4b5a6978"),
		mustHexDecode("9fadf409c00811d00431d67efbd88fba59218d5d6708b1d685863fabbb0e961e" +
			"ea480fd6fb532bfd494b2151015057423ab60a63fe4f55f7a212e2167ccab931" +
			"fbfd29cf7bc1d279eddf25dd316bb8843d6edee0bd1ef121d12fa17cbc2c574c" +
			"ccab5e275167b08bd686f8a09df87ec3ffb35361b94ebfa13fec0e4889d18da5"),
		0,
	},
	testVector{
		mustHexDecode("c46ec1b18ce8a878725a37e780dfb7351f68ed2e194c79fbc6aebee1a667975d"),
		mustHexDecode("1ada31d5cf688221"),
		mustHexDecode("f63a89b75c2271f9368816542ba52f06ed49241792302b00b5e8f80ae9a473af" +
			"c25b218f519af0fdd406362e8d69de7f54c604a6e00f353f110f771bdca8ab92" +
			"e5fbc34e60a1d9a9db17345b0a402736853bf910b060bdf1f897b6290f01d138" +
			"ae2c4c90225ba9ea14d518f55929dea098ca7a6ccfe61227053c84e49a4a3332"),
		0,
	},
}

// stolen from https://github.com/codahale/chacha20/blob/master/chacha20_test.go
var xTestVectors = []testVector{
	testVector{
		[]byte{
			0x1b, 0x27, 0x55, 0x64, 0x73, 0xe9, 0x85, 0xd4,
			0x62, 0xcd, 0x51, 0x19, 0x7a, 0x9a, 0x46, 0xc7,
			0x60, 0x09, 0x54, 0x9e, 0xac, 0x64, 0x74, 0xf2,
			0x06, 0xc4, 0xee, 0x08, 0x44, 0xf6, 0x83, 0x89,
		},
		[]byte{
			0x69, 0x69, 0x6e, 0xe9, 0x55, 0xb6, 0x2b, 0x73,
			0xcd, 0x62, 0xbd, 0xa8, 0x75, 0xfc, 0x73, 0xd6,
			0x82, 0x19, 0xe0, 0x03, 0x6b, 0x7a, 0x0b, 0x37,
		},
		[]byte{
			0x4f, 0xeb, 0xf2, 0xfe, 0x4b, 0x35, 0x9c, 0x50,
			0x8d, 0xc5, 0xe8, 0xb5, 0x98, 0x0c, 0x88, 0xe3,
			0x89, 0x46, 0xd8, 0xf1, 0x8f, 0x31, 0x34, 0x65,
			0xc8, 0x62, 0xa0, 0x87, 0x82, 0x64, 0x82, 0x48,
			0x01, 0x8d, 0xac, 0xdc, 0xb9, 0x04, 0x17, 0x88,
			0x53, 0xa4, 0x6d, 0xca, 0x3a, 0x0e, 0xaa, 0xee,
			0x74, 0x7c, 0xba, 0x97, 0x43, 0x4e, 0xaf, 0xfa,
			0xd5, 0x8f, 0xea, 0x82, 0x22, 0x04, 0x7e, 0x0d,
			0xe6, 0xc3, 0xa6, 0x77, 0x51, 0x06, 0xe0, 0x33,
			0x1a, 0xd7, 0x14, 0xd2, 0xf2, 0x7a, 0x55, 0x64,
			0x13, 0x40, 0xa1, 0xf1, 0xdd, 0x9f, 0x94, 0x53,
			0x2e, 0x68, 0xcb, 0x24, 0x1c, 0xbd, 0xd1, 0x50,
			0x97, 0x0d, 0x14, 0xe0, 0x5c, 0x5b, 0x17, 0x31,
			0x93, 0xfb, 0x14, 0xf5, 0x1c, 0x41, 0xf3, 0x93,
			0x83, 0x5b, 0xf7, 0xf4, 0x16, 0xa7, 0xe0, 0xbb,
			0xa8, 0x1f, 0xfb, 0x8b, 0x13, 0xaf, 0x0e, 0x21,
			0x69, 0x1d, 0x7e, 0xce, 0xc9, 0x3b, 0x75, 0xe6,
			0xe4, 0x18, 0x3a,
		},
		0,
	},
}

func testChaCha20(t *testing.T, newChaCha20 func(key, nonce []byte) (cipher.Stream, error), vectors []testVector) {
	for i, vector := range vectors {
		t.Run(fmt.Sprintf("vector%d", i), func(t *testing.T) {
			c, err := newChaCha20(vector.key, vector.nonce)
			if err != nil {
				t.Fatal(err)
			}

			var block [64]byte
			for i := uint64(0); i < vector.counter; i++ {
				c.XORKeyStream(block[:], block[:])
			}

			dst := make([]byte, len(vector.keyStream))
			c.XORKeyStream(dst, dst)

			if bytes.Equal(vector.keyStream, dst) {
				return
			}

			t.Error("Bad keystream:")
			t.Errorf("\texpected %x", vector.keyStream)
			t.Errorf("\twas      %x", dst)

			for i, v := range vector.keyStream {
				if dst[i] != v {
					t.Logf("\tMismatch at offset %d: %x vs %x", i, v, dst[i])
					break
				}
			}
		})
	}
}

func testChaCha20x64(t *testing.T, newChaCha20 func(key, nonce []byte) (cipher.Stream, error), vectors []testVector) {
	if useRef {
		t.Skip("skipping: do not have x64 implementation")
	}

	oldAVX, oldAVX2 := useAVX, useAVX2
	useAVX, useAVX2 = false, false
	defer func() {
		useAVX, useAVX2 = oldAVX, oldAVX2
	}()

	testChaCha20(t, newChaCha20, vectors)
}

func testChaCha20AVX(t *testing.T, newChaCha20 func(key, nonce []byte) (cipher.Stream, error), vectors []testVector) {
	if !useAVX {
		t.Skip("skipping: do not have AVX implementation")
	}

	oldAVX, oldAVX2 := useAVX, useAVX2
	useAVX, useAVX2 = true, false
	defer func() {
		useAVX, useAVX2 = oldAVX, oldAVX2
	}()

	testChaCha20(t, newChaCha20, vectors)
}

func testChaCha20AVX2(t *testing.T, newChaCha20 func(key, nonce []byte) (cipher.Stream, error), vectors []testVector) {
	if !useAVX2 {
		t.Skip("skipping: do not have AVX2 implementation")
	}

	oldAVX, oldAVX2 := useAVX, useAVX2
	useAVX, useAVX2 = false, true
	defer func() {
		useAVX, useAVX2 = oldAVX, oldAVX2
	}()

	testChaCha20(t, newChaCha20, vectors)
}

func TestRFCChaCha20x64(t *testing.T) {
	testChaCha20x64(t, NewRFC, rfcTestVectors)
}

func TestRFCChaCha20AVX(t *testing.T) {
	testChaCha20AVX(t, NewRFC, rfcTestVectors)
}

func TestRFCChaCha20AVX2(t *testing.T) {
	testChaCha20AVX2(t, NewRFC, rfcTestVectors)
}

func TestRFCChaCha20Go(t *testing.T) {
	testChaCha20(t, ref.NewRFC, rfcTestVectors)
}

func TestDraftChaCha20x64(t *testing.T) {
	testChaCha20x64(t, NewDraft, draftTestVectors)
}

func TestDraftChaCha20AVX(t *testing.T) {
	testChaCha20AVX(t, NewDraft, draftTestVectors)
}

func TestDraftChaCha20AVX2(t *testing.T) {
	testChaCha20AVX2(t, NewDraft, draftTestVectors)
}

func TestDraftChaCha20Go(t *testing.T) {
	testChaCha20(t, ref.NewDraft, draftTestVectors)
}

func TestXChaCha20x64(t *testing.T) {
	testChaCha20x64(t, NewXChaCha, xTestVectors)
}

func TestXChaCha20AVX(t *testing.T) {
	testChaCha20AVX(t, NewXChaCha, xTestVectors)
}

func TestXChaCha20AVX2(t *testing.T) {
	testChaCha20AVX2(t, NewXChaCha, xTestVectors)
}

func TestXChaCha20Go(t *testing.T) {
	testChaCha20(t, ref.NewXChaCha, xTestVectors)
}

func testBadSize(t *testing.T, newChaCha20 func(key, nonce []byte) (cipher.Stream, error), keysize, nonceSize int, expect error) {
	key := make([]byte, keysize)
	nonce := make([]byte, nonceSize)

	_, err := newChaCha20(key, nonce)

	if err != expect {
		t.Errorf("expected error %v, got %v", expect, err)
	}
}

func TestRFCRightSizes(t *testing.T) {
	testBadSize(t, NewRFC, KeySize, RFCNonceSize, nil)
}

func TestDraftRightSizes(t *testing.T) {
	testBadSize(t, NewDraft, KeySize, DraftNonceSize, nil)
}

func TestXRightSizes(t *testing.T) {
	testBadSize(t, NewXChaCha, KeySize, XNonceSize, nil)
}

func TestRFCBadKeySize(t *testing.T) {
	testBadSize(t, NewRFC, 3, DraftNonceSize, ErrInvalidKey)
}

func TestDraftBadKeySize(t *testing.T) {
	testBadSize(t, NewDraft, 3, DraftNonceSize, ErrInvalidKey)
}

func TestXBadKeySize(t *testing.T) {
	testBadSize(t, NewXChaCha, 3, XNonceSize, ErrInvalidKey)
}

func TestRFCBadNonceSize(t *testing.T) {
	testBadSize(t, NewRFC, KeySize, 3, ErrInvalidNonce)
}

func TestDraftBadNonceSize(t *testing.T) {
	testBadSize(t, NewDraft, KeySize, 3, ErrInvalidNonce)
}

func TestXBadNonceSize(t *testing.T) {
	testBadSize(t, NewXChaCha, KeySize, 3, ErrInvalidNonce)
}

func TestNewBadNonceSize(t *testing.T) {
	testBadSize(t, New, KeySize, 3, ErrInvalidNonce)
}

func testEqual(t *testing.T, new1, new2 func(key, nonce []byte) (cipher.Stream, error), noncesize, calls int, label1, label2 string) {
	t.Parallel()

	if err := quick.Check(func(key, nonce, src []byte) bool {
		c1, err := new1(key, nonce)
		if err != nil {
			t.Error(err)
			return false
		}

		c2, err := new2(key, nonce)
		if err != nil {
			t.Error(err)
			return false
		}

		dst1 := make([]byte, len(src))
		dst2 := make([]byte, len(src))

		for i := 0; i < calls; i++ {
			c1.XORKeyStream(dst1, src)
			c2.XORKeyStream(dst2, src)
		}

		if bytes.Equal(dst1, dst2) {
			return true
		}

		t.Error("Bad output:")
		t.Errorf("\t%s: %x", label2, dst2)
		t.Errorf("\t%s: %x", label1, dst1)

		for i, v := range dst2 {
			if dst1[i] != v {
				t.Logf("\tMismatch at offset %d: %x vs %x", i, v, dst1[i])
				break
			}
		}

		return false
	}, &quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			key := make([]byte, KeySize)
			rand.Read(key)
			args[0] = reflect.ValueOf(key)

			nonce := make([]byte, noncesize)
			rand.Read(nonce)
			args[1] = reflect.ValueOf(nonce)

			src := make([]byte, 1+rand.Intn(1024*1024))
			rand.Read(src)
			args[2] = reflect.ValueOf(src)
		},

		MaxCountScale: 0.5,
	}); err != nil {
		t.Error(err)
	}
}

func TestRFCEqualOneShot(t *testing.T) {
	testEqual(t, NewRFC, ref.NewRFC, RFCNonceSize, 1, "tmthrgd/chacha20", "tmthrgd/chacha20/internal/ref")
}

func TestRFCEqualMultiUse(t *testing.T) {
	testEqual(t, NewRFC, ref.NewRFC, RFCNonceSize, 5, "tmthrgd/chacha20", "tmthrgd/chacha20/internal/ref")
}

func TestDraftEqualOneShot(t *testing.T) {
	testEqual(t, NewDraft, codahale.New, DraftNonceSize, 1, "tmthrgd/chacha20", "codahale/chacha20")
}

func TestDraftEqualMultiUse(t *testing.T) {
	testEqual(t, NewDraft, codahale.New, DraftNonceSize, 5, "tmthrgd/chacha20", "codahale/chacha20")
}

func TestDraftEqualOneShotGo(t *testing.T) {
	testEqual(t, ref.NewDraft, codahale.New, DraftNonceSize, 1, "tmthrgd/chacha20/internal/ref", "codahale/chacha20")
}

func TestDraftEqualMultiUseGo(t *testing.T) {
	testEqual(t, ref.NewDraft, codahale.New, DraftNonceSize, 5, "tmthrgd/chacha20/internal/ref", "codahale/chacha20")
}

func TestXEqualOneShot(t *testing.T) {
	testEqual(t, NewXChaCha, codahale.NewXChaCha, XNonceSize, 1, "tmthrgd/chacha20", "codahale/chacha20")
}

func TestXEqualMultiUse(t *testing.T) {
	testEqual(t, NewXChaCha, codahale.NewXChaCha, XNonceSize, 5, "tmthrgd/chacha20", "codahale/chacha20")
}

func TestXEqualOneShotGo(t *testing.T) {
	testEqual(t, ref.NewXChaCha, codahale.NewXChaCha, XNonceSize, 1, "tmthrgd/chacha20/internal/ref", "codahale/chacha20")
}

func TestXEqualMultiUseGo(t *testing.T) {
	testEqual(t, ref.NewXChaCha, codahale.NewXChaCha, XNonceSize, 5, "tmthrgd/chacha20/internal/ref", "codahale/chacha20")
}

func testNewNewVar(t *testing.T, newVariant func(key, nonce []byte) (cipher.Stream, error), nonceSize int) {
	var key [KeySize]byte
	nonce := make([]byte, nonceSize)

	c1, err := New(key[:], nonce)
	if err != nil {
		t.Fatal(err)
	}

	c2, err := newVariant(key[:], nonce)
	if err != nil {
		t.Fatal(err)
	}

	var block1 [64]byte
	c1.XORKeyStream(block1[:], block1[:])

	var block2 [64]byte
	c2.XORKeyStream(block2[:], block2[:])

	if !bytes.Equal(block1[:], block2[:]) {
		t.Error("New returned incorrect cipher")
	}
}

func TestNewNewRFC(t *testing.T) {
	testNewNewVar(t, NewRFC, RFCNonceSize)
}

func TestNewNewDraft(t *testing.T) {
	testNewNewVar(t, NewDraft, DraftNonceSize)
}

func TestNewNewXChaCha(t *testing.T) {
	testNewNewVar(t, NewXChaCha, XNonceSize)
}

func TestXOREmptyKeyStream(t *testing.T) {
	var key [KeySize]byte
	var nonce [RFCNonceSize]byte

	c, err := NewRFC(key[:], nonce[:])
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Errorf("XORKeyStream caused panic on empty src")
		}
	}()

	var out [0]byte
	c.XORKeyStream(out[:], out[:0])
}

func TestXORNoExtraKeyStream(t *testing.T) {
	var key [KeySize]byte
	var nonce [RFCNonceSize]byte

	c, err := NewRFC(key[:], nonce[:])
	if err != nil {
		t.Fatal(err)
	}

	var out, zero [167]byte
	c.XORKeyStream(out[:], out[:153])

	if bytes.Equal(out[:153], zero[:153]) {
		t.Error("XORKeyStream did not update partial block")
	}

	if !bytes.Equal(out[153:], zero[153:]) {
		t.Error("XORKeyStream updated past len(src)")
	}
}

func TestXORKeyStreamBufferEmpty(t *testing.T) {
	var key [KeySize]byte
	var nonce [RFCNonceSize]byte

	c, err := NewRFC(key[:], nonce[:])
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Errorf("XORKeyStream caused panic on empty after buffer fill")
		}
	}()

	var out [256]byte
	c.XORKeyStream(out[:], out[:224])
	c.XORKeyStream(out[224:], out[224:])
}

func ExampleNewRFC() {
	key, err := hex.DecodeString("60143a3d7c7137c3622d490e7dbb85859138d198d9c648960e186412a6250722")
	if err != nil {
		panic(err)
	}

	// A nonce should only be used once. Generate it randomly.
	nonce, err := hex.DecodeString("00000000308c92676fa95973")
	if err != nil {
		panic(err)
	}

	c, err := NewRFC(key, nonce)
	if err != nil {
		panic(err)
	}

	src := []byte("hello I am a secret message")
	dst := make([]byte, len(src))

	c.XORKeyStream(dst, src)

	fmt.Printf("%x\n", dst)
	// Output: a05452ebd981422dcdab2c9cde0d20a03f769e87d3e976ee6d6a11
}

func ExampleNewDraft() {
	key, err := hex.DecodeString("60143a3d7c7137c3622d490e7dbb85859138d198d9c648960e186412a6250722")
	if err != nil {
		panic(err)
	}

	// A nonce should only be used once. Generate it randomly.
	nonce, err := hex.DecodeString("308c92676fa95973")
	if err != nil {
		panic(err)
	}

	c, err := NewDraft(key, nonce)
	if err != nil {
		panic(err)
	}

	src := []byte("hello I am a secret message")
	dst := make([]byte, len(src))

	c.XORKeyStream(dst, src)

	fmt.Printf("%x\n", dst)
	// Output: a05452ebd981422dcdab2c9cde0d20a03f769e87d3e976ee6d6a11
}
