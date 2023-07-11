Output docs/docs/demo/tui-overview.gif
Output docs/docs/demo/tui-overview.webm

Set Padding 0

Set FontSize 16
Set Framerate 60

Set Width 1000
Set Height 600

# Start
Hide
Type@0s "SALT_DEMO=true ./salt-live -ipc-file ./e2e_test/ipc.ignore/master_event_pub.ipc -hard-filter='!ping_master'"
Enter
Sleep 5s
Show

# Admire the output for a bit
Sleep 3s

Down@1s 3

Sleep 1s
