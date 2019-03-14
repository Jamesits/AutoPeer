# AutoPeer

**This project is in development, usage might change every night. Use with caution.**

We automated peering process by pulling peer data from [PeeringDB](https://www.peeringdb.com/), so you are finally free from the boring and error-prone process of manually inputting IPs into your BGP daemon config. You only need to configure your peer's ASN and you are good to go.

[![Build Status](https://dev.azure.com/nekomimiswitch/General/_apis/build/status/AutoPeer?branchName=master)](https://dev.azure.com/nekomimiswitch/General/_build/latest?definitionId=43&branchName=master)

Prebuilt binaries for Linux (amd64) can be retrieved from the CI above, or use my personal [Fast Ring](https://releases.swineson.me/autopeer/autopeer-linux-amd64).

## BGP Daemon Support

* [Bird](https://bird.network.cz/) 1.6

## Usage

1. Create an `autopeer.toml` ([example](doc/examples/autopeer.toml))
2. Run `autopeer --config autopeer.toml --output /etc/bird`
3. Include the config: 
    * In `bird.conf`: `include "autopeer4.conf";`
    * In `bird6.conf`: `include "autopeer6.conf";`
4. `configure` all bird instances
5. Profit

## Thanks

This project is proudly sponsored by [xTom](https://xtom.com/).

![我大哥是Showfom.webp](doc/assets/my_brother.png)

Thank every [NekomimiRouter.com](https://nekomimirouter.com/) operators for their help during the development of this project.