{
  description = "A slack bot; mostly for wasting time, but sometimes is useful";

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  inputs = {
    devenv.url = "github:cachix/devenv";
    devenv.inputs.nixpkgs.follows = "nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = inputs @ {
    self,
    devenv,
    flake-utils,
    nixpkgs,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};

        # Generate a user-friendly version number.
        version = builtins.substring 0 8 self.lastModifiedDate;
      in rec {
        # Provide some binary packages for selected system types.
        packages = flake-utils.lib.flattenTree {
          walbot = pkgs.buildGoModule {
            pname = "walbot";
            version = version;
            # In 'nix develop', we don't need a copy of the source tree
            # in the Nix store.
            src = ./.;
            vendorHash = "sha256-8HfaR3McfaELatGYf1vyQZGBZcuBSnQdahM63muWwPs=";
          };
        };
        defaultPackage = packages.walbot;
        apps.walbot = flake-utils.lib.mkApp {drv = packages.walbot;};
        defaultApp = apps.walbot;
        devShell = devenv.lib.mkShell {
          inherit inputs pkgs;
          modules = [./devenv.nix];
        };
      }
    );
}
