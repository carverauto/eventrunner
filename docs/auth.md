# auth

## Design

We will use Ory Kratos for identity management and Ory Hydra for OAuth2.0. 
A SuperUser is the first user created and has the ability to create other tenants and users.

## Identity

When we create an identity for the user, we will first create a tenant and then create the user. 
The tenant will have a unique ID that will be used to create the user. The user will have a unique ID 
that will be used to create the identity. We will set that tenantID and customerID in the user's traits,
and retrieve that in our app middleware for authorization.

Setting up Ory

### JWKS

Generate keys

```shell
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
```

Use the `jwtKeys` go tool to generate a JWKS file. Supply this JSON to Ory when configuring the identity provider.
Get the address for the `jwks.json` file and set it in the .env file or k8s manifest.

### Custom Identity Schema

```json
{
  "$id": "https://schemas.ory.sh/presets/kratos/identity.email.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Person",
  "type": "object",
  "properties": {
    "traits": {
      "type": "object",
      "properties": {
        "email": {
          "type": "string",
          "format": "email",
          "title": "E-Mail",
          "ory.sh/kratos": {
            "credentials": {
              "password": {
                "identifier": true
              },
              "webauthn": {
                "identifier": true
              },
              "totp": {
                "account_name": true
              },
              "code": {
                "identifier": true,
                "via": "email"
              },
              "passkey": {
                "display_name": true
              }
            },
            "recovery": {
              "via": "email"
            },
            "verification": {
              "via": "email"
            }
          },
          "maxLength": 320
        },
        "tenant_id": {
          "type": "string",
          "title": "Tenant ID",
          "description": "The ID of the tenant this user belongs to",
          "pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
        },
        "customer_id": {
          "type": "string",
          "title": "Customer ID",
          "description": "The ID of the customer this user is associated with",
          "pattern": "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
        },
        "roles": {
          "type": "array",
          "items": {
            "type": "string",
            "enum": ["admin", "user", "superuser"]
          },
          "title": "Roles",
          "description": "The roles assigned to this user"
        }
      },
      "required": [
        "email",
        "tenant_id",
        "roles"
      ],
      "additionalProperties": false
    }
  }
}
```
