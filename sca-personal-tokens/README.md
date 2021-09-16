## Signing SCA requests when using personal tokens
Wise API users in the UK / Europe / EEA need to satisfy a Strong Customer Authentication (SCA) challenge from Wise when they are moving money or reading account statements. This is an additional security check to comply with the 2nd Payment Services Directive (PSD2) regulation. 

How it works is that a standard API call to pull an account statement will fail with HTTP 403 and an 'x-2fa-approval' header with a one-time-token value. The user will need to sign the one-time-token with a private key that corresponds to a public key uploaded to their Wise account. They can then retry the same request including the signed header and the request can succeed.

### Usage
1. Read the [Wise API documentation](https://api-docs.wise.com/#strong-customer-authentication-personal-token), generate your keypair and upload your public key.
2. Add your profile ID, borderless account ID, and private key file path to the script.
3. Add your API token as an env var with key: API_TOKEN
```bash
$ export API_TOKEN=<YOUR API TOKEN HERE>
```
4. Run the script (go, python, java)

### Go implementation

Source code: [get-statements-sca.go](https://github.com/transferwise/digital-signatures-examples/blob/main/sca-personal-tokens/get-statements-sca.go)

```bash
$ go run get-statements-sca.go
```

### Python implementation

Source code: [get-statements-sca.py](https://github.com/transferwise/digital-signatures-examples/blob/main/sca-personal-tokens/get-statements-sca.py)

```bash
$ python3 get-statements-sca.py
```

### Java implementation

Source code: https://github.com/transferwise/digital-signatures

### API Documentation
- [Wise API: Get Account Statements](https://api-docs.transferwise.com/#borderless-accounts-get-account-statement)
- [Wise API: Strong Customer Authentication](https://api-docs.transferwise.com/#strong-customer-authentication)
