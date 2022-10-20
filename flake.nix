{
  description = "A slack bot; mostly for wasting time, but sometimes is useful";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";

    #devshell.inputs.flake-utils.follows = "flake-utils";
    #devshell.inputs.nixpkgs.follows = "nixpkgs";
    #devshell.url = "github:numtide/devshell";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        /*
        pkgs = import nixpkgs {
          inherit system;
          overlays = [devshell.overlay];
        };
        */
        pkgs = nixpkgs.legacyPackages.${system};

        # Generate a user-friendly version number.
        version = builtins.substring 0 8 self.lastModifiedDate;
      in rec {
        # Provide some binary packages for selected system types.
        packages = flake-utils.lib.flattenTree {
          walbot = pkgs.buildGo119Module {
            pname = "walbot";
            inherit version;
            # In 'nix develop', we don't need a copy of the source tree
            # in the Nix store.
            src = ./.;
            vendorSha256 = "sha256-+/zZmn7J8325g7blv0kv3FJg9+LNK0rPaNpXK86QOgQ=";
          };
        };
        defaultPackage = packages.walbot;
        apps.walbot = flake-utils.lib.mkApp {drv = packages.walbot;};
        defaultApp = apps.walbot;
      }
    );
}
