{pkgs, ...}: {
  packages = with pkgs; [
    gofumpt
  ];

  devcontainer.enable = true;
  difftastic.enable = true;

  languages.go.enable = true;
  languages.nix.enable = true;

  pre-commit.hooks = {
    alejandra.enable = true;
    prettier.enable = true;
    shellcheck.enable = true;
    shfmt.enable = true;
    gofumpt = {
      enable = true;
      name = "gofumpt";
      entry = "${pkgs.gofumpt}/bin/gofumpt -d";
      types = ["go"];
    };
  };
}
