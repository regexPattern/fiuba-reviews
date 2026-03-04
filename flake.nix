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

        mkShells = {
          devPackages,
          runtimePackages,
        }: let
          dev = pkgs.mkShell {
            packages = devPackages;
          };
          runtime = pkgs.mkShell {
            packages = runtimePackages;
          };
        in {
          inherit dev runtime;
          full = pkgs.mkShell {
            inputsFrom = [dev runtime];
          };
        };
        
      in {
        devShells.web = mkShells {
          devPackages = with pkgs; [
            prettierd
            svelte-language-server
            tailwindcss-language-server
            typescript-language-server
            vscode-langservers-extracted
          ];
          runtimePackages = with pkgs; [
            nodejs_24
            pnpm
          ];
        };

        devShells.actualizador.cliente = mkShells {
          devPackages = with pkgs; [
            prettierd
            svelte-language-server
            tailwindcss-language-server
            typescript-language-server
            vscode-langservers-extracted
          ];
          runtimePackages = with pkgs; [
            nodejs_24
            pnpm
          ];
        };

        devShells.actualizador.servidor = mkShells {
          devPackages = with pkgs; [
            golangci-lint-langserver
            gopls
          ];
          runtimePackages = with pkgs; [
            air
            go
            golangci-lint
            pgformatter
          ];
        };

        formatter = pkgs.alejandra;
      }
    );
}
