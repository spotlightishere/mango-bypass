# What's this?
It connects to the Wii Shop Channel for you, and proxies GET/POST requests. This allows you to modify what's passed through to a server you control.

# How do I set it up?
1. Do to https://larsenv.github.io/NintendoCerts/index.html
2. Download "Wii NWC Prod 1" wherever you'd like, i put it in this repo (since ignored)
3.
```
openssl pkcs12 -in WII_NWC_1_CERT.p12 -out shop.crt -clcerts -nokeys
openssl pkcs12 -in WII_NWC_1_CERT.p12 -out shop.key -nocerts -nodes
```
4. Create your own CA, or import the DST Root CA X3 from Opera 9 on Windows.
5. Using similar magic, edit `00000002.app`'s U8 and edit the opera include file to have `http://*.base.proxy.domain`
6. Find a way to have this domain route to this server. Good luck!
7. Copy `config.go.example` to `config.go`, using similar as above (`base.proxy.domain`) as the base proxy domain
8. fire this up
9. ???
10. potentially quite literally, profit
