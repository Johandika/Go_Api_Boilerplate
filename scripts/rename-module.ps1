param(
    [Parameter(Mandatory = $true)]
    [string]$OldModule,

    [Parameter(Mandatory = $true)]
    [string]$NewModule
)

$root = Split-Path -Parent $PSScriptRoot
$files = Get-ChildItem -LiteralPath $root -Recurse -File |
    Where-Object {
        $_.FullName -notmatch '\\.git\\' -and
        $_.FullName -notmatch '\\tmp\\' -and
        $_.FullName -notmatch '\\bin\\'
    }

foreach ($file in $files) {
    $content = Get-Content -LiteralPath $file.FullName -Raw
    if ($content -like "*$OldModule*") {
        $content.Replace($OldModule, $NewModule) | Set-Content -LiteralPath $file.FullName -NoNewline
    }
}

go mod edit -module $NewModule
go mod tidy
