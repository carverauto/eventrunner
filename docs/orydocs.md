---
id: api-access-rules
title: API access rules
---

Ory Oathkeeper reaches decisions to allow or deny access by applying Access Rules. Access Rules can be stored on the file system,
set as an environment variable, or fetched from HTTP(s) remotes. These repositories can be configured in the configuration file
(`oathkeeper -c ./path/to/config.yml ...`)

```yaml
# Configures Access Rules
access_rules:
  # Locations (list of URLs) where access rules should be fetched from on boot.
  # It's expected that the documents at those locations return a JSON or YAML Array containing Ory Oathkeeper Access Rules.
  repositories:
    # If the URL Scheme is `file://`, the access rules (an array of access rules is expected) will be
    # fetched from the local file system.
    - file://path/to/rules.json
    # If the URL Scheme is `inline://`, the access rules (an array of access rules is expected)
    # are expected to be a base64 encoded (with padding!) JSON/YAML string (base64_encode(`[{"id":"foo-rule","authenticators":[....]}]`)):
    - inline://W3siaWQiOiJmb28tcnVsZSIsImF1dGhlbnRpY2F0b3JzIjpbXX1d
    # If the URL Scheme is `http://` or `https://`, the access rules (an array of access rules is expected) will be
    # fetched from the provided HTTP(s) location.
    - https://path-to-my-rules/rules.json
    # If the URL Scheme is `s3://`, `gs://` or `azblob://`, the access rules (an array of access rules is expected)
    # will be fetched by an object storage (AWS S3, Google Cloud Storage, Azure Blob Storage).
    #
    # S3 storage also supports S3-compatible endpoints served by Minio or Ceph.
    # See aws.ConfigFromURLParams (https://godoc.org/gocloud.dev/aws#ConfigFromURLParams) for more details on supported URL options for S3.
    - s3://my-bucket-name/rules.json
    - s3://my-bucket-name/rules.json?endpoint=minio.my-server.net
    - gs://gcp-bucket-name/rules.json
    - azblob://my-blob-container/rules.json
  # Determines a matching strategy for the access rules . Supported values are `glob` and `regexp`. Empty string defaults to regexp.
  matching_strategy: glob
```

or by setting the equivalent environment variable:

```sh
export ACCESS_RULES_REPOSITORIES='file://path/to/rules.json,https://path-to-my-rules/rules.json,inline://W3siaWQiOiJmb28tcnVsZSIsImF1dGhlbnRpY2F0b3JzIjpbXX1d'
```

The repository (file, inline, remote) must be formatted either as a JSON or a YAML array containing the access rules:

```sh
cat ./rules.json
[{
    "id": "my-first-rule"
},{
    "id": "my-second-rule"
}]

cat ./rules.yaml
- id: my-first-rule
  version: v0.36.0-beta.4
  authenticators:
    - handler: noop
- id: my-second-rule
  version: v0.36.0-beta.4
  authorizer:
    handler: allow
```

## Access rule format

Access Rules have four principal keys:

- `id` (string): The unique ID of the Access Rule.
- `version` (string): The version of Ory Oathkeeper uses [Semantic Versioning](https://semver.org). Please use
  `vMAJOR.MINOR.PATCH` notation format. Ory Oathkeeper can migrate access rules across versions. If left empty Ory Oathkeeper will
  assume that the rule uses the same tag as the running version. Examples: `v0.1.3` or `v1.2.3`
- `upstream` (object): The location of the server where requests matching this rule should be forwarded to. This only needs to be
  set when using the Ory Oathkeeper Proxy as the Decision API doesn't forward the request to the upstream.
  - `url` (string): The URL the request will be forwarded to.
  - `preserve_host` (bool): If set to `false` (default), the forwarded request will include the host and port of the `url` value.
    If `true`, the host and port of the Ory Oathkeeper Proxy will be used instead:
    - `false`: Incoming HTTP Header `Host: mydomain.com`-> Forwarding HTTP Header `Host: someservice.intranet.mydomain.com:1234`
  - `strip_path` (string): If set, replaces the provided path prefix when forwarding the requested URL to the upstream URL:
    - set to `/api/v1`: Incoming HTTP Request at `/api/v1/users` -> Forwarding HTTP Request at `/users`.
    - unset: Incoming HTTP Request at `/api/v1/users` -> Forwarding HTTP Request at `/api/v1/users`.
- `match` (object): Defines the URL(s) this Access Rule should match.

  - `methods` (string[]): Array of HTTP methods (for example GET, POST, PUT, DELETE, ...).
  - `url` (string): The URL that should be matched. You can use regular expressions or glob patterns in this field to match more
    than one url. The matching strategy (glob or regexp) is defined in the global configuration file as
    `access_rules.matching_strategy`. This matcher ignores query parameters. Regular expressions (or glob patterns) are
    encapsulated in brackets `<` and `>`.

    Regular expressions examples:

    - `https://mydomain.com/` matches `https://mydomain.com/` and doesn't match `https://mydomain.com/foo` or
      `https://mydomain.com`.
    - `<https|http>://mydomain.com/<.*>` matches:`https://mydomain.com/` or `http://mydomain.com/foo`. Doesn't match:
      `https://other-domain.com/` or `https://mydomain.com`.
    - `http://mydomain.com/<[[:digit:]]+>` matches `http://mydomain.com/123` and doesn't match `http://mydomain/abc`.
    - `http://mydomain.com/<(?!protected).*>` matches `http://mydomain.com/resource` and doesn't match
      `http://mydomain.com/protected`

    [Glob](http://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm) patterns examples:

    - `https://mydomain.com/<m?n>` matches `https://mydomain.com/man` and does not match `http://mydomain.com/foo`.
    - `https://mydomain.com/<{foo*,bar*}>` matches `https://mydomain.com/foo` or `https://mydomain.com/bar` and doesn't match
      `https://mydomain.com/any`.

- `authenticators`: A list of authentication handlers that authenticate the provided credentials. Authenticators are checked
  iteratively from index `0` to `n` and the first authenticator to return a positive result will be the one used. If you want the
  rule to first check a specific authenticator before "falling back" to others, have that authenticator as the first item in the
  array. For the full list of available authenticators, click [here](pipeline/authn.md).
- `authorizer`: The authorization handler which will try to authorize the subject ("user") from the previously validated
  credentials making the request. For example, you could check if the subject ("user") is part of the "admin" group or if he/she
  has permission to perform that action. For the full list of available authorizers, click [here](pipeline/authz.md).
- `mutators`: A list of mutation handlers that transform the HTTP request before forwarding it. A common use case is generating a
  new set of credentials (for example JWT) which then will be forwarded to the upstream server. When using Ory Oathkeeper's
  Decision API, it's expected that the API Gateway forwards the mutated HTTP Headers to the upstream server. For the full list of
  available mutators, click [here](pipeline/mutator.md).
- `errors`: A list of error handlers that are executed when any of the previous handlers (for example authentication) fail. Error
  handlers define what to do in case of an error, for example redirect the user to the login endpoint when a unauthorized (HTTP
  Status Code 401) error occurs. If left unspecified, errors will always be handled as JSON responses unless the global
  configuration key `errors.fallback` was changed. For more information on error handlers, click [here](pipeline/error.md).

#### Examples

Rule in JSON format:

```json
{
  "id": "some-id",
  "version": "v0.36.0-beta.4",
  "upstream": {
    "url": "http://my-backend-service",
    "preserve_host": true,
    "strip_path": "/api/v1"
  },
  "match": {
    "url": "http://my-app/some-route/<.*>",
    "methods": ["GET", "POST"]
  },
  "authenticators": [{ "handler": "noop" }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }],
  "errors": [{ "handler": "json" }]
}
```

Rule in YAML format:

```yaml
id: some-id
version: v0.36.0-beta.4
upstream:
  url: http://my-backend-service
  preserve_host: true
  strip_path: /api/v1
match:
  url: http://my-app/some-route/<.*>
  methods:
    - GET
    - POST
authenticators:
  - handler: noop
authorizer:
  handler: allow
mutators:
  - handler: noop
errors:
  - handler: json
```

## Handler configuration

Handlers (Authenticators, Mutators, Authorizers, Errors) sometimes require configuration. The configuration can be defined
globally as well as per Access Rule. The configuration from the Access Rules overrides values from the global configuration.

**oathkeeper.yml**

```yaml
authenticators:
  anonymous:
    enabled: true
    config:
      subject: anon
```

**rule.json**

```json
{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service",
    "preserve_host": true,
    "strip_path": "/api/v1"
  },
  "match": {
    "url": "http://my-app/some-route/<.*>",
    "methods": ["GET", "POST"]
  },
  "authenticators": [{ "handler": "anonymous", "config": { "subject": "anon" } }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }]
}
```

## Scoped credentials

Some credentials are scoped. For example, OAuth 2.0 Access Tokens usually are scoped ("OAuth 2.0 Scope"). Scope validation depends
on the meaning of the scope. Therefore, wherever Ory Oathkeeper validates a scope, these scope strategies are supported:

- `hierarchic`: Scope `foo` matches `foo`, `foo.bar`, `foo.baz` but not `bar`
- `wildcard`: Scope `foo.*` matches `foo`, `foo.bar`, `foo.baz` but not `bar`. Scope `foo` matches `foo` but not `foo.bar` nor
  `bar`
- `exact`: Scope `foo` matches `foo` but not `bar` nor `foo.bar`
- `none`: Scope validation is disabled. If however a scope is configured to be validated, the request will fail with an error
  message.

## Match strategy behavior

With the **Regular expression** strategy, you can use the extracted groups in all handlers where the substitutions are supported
by using the Go [`text/template`](https://golang.org/pkg/text/template/) package, receiving the
[AuthenticationSession](https://github.com/ory/oathkeeper/blob/master/pipeline/authn/authenticator.go#L39) struct:

```go
type AuthenticationSession struct {
  Subject      string
  Extra        map[string]interface{}
  Header       http.Header
  MatchContext MatchContext
}

type MatchContext struct {
  RegexpCaptureGroups []string
  URL                 *url.URL
}
```

If the match URL is `<https|http>://mydomain.com/<.*>` and the request is `http://mydomain.com/foo`, the `MatchContext` field will
contain

- `RegexpCaptureGroups`: ["http", "foo"]
- `URL`: "http://mydomain.com/foo"
---
id: pipeline
title: Access rule pipeline
---

Read more about the [principal components and execution pipeline of access rules](api-access-rules.md) if you haven't already.
This chapter explains the different pipeline handlers available to you:

- [Authentication handlers](pipeline/authn.md) inspect HTTP requests (for example the HTTP Authorization Header) and execute some
  business logic that return true (for authentication ok) or false (for authentication invalid) as well as a subject ("user"). The
  subject is typically the "user" that made the request, but it could also be a machine (if you have machine-2-machine
  interaction) or something different.
- [Authorization handlers](pipeline/authz.md): ensure that a subject ("user") has the right permissions. For example, a specific
  endpoint might only be accessible to subjects ("users") from group "admin". The authorizer handles that logic.
- [Mutation handlers](pipeline/mutator.md): transforms the credentials from incoming requests to credentials that your backend
  understands. For example, the `Authorization: basic` header might be transformed to `X-User: <subject-id>`. This allows you to
  write backends that don't care if the original request was an anonymous one, an OAuth 2.0 Access Token, or some other credential
  type. All your backend has to do is understand, for example, the `X-User:`.
- [Error handlers](pipeline/error.md): are responsible for executing logic after, for example, authentication or authorization
  failed. Ory Oathkeeper supports different error handlers and we will add more as the project progresses.

## Templating

Some handlers such as the [ID Token Mutator](pipeline/mutator.md#id_token) support templating using
[Golang Text Templates](https://golang.org/pkg/text/template/)
([examples](https://blog.gopheracademy.com/advent-2017/using-go-templates/)). The [sprig](http://masterminds.github.io/sprig/) is
also supported, on top of these two functions:

```go
var _ = template.FuncMap{
    "print": func(i interface{}) string {
        if i == nil {
            return ""
        }
        return fmt.Sprintf("%v", i)
    },
    "printIndex": func(element interface{}, i int) string {
        if element == nil {
            return ""
        }

        list := reflect.ValueOf(element)

        if list.Kind() == reflect.Slice && i < list.Len() {
            return fmt.Sprintf("%v", list.Index(i))
        }

        return ""
    },
}
```

## Session

In all configurations supporting [templating](#templating) instructions, it's possible to use the
[`AuthenticationSession`](https://github.com/ory/oathkeeper/blob/master/pipeline/authn/authenticator.go#L39) struct content.

```go
type AuthenticationSession struct {
  Subject      string
  Extra        map[string]interface{}
  Header       http.Header
  MatchContext MatchContext
}

type MatchContext struct {
  RegexpCaptureGroups []string
  URL                 *url.URL
  Method              string
  Header              http.Header
}
```

### RegexpCaptureGroups

### Configuration Examples

To use the subject extract to the token

```json
{ "config_field": "{{ print .Subject }}" }
```

To use any arbitrary header value from the request headers

```json
{ "config_field": "{{ .MatchContext.Header.Get \"some_header\" }}" }
```

To use an embedded value in the `Extra` map (most of the time, it's a JWT token claim)

```json
{ "config_field": "{{ print .Extra.some.arbitrary.data }}" }
```

To use a Regex capture from the request URL Note the usage of `printIndex` to print a value from the array

```json
{
  "claims": "{\"aud\": \"{{ print .Extra.aud }}\", \"resource\": \"{{ printIndex .MatchContext.RegexpCaptureGroups 0 }}\""
}
```

To display a string array to JSON format, we can use the [fmt printf](https://golang.org/pkg/fmt/) function

```json
{
  "claims": "{\"aud\": \"{{ print .Extra.aud }}\", \"scope\": {{ printf \"%+q\" .Extra.scp }}}"
}
```

Note that the `AuthenticationSession` struct has a field named `Extra` which is a `map[string]interface{}`, which receives varying
introspection data from the authentication process. Because the contents of `Extra` are so variable, nested and potentially
non-existent values need special handling by the `text/template` parser, and a `print` FuncMap function has been provided to
ensure that non-existent map values will simply return an empty string, rather than `<no value>`.

If you find that your field contain the string `<no value>` then you have most likely omitted the `print` function, and it's
recommended you use it for all values out of an abundance of caution and for consistency.

In the same way, a `printIndex` FuncMap function is provided to avoid _out of range_ exception to access in a array. It can be
useful for the regexp captures which depend of the request.
---
id: authn
title: Authenticators
---

An authenticator is responsible for authenticating request credentials. Ory Oathkeeper supports different authenticators and we
will add more as the project progresses.

An authenticator inspects the HTTP request (for example the HTTP Authorization Header) and executes some business logic that
returns true (for authentication ok) or false (for authentication invalid) as well as a subject ("user"). The subject is typically
the "user" that made the request, but it could also be a machine (if you have machine-2-machine interaction) or something
different.

Each authenticator has two keys:

- `handler` (string, required): Defines the handler (for example `noop`) to be used.
- `config` (object, optional): Configures the handler. Configuration keys vary per handler. The configuration can be defined in
  the global configuration file, or per access rule.

```json
{
  "authenticators": [
    {
      "handler": "noop",
      "config": {}
    }
  ]
}
```

You can define more than one authenticator in the Access Rule. The first authenticator that's able to handle the credentials will
be consulted and other authenticators will be ignored:

```json
{
  "authenticators": [
    {
      "handler": "a"
    },
    {
      "handler": "b"
    },
    {
      "handler": "c"
    }
  ]
}
```

If handler `a` is able to handle the provided credentials, then handler `b` and `c` will be ignored. If handler `a` can't handle
the provided credentials but handler `b` can, then handler `a` and `c` will be ignored. Handling the provided credentials means
that the authenticator knows how to handle, for example, the `Authorization: basic` header. It doesn't mean that the credentials
are valid! If a handler encounters invalid credentials, then other handlers will be ignored too.

## `noop`

The `noop` handler tells Ory Oathkeeper to bypass authentication, authorization, and mutation. This implies that no authorization
will be executed and no credentials will be issued. It's basically a pass-all authenticator that allows any request to be
forwarded to the upstream URL.

> Using this handler is basically an allow-all configuration. It makes sense when the upstream handles access control itself or
> doesn't need any type of access control.

### `noop` configuration

This handler isn't configurable.

To enable this handler, set:

```yaml
# Global configuration file oathkeeper.yml
authenticators:
  noop:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true
```

### `noop` access rule example

```sh
cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{
    "handler": "noop"
  }]
}]

curl -X GET http://my-app/some-route

HTTP/1.0 200 Status OK
The request has been allowed!
```

## `unauthorized`

The `unauthorized` handler tells Ory Oathkeeper to reject all requests as unauthorized.

### `unauthorized` Configuration

This handler isn't configurable.

To enable this handler, set:

```yaml
# Global configuration file oathkeeper.yml
unauthorized:
  # Set 'enabled' to 'true' if the authenticator should be enabled and 'false' to disable the authenticator. Defaults to 'false'.
  enabled: true
```

### `unauthorized` access rule example

```sh
cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{
    "handler": "unauthorized"
  }]
}]

curl -X GET http://my-app/some-route

HTTP/1.0 401 Unauthorized
```

## `anonymous`

The `anonymous` authenticator checks whether or not an `Authorization` header is set. If not, it will set the subject to
`anonymous`.

### `anonymous` Configuration

- `subject` (string, optional) - Sets the anonymous username. Defaults to "anonymous". Common names include "guest", "anon",
  "anonymous", "unknown".

```yaml
# Global configuration file oathkeeper.yml
authenticators:
  anonymous:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true

    config:
      subject: guest
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authenticators:
  - handler: anonymous
    config:
      subject: guest
```

### `anonymous` access rule example

The following rule allows all requests to `GET http://my-app/some-route` and sets the subject name to the anonymous username, as
long as no `Authorization` header is set in the HTTP request:

```shell
cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{
    "handler": "anonymous"
  }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }]
}]

curl -X GET http://my-app/some-route

HTTP/1.0 200 Status OK
The request has been allowed! The subject is: "anonymous"

curl -X GET -H "Authorization: Bearer foobar" http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because credentials have been provided but only the anonymous
authenticator is enabled for this URL.
```

## `cookie_session`

The `cookie_session` authenticator will forward the request method, path and headers to a session store. If the session store
returns `200 OK` and body `{ "subject": "...", "extra": {} }` then the authenticator will set the subject appropriately. Please
note that Gzipped responses from `check_session_url` are not supported, and will fail silently.

### `cookie_session` configuration

- `check_session_url` (string, required) - The session store to forward request method/path/headers to for validation.
- `only` ([]string, optional) - If set, only requests that have at least one of the set cookies will be forwarded, others will be
  passed to the next authenticator. If unset, all requests are forwarded.
- `preserve_path` (boolean, optional) - If set, any path in `check_session_url` will be preserved instead of replacing the path
  with the path of the request being checked.
- `preserve_query` (boolean, optional) - If unset or true, query parameters in `check_session_url` will be preserved instead of
  replacing them with the query of the request being checked.
- `force_method` (string, optional) - If set uses the given HTTP method when forwarding the request, instead of the original HTTP
  method.
- `extra_from` (string, optional - defaults to `extra`) - A [GJSON Path](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)
  pointing to the `extra` field. This defaults to `extra`, but it could also be `@this` (for the root element), `session.foo.bar`
  for `{ "subject": "...", "session": { "foo": {"bar": "whatever"} } }`, and so on.
- `subject_from` (string, optional - defaults to `subject`) - A
  [GJSON Path](https://github.com/tidwall/gjson/blob/master/SYNTAX.md) pointing to the `subject` field. This defaults to
  `subject`. Example: `identity.id` for `{ "identity": { "id": "1234" } }`.
- `additional_headers` (map[string]string, optional - defaults empty) - If set, you can either add additional headers or override
  existing ones.
- `forward_http_headers` ([]string, optional - defaults ["Authorization", "Cookie"]) - If set, you can specify which headers will
  be forwarded.

```yaml
# Global configuration file oathkeeper.yml
authenticators:
  cookie_session:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true

    config:
      check_session_url: https://session-store-host
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authenticators:
  - handler: cookie_session
    config:
      check_session_url: https://session-store-host
      only:
        - sessionid
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authenticators:
  - handler: cookie_session
    config:
      check_session_url: https://session-store-host
      only:
        - sessionid
      forward_http_headers:
        - Connect
        - Authorization
        - Cookie
        - X-Forwarded-For
```

```yaml
# Some Access Rule Preserving Path: access-rule-2.yaml
id: access-rule-2
# match: ...
# upstream: ...
authenticators:
  - handler: cookie_session
    config:
      check_session_url: https://session-store-host/check-session
      only:
        - sessionid
      preserve_path: true
      preserve_query: true
```

### `cookie_session` access rule example

```shell
cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{
    "handler": "cookie_session"
  }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }]
}]

curl -X GET -b sessionid=abc http://my-app/some-route

HTTP/1.0 200 OK
The request has been allowed! The subject is: "peter"

curl -X GET -b sessionid=def http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because the provided credentials are invalid.
```

## `bearer_token`

The `bearer_token` authenticator will forward the request method, path and headers to a session store. If the session store
returns `200 OK` and body `{ "subject": "...", "extra": {} }` then the authenticator will set the subject appropriately. Please
note that Gzipped responses from `check_session_url` are not supported, and will fail silently.

### `bearer_token` configuration

- `check_session_url` (string, required) - The session store to forward request method/path/headers to for validation.
- `preserve_path` (boolean, optional) - If set, any path in `check_session_url` will be preserved instead of replacing the path
  with the path of the request being checked.
- `preserve_query` (boolean, optional) - If unset or true, query parameters in `check_session_url` will be preserved instead of
  replacing them with the query of the request being checked.
- `force_method` (string, optional) - If set uses the given HTTP method when forwarding the request, instead of the original HTTP
  method.
- `extra_from` (string, optional - defaults to `extra`) - A [GJSON Path](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)
  pointing to the `extra` field. This defaults to `extra`, but it could also be `@this` (for the root element), `session.foo.bar`
  for `{ "subject": "...", "session": { "foo": {"bar": "whatever"} } }`, and so on.
- `subject_from` (string, optional - defaults to `sub`) - A [GJSON Path](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)
  pointing to the `sub` field. This defaults to `sub`. Example: `identity.id` for `{ "identity": { "id": "1234" } }`.
- `token_from` (object, optional) - The location of the bearer token. If not configured, the token will be received from a default
  location - 'Authorization' header. One and only one location (header, query, or cookie) must be specified.
  - `header` (string, required, one of) - The header (case insensitive) that must contain a Bearer token for request
    authentication. It can't be set along with `query_parameter` or `cookie`.
  - `query_parameter` (string, required, one of) - The query parameter (case sensitive) that must contain a Bearer token for
    request authentication. It can't be set along with `header` or `cookie`.
  - `cookie` (string, required, one of) - The cookie (case sensitive) that must contain a Bearer token for request authentication.
    It can't be set along with `header` or `query_parameter`
- `forward_http_headers` ([]string, optional - defaults ["Authorization", "Cookie"]) - If set, you can specify which headers will
  be forwarded.

```yaml
# Global configuration file oathkeeper.yml
authenticators:
  bearer_token:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true

    config:
      check_session_url: https://session-store-host
      token_from:
        header: Custom-Authorization-Header
        # or
        # query_parameter: auth-token
        # or
        # cookie: auth-token
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authenticators:
  - handler: bearer_token
    config:
      check_session_url: https://session-store-host
      token_from:
        query_parameter: auth-token
        # or
        # header: Custom-Authorization-Header
        # or
        # cookie: auth-token
```

```yaml
# Some Access Rule Preserving Path: access-rule-2.yaml
id: access-rule-2
# match: ...
# upstream: ...
authenticators:
  - handler: bearer_token
    config:
      check_session_url: https://session-store-host/check-session
      token_from:
        query_parameter: auth-token
        # or
        # header: Custom-Authorization-Header
        # or
        # cookie: auth-token
      preserve_path: true
      preserve_query: true
      forward_http_headers:
        - Authorization
        - Cookie
        - X-Forwarded-For
```

### `bearer_token` access rule example

```shell
cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{
    "handler": "bearer_token"
  }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }]
}]

curl -X GET -H 'Authorization: Bearer valid-token' http://my-app/some-route

HTTP/1.0 200 OK
The request has been allowed! The subject is: "peter"

curl -X GET -H 'Authorization: Bearer invalid-token' http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because the provided credentials are invalid.
```

## `oauth2_client_credentials`

This `oauth2_client_credentials` uses the username and password from HTTP Basic Authorization
(`Authorization: Basic base64(<username:password>)` to perform the OAuth 2.0 Client Credentials grant in order to detect if the
provided credentials are valid.

This authenticator will use the username from the HTTP Basic Authorization header as the subject for this request.

> If you are unfamiliar with OAuth 2.0 Client Credentials we recommend
> [reading this guide](https://www.oauth.com/oauth2-servers/access-tokens/client-credentials/).

### `oauth2_client_credentials` configuration

- `token_url` (string, required) - The OAuth 2.0 Token Endpoint that will be used to validate the client credentials.
- `retry` (object, optional) - Configures timeout and delay settings for the request against the token endpoint
  - `give_up_after` (string) timeout
  - `max_delay` (string) time to wait between retries
- `cache` (object, optional) - Enables caching of requested tokens
  - `enabled` (bool, optional) - Enable the cache, will use exp time of token to determine when to evict from cache. Defaults to
    false.
  - `ttl` (string) - Can override the default behavior of using the token exp time, and specify a set time to live for the token
    in the cache. If the token exp time is lower than the set value the token exp time will be used instead.
  - `max_tokens` (int) - Max number of tokens to cache. Defaults to 1000.
- `required_scope` ([]string, optional) - Sets what scope is required by the URL and when making performing OAuth 2.0 Client
  Credentials request, the scope will be included in the request:

```yaml
# Global configuration file oathkeeper.yml
authenticators:
  oauth2_client_credentials:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true

    config:
      token_url: https://my-website.com/oauth2/token
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authenticators:
  - handler: oauth2_client_credentials
    config:
      token_url: https://my-website.com/oauth2/token
```

### `oauth2_client_credentials` access rule example

```shell
cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{
    "handler": "oauth2_client_credentials",
    "config": {
      "required_scope": ["scope-a", "scope-b"]
    }
  }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }]
}]

curl -X GET http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because no credentials have been provided.

curl -X GET --user idonotexist:whatever http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because the provided credentials are invalid.

curl -X GET --user peter:somesecret http://my-app/some-route

HTTP/1.0 200 OK
The request has been allowed! The subject is: "peter"
```

In the background, a request to the OAuth 2.0 Token Endpoint (value of `authenticators.oauth2_client_credentials.token_url`) will
be made, using the OAuth 2.0 Client Credentials Grant:

```bash
POST /oauth2/token HTTP/1.1
Host: authorization-server.com

grant_type=client_credentials
&client_id=peter
&client_secret=somesecret
&scope=scope-a+scope-b
```

If the request succeeds, the credentials are considered valid and if the request fails, the credentials are considered invalid.

## `oauth2_introspection`

The `oauth2_introspection` authenticator handles requests that have an Bearer Token in the Authorization Header
(`Authorization: bearer <token>`) or in a different header or query parameter specified in configuration. It then uses OAuth 2.0
Token Introspection to check if the token is valid and if the token was granted the requested scope.

> If you are unfamiliar with OAuth 2.0 Introspection we recommend
> [reading this guide](https://www.oauth.com/oauth2-servers/token-introspection-endpoint/).

### `oauth2_introspection` configuration

- `introspection_url` (string, required) - The OAuth 2.0 Token Introspection endpoint.
- `scope_strategy` (string, optional) - Sets the strategy to be used to validate/match the token scope. Supports "hierarchic",
  "exact", "wildcard", "none". Defaults to "none".
- `required_scope` ([]string, optional) - Sets what scope is required by the URL and when performing OAuth 2.0 Client Credentials
  request, the scope will be included in the request.
- `target_audience` ([]string, optional) - Sets what audience is required by the URL.
- `trusted_issuers` ([]string, optional) - Sets a list of trusted token issuers.
- `pre_authorization` (object, optional) - Enable pre-authorization in cases where the OAuth 2.0 Token Introspection endpoint is
  protected by OAuth 2.0 Bearer Tokens that can be retrieved using the OAuth 2.0 Client Credentials grant.
  - `enabled` (bool, optional) - Enable pre-authorization. Defaults to false.
  - `client_id` (string, required if enabled) - The OAuth 2.0 Client ID to be used for the OAuth 2.0 Client Credentials Grant.
  - `client_secret` (string, required if enabled) - The OAuth 2.0 Client Secret to be used for the OAuth 2.0 Client Credentials
    Grant.
  - `token_url` (string, required if enabled) - The OAuth 2.0 Token Endpoint where the OAuth 2.0 Client Credentials Grant will be
    performed.
  - `audience` (string, optional) - The OAuth 2.0 Audience to be requested during the OAuth 2.0 Client Credentials Grant.
  - `scope` ([]string, optional) - The OAuth 2.0 Scope to be requested during the OAuth 2.0 Client Credentials Grant.
- `token_from` (object, optional) - The location of the bearer token. If not configured, the token will be received from a default
  location - 'Authorization' header. One and only one location (header, query, or cookie) must be specified.
  - `header` (string, required, one of) - The header (case insensitive) that must contain a Bearer token for request
    authentication. It can't be set along with `query_parameter` or `cookie`.
  - `query_parameter` (string, required, one of) - The query parameter (case sensitive) that must contain a Bearer token for
    request authentication. It can't be set along with `header` or `cookie`.
  - `cookie` (string, required, one of) - The cookie (case sensitive) that must contain a Bearer token for request authentication.
    It can't be set along with `header` or `query_parameter`
- `introspection_request_headers` (object, optional) - Additional headers to add to the introspection request.
- `retry` (object, optional) - Configure the retry policy
  - `max_delay` (string, optional, default to 500ms) - Maximum delay to wait before retrying the request
  - `give_up_after` (string, optional, default to 1s) - Maximum delay allowed for retries
- `cache` (object, optional) - Enables caching of incoming tokens
  - `enabled` (bool, optional) - Enable the cache, will use exp time of token to determine when to evict from cache. Defaults to
    false.
  - `ttl` (string) - Can override the default behavior of using the token exp time, and specify a set time to live for the token
    in the cache.
  - `max_cost` (int) - Max cost to cache. Defaults to 100000000.

Please note that caching won't be used if the scope strategy is `none` and `required_scope` isn't empty. In that case, the
configured introspection URL will always be called and is expected to check if the scope is valid or not.

```yaml
# Global configuration file oathkeeper.yml
authenticators:
  oauth2_introspection:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true

    config:
      introspection_url: https://my-website.com/oauth2/introspection
      scope_strategy: exact
      required_scope:
        - photo
        - profile
      target_audience:
        - example_audience
      trusted_issuers:
        - https://my-website.com/
      pre_authorization:
        enabled: true
        client_id: some_id
        client_secret: some_secret
        scope:
          - introspect
        token_url: https://my-website.com/oauth2/token
      token_from:
        header: Custom-Authorization-Header
        # or
        # query_parameter: auth-token
        # or
        # cookie: auth-token
      introspection_request_headers:
        x-forwarded-proto: https
      retry:
        max_delay: 300ms
        give_up_after: 2s
      cache:
        enabled: true
        ttl: 60s
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authenticators:
  - handler: oauth2_introspection
    config:
      introspection_url: https://my-website.com/oauth2/introspection
      scope_strategy: exact
      required_scope:
        - photo
        - profile
      target_audience:
        - example_audience
      trusted_issuers:
        - https://my-website.com/
      pre_authorization:
        enabled: true
        client_id: some_id
        client_secret: some_secret
        scope:
          - introspect
        token_url: https://my-website.com/oauth2/token
      token_from:
        query_parameter: auth-token
        # or
        # header: Custom-Authorization-Header
        # or
        # cookie: auth-token
      introspection_request_headers:
        x-forwarded-proto: https
        x-foo: bar
      retry:
        max_delay: 300ms
        give_up_after: 2s
```

### `oauth2_introspection` access rule example

```shell
cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{
    "handler": "oauth2_introspection",
    "config": {
      "required_scope": ["scope-a", "scope-b"],
      "target_audience": ["example_audience"]
    }
  }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }]
}]

curl -X GET http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because no credentials have been provided.

curl -X GET -H 'Authorization: Bearer invalid-token' http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because the provided credentials are invalid.

curl -X GET -H 'Authorization: Bearer valid.access.token.from.peter' http://my-app/some-route

HTTP/1.0 200 OK
The request has been allowed! The subject is: "peter"
```

In the background, this handler will make a request to the OAuth 2.0 Token Endpoint (configuration value
`authenticators.oauth2_introspection.introspection_url`) to check if the Bearer Token is valid:

```bash
POST /oauth2/introspect HTTP/1.1

token=valid.access.token.from.peter
```

If pre-authorization is enabled, that request will include an Authorization Header:

```bash
POST /oauth2/introspect HTTP/1.1
Authorization: Bearer token-received-by-performing-pre-authorization

token=valid.access.token.from.peter
```

The Token is considered valid if the Introspection response is HTTP 200 OK and includes `{"active":true}` in the response payload.
The subject is extracted from the `username` field.

## `jwt`

The `jwt` authenticator handles requests that have an Bearer Token in the Authorization Header (`Authorization: bearer <token>`)
or in a different header or query parameter specified in configuration. It assumes that the token is a JSON Web Token and tries to
verify the signature of it.

### `jwt` configuration

- `jwks_urls` ([]string, required) - The URLs where Ory Oathkeeper can retrieve JSON Web Keys from for validating the JSON Web
  Token. Usually something like `https://my-keys.com/.well-known/jwks.json`. The response of that endpoint must return a JSON Web
  Key Set (JWKS).
- `jwks_max_wait` (duration, optional) - The maximum time for which the JWK fetcher should wait for the JWK request to complete.
  After the interval passes, the JWK fetcher will return expired or no JWK at all. If the initial JWK request finishes
  successfully, it will still refresh the cached JWKs. Defaults to "1s".
- `jwks_ttl` (duration, optional) - The duration for which fetched JWKs should be cached internally. Defaults to "30s".
- `scope_strategy` (string, optional) - Sets the strategy to be used to validate/match the scope. Supports "hierarchic", "exact",
  "wildcard", "none". Defaults to "none".
- If `trusted_issuers` ([]string) is set, the JWT must contain a value for claim `iss` that matches _exactly_ (case-sensitive) one
  of the values of `trusted_issuers`. If no values are configured, the issuer will be ignored.
- If `target_audience` ([]string) is set, the JWT must contain all values (exact, case-sensitive) in the claim `aud`. If no values
  are configured, the audience will be ignored.
- Value `allowed_algorithms` ([]string) sets what signing algorithms are allowed. Defaults to `RS256`.
- Value `required_scope` ([]string) validates the scope of the JWT. It will checks for claims `scp`, `scope`, `scopes` in the JWT
  when validating the scope as that claim isn't standardized.
- `token_from` (object, optional) - The location of the bearer token. If not configured, the token will be received from a default
  location - 'Authorization' header. One and only one location (header, query, or cookie) must be specified.
  - `header` (string, required, one of) - The header (case insensitive) that must contain a Bearer token for request
    authentication. It can't be set along with `query_parameter` or `cookie`.
  - `query_parameter` (string, required, one of) - The query parameter (case sensitive) that must contain a Bearer token for
    request authentication. It can't be set along with `header` or `cookie`.
  - `cookie` (string, required, one of) - The cookie (case sensitive) that must contain a Bearer token for request authentication.
    It can't be set along with `header` or `query_parameter`

```yaml
# Global configuration file oathkeeper.yml
authenticators:
  jwt:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true

    config:
      jwks_urls:
        - https://my-website.com/.well-known/jwks.json
        - https://my-other-website.com/.well-known/jwks.json
        - file://path/to/local/jwks.json
      scope_strategy: none
      required_scope:
        - scope-a
        - scope-b
      target_audience:
        - https://my-service.com/api/users
        - https://my-service.com/api/devices
      trusted_issuers:
        - https://my-issuer.com/
      allowed_algorithms:
        - RS256
      token_from:
        header: Custom-Authorization-Header
        # or
        # query_parameter: auth-token
        # or
        # cookie: auth-token
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authenticators:
  - handler: jwt
    config:
      jwks_urls:
        - https://my-website.com/.well-known/jwks.json
        - https://my-other-website.com/.well-known/jwks.json
        - file://path/to/local/jwks.json
      scope_strategy: none
      required_scope:
        - scope-a
        - scope-b
      target_audience:
        - https://my-service.com/api/users
        - https://my-service.com/api/devices
      trusted_issuers:
        - https://my-issuer.com/
      allowed_algorithms:
        - RS256
      token_from:
        query_parameter: auth-token
        # or
        # header: Custom-Authorization-Header
        # or
        # cookie: auth-token
```

#### `jwt` validation example

```json
{
  "handler": "jwt",
  "config": {
    "required_scope": ["scope-a", "scope-b"],
    "target_audience": ["https://my-service.com/api/users", "https://my-service.com/api/devices"],
    "trusted_issuers": ["https://my-issuer.com/"],
    "allowed_algorithms": ["RS256", "RS256"]
  }
}
```

That exemplary Access Rule consider the following (decoded) JSON Web Token as valid:

```json
{
  "alg": "RS256"
}
{
  "iss": "https://my-issuer.com/",
  "aud": ["https://my-service.com/api/users", "https://my-service.com/api/devices"],
  "scp": ["scope-a", "scope-b"]
}
```

And this token as invalid (audience is missing, issuer isn't matching, scope is missing, wrong algorithm):

```json
{
  "alg": "HS256"
}
{
  "iss": "https://not-my-issuer.com/",
  "aud": ["https://my-service.com/api/users"],
  "scp": ["not-scope-a", "scope-b"]
}
```

### `jwt` Scope

JSON Web Tokens can be scoped. However, that feature isn't standardized and several claims that represent the token scope have
been seen in the wild: `scp`, `scope`, `scopes`. Additionally, the claim value can be a string (`"scope-a"`), a space-delimited
string (`"scope-a scope-b"`) or a JSON string array (`["scope-a", "scope-b"]`). Because of this ambiguity, all of those claims are
checked and parsed and will be available as `scp` (string array) in the authentication session (`.Extra["scp"]`).

### `jwt` access rule example

```shell
cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{
    "handler": "jwt",
    "config": {
      "required_scope": ["scope-a", "scope-b"],
      "target_audience": ["aud-1"],
      "trusted_issuers": ["iss-1"]
    }
  }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }]
}]

curl -X GET http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because no credentials have been provided.

curl -X GET -H 'Authorization: Bearer invalid-token' http://my-app/some-route

HTTP/1.0 401 Status Unauthorized
The request isn't authorized because the provided credentials are invalid.

curl -X GET -H 'Authorization: Bearer valid.jwtfrom.peter' http://my-app/some-route

HTTP/1.0 200 OK
The request has been allowed! The subject is: "peter"
```

In the background, this handler will fetch all JSON Web Key Sets provided by configuration key `authenticators.jwt.jwks_urls` and
use those keys to verify the signature. If the signature can't be verified by any of those keys, the JWT is considered invalid.
---
id: authz
title: Authorizers
---

An "authorizer" is responsible for properly permissioning a subject. Ory Oathkeeper supports different kinds of authorizers. The
list of authorizers increases over time due to new features and requirements.

Authorizers assure that a subject, for instance a "user", has the permissions necessary to access or perform a particular service.
For example, an authorizer can permit access to an endpoint or URL for specific subjects or "users" from a specific group "admin".
The authorizer permits the subjects the desired access to the endpoint.

Each authorizer has two keys:

- `handler` (string, required): Defines the handler, for example `noop`, to be used.
- `config` (object, optional): Configures the handler. Configuration keys can vary for each handler.s

```json
{
  "authorizer": {
    "handler": "noop",
    "config": {}
  }
}
```

There is a 1:1 mandatory relationship between an authoriser and an access rule. It isn't possible to configure more than one
authorizer per Access Rule.

## `allow`

This authorizer permits every action allowed.

### `allow` configuration

This handler isn't configurable.

To enable this handler, set as follows:

```yaml
# Global configuration file oathkeeper.yml
authorizers:
  allow:
    # Set enabled to "true" to enable the authenticator, and "false" to disable the authenticator. Defaults to "false".
    enabled: true
```

### `allow` access rule example

```sh
$ cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{ "handler": "anonymous" }],
  "authorizer": { "handler": "allow" },
  "mutators": [{ "handler": "noop" }]
}]

$ curl -X GET http://my-app/some-route

HTTP/1.0 200 Status OK
The request has been allowed!
```

## `deny`

This authorizer considers every action unauthorized therefore "forbidden" or "disallowed".

### `deny` configuration

This handler isn't configurable.

To enable this handler, set:

```yaml
# Global configuration file oathkeeper.yml
authorizers:
  deny:
    # Set enabled to "true" to enable the authenticator, and "false" to disable the authenticator. Defaults to "false".
    enabled: true
```

### `deny` access rule example

```sh
$ cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/some-route",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [{ "handler": "anonymous" }],
  "authorizer": { "handler": "deny" },
  "mutators": [{ "handler": "noop" }]
}]

$ curl -X GET http://my-app/some-route

HTTP/1.0 403 Forbidden
The request is forbidden!
```

## `keto_engine_acp_ory`

This authorizer uses the Ory Keto API to carry out access control using "Ory-flavored" Access Control Policies. The conventions
used in the Ory Keto project are located on [GitHub Ory Keto](https://github.com/ory/keto) for consultation prior to using this
authorizer.

### `keto_engine_acp_ory` configuration

- `base_url` (string, required) - The base URL of Ory Keto, typically something like `https://hostname:port/`
- `required_action` (string, required) - See section below.
- `required_resource` (string, required) - See section below.
- `subject` (string, optional) - See section below.
- `flavor` (string, optional) - See section below.

#### Resource, relation (action), subject

> Actions were renamed to relations. Read the
> [Ory Keto policy migration guide](https://www.ory.sh/docs/keto/guides/migrating-legacy-policies#rewriting-it-to-relationships)
> for more details.

This authorizer has four configuration options, `required_action`, `required_resource`, `subject`, and `flavor`:

```json
{
  "handler": "keto_engine_acp_ory",
  "config": {
    "required_action": "...",
    "required_resource": "...",
    "subject": "...",
    "flavor": "..."
  }
}
```

All configuration options except `flavor` support Go [`text/template`](https://golang.org/pkg/text/template/). For example in the
following match configuration:

```json
{
  "match": {
    "url": "http://my-app/api/users/<[0-9]+>/<[a-zA-Z]+>",
    "methods": ["GET"]
  }
}
```

The following example shows how to reference the values matched by or resulting from the two regular expressions, `<[0-9]+>` and
`<[a-zA-Z]+>`. using the `AuthenticationSession` struct:

```json
{
  "handler": "keto_engine_acp_ory",
  "config": {
    "required_action": "my:action:{{ printIndex .MatchContext.RegexpCaptureGroups 0 }}",
    "required_resource": "my:resource:{{ printIndex .MatchContext.RegexpCaptureGroups 1 }}:foo:{{ printIndex .MatchContext.RegexpCaptureGroups 0 }}"
  }
}
```

Assuming a request to `http://my-api/api/users/1234/foobar` was made, the config from above would expand to:

```json
{
  "handler": "keto_engine_acp_ory",
  "config": {
    "required_action": "my:action:1234",
    "required_resource": "my:resource:foobar:foo:1234"
  }
}
```

The `subject` field configures the subject that passes to the Ory Keto endpoint. If `subject` isn't specified it will default to
`AuthenticationSession.Subject`.

For more details about supported Go template substitution, see. [How to use session variables](../pipeline.md#session)

#### `keto_engine_acp_ory` example

```yaml
# Global configuration file oathkeeper.yml
authorizers:
  keto_engine_acp_ory:
    # Set enabled to "true" to enable the authenticator, and "false" to disable the authenticator. Defaults to "false".
    enabled: true

    config:
      base_url: http://my-keto/
      required_action: ...
      required_resource: ...
      subject: ...
      flavor: ...
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authorizers:
  - handler: keto_engine_acp_ory
    config:
      base_url: http://my-keto/
      required_action: ...
      required_resource: ...
      subject: ...
      flavor: ...
```

### `keto_engine_acp_ory` access rule example

```shell
$ cat ./rules.json

[{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/api/users/<[0-9]+>/<[a-zA-Z]+>",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [
    {
      "handler": "anonymous"
    }
  ],
  "authorizer": {
    "handler": "keto_engine_acp_ory",
    "config": {
      "required_action": "my:action:$1",
      "required_resource": "my:resource:$2:foo:$1"
      "subject": "{{ .Extra.email }}",
      "flavor": "exact"
    }
  }
  "mutators": [
    {
      "handler": "noop"
    }
  ]
}]
```

## `remote`

This authorizer performs authorization using a remote authorizer. The authorizer makes a HTTP POST request to a remote endpoint
with the original body request as body. If the endpoint returns a "200 OK" response code, the access is allowed, if it returns a
"403 Forbidden" response code, the access is denied.

### `remote` configuration

- `remote` (string, required) - The remote authorizer's URL. The remote authorizer is expected to return either "200 OK" or "403
  Forbidden" to allow/deny access.
- `headers` (map of strings, optional) - The HTTP headers sent to the remote authorizer. The values will be parsed by the Go
  [`text/template`](https://golang.org/pkg/text/template/) package and applied to an
  [`AuthenticationSession`](https://github.com/ory/oathkeeper/blob/master/pipeline/authn/authenticator.go#L40) object. See
  [Session](../pipeline.md#session) for more details.
- `forward_response_headers_to_upstream` (slice of strings, optional) - The HTTP headers that will be allowed from remote
  authorizer responses. If returned, headers on this list will be forward to upstream services.
- `retry` (object, optional) - Configures timeout and delay settings for the request against the token endpoint
  - `give_up_after` (string) max delay duration of retry. The value will be parsed by the Go
    [duration parser](https://pkg.go.dev/time#ParseDuration).
  - `max_delay` (string) time to wait between retries and max service response time. The value will be parsed by the Go
    [duration parser](https://pkg.go.dev/time#ParseDuration).

#### `remote` example

```yaml
# Global configuration file oathkeeper.yml
authorizers:
  remote:
    # Set enabled to "true" to enable the authenticator, and "false" to disable the authenticator. Defaults to "false".
    enabled: true

    config:
      remote: http://my-remote-authorizer/authorize
      headers:
        X-Subject: "{{ print .Subject }}"
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authorizers:
  - handler: remote
    config:
      remote: http://my-remote-authorizer/authorize
      headers:
        X-Subject: "{{ print .Subject }}"
```

### `remote` access rule example

```shell
{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/api/<.*>",
    "methods": ["GET"]
  },
  "authenticators": [
    {
      "handler": "anonymous"
    }
  ],
  "authorizer": {
    "handler": "remote",
    "config": {
      "remote": "http://my-remote-authorizer/authorize",
      "headers": {
        "X-Subject": "{{ print .Subject }}"
      },
      "forward_response_headers_to_upstream": [
        "X-Foo"
      ]
    }
  }
  "mutators": [
    {
      "handler": "noop"
    }
  ]
}
```

## `remote_json`

This authorizer performs authorization using a remote authorizer. The authorizer makes a HTTP POST request to a remote endpoint
with a JSON body. If the endpoint returns a "200 OK" response code, the access is allowed, if it returns a "403 Forbidden"
response code, the access is denied.

### `remote_json` configuration

- `remote` (string, required) - The remote authorizer's URL. The remote authorizer is expected to return either "200 OK" or "403
  Forbidden" to allow/deny access.
- `payload` (string, required) - The request's JSON payload sent to the remote authorizer. The string will be parsed by the Go
  [`text/template`](https://golang.org/pkg/text/template/) package and applied to an
  [`AuthenticationSession`](https://github.com/ory/oathkeeper/blob/master/pipeline/authn/authenticator.go#L40) object. See
  [Session](../pipeline.md#session) for more details.
- `headers` (map of strings, optional) - The HTTP headers sent to the remote authorizer. The values will be parsed by the Go
  [`text/template`](https://golang.org/pkg/text/template/) package and applied to an
  [`AuthenticationSession`](https://github.com/ory/oathkeeper/blob/master/pipeline/authn/authenticator.go#L40) object. See
  [Session](../pipeline.md#session) for more details.
- `forward_response_headers_to_upstream` (slice of strings, optional) - The HTTP headers that will be allowed from remote
  authorizer responses. If returned, headers on this list will be forward to upstream services.
- `retry` (object, optional) - Configures timeout and delay settings for the request against the token endpoint
  - `give_up_after` (string) max delay duration of retry. The value will be parsed by the Go
    [duration parser](https://pkg.go.dev/time#ParseDuration).
  - `max_delay` (string) time to wait between retries and max service response time. The value will be parsed by the Go
    [duration parser](https://pkg.go.dev/time#ParseDuration).

#### `remote_json` example

```yaml
# Global configuration file oathkeeper.yml
authorizers:
  remote_json:
    # Set enabled to "true" to enable the authenticator, and "false" to disable the authenticator. Defaults to "false".
    enabled: true

    config:
      remote: http://my-remote-authorizer/authorize
      headers:
        Y-Api-Key: '{{ .MatchContext.Header.Get "X-Api-Key" }}'
      payload: |
        {
          "subject": "{{ print .Subject }}",
          "resource": "{{ printIndex .MatchContext.RegexpCaptureGroups 0 }}"
        }
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
authorizers:
  - handler: remote_json
    config:
      remote: http://my-remote-authorizer/authorize
      headers:
        Y-Api-Key: '{{ .MatchContext.Header.Get "X-Api-Key" }}'
      payload: |
        {
          "subject": "{{ print .Subject }}",
          "resource": "{{ printIndex .MatchContext.RegexpCaptureGroups 0 }}"
        }
```

### `remote_json` access rule example

```shell
{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/api/<.*>",
    "methods": ["GET"]
  },
  "authenticators": [
    {
      "handler": "anonymous"
    }
  ],
  "authorizer": {
    "handler": "remote_json",
    "config": {
      "headers": {
         "Y-Api-Key": "{{ .MatchContext.Header.Get \"X-Api-Key\" }}"
      },
      "remote": "http://my-remote-authorizer/authorize",
      "payload": "{\"subject\": \"{{ print .Subject }}\", \"resource\": \"{{ printIndex .MatchContext.RegexpCaptureGroups 0 }}\"}"
    },
    "forward_response_headers_to_upstream": [
      "X-Foo"
    ]
  }
  "mutators": [
    {
      "handler": "noop"
    }
  ]
}
```
---
id: mutator
title: Mutators
---

A mutator transforms the credentials from incoming requests to credentials that your backend understands. For example, the
`Authorization: basic` header might be transformed to `X-User: <subject-id>`. This allows you to write backends that don't care if
the original request was an anonymous one, an OAuth 2.0 Access Token, or some other credential type. All your backend has to do is
understand, for example, the `X-User:`.

The Access Control Decision API will return the mutated result as the HTTP Response.

## `noop`

This mutator doesn't transform the HTTP request and simply forwards the headers as-is. This is useful if you don't want to
replace, for example, `Authorization: basic` with `X-User: <subject-id>`.

### `noop` configuration

```yaml
# Global configuration file oathkeeper.yml
mutators:
  noop:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
mutators:
  - handler: noop
```

### `noop` access rule example

```shell
cat ./rules.json
{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/api/users/<[0-9]+>/<[a-zA-Z]+>",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [
    {
      "handler": "anonymous"
    }
  ],
  "authorizer": {
    "handler": "allow"
  },
  "mutators": [
    {
      "handler": "noop"
    }
  ]
}

curl -X GET http://my-app/some-route

HTTP/1.0 200 Status OK
The request has been allowed! The original HTTP Request hasn't been modified.
```

## `id_token`

This mutator takes the authentication information (such as subject) and transforms it to a signed JSON Web Token, and more
specifically to an OpenID Connect ID Token. Your backend can verify the token by fetching the (public) key from the
`/.well-known/jwks.json` endpoint provided by the Ory Oathkeeper API.

Let's say a request is made to a resource protected by Ory Oathkeeper using Basic Authorization:

```bash
GET /api/resource HTTP/1.1
Host: www.example.com
Authorization: Basic Zm9vOmJhcg==
```

Assuming that Ory Oathkeeper is granting the access request, `Basic Zm9vOmJhcg==` will be replaced with a cryptographically signed
JSON Web Token:

```bash
GET /api/resource HTTP/1.1
Host: internal-api-endpoint-dns
Authorization: Bearer <jwt-signed-id-token>
```

Now, the protected resource is capable of decoding and validating the JSON Web Token using the public key supplied by Ory
Oathkeeper's API. The public key for decoding the ID token is available at Ory Oathkeeper's `/.well-known/jwks.json` endpoint:

```bash
http://oathkeeper:4456/.well-known/jwks.json
```

The related flow diagram looks like this:

Let's say the `oauth2_client_credentials` authenticator successfully authenticated the credentials `client-id:client-secret`. This
mutator will craft an ID Token (JWT) with the following exemplary claims:

```json
{
  "iss": "https://server.example.com",
  "sub": "client-id",
  "aud": "s6BhdRkqt3",
  "jti": "n-0S6_WzA2Mj",
  "exp": 1311281970,
  "iat": 1311280970
}
```

The ID Token Claims are as follows:

- `iss`: Issuer Identifier for the Issuer of the response. The iss value is a case sensitive URL using the https scheme that
  contains scheme, host, and optionally, port number and path components and no query or fragment components. Typically, this is
  the URL of Ory Oathkeeper, for example: `https://oathkeeper.myapi.com`.
- `sub`: Subject Identifier. A locally unique and never reassigned identifier within the Issuer for the End-User, which is
  intended to be consumed by the Client, for example, 24400320 or AItOawmwtWwcT0k51BayewNvutrJUqsvl6qs7A4. It must not exceed 255
  ASCII characters in length. The sub value is a case sensitive string. The End-User might also be an OAuth 2.0 Client, given that
  the access token was granted using the OAuth 2.0 Client Credentials flow.
- `aud`: Audience(s) that this ID Token is intended for. It MUST contain the OAuth 2.0 client_id of the Relying Party as an
  audience value. It MAY also contain identifiers for other audiences. In the general case, the aud value is an array of case
  sensitive strings.
- `exp`: Expiration time on or after which the ID Token MUST NOT be accepted for processing. The processing of this parameter
  requires that the current date/time MUST be before the expiration date/time listed in the value. Its value is a JSON number
  representing the number of seconds from 1970-01-01T0:0:0Z as measured in UTC until the date/time. See RFC 3339 [RFC3339] for
  details regarding date/times and UTC in particular.
- `iat`: Time at which the JWT was issued. Its value is a JSON number representing the number of seconds from 1970-01-01T0:0:0Z as
  measured in UTC until the date/time.
- `jti`: A cryptographically strong random identifier to ensure the ID Token's uniqueness.

### `id_token` configuration

- `issuer_url` (string, required) - Sets the "iss" value of the ID Token.
- `jwks_url` (string, required) - Sets the URL where keys should be fetched from. Supports remote locations (http, https, s3, gs,
  azblob) as well as local filesystem paths.
- `ttl` (string, optional) - Sets the time-to-live of the ID token. Defaults to one minute. Valid time units are: s (second), m
  (minute), h (hour).
- `claims` (string, optional) - Allows you to customize the ID Token claims and support Go Templates. For more information, check
  section [Claims](#id_token-claims)

```yaml
# Global configuration file oathkeeper.yml
mutators:
  id_token:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true
    config:
      issuer_url: https://my-oathkeeper/
      jwks_url: https://fetch-keys/from/this/location.json
      # jwks_url: file:///from/this/absolute/location.json
      # jwks_url: file://../from/this/relative/location.json
      ttl: 60s
      claims: '{"aud": ["https://my-backend-service/some/endpoint"],"def": "{{ print .Extra.some.arbitrary.data }}"}'
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
mutators:
  - handler: id_token
    config:
      issuer_url: https://my-oathkeeper/
      jwks_url: https://fetch-keys/from/this/location.json
      # jwks_url: file:///from/this/absolute/location.json
      # jwks_url: file://../from/this/relative/location.json
      ttl: 60s
      claims: '{"aud": ["https://my-backend-service/some/endpoint"],"def": "{{ print .Extra.some.arbitrary.data }}"}'
```

The first private key found in the JSON Web Key Set defined by `mutators.id_token.jwks_url` will be used for signing the JWT:

- If the first key found is a symmetric key (`HS256` algorithm), that key will be used. That key **won't** be broadcasted at
  `/.well-known/jwks.json`. You must manually configure the upstream to be able to fetch the key (for example from an environment
  variable).
- If the first key found is an asymmetric private key (for example `RS256`, `ES256`, ...), that key will be used. The related
  public key will be broadcasted at `/.well-known/jwks.json`.

#### `id_token` Claims

This mutator allows you to specify custom claims, like the audience of ID tokens, via the `claims` field of the mutator's `config`
field. The keys represent names of claims and the values are arbitrary data structures which will be parsed by the Go
[text/template](https://golang.org/pkg/text/template/) package for value substitution, receiving the `AuthenticationSession`
struct.

For more details please check [Session variables](../pipeline.md#session)

The claims configuration expects a string which is expected to be valid JSON:

```json
{
  "handler": "id_token",
  "config": {
    "claims": "{\"aud\": [\"https://my-backend-service/some/endpoint\"],\"def\": \"{{ print .Extra.some.arbitrary.data }}\"}"
  }
}
```

Please keep in mind that certain keys (such as the `sub`) claim **can't** be overwritten!

### `id_token` access rule example

```shell
cat ./rules.json
{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/api/users/<[0-9]+>/<[a-zA-Z]+>",
    "methods": [
      "GET"
    ]
  },
  "authenticators": [
    {
      "handler": "anonymous"
    }
  ],
  "authorizer": {
    "handler": "allow"
  },
  "mutators": [
    {
      "handler": "id_token",
      "config": {
        "aud": [
          "audience-1",
          "audience-2"
        ],
        "claims": "{\"abc\": \"{{ print .Subject }}\",\"def\": \"{{ print .Extra.some.arbitrary.data }}\"}"
      }
    }
  ]
}
```

## `header`

This mutator will transform the request, allowing you to pass the credentials to the upstream application via the headers. This
will augment, for example, `Authorization: basic` with `X-User: <subject-id>`.

### `header` configuration

- `headers` (object (`string: string`), required) - A keyed object (`string:string`) representing the headers to be added to this
  request, see section [headers](#headers).

```yaml
# Global configuration file oathkeeper.yml
mutators:
  header:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true
    config:
      headers:
        X-User: "{{ print .Subject }}"
        X-Some-Arbitrary-Data: "{{ print .Extra.some.arbitrary.data }}"
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
mutators:
  - handler: header
    config:
      headers:
        X-User: "{{ print .Subject }}"
        X-Some-Arbitrary-Data: "{{ print .Extra.some.arbitrary.data }}"
```

#### Headers

The headers are specified via the `headers` field of the mutator's `config` field. The keys are the header name and the values are
a string which will be parsed by the Go [`text/template`](https://golang.org/pkg/text/template/) package for value substitution,
receiving the `AuthenticationSession` struct.

For more details please check [Session variables](../pipeline.md#session)

### `header` access rule example

```json
{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/api/<.*>",
    "methods": ["GET"]
  },
  "authenticators": [
    {
      "handler": "anonymous"
    }
  ],
  "authorizer": {
    "handler": "allow"
  },
  "mutators": [
    {
      "handler": "header",
      "config": {
        "headers": {
          "X-User": "{{ print .Subject }}",
          "X-Some-Arbitrary-Data": "{{ print .Extra.some.arbitrary.data }}"
        }
      }
    }
  ]
}
```

## `cookie`

This mutator will transform the request, allowing you to pass the credentials to the upstream application via the cookies.

### `cookie` configuration

- `cookies` (object (`string: string`), required) - A keyed object (`string:string`) representing the cookies to be added to this
  request, see section [cookies](#cookies).

```yaml
# Global configuration file oathkeeper.yml
mutators:
  cookie:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true
    config:
      cookies:
        user: "{{ print .Subject }}",
        some-arbitrary-data: "{{ print .Extra.some.arbitrary.data }}"
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
mutators:
  - handler: cookie
    config:
      cookies:
        user: "{{ print .Subject }}",
        some-arbitrary-data: "{{ print .Extra.some.arbitrary.data }}"
```

### Cookies

The cookies are specified via the `cookies` field of the mutators `config` field. The keys are the cookie name and the values are
a string which will be parsed by the Go [`text/template`](https://golang.org/pkg/text/template/) package for value substitution,
receiving the `AuthenticationSession` struct.

For more details please check [Session variables](../pipeline.md#session)

#### `cookie` example

```json
{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/api/<.*>",
    "methods": ["GET"]
  },
  "authenticators": [
    {
      "handler": "anonymous"
    }
  ],
  "authorizer": {
    "handler": "allow"
  },
  "mutators": [
    {
      "handler": "cookie",
      "config": {
        "cookies": {
          "user": "{{ print .Subject }}",
          "some-arbitrary-data": "{{ print .Extra.some.arbitrary.data }}"
        }
      }
    }
  ]
}
```

## `hydrator`

This mutator allows for fetching additional data from external APIs, which can be then used by other mutators. It works by making
an upstream HTTP call to an API specified in the **Per-Rule Configuration** section below. The request is a POST request and it
contains JSON representation of
[AuthenticationSession](https://github.com/ory/oathkeeper/blob/master/pipeline/authn/authenticator.go#L39) struct in body, which
is:

```json
{
  "subject": String,
  "extra": Object,
  "header": Object,
  "match_context": {
    "regexp_capture_groups": Object,
    "url": Object
  }
}
```

As a response the mutator expects similar JSON object, but with `extra` or `header` fields modified.

Example request/response payload:

```json
{
  "subject": "anonymous",
  "extra": {
    "foo": "bar"
  },
  "header": {
    "foo": ["bar1", "bar2"]
  },
  "match_context": {
    "regexp_capture_groups": ["http", "foo"],
    "url": "http://domain.com/foo"
  }
}
```

:::note

Ory Oathkeeper is case-insensitive when accepting custom request headers. However, all incoming custom headers are converted to
the canonical format of the MIME header key. This means that the first letter of the incoming header, as well as any letter that
follows a hyphen, is converted into upper case and the rest of the letters are converted into lowercase. For example, the incoming
header `x-user-company` is converted and returned by Oathkeeper as `X-User-Company`.

:::

The AuthenticationSession from this object replaces the original one and is passed to the next mutator, where it can be used to
set a particular cookie to the value received from an API.

Setting `extra` field doesn't transform the HTTP request, whereas headers set in the `header` field will be added to the final
request headers.

### Cache

This handler supports caching. If caching is enabled, the `api.url` configuration value and the full `AuthenticationSession`
payload.

:::info

Because the cache key is quite complex, the caching handler has a higher chance of cache misses. This will be improved in future
versions.

:::

### `hydrator` configuration

- `api.url` (string - required) - The API URL.
- `api.auth.basic.*` (optional) - Enables HTTP Basic Authorization.
- `api.auth.retry.*` (optional) - Configures the retry logic.
- `cache.ttl` (optional) - Configures how long to cache hydrate requests

```yaml
# Global configuration file oathkeeper.yml
mutators:
  hydrator:
    # Set enabled to true if the authenticator should be enabled and false to disable the authenticator. Defaults to false.
    enabled: true
    config:
      api:
        url: http://my-backend-api
        auth:
          basic:
            username: someUserName
            password: somePassword
        retry:
          give_up_after: 2s
          max_delay: 100ms
      cache:
        ttl: 60s
```

```yaml
# Some Access Rule: access-rule-1.yaml
id: access-rule-1
# match: ...
# upstream: ...
mutators:
  - handler: hydrator
    config:
      api:
        url: http://my-backend-api
        auth:
          basic:
            username: someUserName
            password: somePassword
        retry:
          give_up_after: 2s
          max_delay: 100ms
      cache:
        ttl: 60s
```

### `hydrator` access rule example

```json
{
  "id": "some-id",
  "upstream": {
    "url": "http://my-backend-service"
  },
  "match": {
    "url": "http://my-app/api/<.*>",
    "methods": ["GET"]
  },
  "authenticators": [
    {
      "handler": "anonymous"
    }
  ],
  "authorizer": {
    "handler": "allow"
  },
  "mutators": [
    {
      "handler": "hydrator",
      "config": {
        "api": {
          "url": "http://my-backend-api"
        }
      }
    },
    {
      "handler": "cookie",
      "config": {
        "cookies": {
          "some-arbitrary-data": "{{ print .Extra.cookie }}"
        }
      }
    }
  ]
}
```

