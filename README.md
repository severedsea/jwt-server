# jwt-server
Sample application that demonstrates access token issuance and verification

## Pre-requisites

- Go 1.21+ - <https://golang.org/doc/go1.21>
- Docker - <https://www.docker.com>

## Getting started

### Setup

Sets up the local environment and generates JWT keys

```
$ make setup
```

### Run

Runs the server at port 3000. Environment variables can be configured in `.env.local`

```
$ make run
```

### Test

Executes all tests. Test environment variables can be configured in `.env.test`

```
$ make test
```

## How it works

### Generate access token 
```
GET /v1/login?subject={uid}
```

Access token will be returned as:
- JSON response
- Cookie

--- 

Logic: 
1. Generates the claims based on the `subject` provided
1. Signs the claims to generate an `access_token`
1. Saves the token in Redis for session management
1. Returns the token as a cookie and body in the HTTP response

### Verify access token
```
POST /v1/verify
```

Access token can be provided either by: 
- Authorization: Bearer {access_token}
- Cookie

--- 

Logic: 
1. Checks if the `access_token` is in the HTTP request (Authorization header or cookie)
1. Parse the `access_token` using the JWT public key
1. Validate expiry
1. Validate issuer
1. Check if token is in Redis

### Invalidate access token
```
GET /v1/logout
```

Access token can be provided either by: 
- Authorization: Bearer {access_token}
- Cookie

--- 

Logic: 
1. Perform the `Verify` logic
1. Delete the token in Redis
1. Invalidate the cookie
