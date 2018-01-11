package main

import (
	"net/http"
	"log"
	"fmt"
	"io"
	"crypto/tls"
	"github.com/logrusorgru/aurora"
)

var cert tls.Certificate

func main() {
	// Add shop cert
	shopCert, err := tls.LoadX509KeyPair("./shop.crt", "./shop.key")
	if err != nil {
		log.Fatal(err)
	}
	cert = shopCert

	http.HandleFunc("/", unknownHandler)
	log.Println("Running...")
	http.ListenAndServe(":80", logRequest(http.DefaultServeMux))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", aurora.Blue(r.Method), aurora.Red(r.URL))
		// TODO: remove header dumping
		for name, value := range r.Header {
			log.Print(name, " ", aurora.Green("=>"), " ", value)
		}
		handler.ServeHTTP(w, r)
	})
}

func unknownHandler(w http.ResponseWriter, r *http.Request) {
	// Proxy from Nintendo.

	// First, add cert to tls handshake...
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: tr,
	}
	// then, mirror current server path...
	response, err := client.Get(fmt.Sprintf("https://oss-auth.shop.wii.com%s", r.URL))
	if err != nil {
		panic(err)
	}
	// and disguise ourselves as the shop channel...
	response.Header.Add("User-Agent", "Opera/9.30 (Nintendo Wii; U; ; 2077-4; Wii Shop Channel/21.0(A); en)")

	defer response.Body.Close()
	// then, literally mirror back!
	_, err = io.Copy(w, response.Body)
	if err != nil {
		panic(err)
	}
}
