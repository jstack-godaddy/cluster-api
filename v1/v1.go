package v1

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	docs "dbs-api/docs"
	endpoints "dbs-api/v1/endpoints"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"

	gin "github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var metdataurl = "https://sso.godaddy.com/metadata" //Metadata of the IDP
var sessioncert = "./sessioncert"                   //Key pair used for creating a signed session
var sessionkey = "./sessionkey"
var serverkey = "./serverkey" //Server TLS
var servercert = "./servercert"
var serverurl = "https://cluster.api.databaseservices.int.gdcorp.tools" // base url of this service
var entityId = serverurl                                                //Entity ID uniquely identifies your service for IDP (does not have to be server url)
var listenAddr = "0.0.0.0:443"

// @BasePath /api/v1
func Initialize() {
	r := gin.Default()

	keyPair, err := tls.LoadX509KeyPair(sessioncert, sessionkey)
	panicIfError(err)
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	panicIfError(err)
	idpMetadataURL, err := url.Parse(metdataurl)
	panicIfError(err)
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient,
		*idpMetadataURL)
	panicIfError(err)
	rootURL, err := url.Parse(serverurl)
	panicIfError(err)
	samlSP, _ := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata, // you can also have Metadata XML instead of URL
		EntityID:    entityId,
	})
	app := http.HandlerFunc(hello)
	http.Handle("/hello", samlSP.RequireAccount(app))
	http.Handle("/saml/", samlSP)
	panicIfError(http.ListenAndServeTLS(listenAddr, servercert, serverkey, nil))

	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")
	{
		endpoints.Initialize(v1)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}

func hello(w http.ResponseWriter, r *http.Request) {
	s := samlsp.SessionFromContext(r.Context())
	if s == nil {
		return
	}
	sa, ok := s.(samlsp.SessionWithAttributes)
	if !ok {
		return
	}

	fmt.Fprintf(w, "Token contents, %+v!", sa.GetAttributes())
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
