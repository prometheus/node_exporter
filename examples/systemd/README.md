# Systemd Unit

If you are using distribution packages or the copr repository, you don't need to deal with these files!

The unit file in this directory is to be put into `/etc/systemd/system`.
It needs a user named `node_exporter`, whose shell should be `/sbin/nologin` and should not have any special privileges.
It needs a sysconfig file in `/etc/sysconfig/node_exporter`.
It needs a directory named `/var/lib/node_exporter/textfile_collector`, whose owner should be `node_exporter`:`node_exporter`.
A sample file can be found in `sysconfig.node_exporter`.
