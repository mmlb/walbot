{
  description = "A slack bot; mostly for wasting time, but sometimes is useful";

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";

    devenv.url = "github:cachix/devenv/latest";
    devenv.inputs.nixpkgs.follows = "nixpkgs";
    devshell.inputs.flake-utils.follows = "flake-utils";
    devshell.inputs.nixpkgs.follows = "nixpkgs";
    devshell.url = "github:numtide/devshell";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    devenv,
    devshell,
    flake-utils,
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [devshell.overlays.default];
        };

        # Generate a user-friendly version number.
        version = builtins.substring 0 8 self.lastModifiedDate;
      in rec {
        # Provide some binary packages for selected system types.
        packages = flake-utils.lib.flattenTree {
          walbot = pkgs.buildGoModule {
            pname = "walbot";
            inherit version;
            # In 'nix develop', we don't need a copy of the source tree
            # in the Nix store.
            src = ./.;
            vendorSha256 = "sha256-ur2iBQayIBEdrEn4PvLOhuiEm9RFugVywNaXnYHYjZQ=";
          };
        };
        defaultPackage = packages.walbot;
        apps.walbot = flake-utils.lib.mkApp {drv = packages.walbot;};
        defaultApp = apps.walbot;
        devShell = pkgs.devshell.mkShell {
          motd = "";
          packages = [devenv.packages.${system}.devenv];
        };
      }
    );
}
