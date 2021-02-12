$streamdeckPiholeSrcDir = ".\cmd\streamdeck-pihole"

$installDir = "${env:APPDATA}\Elgato\StreamDeck\Plugins\com.craiggwilson.streamdeck.pihole.sdPlugin"
echo "Installing to $installDir"

go build -o "$installDir\streamdeck-pihole.exe" $streamdeckPiholeSrcDir
Copy-Item "$streamdeckPiholeSrcDir\*.json" $installDir
Copy-Item "$streamdeckPiholeSrcDir\*.html" $installDir
Copy-Item "$streamdeckPiholeSrcDir\*.css" $installDir
