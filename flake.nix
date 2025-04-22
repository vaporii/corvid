{
  inputs = rec {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{self, nixpkgs, ... }: let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
  in {
    devShells.${system}.default = pkgs.mkShell {
      packages = with pkgs; [
        go
      ];
    };

    packages.${system} = {
      default = self.packages.${system}.corvid;
      corvid = pkgs.buildGoModule {
        pname = "corvid";
        version = "v1.0.0";

        src = ./.;

        vendorHash = "sha256-WUTGAYigUjuZLHO1YpVhFSWpvULDZfGMfOXZQqVYAfs=";
      };
    };

    overlays.default = self.overlays.corvid;
    overlays.corvid = final: prev: {
      corvid = self.packages.${system}.corvid;
    };
  };
}