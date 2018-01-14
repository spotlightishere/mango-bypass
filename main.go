package main

import (
	"net/http"
	"log"

	"crypto/tls"
	"github.com/logrusorgru/aurora"
	"fmt"
	"io/ioutil"
	//"strings"
	"path"
	"os"
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
		//// TODO: remove header dumping
		//for name, value := range r.Header {
		//	log.Print(name, " ", aurora.Green("=>"), " ", value)
		//}
		log.Printf("Accessing from: %s", aurora.Blue(r.Host))
		handler.ServeHTTP(w, r)
	})
}

func unknownHandler(wiiWriter http.ResponseWriter, wiiRequest *http.Request) {
	// First, we need to check if stuff has been "cached".
	base := path.Ext(wiiRequest.URL.Path)
	filename := path.Base(wiiRequest.URL.Path)
	cachedPath := "cache/" + filename
	switch base {
	case ".png", ".gif", ".js", ".css":
		if _, err := os.Stat(cachedPath); os.IsNotExist(err) {
			// not cached! we'll get 'er next time
			log.Printf("Looks like %s isn't cached, requesting.", aurora.Blue(filename))
			break
		} else {
			// ok we're good lol
			coolFile, err := ioutil.ReadFile(cachedPath)
			if err != nil {
				panic(err)
			}

			wiiWriter.Write(coolFile)
			return
		}
	}

	// Not cached, so "proxy" from Nintendo.

	// First, add cert to tls handshake...
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_RC4_128_SHA,
			},
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: tr,
	}
	// then, mirror current server path...
	url := fmt.Sprintf("https://oss-auth.shop.wii.com%s", wiiRequest.URL)
	//url := "https://fl0.co/agent.php"
	log.Print("Mirroring ", aurora.Green(url))
	shopRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	// Set headers for shop request.
	// and disguise ourselves as the shop channel, passing on some headers...
	shopRequest.Header = wiiRequest.Header
	//shopRequest.Header.Add("Te", "deflate, gzip, chunked, identity, trailers")
	//shopRequest.Header.Add("Accept", "text/html, application/xml;q=0.9, application/xhtml+xml, image/png, image/jpeg, image/gif, image/x-xbitmap, */*;q=0.1")
	//shopRequest.Header.Add("Accept-Language", "en")
	//shopRequest.Header.Add("Accept-Charset", "iso-8859-1, utf-8, utf-16, *;q=0.1")
	//shopRequest.Header.Add("Accept-Encoding", "deflate, gzip, x-gzip, identity, *;q=0")
	//shopRequest.Header.Add("Connection", "Keep-Alive, TE")
	////shopRequest.Header.Add("Connection", "Keep-Alive")
	shopRequest.Header.Set("User-Agent", "Opera/9.30 (Nintendo Wii; U; ; 2077-4; Wii Shop Channel/21.0(A); en)")

	shopResponse, err := client.Do(shopRequest)
	if err != nil {
		panic(err)
	}

	defer shopResponse.Body.Close()
	shopBody, err := ioutil.ReadAll(shopResponse.Body)
	if err != nil {
		panic(err)
	}

	// Check what we need to do.
	switch base {
	case ".jsp":
		// grab our cool api
		api, err := ioutil.ReadFile("cool-api.js")
		if err != nil {
			panic(err)
		}

		injectedBody := shopBody
		injectedBody = append(injectedBody, []byte("<script>\n")...)
		injectedBody = append(injectedBody, api...)
		injectedBody = append(injectedBody, []byte("\n</script>")...)
		injectedBody = append(injectedBody, shopBody...)
		_, err = wiiWriter.Write(injectedBody)
		if err != nil {
			panic(err)
		}
		return
	case ".png", ".gif", ".js", ".css":
		// Cache them, by all means!
		err := ioutil.WriteFile(cachedPath, shopBody, 777)
		if err != nil {
			panic(err)
		}
		break
	}

	// then, literally mirror back!
	_, err = wiiWriter.Write(shopBody)
	if err != nil {
		panic(err)
	}
}
