Title: Download DapperDox
Description: Download the latest pre-built release of DapperDox
Keywords: Download, linux, mac, osx, windows, binary, tar, zip

# Download DapperDox

We provide the lastest release builds for the most common operating systems and architectures.
If yours is not listed here, then clone the source from [GitHub](http://github.com/dapperdox/dapperdox) and follow the [build instructions](#building-from-source).

## Precompiled releases

**1.0.1 (2016-11-22)**

| Filename | OS   | Arch | Size | Checksum |
| -------- | ---- | ---- | ---- | -------- |
[dapperdox-1.0.1.darwin-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.1/dapperdox-1.0.1.darwin-amd64.tgz) | darwin | amd64 | 3.9M | bf1b5bc402e3d1fa68bece7368f23a0a8b6fa42dae80daa2e31bd854a5b8081b |
[dapperdox-1.0.1.linux-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.1/dapperdox-1.0.1.linux-amd64.tgz) | linux | amd64 | 3.9M | 1b189caa195932653b351fdeb2ed3d0a3005941418c798211915c0f2a5de78f1 |
[dapperdox-1.0.1.linux-arm.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.1/dapperdox-1.0.1.linux-arm.tgz) | linux | arm | 3.5M | 1250e273abc4dd7ee352b4e52cedc9e151114b0ae1f9ae967fac665d02805aea |
[dapperdox-1.0.1.linux-arm64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.1/dapperdox-1.0.1.linux-arm64.tgz) | linux | arm64 | 3.6M | f211911bd02ee1abe2e2959045779e18b75b72e8c1eca40c77afd084601d8ef0 |
[dapperdox-1.0.1.linux-x86.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.1/dapperdox-1.0.1.linux-x86.tgz) | linux | x86 | 3.7M | 6808b00a0711d1de0868ddf7ae47823839aed53d69b3f17787214f9921e7e98c |

## Building from source

To build from source, clone the [GitHub repo](https://github.com/dapperdox/dapperdox):

```bash
> git clone https://github.com/dapperdox/dapperdox
```

Now build dapperdox (this assumes that you have your [golang](https://golang.org/doc/install) environment configured correctly):

```
go get && go build
```
