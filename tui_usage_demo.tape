Output docs/docs/demo/tui-usage.gif
Output docs/docs/demo/tui-usage.webm

Set Padding 0

Set FontSize 16
Set Framerate 60

Set Width 1000
Set Height 600

# Start
Hide
Type@0s "SALT_DEMO=true ./salt-live -ipc-file ./e2e_test/ipc.ignore/master_event_pub.ipc -hard-filter='!ping_master'"
Enter
Sleep 100ms
Type "$"
Sleep 100ms

Show
Type@50ms "Simply start salt-live"
Sleep 2s
Type "$"

# Admire the output for a bit
Sleep 6s

# Navigate
Type@60ms "$Use the arrows to navigate in the event list"
Sleep 3s
Type "$"
Sleep 1s
Down@1s 3
Sleep 1s
Up 1
Sleep 2s

Type@60ms "$While navigating, the list is frozen. Press 'f' to get the events in real time again."
Sleep 3s
Type "$"
Sleep 1s
Type "f"
Sleep 1s

# Switch side view output format
Type@60ms "$Press 'm' to change the output format of the event details on the right"
Sleep 3s
Type "$"
Hide
Down 1
Up 1
Show
Sleep 1s
Type@2s "mmmm"
Sleep 2s

# Filter
Type@60ms "$Press '/' to filter the events"
Sleep 3s
Type "$"
Sleep 1s
Type '/'
Type@10ms 'foo state.sls'
Sleep 2s
Enter

Sleep 5s
