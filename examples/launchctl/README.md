# MacOS LaunchDaemon

If you're installing through a package manager, you probably don't need to deal
with this file.

The `plist` file should be put in `/Library/LaunchDaemons/` (user defined daemons), and the binary installed at
`/usr/local/bin/node_exporter`.

Ex. install globally by

    sudo cp -n node_exporter /usr/local/bin/
    sudo cp -n examples/launchctl/io.prometheus.node_exporter.plist /Library/LaunchDaemons/
    sudo launchctl bootstrap system/ /Library/LaunchDaemons/io.prometheus.node_exporter.plist

    # Optionally configure by dropping CLI arguments in a file
    echo -- '--web.listen-address=:9101' | sudo tee /usr/local/etc/node_exporter.args

    # Check it's running
    sudo launchctl list | grep node_exporter

    # See full process state
    sudo launchctl print system/io.prometheus.node_exporter

    # View logs
    sudo tail /tmp/node_exporter.log
