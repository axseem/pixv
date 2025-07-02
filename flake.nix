{
  description = "A CLI tool for converting pixel art images into SVG";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, ... }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          pixv = pkgs.buildGoModule {
            pname = "pixv";
            version = "0.3.0";
            src = ./.;
            vendorHash = "sha256-kYbVtbVorI+cCMAqizYfSm2BcSjyKmICkcKahi9Llx0=";
          };
          default = pixv;
        });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.go
              pkgs.gotools
            ];
          };
        });
    };
}