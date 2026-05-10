{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
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
}
