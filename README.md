# the heck is this
good question lol

# so what do???
if you really wanna hax0r that ninty server, here's a cool how2:
1. go to https://larsenv.github.io/NintendoCerts/index.html
2. download "Wii NWC Prod 1"
3. copypasta wherever you'd like, i put it in this repo (since ignored)
4.
```
openssl pkcs12 -in WII_NWC_1_CERT.p12 -out shop.crt -clcerts -nokeys
openssl pkcs12 -in WII_NWC_1_CERT.p12 -out shop.key -nocerts -nodes
```
5. use magic, edit wii shop dol to have `https://oss-auth.shop.wii.com/` -> `http://oss-auth.shop.wii.com/`
6. using similar magic, edit `00000002.app`'s U8 and edit the opera include file to have `http://*.shop.wii.com`
7. grab yourself a dns server, bind `oss-auth.shop.wii.com` to wherever you wanna run this
8. fire this up
9. ???
10. potentially quite literally, profit
