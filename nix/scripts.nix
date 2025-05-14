{ pkgs
, config
, mrmintchain ? (import ../. { inherit pkgs; })
}: rec {
  start-mrmintchain = pkgs.writeShellScriptBin "start-mrmintchain" ''
    # rely on environment to provide mrmintchaind
    export PATH=${pkgs.test-env}/bin:$PATH
    ${../scripts/start-mrmintchain.sh} ${config.mrmintchain-config} ${config.dotenv} $@
  '';
  start-geth = pkgs.writeShellScriptBin "start-geth" ''
    export PATH=${pkgs.test-env}/bin:${pkgs.go-ethereum}/bin:$PATH
    source ${config.dotenv}
    ${../scripts/start-geth.sh} ${config.geth-genesis} $@
  '';
  start-scripts = pkgs.symlinkJoin {
    name = "start-scripts";
    paths = [ start-mrmintchain start-geth ];
  };
}
