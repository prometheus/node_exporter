# Contributing

Prometheus uses GitHub to manage reviews of pull requests.

* If you have a trivial fix or improvement, go ahead and create a pull request,
  addressing (with `@...`) the maintainer of this repository (see
  [MAINTAINERS.md](MAINTAINERS.md)) in the description of the pull request.

* If you plan to do something more involved, first discuss your ideas
  on our [mailing list](https://groups.google.com/forum/?fromgroups#!forum/prometheus-developers).
  This will avoid unnecessary work and surely give you and us a good deal
  of inspiration.

* Relevant coding style guidelines are the [Go Code Review
  Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments)
  and the _Formatting and style_ section of Peter Bourgon's [Go: Best
  Practices for Production
  Environments](http://peter.bourgon.org/go-in-production/#formatting-and-style).

* Sign your work to certify that your changes were created by yourself or you
  have the right to submit it under our license. Read
  https://developercertificate.org/ for all details and append your sign-off to
  every commit message like this:

        Signed-off-by: Random J Developer <example@example.com>


## Collector Implementation Guidelines

The Node Exporter is not a general monitoring agent. Its sole purpose is to
expose machine metrics, as oppose to service metrics, with the only exception
being the textfile collector.

The metrics should not get transformed in a way that is hardware specific and
would require maintaining any form of vendor based mappings or conditions. If
for example a proc file contains the magic number 42 as some identifier, the
Node Exporter should expose it as it is and not keep a mapping in code to make
this human readable. Instead, the textfile collector can be used to add a static
metric which can be joined with the metrics exposed by the exporter to get human
readable identifier.

A Collector may only read `/proc` or `/sys` files, use system calls or local
sockets to retrieve metrics. It may not require root privileges. Running
external commands is not allowed for performance and reliability reasons. Use a
dedicated exporter instead or gather the metrics via the textfile collector.

The Node Exporter tries to support the most common machine metrics. For more
exotic metrics, use the textfile collector or a dedicated Exporter.
