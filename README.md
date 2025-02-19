# go-ethtool-metrics

A library for parsing ethtool metrics. Check out the [complete exporter](https://github.com/newrushbolt/go-ethtool-exporter) built on top.

## Motivation and origins

This library is a logical extension of [my fork](https://github.com/newrushbolt/prometheus-ethtool-exporter) of [Showmax/prometheus-ethtool-exporter](https://github.com/Showmax/prometheus-ethtool-exporter), which "drew some inspiration from [adeverteuil/ethtool_exporter](https://github.com/adeverteuil/ethtool_exporter)".  
Kudos to all the persons contributed to them.  
My fork was pretty ugly, but it had one important difference: tests with data, collected from real NICs and SFPs.

Once I decided to make a new shiny go-ethtool-exporter, I realized that monitoring formats come and go,  
and it's better to build a flexible parsing **library**, still having the test data, and build exporter on top of that.

So when time time comes the new monitoring format will kick in (openTelemetry or whatever),  
someone could just slap on 50-lines of monitoring agent on top of the library.  
Hope Golang will still be a thing :)

That's said, I clearly understand that monitoring hardware and especially the NIC\SFP metrics is VERY niche.  
A few companies host their hardware on premise nowadays, and even fewer cares about anything lower than TCP metrics in their OS.  
The library is build with the background of hosting private cloud for internal customers, with the following default settings.

This library does not have any compatibility with the latter exporters in terms of metric naming.  
This is because those exporters were mostly transparent wrappers around the ethtool itself.  

But ethtool does not really provides any naming standarts for metrics, and is loosely coupled with the kernel drivers ([example](https://github.com/torvalds/linux/blob/v5.19/drivers/net/ethernet/intel/i40e/i40e_ethtool.c)), so metrics can differ very much even within one vendor.  

So you can pretty much say that this library aims to **give some abstraction level over various vendors and models**, allowing you to be more vendor-agnostic.

## Adding new testdata

[Here](testdata/README.md)
