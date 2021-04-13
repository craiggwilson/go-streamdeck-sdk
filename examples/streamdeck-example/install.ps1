$srcDir = $PSScriptRoot

$installDir = "${env:APPDATA}\Elgato\StreamDeck\Plugins\com.craiggwilson.streamdeck.example.sdPlugin"
echo "Installing to $installDir"

if (Test-Path $installDir)
{
    Remove-Item -Force -Recurse $installDir
}

go build -o "$installDir\go-streamdeck-sdk-example.exe" $srcDir
Copy-Item "$srcDir\*.json" $installDir
