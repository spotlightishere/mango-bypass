echo Downloading key....
wget http://yuzu.api.6100m.ga/WII_NWC_1_CERT.p12
echo Executing patching stage 1...
openssl pkcs12 -in WII_NWC_1_CERT.p12 -out shop.crt -clcerts -nokeys
echo Executing patching stage 2...
openssl pkcs12 -in WII_NWC_1_CERT.p12 -out shop.key -nocerts -nodes
echo Downloading CA....
wget http://yuzu.api.6100m.ga/dst.pem
echo Please specify the name of the repacked 00000002.app file.
read appfile
echo Please specify the name of the original 00000002.app file.
read ogfile
echo Generating patch...
diff -u $ogfile $appfile > tmp.patch
echo Cleaning up....
rm $appfile
echo Patching....
patch $ogfile tmp.patch
echo Patched!
sleep 1
echo Executing the Repacking WAD stage 1....
mono Sharpii.exe WAD -u WSC.wad WSC/ #Thanks KCR and Spotlight for the IOS-Patcher, I'm basing it off of it a bit.
echo Executing the Repacking WAD stage 2....
rm WSC/00000002.app
echo Executing the Repacking WAD stage 3....
cp 00000002.app WSC/00000002.app
echo Executing the Repacking WAD stage 4....
mkdir WAD
echo Executing the Repacking WAD stage 5....
mono Sharpii.exe WAD -p WSC/ ./WAD/WSC-Patched.wad
echo Patched version is now avaibile at the following location: $PWD/WAD/WSC-Patched.wad
sleep 1
echo Executing setup stage....
mkdir public_html
echo Running backup stage 1...
mkdir old
echo Running backup stage 2...
package main
echo Please specify the name of the Offical Domain you want to use. Make sure it is in quotation marks!
read offical
echo Please specify the name of the Proxying Domain you want to use. Make sure it is in quotation marks!
read proxying
cat <<EOT >> config.go.example
const (
	// OfficialDomain is the base subdomain this proxy needs to connect to.
	OfficialDomain = $offical
	// ProxyingDomain is the base subdomain the Wii is connecting to, to which we need to edit redirects to.
	ProxyingDomain = $proxying
)
EOT
mv config.go.example old/config.go.example
echo Installing config....
cp old/config.go.example $PWD/config.go
echo Done!
sleep 1
echo Exiting in 12...
sleep 1  # Waits 1 second(s).
echo Exiting in 11...
sleep 1  # Waits 1 second(s).
echo Exiting in 10...
sleep 1  # Waits 1 second(s).
echo Exiting in 9...
sleep 1  # Waits 1 second(s).
echo Exiting in 8...
sleep 1  # Waits 1 second(s).
echo Exiting in 7...
sleep 1  # Waits 1 second(s).
echo Exiting in 6...
sleep 1  # Waits 1 second(s).
echo Exiting in 5...
sleep 1  # Waits 1 second(s).
echo Exiting in 4...
sleep 1  # Waits 1 second(s).
echo Exiting in 3...
sleep 1  # Waits 1 second(s).
echo Exiting in 2...
sleep 1  # Waits 1 second(s).
echo Exiting in 1...
sleep 1  # Waits 1 second(s).
exit 0
