# Authentication through Facebook

## Overview

One of the login options for dex uses the Facebook OAuth2 flow to identify the end user through their Facebook account.

When a client redeems a refresh token through dex, dex will re-query Facebook to update user information in the ID Token. To do this, __dex stores a readonly Facebook access token in its backing datastore.__ Users that reject dex's access through Facebook will also revoke all dex clients which authenticated them through Facebook.

## Configuration

Register a new application via `Add new application` at [Facebook for developers](https://developers.facebook.com/apps/) ensuring the callback URL is `(dex issuer)/callback`. For example if dex is listening at the non-root path `https://auth.example.com/dex` the callback would be `https://auth.example.com/dex/callback`.

The following is an example of a configuration for `examples/config-dev.yaml`:

```yaml
connectors:
  - type: facebook
    # Required field for connector id.
    id: facebook
    # Required field for connector name.
    name: Facebook
    config:
      # Credentials can be string literals or pulled from the environment.
      clientID: $FACEBOOK_APPLICATION_ID
      clientSecret: $FACEBOOK_CLIENT_SECRET
      redirectURI: http://127.0.0.1:5556/dex/callback
```