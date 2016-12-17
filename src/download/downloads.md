Title: Download DapperDox
Description: Download the latest pre-built release of DapperDox
Keywords: Download, linux, mac, osx, windows, binary, tar, zip

# Download DapperDox

We provide the lastest release builds for the most common operating systems and architectures.
If yours is not listed here, then clone the source from [GitHub](http://github.com/dapperdox/dapperdox) and follow the [build instructions](#building-from-source).

## Precompiled releases

**1.0.2 (2016-12-16)**

| Filename | OS   | Arch | Size | Checksum |
| -------- | ---- | ---- | ---- | -------- |
[dapperdox-1.0.2.darwin-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.2/dapperdox-1.0.2.darwin-amd64.tgz) | darwin | amd64 | 3.8M | e3a36d8b708df22e8671ab9a83c71927e0408eb2d700ec9c29b40bd9492046a6 |
[dapperdox-1.0.2.linux-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.2/dapperdox-1.0.2.linux-amd64.tgz) | linux | amd64 | 3.8M | aef411e414ef28e2344b73f2555cb6391464106e288057cec31c0a4be5d7d689 |
[dapperdox-1.0.2.linux-arm.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.2/dapperdox-1.0.2.linux-arm.tgz) | linux | arm | 3.4M | d63ab609364dc41198c45330c12d7c8cd86e8604ed27b007afdee90fb7f05691 |
[dapperdox-1.0.2.linux-arm64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.2/dapperdox-1.0.2.linux-arm64.tgz) | linux | arm64 | 3.5M | 090c709a1ccbbf8cae330ddcde9b43be98db13308a5e70b3a1b7fd8a2ea8dc63 |
[dapperdox-1.0.2.linux-x86.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.0.2/dapperdox-1.0.2.linux-x86.tgz) | linux | x86 | 3.5M | 82bde928c83f0fb884dc3c29f8af7fdbf336eb0a4de6f6ee6f0ff1ae2543191a |

## Building from source

To build from source, clone the [GitHub repo](https://github.com/dapperdox/dapperdox):

```bash
> git clone https://github.com/dapperdox/dapperdox
```

Now build dapperdox (this assumes that you have your [golang](https://golang.org/doc/install) environment configured correctly):

```
go get && go build
```
