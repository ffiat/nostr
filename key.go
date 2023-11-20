package nostr

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/bech32"
)

const Prefix = "nostr:"

const (
	UriPub     = "npub"
	UriProfile = "nprofile"
	UriEvent = "nevent"
)

// NIP-01
func GeneratePrivateKey() string {
	params := btcec.S256().Params()
	one := new(big.Int).SetInt64(1)

	b := make([]byte, params.BitSize/8+8)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	k := new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)

	return hex.EncodeToString(k.Bytes())
}

// NIP-01
func GetPublicKey(sk string) (string, error) {
	b, err := hex.DecodeString(sk)
	if err != nil {
		return "", err
	}

	_, pk := btcec.PrivKeyFromBytes(b)
	return hex.EncodeToString(schnorr.SerializePubKey(pk)), nil
}

// NIP-01
func IsValidPublicKeyHex(pk string) bool {
	if strings.ToLower(pk) != pk {
		return false
	}
	dec, _ := hex.DecodeString(pk)
	return len(dec) == 32
}

// NIP-19
func DecodeBech32(bech32string string) (string, string, error) {

	prefix, bits5, err := bech32.DecodeNoLimit(bech32string)
	if err != nil {
		return "", "", err
	}

	data, err := bech32.ConvertBits(bits5, 5, 8, false)
	if err != nil {
		return prefix, "", fmt.Errorf("failed translating data into 8 bits: %s", err.Error())
	}

	switch prefix {
	case "npub", "nsec", "note":
		if len(data) < 32 {
			return prefix, "", fmt.Errorf("data is less than 32 bytes (%d)", len(data))
		}
		return prefix, hex.EncodeToString(data[0:32]), nil
	case "nprofile":
		log.Fatalln("nostr: nprofile decoding not implemented")
	case "nevent":
		log.Fatalln("nostr: nevent decoding not implemented")
	case "naddr":
		log.Fatalln("nostr: naddr decoding not implemented")
	}

	return prefix, string(data), fmt.Errorf("unknown tag %s", prefix)
}

// NIP-19
func EncodePrivateKey(privateKeyHex string) (string, error) {
	b, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key hex: %w", err)
	}

	bits5, err := bech32.ConvertBits(b, 8, 5, true)
	if err != nil {
		return "", err
	}

	return bech32.Encode("nsec", bits5)
}

// NIP-19
func EncodePublicKey(publicKeyHex string) (string, error) {
	b, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key hex: %w", err)
	}

	bits5, err := bech32.ConvertBits(b, 8, 5, true)
	if err != nil {
		return "", err
	}

	return bech32.Encode("npub", bits5)
}

func EncodeNote(eventIDHex string) (string, error) {

	b, err := hex.DecodeString(eventIDHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode event id hex: %w", err)
	}

	bits5, err := bech32.ConvertBits(b, 8, 5, true)
	if err != nil {
		return "", err
	}

	return bech32.Encode("note", bits5)
}
