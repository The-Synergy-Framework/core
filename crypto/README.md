# Crypto - Cryptographic Utilities

Cryptographic utilities for the Synergy Framework.

## Overview

The crypto package provides common cryptographic operations used across the Synergy Framework. It focuses on providing safe, easy-to-use utilities for common crypto tasks.

## Features

### RSA Key Parsing
- `ParseRSAPrivateKey(keyData string)` - Parse PEM-encoded RSA private keys
- `ParseRSAPublicKey(keyData string)` - Parse PEM-encoded RSA public keys

Supports both PKCS1 and PKCS8 formats for private keys.

## Example

```go
package main

import (
    "core/crypto"
    "fmt"
)

func main() {
    // Parse RSA private key
    privateKeyPEM := `-----BEGIN RSA PRIVATE KEY-----
    MIIEpAIBAAKCAQEA...
    -----END RSA PRIVATE KEY-----`
    
    privateKey, err := crypto.ParseRSAPrivateKey(privateKeyPEM)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Private key size: %d bits\n", privateKey.Size()*8)
}
```

## Design Principles

- **Safe defaults**: All utilities use secure defaults
- **Error handling**: Clear error messages for debugging
- **Format support**: Support multiple standard formats where applicable
- **Zero dependencies**: Only uses Go standard library 