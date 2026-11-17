---
title: info.json can pass `check` but still be non-canonical
date: 2026-11-17
area: asset-metadata
files:
  - blockchains/smartchain/assets/0x10BA136F3d56c30D271D77f09d7dc251992b0F14/info.json
  - internal/info/model.go
  - internal/info/kank_asset_test.go
symptom: A hand-written info.json passes `make check` but `make fix` rewrites it (indent, key order, empty fields dropped).
tags: [info-json, formatting, check, fix, hand-edited, bep20]
summary: `check` validates content, not formatting; hand-written info.json must use 4-space indent, AssetModel key order, and no empty research/tags/links.
---

## Symptom
The KANK token (`smartchain/0x10BA…0F14`) was added by hand with 2-space indent,
an arbitrary key order and empty `"research": ""`, `"tags": []`, `"links": []`.
`bin/assets check` reported 0 errors, yet `bin/assets fix` rewrote the file.

## Root cause
`check` only runs the content validators (`internal/info/ValidateAsset`: required keys,
id checksum, decimals, description, links). It does not compare the file to its
canonical serialization. `fix` re-encodes through `AssetModel` with 4-space indent;
field order comes from the struct, and `omitempty` drops empty strings/slices.

## Rule
- Prefer `make add-token` over hand-writing info.json.
- If hand-writing, match `AssetModel` order (name, symbol, type, decimals,
  description, website, explorer, status, id, links, tags), 4-space indent, and omit empty fields.
- Don't run `bin/assets fix` repo-wide to normalize one file: it rewrites thousands
  of unrelated assets. Run it, keep only the target file's diff, and revert the rest.

## Regression guard
`internal/info/kank_asset_test.go` validates the KANK asset and asserts its file equals
the canonical `AssetModel` encoding.
