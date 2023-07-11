{
  inputs = {
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };

    flake-utils = {
      url = "github:numtide/flake-utils";
    };
  };

  outputs = { self, nixpkgs, ... }@inputs:
    inputs.flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        rec {
          packages.gcov2lcov = pkgs.buildGoModule rec {
            pname = "gcov2lcov";
            version = "1.0.5";

            src = pkgs.fetchFromGitHub {
              owner = "jandelgado";
              repo = pname;
              rev = "v${version}";
              sha256 = "sha256-HDG7diLPc3CvFwmEUb+W1AAmNxiJC5vJclRME4svMUA=";
            };

            vendorSha256 = "sha256-ReVaetR3zkLLLc3d0EQkBAyUrxwBn3iq8MZAGzkQfeY=";

            # https://github.com/jandelgado/gcov2lcov/issues/15
            doCheck = false;
          };

          devShells.default = pkgs.mkShell {
            name = "Development environment";

            nativeBuildInputs = with pkgs; [
              self.packages.${system}.gcov2lcov
              ginkgo
              go
              golangci-lint
              mockgen
            ];
          };

          devShells.ci = pkgs.mkShell {
            name = "CI environment";

            nativeBuildInputs = with pkgs; devShells.default.nativeBuildInputs ++ [
              lcov
            ];
          };
        }
      );
}
