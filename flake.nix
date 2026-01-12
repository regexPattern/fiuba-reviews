{
  description = "FIUBA Reviews";

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
        webDevPkgs = with pkgs; [
          bun
          prettierd
          svelte-language-server
          tailwindcss-language-server
          typescript-language-server
          vscode-langservers-extracted
        ];
      in {
        devShells.web = pkgs.mkShell {
          packages = webDevPkgs;
        };
        devShells.actualizador = {
          cliente = pkgs.mkShell {
            packages = webDevPkgs;
          };
          servidor = pkgs.mkShell {
            packages = with pkgs; [
              air
              go
              golangci-lint
              golangci-lint-langserver
              gopls
              pgformatter
            ];
          };
        };
        formatter = pkgs.alejandra;
      }
    );
}
