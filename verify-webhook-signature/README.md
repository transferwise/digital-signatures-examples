## Webhook signature verification

Each outgoing webhook request is signed. You should verify that any request you handle was sent by Wise and has not been forged or tampered with. You should not process any requests with signatures that fail verification.

Signatures are generated using an RSA key and SHA256 digest of the message body. They are transmitted using the X-Signature-SHA256 request header and are Base64 encoded.

The code examples use the sandbox public key.

### Java implementation

Source code: [verify-signature.java](https://github.com/transferwise/digital-signatures-examples/blob/main/verify-webhook-signature/verify-signature.java)

### Ruby implementation

Source code: [verify-signature.rb](https://github.com/transferwise/digital-signatures-examples/blob/main/verify-webhook-signature/verify-signature.rb)

```bash
$ ruby verify-signature.rb
```

### Node implementation

Source code: [verify-signature.js](https://github.com/transferwise/digital-signatures-examples/blob/main/verify-webhook-signature/verify-signature.js)

```bash
$ node verify-signature.js
```

### PHP implementation

Source code: [verify-signature.php](https://github.com/transferwise/digital-signatures-examples/blob/main/verify-webhook-signature/verify-signature.php)

```bash
$ php verify-signature.php
```

### API Documentation
- [Wise API: Webhooks](https://api-docs.wise.com/#webhook-events-event-http-requests)
