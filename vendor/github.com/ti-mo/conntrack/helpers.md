# Conntrack Helpers

Support was planned originally for creating helper entries from userspace. After some experimentation,
creating expectations from userspace does not seem possible, even with upstream tools. (`conntrack`)
The library code to assemble and marshal expectations will be kept around in case this becomes possible
later on in the kernel. (probably not)

Expectations follow a specific pattern, and can be created as follows (simple example using FTP server)
w/ client in passive mode.

```
sudo modprobe nf_conntrack_ftp
echo "1" | sudo tee /proc/sys/net/netfilter/nf_conntrack_helper
sudo iptables -A INPUT -m conntrack --ctstate RELATED -m helper --helper ftp -d 127.0.0.1 -p tcp --dport 30000:30009 -j ACCEPT

docker run -d --rm -e FTP_USER_NAME=bob -e FTP_USER_PASS=12345 -e FTP_USER_HOME=/home/bob -p 21:21 -p 30000-30009:30000-30009 stilliard/pure-ftpd

ftp 127.0.0.1 21 -p
```

Log in with: `bob/12345` and send a file to the server.

This should yield records like follows:

```
[EventExpNew]
    Timeout: 300,
    Master: <tcp, Src: 127.0.0.1:42706, Dst: 127.0.0.1:21>,
    Tuple: <tcp, Src: 127.0.0.1:0, Dst: 127.0.0.1:30000>,
    Mask: <tcp, Src: 255.255.255.255:0, Dst: 255.255.255.255:65535>,
    Zone: {0 0}, Helper: ftp, Class: 0x30

[EventExpDestroy] Timeout: 300, Master: <tcp, Src: 127.0.0.1:42706, Dst: 127.0.0.1:21>, Tuple: <tcp, Src: 127.0.0.1:0, Dst: 127.0.0.1:30000>, Mask: <tcp, Src: 255.255.255.255:0, Dst: 255.255.255.255:65535>, Zone: {0 0}, Helper: ftp, Class: 0x30
```
