qdisc [![Build Status](https://travis-ci.org/ema/qdisc.svg?branch=master)](https://travis-ci.org/ema/qdisc)
=====

Package `qdisc` allows to get queuing discipline information via netlink,
similarly to what `tc -s qdisc show` does.

Example usage
-------------

    package main

    import (
        "fmt"

        "github.com/ema/qdisc"
    )

    func main() {
        info, err := qdisc.Get()

        if err == nil {
            for _, msg := range info {
                fmt.Printf("%+v\n", msg)
            }
        }
    }
