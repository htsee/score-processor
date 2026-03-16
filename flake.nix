{
  description = "Flake for score-processor";

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
        version = "0.1.0";
        src = self;
        goSum = ./go.sum;
        vendorHash = "sha256-UtCtuypTlNntJu4W2oc6WLxdBR+QOZwz2DUZgBeiXXI=";
        nativeBuildInputs = with pkgs; [
          pkg-config
          makeWrapper
        ];
        buildInputs = with pkgs; [
          opencv
          mupdf-headless
        ];
        postInstall = ''
                    mv $out/bin/score-processor $out/bin/sp
          					wrapProgram $out/bin/sp \
          					--prefix PATH : ${pkgs.mupdf-headless}/bin
          				'';
      };
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
          golangci-lint
          opencv
          mupdf-headless
          pkg-config
        ];
      };
    };
}
