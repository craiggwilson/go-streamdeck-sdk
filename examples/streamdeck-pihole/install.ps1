$srcDir = $PSScriptRoot

$installDir = "${env:APPDATA}\Elgato\StreamDeck\Plugins\com.craiggwilson.streamdeck.pihole.sdPlugin"
echo "Installing to $installDir"

if (Test-Path $installDir)
{
    Remove-Item -Force -Recurse $installDir
}

go build -o "$installDir\streamdeck-pihole.exe" $srcDir
Copy-Item "$srcDir\*.json" $installDir
Copy-Item "$srcDir\*.html" $installDir
Copy-Item "$srcDir\*.css" $installDir
Copy-Item "$srcDir\images" "$installDir\images" -Recurse
