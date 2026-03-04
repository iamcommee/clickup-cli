# Fix Hostname-Based Encryption Key Instability

**Status: Completed**

## Overview
macOS `os.Hostname()` returns varying values (with/without `.local` suffix) depending on network state. This changes the AES-GCM encryption key used to encrypt/decrypt the API token, making the token undecryptable when the hostname changes.

## Tasks
- [x] Write tests for hostname normalization
- [x] Fix `getEncryptionKey()` to normalize hostname
- [x] Add migration: try old key if new key fails decryption
- [x] Write migration tests
- [x] Run all tests
- [x] Rebuild binary

## Progress Log

| DateTime | Task | Status | Notes |
|----------|------|--------|-------|
| 2026-03-04 11:47 | Started investigation | Complete | Confirmed hostname `.local` suffix causes key mismatch |
| 2026-03-04 11:48 | Write tests (TDD) | Complete | 3 new tests: normalizeHostname, DecryptWithFallback, DecryptWithFallback_OldKey |
| 2026-03-04 11:49 | Implement fix | Complete | Normalized hostname, added fallback decryption, auto re-encryption |
| 2026-03-04 11:50 | All tests passing | Complete | 19/19 tests pass |
| 2026-03-04 11:50 | Build & verify | Complete | Binary rebuilt, config loads correctly |
