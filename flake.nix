{
  description = "FIUBA Reviews";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  # Estructura de las shells
  #
  # Cada proyecto tiene dos shells: una para desarrollo y otra que incluye el
  # runtime del mismo.
  #
  # La shell de desarrollo está pensada para ser utilizada en entornos donde el
  # editor de código busca los binarios de las herramientas de desarrollo en el
  # PATH. De esta forma, el desarrollador tiene acceso a las versiones pineadas
  # de esos binarios directamente desde este flake y su respectivo lockfile. Es
  # decir, por ejemplo, si el editor de código a elección es Neovim, asumiendo
  # que este está correctamente configurado, el editor va a poder encontrar los
  # language servers necesarios en el PATH, una vez se haya cargado la shell,
  # ya sea con el comando `nix develop` o con alguna utilidad como direnv. Si
  # en cambio utilizas algo como VSCode, la mayoría de extensiones ya traen su
  # propio servidor incluido, por lo que capaz te resulte más útil simplemente
  # aceptar las extensiones sugeridas para cada proyecto.
  #
  # La shell de runtime está pensada para demás dependencias que no sean
  # meramente de desarrollo, sino que más bien son las dependecias para la
  # compilación, ejecución, testing, etc. del proyecto correspondiente.

  outputs = {
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.web = {
          dev = pkgs.mkShell {
            packages = with pkgs; [
              prettierd
              svelte-language-server
              tailwindcss-language-server
              typescript-language-server
              vscode-langservers-extracted
            ];
          };
          runtime = pkgs.mkShell {
            packages = with pkgs; [
              nodejs_24
              pnpm
            ];
          };
        };

        devShells.actualizador.cliente = {
          dev = pkgs.mkShell {
            packages = with pkgs; [
              prettierd
              svelte-language-server
              tailwindcss-language-server
              typescript-language-server
              vscode-langservers-extracted
            ];
          };
          runtime = pkgs.mkShell {
            packages = with pkgs; [
              nodejs_24
              pnpm
            ];
          };
        };

        devShells.actualizador.servidor = {
          dev = pkgs.mkShell {
            packages = with pkgs; [
              golangci-lint-langserver
              gopls
            ];
          };
          runtime = pkgs.mkShell {
            packages = with pkgs; [
              air
              go
              golangci-lint
              pgformatter
            ];
          };
        };

        formatter = pkgs.alejandra;
      }
    );
}
