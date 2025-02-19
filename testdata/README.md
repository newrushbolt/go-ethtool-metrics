# Testdata

## Structure

```bash
testdata
├── README.md
└── intel # vendor name
    └── i40e # driver name, `ethtool -i eth0 | grep 'driver:'`
        └── 00_sfp_10g_sr85 # index and additional info, like module type
            ├── results # Expected outputs in JSON for different option sets (default, full, etc)
            │   ├── driver_info.default.json
            │   ├── driver_info.full.json
            │   ├── generic_info.default.json
            │   ├── module_info.default.json
            │   ├── module_info.full.json
            │   └── statistics.default.json
            └── src # Plain output of ethtool command
                ├── driver_info # `ethtool -i eth0`
                ├── generic_info # `ethtool eth0`
                ├── module_info # `ethtool -m eth0`
                └── statistics # `ethtool -S eth0`
```

## Adding new testdata

### Driver info, generic info

For `driver_info` and `generic_info` adding new data is pretty straightforward:

* create new folders with incremented index

  ```bash
  mkdir -p testdata/intel/i40e/02_sfp_10_or_25g_sr/src
  mkdir -p testdata/intel/i40e/02_sfp_10_or_25g_sr/results
  ```

* fill files in `src` with ethtool output

  ```bash
  ssh server-1.mycompany.local sudo ethtool -i eth0 2>/dev/null \
    > testdata/intel/i40e/02_sfp_10_or_25g_sr/src/driver_info
  
  ssh server-1.mycompany.local sudo ethtool eth0 2>/dev/null \
    > testdata/intel/i40e/02_sfp_10_or_25g_sr/src/generic_info
  ```

* add new testdata to `GetFixtureList()` in [metrics_test.go](../metrics_test.go)

  ```diff
  --- a/metrics_test.go
  +++ b/metrics_test.go
  @@ -38,6 +38,7 @@ func GetFixtureList() []string {
          return []string{
                  "intel/i40e/00_sfp_10g_sr85",
                  "intel/i40e/01_int_tp",
  +               "intel/i40e/02_sfp_10_or_25g_sr",
                  "intel/igb/00_int_tp",
                  "broadcom/bnxt_en/00_sfp_10gsr85",
          }
  ```

* create empty result-files in `testdata/intel/i40e/02_sfp_10_or_25g_sr/results` for each mode

  ```bash
  touch testdata/intel/i40e/02_sfp_10_or_25g_sr/results/module_info.default.json
  touch testdata/intel/i40e/02_sfp_10_or_25g_sr/results/driver_info.default.json
  touch testdata/intel/i40e/02_sfp_10_or_25g_sr/results/generic_info.default.json
  touch testdata/intel/i40e/02_sfp_10_or_25g_sr/results/module_info.full.json
  touch testdata/intel/i40e/02_sfp_10_or_25g_sr/results/driver_info.full.json
  touch testdata/intel/i40e/02_sfp_10_or_25g_sr/results/generic_info.full.json
  ```

* run related tests

  ```bash
  go test -v -test.run 'TestDriverInfo*'
  go test -v -test.run 'TestGenericInfo*'
  ...
    === RUN   TestDriverInfoDefault/intel/i40e/02_sfp_10_or_25g_sr
        metrics_test.go:90:
              Error Trace:  /Users/newrushbolt/projects/personal/go-ethtool-metrics/metrics_test.go:90
              Error:        Not equal:
              ...
                            Diff:
                            --- Expected
                            +++ Actual
                            @@ -1,12 +1 @@
                            -{
                            -    "DriverName": "ice",
                            -    "DriverVersion": "4.18.0-536.el8.x86_64",
                            -    "FirmwareVersion": "4.40 0x8001ba1d 22.5.7",
                            -    "FirmwareVersionParts": [
                            -        "4.40",
                            -        "0x8001ba1d",
                            -        "22.5.7"
                            -    ],
                            -    "BusAddress": "0000:22:00.3",
                            -    "Features": null
                            -}
                            +
              Test:         TestDriverInfoDefault/intel/i40e/02_sfp_10_or_25g_sr
  ...
  ```

* get parsed metrics from failed asserts, and fill the required json-files in `testdata/intel/i40e/02_sfp_10_or_25g_sr/results`.  
  Go one failed test at a time, until all tests of `TestGenericInfo*` and `TestDriverInfo*` pass.

* please check that `driver_info.firmware-version` correctly converted to array `FirmwareVersionParts`, eg:

  ```yaml
  firmware-version: 4.40 0x8001ba1d 22.5.7
  # ->
  "FirmwareVersionParts": [
      "4.40",
      "0x8001ba1d",
      "22.5.7"
  ]
  ```

  This helps to build better dashboards and alerts about firmware versions.

### Module info

This is where things gets harder. Although driver could announce that we should be able to get diagnostics from the module, it doesn't mean we will get any.

```bash
# Should be able to show some diagnostics data
ethtool -i eth0
...
supports-test: yes

# But shows no diagnostic information
ethtool -m eth0 | grep -i threshold
```

In my experience this happens because of incompatible driver or firmware versions,  
so you may need to get tuple kernel_version:driver_version:nic_firmware_version right.

After fixing that, and ensuring `ethtool -m eno2` shows you proper diagnostics, the steps are simple:

* fill files in `src` with ethtool output

  ```bash
  ssh server-1.mycompany.local sudo ethtool -m eth0 2>/dev/null \
    > testdata/intel/i40e/02_sfp_10_or_25g_sr/src/module_info
  ```

* run related tests

  ```bash
  go test -v -test.run 'TestDriverInfo*'
  ```

* get parsed metrics from failed asserts, and fill the required json-files in `testdata/intel/i40e/02_sfp_10_or_25g_sr/results`

* make sure related tests pass

### Statistics

This where it will require a little bit of data modification, especially if your interfaces haven't yet experienced a lot of error-like events.

* lets gather the metrics first

  ```bash
  ssh server-1.mycompany.local sudo ethtool -S eth0 2>/dev/null \
    > testdata/intel/i40e/02_sfp_10_or_25g_sr/src/statistics

  touch testdata/intel/i40e/02_sfp_10_or_25g_sr/results/statistics.default.json
  ```

* and have a look at the metrics

  ```bash
  $ go test -v -test.run 'TestStatistics*'
  === RUN   TestStatistics/intel/i40e/02_sfp_10_or_25g_sr
      metrics_test.go:120:
            Error:
                          Diff:
                          --- Expected
                          +++ Actual
                          @@ -1,12 +1 @@
                          -{
                          -    "General": {
                          -        "TxBytes": 24983496586400,
                          -        "RxBytes": 14645706696119,
                          -        "RxErrors": 0,
                          -        "TxErrors": 0,
                          -        "RxDiscards": 0,
                          -        "TxDiscards": 0,
                          -        "TxCollisions": 0,
                          -        "RxCrcErrors": 0
                          -    }
                          -}
                          +
            Test:         TestStatistics/intel/i40e/02_sfp_10_or_25g_sr
  ```

As you can see, all error-like metrics returned as zeroes. It looks good from the monitoring perspective, but in tests we cannot differentiate "real" zeroes from "no such metric found or parsed" zeroes.

So it would be a nice touch to:

* check if error metrics are presented in ethtool output, and have the same names, as stated in the [GeneralStatisticsInfo](pkg/metrics/statistics/statistics_structs.go).  
  If names don't match you may need to add new alias for this metric:

  ```diff
  --- a/pkg/metrics/statistics/statistics_structs.go
  +++ b/pkg/metrics/statistics/statistics_structs.go
  @@ -14,7 +14,7 @@ type GeneralStatisticsInfo struct {
          RxDiscards   uint64 `general_statistics:"rx_discards,veb.rx_discards,rx_stat_discard"`
          TxDiscards   uint64 `general_statistics:"tx_discards,veb.tx_discards,tx_stat_discard"`
          TxCollisions uint64 `general_statistics:"tx_collisions,tx_total_collisions,collisions"`
  -       RxCrcErrors  uint64 `general_statistics:"rx_crc_errors"` // Only exists in Intel
  +       RxCrcErrors  uint64 `general_statistics:"rx_crc_errors,rx_crc_errors.nic"` // Only exists in Intel
  }
  ```

* you may also manually set error-like metrics in ethtool output to non-zeroes, just to make sure they are parsed correctly
