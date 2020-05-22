package main

import (
	"bytes"
	"log"
	"net/http"

	"crypto/tls"
	"fmt"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
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
	http.HandleFunc("/verysecrethandler", secretHandler)
	log.Println("Running...")
	http.ListenAndServe("192.168.3.135:8080", logRequest(http.DefaultServeMux))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", aurora.Blue(r.Method), aurora.Red(r.URL))
		// TODO: remove header dumping
		//for name, value := range r.Header {
		//	log.Print(name, " ", aurora.Green("=>"), " ", value)
		//}
		//log.Printf("Accessing from: %s", aurora.Blue(r.Host))
		handler.ServeHTTP(w, r)
	})
}

func unknownHandler(wiiWriter http.ResponseWriter, wiiRequest *http.Request) {
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

	url := fmt.Sprintf("https://%s%s", ProxyingDomain, wiiRequest.URL)
	log.Print("Mirroring ", aurora.Green(url))

	var shopRequest *http.Request
	var err error

	// then, mirror current server path...
	if wiiRequest.Method == "GET" {
		shopRequest, err = http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
	} else if wiiRequest.Method == "POST" {
		shopRequest, err = http.NewRequest("POST", url, wiiRequest.Body)
	}

	// Set headers for shop request.
	// and disguise ourselves as the shop channel, passing on some headers...
	shopRequest.Header = wiiRequest.Header
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

	if bytes.Contains(shopBody, []byte(ProxyingDomain)) {
		shopBody = bytes.ReplaceAll(shopBody, []byte(ProxyingDomain), []byte(ConnectingDomain))
	}

	if bytes.Contains(shopBody, []byte("</body>")) {
		shopBody = bytes.ReplaceAll(shopBody, []byte("</body>"), []byte(`
	<script>
		var elem = document.createElement('textarea');
		elem.innerHTML = ec.getLog();
		elem.setAttribute("rows", "15");
		elem.setAttribute("cols", "75");
		document.body.appendChild(elem);
	</script>
</body>
`))
	}

	// then, literally mirror back!
	_, err = wiiWriter.Write(shopBody)
	if err != nil {
		panic(err)
	}
}

func secretHandler(wiiWriter http.ResponseWriter, wiiRequest *http.Request) {

}