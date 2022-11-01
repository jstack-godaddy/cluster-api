package main

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	v1 "dbs-api/v1"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/crewjam/saml/samlsp"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
)

const (
	server_base_url = "https://cluster.api.databaseservices.gdcorp.tools"
	server_cert     = "./certs/server.cert"
	server_key      = "./certs/server.key"
	sso_base_url    = "https://godaddy.okta.com"
	//sso_idp_uri      = "/api/v1/apps/0oaxcfayu8YY96pYa0x7/sso/saml/metadata"
	sso_idp_metadata = "auth/metadata.xml"
)

func main() {
	r := gin.Default()

	// Set up logging
	gin.DisableConsoleColor()
	f, _ := os.Create("/tmp/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Set up SAML SSO
	keyPair, err := tls.LoadX509KeyPair(server_cert, server_key)
	if err != nil {
		panic(err) // TODO handle error
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err) // TODO handle error
	}
	//idpMetadataURL, err := url.Parse(sso_base_url + sso_idp_uri)
	//if err != nil {
	//	panic(err) // TODO handle error
	//}
	metadata, _ := os.ReadFile(sso_idp_metadata)
	idpMetadata, err := samlsp.ParseMetadata(metadata)
	if err != nil {
		panic(err) // TODO handle error
	}
	rootURL, err := url.Parse(server_base_url)
	if err != nil {
		panic(err) // TODO handle error
	}
	samlSP, _ := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
	})
	r.Any("/saml/*action", gin.WrapH(samlSP))
	r.Use(adapter.Wrap(samlSP.RequireAccount))

	// Set up APIv1
	v1.Initialize(r)

	// Redirect the root URL to V1.
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://cluster.api.databaseservices.gdcorp.tools/v1/swagger/index.html")
	})

	// Start the webserver and secure it.
	r.RunTLS(":443", server_cert, server_key)
}
