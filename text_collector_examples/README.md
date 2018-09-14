# Text collector example scripts

These scripts are examples to be used with the Node Exporter Textfile
Collector.

To use these scripts, we recommend using a `sponge` to atomically write the output.

   <collector_script> | sponge <output_file>

Sponge comes from [moreutils](https://joeyh.name/code/moreutils/)
* [brew install moreutils](http://brewformulas.org/Moreutil)
* [apt install moreutils](https://packages.debian.org/search?keywords=moreutils)
* [pkg install moreutils](https://www.freshports.org/sysutils/moreutils/)        

For more information see:
https://github.com/prometheus/node_exporter#textfile-collector
