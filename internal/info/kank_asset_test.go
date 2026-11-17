package info

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/trustwallet/go-primitives/coin"
)

const kankAddr = "0x10BA136F3d56c30D271D77f09d7dc251992b0F14"

func kankInfoPath() string {
	return filepath.Join("..", "..", "blockchains", "smartchain", "assets", kankAddr, "info.json")
}

// TestKANKAssetValid runs the same validator `assets check` uses.
func TestKANKAssetValid(t *testing.T) {
	raw, err := os.ReadFile(kankInfoPath())
	if err != nil {
		t.Fatalf("read info.json: %v", err)
	}

	var a AssetModel
	if err := json.Unmarshal(raw, &a); err != nil {
		t.Fatalf("unmarshal info.json: %v", err)
	}

	if err := ValidateAsset(a, coin.Coins[coin.SMARTCHAIN], kankAddr); err != nil {
		t.Fatalf("ValidateAsset: %v", err)
	}

	if *a.Symbol != "KANK" || *a.Decimals != 18 || *a.Type != "BEP20" {
		t.Errorf("unexpected token fields: symbol=%s decimals=%d type=%s", *a.Symbol, *a.Decimals, *a.Type)
	}
}

// TestKANKAssetCanonicalFormat guards against the original bug: a file that
// passes `check` but is rewritten by `fix` (wrong indent, key order, empty fields).
func TestKANKAssetCanonicalFormat(t *testing.T) {
	raw, err := os.ReadFile(kankInfoPath())
	if err != nil {
		t.Fatalf("read info.json: %v", err)
	}

	var a AssetModel
	if err := json.Unmarshal(raw, &a); err != nil {
		t.Fatalf("unmarshal info.json: %v", err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	if err := enc.Encode(a); err != nil {
		t.Fatalf("encode: %v", err)
	}

	if !bytes.Equal(buf.Bytes(), raw) {
		t.Errorf("info.json is not canonically formatted; run `make fix`.\nwant:\n%s\ngot:\n%s", buf.String(), raw)
	}
}
