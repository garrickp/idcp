# TODO

## idcp

- Add option to create parent directories if they do not already exist.
- Make mode and ownership changes idempotant.
- Add the option for recursive directory copying.

## idcmd

Run a provided command, with a few "idempotant" trigger options, such as
"creates/removes file X", or "if this other command return 0/1"

## idget

Tool to download a file from a given URL, using and respecting
'If-Modified-Since'. Should offer both custom header specification, cookies,
and basic auth. Support hashes for uniquely identifying the content of the
file.

## idonchange

Accepts as input the output of other idtools, and exits '0' if any input
includes a change, and '1' if not.

Note: This should be similar to simply running `egrep 'change +true' && idtool` on the
output.

## idpkg

Ensure a package is installed. Support specification of the repo as part of the
arguments, including any signing keys.

## idrm

Ensure a file does not exist.

## idsvc

Abstract away the thirty ways to manage a service.

## idtar

Idempotantly untar a provided file to a provided location. Support trimming
directories from the tarfile (to support tgzs like those provided by node).
Support recursive ownership and mode setting.

## idtouch

Do a touch on a file, and update its ownership and mode if needed.

## idtemplate

Take an templated input file, and a key/value pair file, and populate the
template file into an output file.

## idlink

Create hardlinks and symlinks.

## man pages

Create man pages for all tools.

## make pkg

Set up package creation for deb, rpm, and tgz.

## idbtget

Use the Bittorrent protocol to download a file.  Input should be a .torrent
file. Exits after all pieces retrieved and in-progress pieces are sent. Should
also serve to peers while downloading. Should support the "WebSeeding"
protocol, and optionally DHT.

https://wiki.theory.org/BitTorrentSpecification

## idbtserve

Use the Bittorrent protocol to seed a given .torrent. Optionally support "Fast Peers" extension.

