{ pkgs ? import (./tools/nix/pinned-nixpkgs.nix) {} }:

pkgs.mkShell {
    buildInputs =
        (import ./tools/nix/general-build-tools.nix { inherit pkgs; }) ++
        (import ./tools/nix/go.nix { inherit pkgs; }) ++
        (import ./tools/nix/npm.nix { inherit pkgs; }) ++
        (import ./tools/nix/goss.nix { inherit pkgs; }) ++
        (import ./tools/nix/dev-tools.nix { inherit pkgs; });

    # Use the SSH client provided by the system (FHS only) to avoid issues with Fedora default settings
    GIT_SSH = if pkgs.lib.pathExists "/usr/bin/ssh" then "/usr/bin/ssh" else "ssh";
}
