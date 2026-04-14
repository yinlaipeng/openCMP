# Go Module Initialization Instructions

## Issue
After adding the LDAP library dependency to the go.mod file, the go.sum file needs to be updated to include the checksums for the new dependencies. This is why you're seeing the error:

```
missing go.sum entry for module providing package github.com/go-ldap/ldap/v3
```

## Solution
The go.sum file needs to be regenerated to include the new LDAP library and its dependencies. This is typically done with the `go mod tidy` command.

## Steps to Resolve

1. **Ensure Go is installed** on your system
2. **Navigate to the backend directory:**
   ```bash
   cd /Users/aurora/Desktop/xtwork/git/openCMP/backend
   ```

3. **Run the following command to update module dependencies:**
   ```bash
   go mod tidy
   ```

4. **Then you can run the server:**
   ```bash
   go run cmd/server/main.go
   ```

## What go mod tidy does:
- Downloads the new dependencies (including github.com/go-ldap/ldap/v3)
- Calculates and records the cryptographic checksums in go.sum
- Cleans up unused dependencies
- Ensures the module graph is consistent

## Dependencies Added:
- `github.com/go-ldap/ldap/v3 v3.4.4` - Primary LDAP library for OpenLDAP/AD integration
- `github.com/go-asn1-ber/asn1-ber v1.5.5` - ASN.1 BER encoding library (used by LDAP library)

## Verification
After running `go mod tidy`, you can verify the module integrity with:
```bash
go mod verify
```

## Expected Result
Once the go.sum file is properly updated, the server should start without errors and the OpenLDAP authentication functionality will be available.