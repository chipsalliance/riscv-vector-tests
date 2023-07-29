{
  description = "Provide multiple version of riscv-vector-tests";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
    {
      packages.sequencer-vector-case-generator = pkgs.buildGoModule {
        pname = "riscv-vector-test";
        version = "unstable-2023-04-12";
        src = pkgs.fetchFromGitHub {
          owner = "ksco";
          repo = "riscv-vector-tests";
          rev = "3fb992b1dc7f89231b27ae4a1e8d50492dde4f5b";
          hash = "sha256-BNbK8+KUwhqk3XfFgKsaeUpX0NuApl8mN/URKwYTYtE=";
        };
        doCheck = false;
        vendorHash = "sha256-9cQlivpHg6IDYpmgBp34n6BR/I0FIYnmrXCuiGmAhNE=";
        # Get rid of copying the whole source
        postInstall = ''
          cp -r $src/configs $out/configs
          mkdir $out/include
          cp $src/macros/sequencer-vector/* $out/include
          cp $src/env/sequencer-vector/* $out/include

          cp $src/env/encoding.h $out/include
          # Let riscv_test.h find the encoding.h file correctly
          sed -i 's/^#include "..\/encoding.h"$/#include "encoding.h"/' $out/include/riscv_test.h
        '';
      };
    });
 }
