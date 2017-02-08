Title: Download DapperDox
Description: Download the latest pre-built release of DapperDox
Keywords: Download, linux, mac, osx, windows, binary, tar, zip

# Download DapperDox

We provide the lastest release builds for the most common operating systems and architectures.
If yours is not listed here, then clone the source from [GitHub](http://github.com/dapperdox/dapperdox) and follow the [build instructions](#building-from-source).

## Precompiled releases

**1.1.0 (2017-01-08)**

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
