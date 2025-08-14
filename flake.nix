{
	description = "Flake for Studying Go";

	inputs = {
		nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
	};

	outputs = { self, nixpkgs }@inputs:
	let
		system = "x86_64-linux";
		
		pkgs = import inputs.nixpkgs{
			inherit system;
			config.allowUnfree = true;
		};
	in 
	{
		devShells.${system}.default = pkgs.mkShell rec {
			name = "Pyxis Service";
			packages = with pkgs; [
				go
				gopls
			];
			shellHook = "tmux -L PyxisService new-session -A -s Pyxis-Service";
			};
		};
	}
