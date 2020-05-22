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
4. use magic, edit wii shop dol to have `https://oss-auth.shop.wii.com/` -> `http://oss-auth.shop.wii.com/`
5. using similar magic, edit `00000002.app`'s U8 and edit the opera include file to have `http://*.shop.wii.com`
6. grab yourself a dns server, bind `oss-auth.shop.wii.com` to wherever you wanna run this
7. copy `config.go.example` to `config.co`, using similar as above (`oss-auth.shop.wii.com`) if you're editing DNS, or if you've changed CAs your own domain
8. fire this up
9. ???
10. potentially quite literally, profit
