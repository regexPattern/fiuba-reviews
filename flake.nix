{
  description = "FIUBA Reviews";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells = {
          actualizador = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              golangci-lint
              gopls
            ];
          };

          resumidor = pkgs.mkShell {
            buildInputs = with pkgs; [
              rustc
              cargo
              rust-analyzer
            ];
          };

          web = pkgs.mkShell {
            buildInputs = with pkgs; [
              nodejs
              pnpm
              typescript
              nodePackages.svelte-language-server
            ];
          };
        };
      }
    );
}

