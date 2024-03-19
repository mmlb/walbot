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

        packages = flake-utils.lib.flattenTree {
          walbot = pkgs.buildGoModule {
            pname = "walbot";
            version = builtins.substring 0 8 self.lastModifiedDate;
            # In 'nix develop', we don't need a copy of the source tree
            # in the Nix store.
            src = ./.;
            vendorHash = "sha256-8HfaR3McfaELatGYf1vyQZGBZcuBSnQdahM63muWwPs=";
          };
        };
      in {
        # Provide some binary packages for selected system types.
        defaultApp = flake-utils.lib.mkApp {drv = packages.walbot;};
        defaultPackage = packages.walbot;
        devShell = devenv.lib.mkShell {
          inherit inputs pkgs;
          modules = [./devenv.nix];
        };
      }
    );
}
