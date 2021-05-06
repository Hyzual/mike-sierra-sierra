{ pkgs ? (import ./pinned-nixpkgs.nix) {} }:

pkgs.buildEnv {
    name = "go-build-tools";
    paths =
        (import ./go.nix { inherit pkgs; }) ++
        (import ./general-build-tools.nix { inherit pkgs; });
}
