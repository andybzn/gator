{
  description = "A Blog Aggregator";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {

        packages.default = pkgs.postgresql_17_jit;

        apps.default = {
          type = "app";
          program = "${pkgs.writeShellScriptBin "start-postgres" ''
            mkdir -p ./data
            chmod 700 ./data

            chmod 777 /tmp

            if [ ! -f ./data/PG_VERSION ]; then
              ${pkgs.postgresql_17_jit}/bin/initdb -D ./data
            fi

            ${pkgs.postgresql_17_jit}/bin/postgres -D ./data
          ''}/bin/start-postgres";
        };

        devShells.default = pkgs.mkShell {
          packages = [ pkgs.zsh ];
          buildInputs = [
            pkgs.postgresql_17_jit
          ];

          shellHook = ''
            export SHELL=${pkgs.zsh}/bin/zsh

            if [[ $SHELL != ${pkgs.zsh}/bin/zsh ]]; then
              exec ${pkgs.zsh}/bin/zsh
            fi
          '';
        };
      }
    );
}
