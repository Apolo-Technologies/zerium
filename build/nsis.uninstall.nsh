Section "Uninstall"
  # uninstall for all users
  setShellVarContext all

  # Delete (optionally) installed files
  {{range $}}Delete $INSTDIR\{{.}}
  {{end}}
  Delete $INSTDIR\uninstall.exe

  # Delete install directory
  rmDir $INSTDIR

  # Delete start menu launcher
  Delete "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk"
  Delete "$SMPROGRAMS\${APPNAME}\Attach.lnk"
  Delete "$SMPROGRAMS\${APPNAME}\Uninstall.lnk"
  rmDir "$SMPROGRAMS\${APPNAME}"

  # Firewall - remove rules if exists
  SimpleFC::AdvRemoveRule "Gabt incoming peers (TCP:30303)"
  SimpleFC::AdvRemoveRule "Gabt outgoing peers (TCP:30303)"
  SimpleFC::AdvRemoveRule "Gabt UDP discovery (UDP:30303)"

  # Remove IPC endpoint (https://github.com/apolo-technologies/EIPs/issues/147)
  ${un.EnvVarUpdate} $0 "ZERIUM_SOCKET" "R" "HKLM" "\\.\pipe\gabt.ipc"

  # Remove install directory from PATH
  Push "$INSTDIR"
  Call un.RemoveFromPath

  # Cleanup registry (deletes all sub keys)
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}"
SectionEnd
