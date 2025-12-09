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
            svelte-language-server
            tailwindcss-language-server
            typescript-language-server
            vscode-langservers-extracted
          ];
        };
        devShells.servidor = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            golangci-lint
            golangci-lint-langserver
            pgformatter
          ];
        };
        formatter = pkgs.alejandra;
      }
    );
}
