# NICCI SDK Golang

## NICCI Profile

Usage:

### To generate the auth URI you can redirect your user to (will end up on the login system on NICCI Profile)

```go
package main

import (
    "github.com/newestindustry/nicci-sdk-golang/profile"
)

func main() {
    clientConfig := profile.ClientConfig{
        ClientID:     "clientID",
        ClientSecret: "clientSecret",
        RedirectURI:  "https://redirectURI/path",
    }
    
    authURI, err := clientConfig.GenerateAuthURI("https://auth.nicci.io", []string{"profile/basic"})
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
    
    fmt.Println(authURI.String())
}

```

### To exchange a code for an access_token

```go
package main

import (
    "github.com/newestindustry/nicci-sdk-golang/profile"
)

func main() {
    clientConfig := profile.ClientConfig{
        ClientID:     "clientID",
        ClientSecret: "clientSecret",
        RedirectURI:  "https://redirectURI/path",
    }
    
    code := "codeReceivedFromAuthFlow"
    
    at, err := clientConfig.ExchangeCode("https://auth.nicci.io", code, []string{"profile/basic"})
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
    
    fmt.Println(at)
}

```