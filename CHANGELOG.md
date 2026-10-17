# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-10-17

First public release.

### Added

- `idcard`: Thai national ID validation (`Validate`, `Check`), formatting
  (`Format`), normalization (`Normalize`) and test-data generation
  (`Generate`).
- `taxid`: Thai tax ID validation, formatting, normalization and juristic-person
  detection (`IsJuristic`).
- `phone`: validation, normalization, E.164 conversion and display formatting
  for Thai mobile and landline numbers, plus `IsMobile` / `IsLandline`.
- `promptpay`: EMVCo-compliant PromptPay QR payload generation for phone,
  national ID / tax ID and e-wallet targets, with optional amount; exported
  `CRC16` and `VerifyChecksum` helpers. Output matches the `promptpay-qr`
  npm package byte for byte.
- `thaidate`: Buddhist Era conversion (`ToBE`, `ToCE`), Thai date formatting
  with standard Go layouts (`Format` and the `Date*` layout constants), Thai
  month and weekday names, DD/MM/YYYY (BE) parsing (`Parse`,
  `ParseInLocation`) and Arabic ↔ Thai digit conversion.
- `bahttext`: amounts in Thai words (`Format`, `FormatString`), following the
  same rules as Excel's `BAHTTEXT`.
- Exported sentinel errors in every package for use with `errors.Is`.
- Runnable programs under `examples/`, CI on Linux, macOS and Windows with
  Go 1.22 and stable, golangci-lint configuration and a Makefile.

[Unreleased]: https://github.com/tonklasocute/thai-utils/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/tonklasocute/thai-utils/releases/tag/v0.1.0
