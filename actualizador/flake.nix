{
  description = "Empty Template";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.cliente = pkgs.mkShell {
          packages = with pkgs; [
            bun
            prettierd
            svelte-language-server
            tailwindcss-language-server
            typescript-language-server
            vscode-langservers-extracted
          ];
        };
        devShells.servidor = pkgs.mkShell {
          packages = with pkgs; [
            air
            go
            golangci-lint
            golangci-lint-langserver
            gopls
            pgformatter
          ];
        };
        formatter = pkgs.alejandra;
      }
    );
}
