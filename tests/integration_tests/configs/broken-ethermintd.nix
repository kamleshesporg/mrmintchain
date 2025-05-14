{ pkgs ? import ../../../nix { } }:
let mrmintchaind = (pkgs.callPackage ../../../. { });
in
mrmintchaind.overrideAttrs (oldAttrs: {
  patches = oldAttrs.patches or [ ] ++ [
    ./broken-mrmintchaind.patch
  ];
})
