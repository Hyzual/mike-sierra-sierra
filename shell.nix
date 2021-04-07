{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/4cb48cc25622334f17ec6b9bf56e83de0d521fb7.tar.gz") {} }:

pkgs.mkShell {
    buildInputs = [
        pkgs.gnumake
        pkgs.gitMinimal
        pkgs.nodejs-slim-14_x
        pkgs.nodePackages.npm
        pkgs.go
        pkgs.mkcert
    ];

    # Use the SSH client provided by the system (FHS only) to avoid issues with Fedora default settings
    GIT_SSH = if pkgs.lib.pathExists "/usr/bin/ssh" then "/usr/bin/ssh" else "ssh";
}
