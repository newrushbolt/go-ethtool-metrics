# Changelog

## [0.0.9](https://github.com/newrushbolt/go-ethtool-metrics/compare/v0.0.8...v0.0.9) (2025-08-26)


### Bug Fixes

* add google gve virtual nic support ([c3f1268](https://github.com/newrushbolt/go-ethtool-metrics/commit/c3f126819bc0959514c23637a957401e85c20157))
* add per-queue drops ([79fd09a](https://github.com/newrushbolt/go-ethtool-metrics/commit/79fd09a8b2616ac6144492f64e34e3fe9f390d35))
* add virtio support ([009937f](https://github.com/newrushbolt/go-ethtool-metrics/commit/009937f579fd450d5f5de0545bac1fa5e96bc52b))
* better per-queue regexps ([c2eeaaf](https://github.com/newrushbolt/go-ethtool-metrics/commit/c2eeaafe62b4262b2f64df0551a99fd17a56dbad))
* better testing and logging for queue parse errors ([882ff73](https://github.com/newrushbolt/go-ethtool-metrics/commit/882ff7317237d6c466ae3f7cd1ad0f2a7f36b700))
* calculate total per-queue bytes for broadcom ([9c84000](https://github.com/newrushbolt/go-ethtool-metrics/commit/9c84000a5e0736e721630447958639d09cc532d1))
* pre-compile queue regexps once for a package ([f289ac8](https://github.com/newrushbolt/go-ethtool-metrics/commit/f289ac821a635bdd219ed63ba08de4c86151ce0f))
* separate dropped packets from discarded ones ([c7b58fe](https://github.com/newrushbolt/go-ethtool-metrics/commit/c7b58fe249a5ca38b587c9bbd8367ca3a9e73a5c))
* split per-queue metrics by groups ([2bbd036](https://github.com/newrushbolt/go-ethtool-metrics/commit/2bbd036090bf8d163f49169ba28f21ef3f2b4587))

## [0.0.8](https://github.com/newrushbolt/go-ethtool-metrics/compare/v0.0.7...v0.0.8) (2025-08-08)


### Bug Fixes

* filter out well-known empty values for strings and slices ([6cc1010](https://github.com/newrushbolt/go-ethtool-metrics/commit/6cc10105a989053d4d0b46b02f44656c94004296))

## [0.0.7](https://github.com/newrushbolt/go-ethtool-metrics/compare/v0.0.6...v0.0.7) (2025-08-06)


### Bug Fixes

* dont panic on empty expected result ([277ceb6](https://github.com/newrushbolt/go-ethtool-metrics/commit/277ceb62f6ce705ab81566cbb1a03b2daa67584d))

## [0.0.6](https://github.com/newrushbolt/go-ethtool-metrics/compare/v0.0.5...v0.0.6) (2025-07-31)


### Bug Fixes

* correct speed units ([f009539](https://github.com/newrushbolt/go-ethtool-metrics/commit/f00953928820d0157a71aae0dae00511460c0f32))

## [0.0.5](https://github.com/newrushbolt/go-ethtool-metrics/compare/v0.0.4...v0.0.5) (2025-07-23)


### Bug Fixes

* add coverage threshold, better parsers for simple types ([95ea966](https://github.com/newrushbolt/go-ethtool-metrics/commit/95ea966454a29951fee3e220aee7a61fa43bef45))
* better collect options for driver_info ([2fae2ca](https://github.com/newrushbolt/go-ethtool-metrics/commit/2fae2ca09e8e8ac56373e4534d379e049f3014a0))
* better slices parse ([c7a5c3d](https://github.com/newrushbolt/go-ethtool-metrics/commit/c7a5c3d090b684d0ab8117eb2880693530b4044a))
* revert to global logger ([9c08c1e](https://github.com/newrushbolt/go-ethtool-metrics/commit/9c08c1e51c94f6764c967e9bfd4ffc7caa52dcb2))


### Miscellaneous Chores

* test empty driver_info ([608da33](https://github.com/newrushbolt/go-ethtool-metrics/commit/608da33eeeefc0155d4263fdc2205398bbc059ef))

## [0.0.4](https://github.com/newrushbolt/go-ethtool-metrics/compare/v0.0.3...v0.0.4) (2025-07-10)


### Bug Fixes

* use separate-logger ([c978b08](https://github.com/newrushbolt/go-ethtool-metrics/commit/c978b081161cbc69b812eb5affbf6a1e87553f3b))

## [0.0.3](https://github.com/newrushbolt/go-ethtool-metrics/compare/v0.0.2...v0.0.3) (2025-07-05)


### Bug Fixes

* minor grammar fixes ([3299ef0](https://github.com/newrushbolt/go-ethtool-metrics/commit/3299ef06d95883ef5e837af8e4f6faa664560120))
* support per-queue metrics ([7f876a8](https://github.com/newrushbolt/go-ethtool-metrics/commit/7f876a8570e4f882744589294a054b183e1773b4))
* test coverage ([7bca81e](https://github.com/newrushbolt/go-ethtool-metrics/commit/7bca81e3e435f7c34fb48ee546142507d3570ddf))
* use `nil` values for missing metrics ([8ab1c20](https://github.com/newrushbolt/go-ethtool-metrics/commit/8ab1c209b45872d2a3bc254d50b37590362a0faa))

## [0.0.2](https://github.com/newrushbolt/go-ethtool-metrics/compare/v0.0.1...v0.0.2) (2025-03-23)


### Bug Fixes

* better config naming ([d1ca28e](https://github.com/newrushbolt/go-ethtool-metrics/commit/d1ca28eb19e3e803cdc1569f58db0ed05c35401b))
