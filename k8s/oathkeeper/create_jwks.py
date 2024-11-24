import json
import base64
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.backends import default_backend

# Read the private key
with open('private.key', 'rb') as f:
    private_key = serialization.load_pem_private_key(
        f.read(),
        password=None,
        backend=default_backend()
    )

# Get the private and public numbers
private_numbers = private_key.private_numbers()
public_numbers = private_key.public_key().public_numbers()

def int_to_base64url(value):
    """Convert an integer to a base64url-encoded string."""
    value_hex = format(value, 'x')
    if len(value_hex) % 2 == 1:
        value_hex = '0' + value_hex
    value_bytes = bytes.fromhex(value_hex)
    return base64.urlsafe_b64encode(value_bytes).decode('utf-8').rstrip('=')

# Create the JWK
jwk = {
    "keys": [{
        "kty": "RSA",
        "kid": "oathkeeper-key",
        "use": "sig",
        "alg": "RS256",
        # Public key components
        "n": int_to_base64url(public_numbers.n),
        "e": int_to_base64url(public_numbers.e),
        # Private key components
        "d": int_to_base64url(private_numbers.d),
        "p": int_to_base64url(private_numbers.p),
        "q": int_to_base64url(private_numbers.q),
        "dp": int_to_base64url(private_numbers.dmp1),
        "dq": int_to_base64url(private_numbers.dmq1),
        "qi": int_to_base64url(private_numbers.iqmp)
    }]
}

# Write the JWKS
with open('jwks.json', 'w') as f:
    json.dump(jwk, f, indent=2)
