{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      packages.${system}.default = pkgs.buildGoModule {
        name = "score-processor";
        src = self;
        goSum = ./go.sum;
        vendorHash = "sha256-+6G2OJPG2aVPqJgrJP6le1+fMtuBs1S59xZ0/wwzuX4=";
        nativeBuildInputs = with pkgs; [ pkg-config ];
        buildInputs = with pkgs; [
          opencv
        ];
      };
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
          opencv
          mupdf-headless
          pkg-config
        ];
      };
    };
}
