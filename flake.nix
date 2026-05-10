{
  inputs = {
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixos-unstable";
    };
    utils = {
      url = "github:numtide/flake-utils";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages = {
          # The actual package
          default = pkgs.buildGoModule {
            pname = "input-display";
            version = "1";
            src = ./.;
            vendorHash = "sha256-elkbVg/j2lVR/ldNakxmoPgDvR/5l7eQDyS8WJMB58Q=";
          };
        };
        devShells = {
          # Development environment
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              pkg-config
              mesa
              libGL
              libX11
              libXrandr
              libXinerama
              libXcursor
              libXi
              libXxf86vm
            ];
          };
        };
      }
    );
}
