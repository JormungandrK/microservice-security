SAML Security
============

This package contains functions and types for enabling SAML as a security mechanism
in a microservice.

# Setting up a SAML security

There are a couple of things you need to do to enable the SAML security middleware.
For details on SAML you can find many resources on the official site: http://saml.xml.org/.

## Setting up the secret keys

Create a directory in which you'll keep your key-pair:

```bash
mkdir saml-keys
cd saml-keys
```

Each service provider must have an self-signed X.509 key pair established. You can generate your own with something like this:

```bash
openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=myservice.example.com"
```

**NOTE:** Make sure you have myservice.key and myservice.cert files in the saml-keys directory

## Set up SAML in Goa

Create a security file app/security.go with the following content:

```go
package app

import (
	"github.com/crewjam/saml/samlsp"
	"net/url"
	"crypto/tls"
	"crypto/x509"
	"crypto/rsa"
)

// NewSAMLSecurity creates a SAML security definition.
func NewSAMLSecurityMiddleware(cert string, key string) *samlsp.Middleware {
	keyPair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
    if err != nil {
        panic(err)
    }

	rootURL, err := url.Parse("http://localhost:8082")
	if err != nil {
		panic(err)
	}

   	idpMetadataURL, err := url.Parse("https://www.testshib.org/metadata/testshib-providers.xml")
	if err != nil {
		panic(err)
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		IDPMetadataURL: idpMetadataURL,
		URL:            *rootURL,
		Key:            keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:    keyPair.Leaf,
	})

	return samlSP
}
```

More details on how to configure the SAML security are available on the following
site:
 * Example: https://github.com/crewjam/saml

## Setting up a SecurityChain

You need to set up a security chain for the microservice.

In the ```main.go``` file of your microservice, set up the SAML Security Chain
middleware and add it to the security chain.

```go

import (
	"github.com/JormungandrK/microservice-security/saml"
	"github.com/JormungandrK/microservice-security/chain"
)

func main() {
	// Create new SAML security chain
  	// "saml-keys" is the directory containing the keys
  	spMiddleware := app.NewSAMLSecurityMiddleware("saml-keys/myservice.cert", "saml-keys/myservice.key")
	SAMLMiddleware := saml.NewSAMLSecurity(spMiddleware)
	sc := chain.NewSecurityChain().AddMiddleware(SAMLMiddleware)

    // other initializations here...

    service.Use(chain.AsGoaMiddleware(sc)) // attach the security chain as Goa middleware
}

```

# Testing the setup

To test the setup, you'll need to generate and sign a JWT token.

Example of the JWT token:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwOi8vbG9jYWxob3N0OjgwODIvc2FtbC9tZXRhZGF0YSIsImF0dHIiO
nsib3JnYW5pemF0aW9ucyI6Ik96cmcxLCBPcmcyIiwicm9sZXMiOiJ1c2VyLCBhZG1pbiIsInVzZXJJZCI6IjU5YTAwNmFlMDAwMDAwMDA
wMDAwMDAwMCIsInVzZXJuYW1lIjoidGVzdC11c2VyIn19.vLl5hWsbYDSybhokeA4sFKJnZznesiUje5tzsCYZzl4
```

(Note that the JWT is actually one line. For readability purposes it is displayed here
  in multiple lines.)

Then you will need to set up the cookie named "saml_token":
```bash
 curl -v -b "saml_token=eyJhbGciO...<full token here>...zsCYZzl4" http://localhost:8082/profiles/m
``` 