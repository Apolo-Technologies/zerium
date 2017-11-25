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
  SimpleFC::AdvRemoveRule "Gzrm incoming peers (TCP:32310)"
  SimpleFC::AdvRemoveRule "Gzrm outgoing peers (TCP:32310)"
  SimpleFC::AdvRemoveRule "Gzrm UDP discovery (UDP:32310)"

  # Remove IPC endpoint (https://github.com/abt/EIPs/issues/147)
  ${un.EnvVarUpdate} $0 "ZERIUM_SOCKET" "R" "HKLM" "\\.\pipe\gzrm.ipc"

  # Remove install directory from PATH
  Push "$INSTDIR"
  Call un.RemoveFromPath

  # Cleanup registry (deletes all sub keys)
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}"
SectionEnd
