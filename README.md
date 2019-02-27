# AutoPeer

[Bird](https://bird.network.cz/) 1.x config generator for BGP players.

[![Build Status](https://dev.azure.com/nekomimiswitch/General/_apis/build/status/AutoPeer?branchName=master)](https://dev.azure.com/nekomimiswitch/General/_build/latest?definitionId=43&branchName=master)

## Usage

1. Create an `autopeer.toml` ([example](doc/examples/autopeer.toml))
2. Run `autopeer --config autopeer.toml --output /etc/bird`
3. Include the config: 
  * In `bird.conf`: `include "autopeer4.conf"`
  * In `bird6.conf`: `include "autopeer6.conf"`

`configure` all bird instances and you are good to go.

## Thanks

This project is proudly sponsored by [xTom](https://xtom.com/).

![我大哥是Showfom.webp](doc/assets/my_brother.png)

Thank every [NekomimiRouter.com](https://nekomimirouter.com/) operators for their help during the development of this project.