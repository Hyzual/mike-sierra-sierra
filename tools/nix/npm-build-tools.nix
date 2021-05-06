{ pkgs ? import (./pinned-nixpkgs.nix) {}}:

pkgs.buildEnv {
    name = "npm-build-tools";
    paths =
        (import ./npm.nix) { inherit pkgs; } ++
        (import ./general-build-tools.nix { inherit pkgs; });
}
