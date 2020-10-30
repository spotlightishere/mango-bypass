# What's this?
It connects to the Wii Shop Channel for you, and proxies GET/POST requests. This allows you to modify what's passed through to a server you control.

# How do I set it up?
1. Using magic, edit `00000002.app`'s U8 and edit the opera include file to have `http://*.base.proxy.domain`
2. Find a way to have this domain route to your server. Good luck!
3. Edit config.go.example with the domain being similar to the domain in Step 1. (`base.proxy.domain`)
4. Run sudo chmod 777 auto.sh.
5. Run ./auto.sh and follow the prompts exactly.
6. fire this up
7. ???
8. potentially quite literally, profit
