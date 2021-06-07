{ pkgs ? (import ./pinned-nixpkgs.nix) {} }:

pkgs.buildEnv {
    name = "docker-build-tools";
    paths =
        (import ./goss.nix { inherit pkgs; }) ++
        (import ./general-build-tools.nix { inherit pkgs; });
}
