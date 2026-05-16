; Switchy Installer — Inno Setup 6
; Run from the switchy\ directory: ISCC.exe build\windows\installer\installer.iss
; or via scripts\build.ps1 at the switchy root.

#define MyAppName "Switchy"
#define NameSmall "switchy"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "ifelse.one"
#define MyAppURL "https://github.com/jonaskahn/switchy"
#define ExeName "switchy.exe"
#define AppDescription "Choose which browser opens your links"
#define ProgID "SwitchyHTML"

[Setup]
AppId={{56C63D05-9D83-492A-ABDD-618FE36ACBFB}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
VersionInfoVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}/issues
AppUpdatesURL={#MyAppURL}/releases
AppReadmeFile={#MyAppURL}#readme
SetupIconFile=..\icon.ico
LicenseFile=..\..\..\LICENSE

UsePreviousAppDir=yes
DefaultDirName={autopf64}\{#MyAppName}
DisableProgramGroupPage=yes
PrivilegesRequired=admin
PrivilegesRequiredOverridesAllowed=dialog

OutputDir=..\..\..\build\bin
OutputBaseFilename=Switchy_Installer

ArchitecturesAllowed=x64compatible
Compression=lzma2
SolidCompression=yes

WizardStyle=modern
SetupLogging=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"
Name: "registerdefault"; Description: "Register Switchy as a default browser candidate (opens Windows Default Apps)"; GroupDescription: "Optional post-install"
Name: "freshstart"; Description: "Fresh start (remove all user data and configuration)"; GroupDescription: "Additional settings"; Flags: unchecked; Check: HasExistingData

[CustomMessages]
UninstallRemoveConfig=Remove user data and configuration files
UninstallConfigDescription=This will delete your browser profiles, rulesets, and preferences.
LaunchProgram=Launch Switchy

[InstallDelete]

[Files]
Source: "..\..\..\build\bin\{#ExeName}"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\{#MyAppName}\{#MyAppName}";          Filename: "{app}\{#ExeName}"; Comment: "{#AppDescription}"
Name: "{autoprograms}\{#MyAppName}\{#MyAppName} Settings"; Filename: "{app}\{#ExeName}"; Parameters: "--settings"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#ExeName}"; Comment: "{#AppDescription}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#ExeName}"; Parameters: "--register"; \
  StatusMsg: "Registering as default browser…"; Tasks: registerdefault; Flags: runhidden waituntilterminated
Filename: "{app}\{#ExeName}"; Description: "{cm:LaunchProgram}"; Flags: nowait postinstall skipifsilent

[UninstallRun]
Filename: "{app}\{#ExeName}"; Parameters: "--unregister"; Flags: runhidden waituntilterminated

[Code]
var
  RemoveConfig: Boolean;

function HasExistingData(): Boolean;
var
  ConfigPath: string;
begin
  ConfigPath := ExpandConstant('{userappdata}\Switchy');
  Result := DirExists(ConfigPath);
end;

function InitializeUninstall(): Boolean;
begin
  if MsgBox('Do you want to remove all user data and configuration files?' + #13 + #10 + #13 + #10 + 'This will delete your browser profiles, rulesets, and preferences.', mbConfirmation, MB_YESNO) = idYes then
    RemoveConfig := True
  else
    RemoveConfig := False;
  Result := True;
end;

procedure CurStepChanged(CurStep: TSetupStep);
var
  ConfigPath: string;
  InstallPath: string;
begin
  if CurStep = ssInstall then
  begin
    InstallPath := ExpandConstant('{app}');
    if DirExists(InstallPath) then
    begin
      DelTree(InstallPath, True, True, True);
    end;
    
    if WizardIsTaskSelected('freshstart') then
    begin
      ConfigPath := ExpandConstant('{userappdata}\Switchy');
      if DirExists(ConfigPath) then
      begin
        DelTree(ConfigPath, True, True, True);
      end;
    end;
  end;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
var
  ConfigPath: string;
  InstallPath: string;
begin
  if CurUninstallStep = usUninstall then
  begin
    if RemoveConfig then
    begin
      ConfigPath := ExpandConstant('{userappdata}\Switchy');
      if DirExists(ConfigPath) then
      begin
        DelTree(ConfigPath, True, True, True);
      end;
    end;
  end;
  if CurUninstallStep = usPostUninstall then
  begin
    InstallPath := ExpandConstant('{app}');
    if DirExists(InstallPath) then
    begin
      DelTree(InstallPath, True, True, True);
    end;
  end;
end;

[Registry]
; ── Browser client registration (HKLM) ──────────────────────────────────────
Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}"; \
  ValueType: string; ValueName: ""; ValueData: "{#MyAppName}"; Flags: uninsdeletekey

Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\Capabilities"; \
  ValueType: string; ValueName: "ApplicationName"; ValueData: "{#MyAppName}"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\Capabilities"; \
  ValueType: string; ValueName: "ApplicationDescription"; ValueData: "{#AppDescription}"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\Capabilities"; \
  ValueType: string; ValueName: "ApplicationIcon"; ValueData: "{app}\{#ExeName},0"; Flags: uninsdeletekey

Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\Capabilities\URLAssociations"; \
  ValueType: string; ValueName: "http";  ValueData: "{#ProgID}"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\Capabilities\URLAssociations"; \
  ValueType: string; ValueName: "https"; ValueData: "{#ProgID}"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\Capabilities\FileAssociations"; \
  ValueType: string; ValueName: ".html"; ValueData: "{#ProgID}"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\Capabilities\FileAssociations"; \
  ValueType: string; ValueName: ".htm";  ValueData: "{#ProgID}"; Flags: uninsdeletekey

Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\DefaultIcon"; \
  ValueType: string; ValueName: ""; ValueData: "{app}\{#ExeName},0"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\shell\open\command"; \
  ValueType: string; ValueName: ""; ValueData: """{app}\{#ExeName}"" ""%1"""; Flags: uninsdeletekey

; ── ProgID (HKLM) ────────────────────────────────────────────────────────────
Root: HKLM; Subkey: "SOFTWARE\Classes\{#ProgID}"; \
  ValueType: string; ValueName: ""; ValueData: "{#MyAppName} URL"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\Classes\{#ProgID}\DefaultIcon"; \
  ValueType: string; ValueName: ""; ValueData: "{app}\{#ExeName},0"
Root: HKLM; Subkey: "SOFTWARE\Classes\{#ProgID}\shell\open\command"; \
  ValueType: string; ValueName: ""; ValueData: """{app}\{#ExeName}"" ""%1"""; Flags: uninsdeletekey

Root: HKLM; Subkey: "SOFTWARE\Classes\.html\OpenWithProgids"; \
  ValueType: string; ValueName: "{#ProgID}"; ValueData: ""; Flags: uninsdeletevalue
Root: HKLM; Subkey: "SOFTWARE\Classes\.htm\OpenWithProgids"; \
  ValueType: string; ValueName: "{#ProgID}"; ValueData: ""; Flags: uninsdeletevalue

; ── RegisteredApplications ───────────────────────────────────────────────────
Root: HKLM; Subkey: "SOFTWARE\RegisteredApplications"; \
  ValueType: string; ValueName: "{#MyAppName}"; \
  ValueData: "SOFTWARE\Clients\StartMenuInternet\{#MyAppName}\Capabilities"; \
  Flags: uninsdeletevalue

