{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "weecal";
          version = "unversioned";
          src = ./.;
          ldflags = [ "-s" "-w" "-X main.version=dev" "-X main.builtBy=flake" ];
          doCheck = false;
          vendorHash = "";
        };

        devShells.default = pkgs.mkShell {
	  packages = with pkgs; [
	    go
	    air
	    tailwindcss
	    templ
	    nodejs_18
	    openssl
	  ];
          shellHook = "go mod tidy";
        };

        # nix develop .#dev
        devShells.dev = pkgs.mkShell {
          packages = with pkgs; [
            go
          ];
        };
      }
    );
}

