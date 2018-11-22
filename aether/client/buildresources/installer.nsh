!macro customInstall
  DetailPrint "Register Aether URI Handler"
  DeleteRegKey HKCR "aether"
  WriteRegStr HKCR "aether" "" "URL:aether"
  WriteRegStr HKCR "aether" "URL Protocol" ""
  WriteRegStr HKCR "aether\DefaultIcon" "" "$INSTDIR\${APP_EXECUTABLE_FILENAME}"
  WriteRegStr HKCR "aether\shell" "" ""
  WriteRegStr HKCR "aether\shell\Open" "" ""
  WriteRegStr HKCR "aether\shell\Open\command" "" "$INSTDIR\${APP_EXECUTABLE_FILENAME} %1"
!macroend