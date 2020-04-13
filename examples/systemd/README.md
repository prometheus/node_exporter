# Systemd Unit

If you are using distribution packages or the copr repository, you don't need to deal with these files!

The unit file in this directory is to be put into `/etc/systemd/system`.
It needs a user named `node_exporter`, whose shell should be `/sbin/nologin` and should not have any special privileges.
It needs a sysconfig file in `/etc/sysconfig/node_exporter`.
It needs a directory named `/var/lib/node_exporter/textfile_collector`, whose owner should be `node_exporter`:`node_exporter`.
A sample file can be found in `sysconfig.node_exporter`.

# Systemd User Unit

An example configuration file for running node_exporter as a user and connecting to the user systemd instance is given in `user/node_exporter.service`. To use it, copy the file to one of the user unit search paths, for example `/usr/lib/systemd/user` or  `~/.config/systemd/user`.

Edit the line `ExecStart` and add any additional options.

On Debian based systems, install the `dbus-user-session` package for headless systems or `dbus-x11` for desktop systems.

Reload the systemd configuration `systemctl --user daemon-reload`, enable the unit `systemctl --user enable node_exporter` and start it `systemctl --user start node_exporter`.