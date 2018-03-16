# MacOS LaunchAgent

If you're installing through a package manager, you probably don't need to deal
with this file.

The `plist` file should be put in `~/Library/LaunchAgents/` (user-install) or
`/Library/LaunchAgents/` (global install), and the binary installed at
`/usr/local/bin/node_exporter`.

Ex. install globally by

    sudo cp -n node_exporter /usr/local/bin/
    sudo cp -n examples/launchctl/node_exporter.plist /Library/LaunchAgents/
    sudo launchctl bootstrap system/ /Library/LaunchAgents/node_exporter.plist
    sudo launchctl start node_exporter

    # Check it's running
    sudo launchctl list | grep node_exporter

    # See full process state
    sudo launchctl print system/node_exporter

    # View logs
    sudo tail /usr/local/var/logs/node_exporter/output.log
