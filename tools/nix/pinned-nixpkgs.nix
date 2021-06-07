{}:

let pinnedNixPkgs = import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/ffb7cfcfad15e3bff9d05336767e59ee6ee24cb6.tar.gz") {};
in pinnedNixPkgs
