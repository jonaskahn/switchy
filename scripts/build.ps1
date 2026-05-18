$INNO_SETUP_COMPILER = "C:/Program Files (x86)/Inno Setup 6/ISCC.exe"

# Always run from the project root (switchy/) regardless of invocation location
Set-Location (Split-Path -Parent $PSScriptRoot)

if (-not (Test-Path $INNO_SETUP_COMPILER)) {
    Write-Error "Inno Setup 6 compiler not found at: $INNO_SETUP_COMPILER"
    exit 1
}

# Ensure Go and the Go bin directory (where wails lives) are on PATH
$env:PATH = [System.Environment]::GetEnvironmentVariable("PATH", "Machine") + ";" +
            [System.Environment]::GetEnvironmentVariable("PATH", "User") + ";" +
            "$env:USERPROFILE\go\bin"

Write-Output "Building Switchy (Go + Wails)..."
wails build -platform windows/amd64
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Output "Building Installer..."
& $INNO_SETUP_COMPILER ./build/windows/installer/installer.iss
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
