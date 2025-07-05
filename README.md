# go-ethtool-metrics

A library for parsing ethtool metrics.  
Check out the [complete prometheus exporter](https://github.com/newrushbolt/go-ethtool-exporter) built on top.

## Target of the project

This library aims to provide some level of abstraction over ethtool metrics over different NIC drivers.  
And for vendor\NIC-specific metrics, better naming and well-defined units.  
Example:

```text
// From this
Vendor name                               : FS
Vendor OUI                                : 00:1b:21
Vendor PN                                 : SFP-10GSR-85
Vendor rev                                : A
Vendor SN                                 : F2020533679
Date code                                 : 200715
Laser bias current                        : 6.090 mA
Laser output power                        : 0.6367 mW / -1.96 dBm
Receiver signal average optical power     : 0.6486 mW / -1.88 dBm
Module temperature                        : 39.18 degrees C / 102.52 degrees F
Module voltage                            : 3.3544 V
```

```json
// To this
"Vendor": {
    "Name": "FS",
    "OUI": "00:1b:21",
    "PartNumber": "SFP-10GSR-85",
    "Revision": "A",
    "SerialNumber": "F2020533679"
},
"Diagnostics": {
    "Values": {
        "BiasMilliAmps": 6.09,
        "OutputPowerMilliWatts": 0.6367,
        "InputPowerMilliWatts": 0.6486,
        "TemperatureCelsius": 39.18,
        "Voltage": 3.3544
    }
}
```

## Origins

This library is a logical extension of [my fork](https://github.com/newrushbolt/prometheus-ethtool-exporter) of [Showmax/prometheus-ethtool-exporter](https://github.com/Showmax/prometheus-ethtool-exporter), which "drew some inspiration from [adeverteuil/ethtool_exporter](https://github.com/adeverteuil/ethtool_exporter)". Kudos to all the persons contributed to them.  

## Motivation for "library instead of exporter"

[My fork](https://github.com/newrushbolt/prometheus-ethtool-exporter) was pretty ugly, but it had one important difference: tests with data, collected from real NICs and SFPs.

Once I decided to make a new shiny go-ethtool-exporter, I realized that monitoring formats come and go,  
and it's better to build a flexible parsing **library**, preserving and extending test data, and later build exporter on top of that.

So when the time comes for the new monitoring format (openTelemetry or whatever), someone could just slap on 50-lines of monitoring agent on top of the library.  

That's said, I clearly understand that monitoring hardware, especially the NIC\SFP metrics is VERY niche.  
A few companies host their hardware on premise nowadays, and even fewer cares about those kind of metrics.  
The library is built with the background of hosting private cloud for internal customers, with the following default settings.

This library does not have any compatibility with the [latter exporters](#origins) in terms of metric naming, because those exporters were mostly transparent wrappers around the ethtool.  
However, ethtool itself is mostly a transparent wrapper around metrics from ethernet drivers, especially for the "statistics" part (you can have a look at the driver's code, [i40e for example](https://github.com/torvalds/linux/blob/v5.19/drivers/net/ethernet/intel/i40e/i40e_ethtool.c)), and doesn't seem to provide any standards for metrics.

Overall, transparent metrics are harder to work with in big and diverse environments, and why I decided to create this library and the following exporter.

## Adding new testdata

[Here](testdata/README.md)

## Missing metrics

If metric is not found, then it is equal `*float64 nil`. In JSON testdata this translates to `null` value, which is easy to use.  
It's possible to use `float64 Nan` instead, but [JSON doesn't not support Nan](https://github.com/golang/go/blob/34c8b14ca9f4096383d658fbd748322a993a2bd2/src/encoding/json/encode.go#L47), and that would require custom JSON encoder for testdata.

Suggested behavior for exporter is to drop all metrics with `*float64 nil`, however, they could also be converted to Nan (which is supported in Prometheus, for example) and preserved if needed.

## Logging

Since library is split to modules, there's no centralized logger.  
However, every module does `slog.SetLogLoggerLevel(internal.GetLogLevel())` on start, and you can set preferred logging level using `GO_ETHTOOL_METRICS_LOG_LEVEL` env.  
Check the possible values [in slog library](https://cs.opensource.google/go/go/+/refs/tags/go1.24.3:src/log/slog/level.go;drc=6d7760cb4292bd156193fb2a3e377f24637f8795;l=67).
