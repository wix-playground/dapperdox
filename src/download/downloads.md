Title: Download DapperDox
Description: Download the latest pre-built release of DapperDox
Keywords: Download, linux, mac, osx, windows, binary, tar, zip

# Download DapperDox

We provide the lastest release builds for the most common operating systems and architectures.
If yours is not listed here, then clone the source from [GitHub](http://github.com/dapperdox/dapperdox) and follow the [build instructions](#building-from-source).

## Precompiled releases

**1.1.1 (06-04-2017)**

| Filename | OS   | Arch | Size | Checksum |
| -------- | ---- | ---- | ---- | -------- |
[dapperdox-1.1.1.darwin-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.darwin-amd64.tgz) | darwin | amd64 | 4.0M | f80297c68efa43502c1e98e6eef508ffe9df91854ce93127e62781d9b0617919 |
[dapperdox-1.1.1.linux-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.linux-amd64.tgz) | linux | amd64 | 4.0M | 3e959b0d972bd4035a46a45810b3afc3faf64d6c7a8173aed4ba7dd1a7a1e846 |
[dapperdox-1.1.1.linux-arm.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.linux-arm.tgz) | linux | arm | 3.6M | 8ceaba3296a865c2b734534ffa19403c07b6b0548f8b1cf5d6f25c32ede522d5 |
[dapperdox-1.1.1.linux-arm64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.linux-arm64.tgz) | linux | arm64 | 3.7M | 7bfd70731fb1ad250872be854b99330f93a66c6b38b69b6b235a5dde5f341240 |
[dapperdox-1.1.1.linux-x86.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.linux-x86.tgz) | linux | x86 | 3.8M | 8576e869d66ffb7ce1bdfaf3e267fcc42471c2ace2095bcf63ae9b31c99b7b46 |
[dapperdox-1.1.1.windows-amd64.zip](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.windows-amd64.zip) | windows | amd64 | 4.0M | 8462c2a9d82747c43760d356a912480297a434a73d9aaecc3b7204c616a1c172 |
[dapperdox-1.1.1.windows-x86.zip](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.windows-x86.zip) | windows | x86 | 3.8M | 8052904d8eb209ae28e03b41b6cbdeb0bff0326351820dcc5a97d597efa2726f |

**1.1.0 (08-01-2017)**

| Filename | OS   | Arch | Size | Checksum |
| -------- | ---- | ---- | ---- | -------- |
[dapperdox-1.1.0.darwin-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.0/dapperdox-1.1.0.darwin-amd64.tgz) | darwin | amd64 | 3.8M | 7231b1743249263612e4fb9c4d5d7340d8acc875e249fd1177c275b14bb666f4 |
[dapperdox-1.1.0.linux-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.0/dapperdox-1.1.0.linux-amd64.tgz) | linux | amd64 | 3.8M | 27f714eeee2348a2f3be60828002f7de52469d20a97cdc2b8059b99b7e5eb62c |
[dapperdox-1.1.0.linux-arm.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.0/dapperdox-1.1.0.linux-arm.tgz) | linux | arm | 3.4M | bd3501632201e8af3d624a14210a4e654afd2130a47431eafbca71cd12b3084f |
[dapperdox-1.1.0.linux-arm64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.0/dapperdox-1.1.0.linux-arm64.tgz) | linux | arm64 | 3.5M | 7fac026649dd84f3f36dd112e04f1359119308e4e29deb501c892c5f843956e0 |
[dapperdox-1.1.0.linux-x86.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.0/dapperdox-1.1.0.linux-x86.tgz) | linux | x86 | 3.5M | 1f540f4f0b978d55a35d796f8f875576476de4d96b477054c086fd17337c8483 |

## Building from source

To build from source, clone the [GitHub repo](https://github.com/dapperdox/dapperdox):

```bash
> git clone https://github.com/dapperdox/dapperdox
```

Now build dapperdox (this assumes that you have your [golang](https://golang.org/doc/install) environment configured correctly):

```
go get && go build
```
