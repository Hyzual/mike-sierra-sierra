{}:

let pinnedNixPkgs = import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/5dbd28d75410738ee7a948c7dec9f9cb5a41fa9d.tar.gz") {};
in pinnedNixPkgs
