{ pkgs ? import <nixpkgs> {} }:

with pkgs;

mkShell {
  name = "prometheus-moto-exporter-dev";
  buildInputs = [ go gotools-unstable ];
}
