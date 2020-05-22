package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	http.HandleFunc("/eval.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./eval.html")
	})
	log.Println("Running...")
	http.ListenAndServe("192.168.3.135:8080", logRequest(http.DefaultServeMux))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s via %s", aurora.Blue(r.Method), aurora.Red(r.URL), aurora.Cyan(r.Host))
		// TODO: remove header dumping
		//for name, value := range r.Header {
		//	log.Print(name, " ", aurora.Green("=>"), " ", value)
		//}
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

	// Given the official format of https://*.shop.wii.com, we want to turn that into our proxying domain.
	// If we've done our job right, however, inbound requests will already be to our proxy, so we need to rewrite.
	properDomain := strings.ReplaceAll(wiiRequest.Host, ProxyingDomain, OfficialDomain)

	properPath := wiiRequest.URL.String()
	url := fmt.Sprintf("https://%s%s", properDomain, properPath)
	log.Print("Proxying from ", aurora.Green(url))

	var shopRequest *http.Request
	var err error

	shouldLog := false
	// If this is a java servlet, we would definitely like to know :)
	// Unless you're the same-old check-in pages. Then we would not.
	if strings.Contains(wiiRequest.URL.Path, ".jsp") && !strings.Contains(wiiRequest.URL.Path, "Check") {
		shouldLog = true
	}

	// then, mirror current server path...
	if wiiRequest.Method == "GET" {
		shopRequest, err = http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
	} else if wiiRequest.Method == "POST" {
		// I am sorry! We're logging this request, no matter how long it takes.
		post, err := ioutil.ReadAll(wiiRequest.Body)
		if err != nil {
			panic(err)
		}
		// Log what's being posted.
		log.Println(string(post))
		// Put the original content back to a buffer to be consumed for our proxied request.
		finalized := bytes.NewBuffer(post)
		shopRequest, err = http.NewRequest("POST", url, finalized)

		// We'd like to see the result as well, please.
		shouldLog = true
	}

	// Set headers for shop request.
	// and disguise ourselves as the shop channel, passing on the rest of the headers...
	shopRequest.Header = wiiRequest.Header
	shopRequest.Header.Set("User-Agent", "Opera/9.30 (Nintendo Wii; U; ; 2077-4; Wii Shop Channel/21.0(A); en)")

	// Cache somehow messes up proxying. TODO: Investigate - is it because of a change applied?
	shopRequest.Header.Del("If-Modified-Since")
	shopRequest.Header.Del("If-None-Match")

	shopResponse, err := client.Do(shopRequest)
	if err != nil {
		panic(err)
	}

	defer shopResponse.Body.Close()
	shopBody, err := ioutil.ReadAll(shopResponse.Body)
	if err != nil {
		panic(err)
	}

	// Replace "base domains" in the document.
	shopBody = bytes.ReplaceAll(shopBody, []byte(OfficialDomain), []byte(ProxyingDomain))

	// Make life a little easier for "progress object" errors.
	shopBody = bytes.ReplaceAll(shopBody, []byte("var errInfo = progress.errInfo;"), []byte(`
	var errInfo = progress.errInfo;
	trace("status: " + progress.status);
	trace("error code: " + progress.errCode);
	trace("error info: " + progress.errInfo);
	`))

	// And inject our text area in as well.
	shopBody = bytes.ReplaceAll(shopBody, []byte("</body>"), []byte(`
		<script>
			var a = document.createElement('a');
			a.setAttribute("href", "/eval.html");
			var button = document.createElement('button');
			button.setAttribute("type", "button");
			button.innerHTML = "Eval";
			a.appendChild(button);
			document.body.appendChild(a);
		</script>
	</body>
	`))

	if shouldLog {
		log.Println(string(shopBody))
	}

	// then, literally mirror back!
	// + headers too!
	for name, values := range shopResponse.Header {
		if name == "Content-Length" {
			continue
		}
		for _, value := range values {
			wiiWriter.Header().Add(name, value)
		}
	}
	_, err = wiiWriter.Write(shopBody)
	if err != nil {
		panic(err)
	}
	log.Printf("Returned proxied error code %d", shopResponse.StatusCode)
}
