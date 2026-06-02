{
  description = "Personal website for Parthiv Krishna";

  inputs = {
    flake-utils = {
      url = "github:numtide/flake-utils";
    };
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixos-unstable";
    };
    nix-colors = {
      url = "github:misterio77/nix-colors";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    { self, nixpkgs, ... }@inputs:
    inputs.flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        inherit (nixpkgs) lib;
        oneDark = inputs.nix-colors.colorSchemes.onedark.palette;
        themeSource = pkgs.writeTextDir "theme/theme.toml" ''
          [name]
          id = "personal-onedark"
          display = "Personal OneDark"

          [colors]
          bg = "#${oneDark.base00}"
          fg = "#${oneDark.base05}"
          muted = "#${oneDark.base04}"
          subtle = "#${oneDark.base03}"

          accent = "#${oneDark.base0D}"
          link = "#${oneDark.base0D}"
          link_hover = "#${oneDark.base0C}"

          heading_marker = "#${oneDark.base03}"
          heading_1 = "#${oneDark.base0D}"
          heading_2 = "#${oneDark.base0E}"
          heading_3 = "#${oneDark.base0C}"
          heading_4 = "#${oneDark.base0B}"
          heading_5 = "#${oneDark.base0A}"
          heading_6 = "#${oneDark.base08}"

          code_fg = "#${oneDark.base05}"
          code_bg = "#${oneDark.base01}"
          code_border = "#${oneDark.base03}"

          quote_fg = "#${oneDark.base0C}"
          quote_marker = "#${oneDark.base03}"
          quote_bg = "#${oneDark.base01}"

          rule = "#${oneDark.base03}"
          selection = "#${oneDark.base02}"

          [typography]
          body_font = "ZedMono Nerd Font, ZedMono Nerd Font Mono, ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace"
          heading_font = "ZedMono Nerd Font, ZedMono Nerd Font Mono, ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace"
          line_height = "1.65"

          [markdown]
          show_heading_markers = true
          heading_marker_weight = "normal"
          heading_text_weight = "700"
          show_link_brackets = false
          show_blockquote_marker = true

          [terminal]
          max_width = 96
          h1_margin_top = 1
          h1_margin_bottom = 1
          h2_margin_top = 1
          h2_margin_bottom = 0
          paragraph_margin_bottom = 1
        '';

        treefmtEval = inputs.treefmt-nix.lib.evalModule pkgs {
          projectRootFile = "flake.nix";
          programs.mdformat.enable = true;
          programs.gofmt.enable = true;
          programs.nixfmt.enable = true;
          programs.taplo.enable = true;
        };

        tuiPackage = pkgs.buildGoModule {
          pname = "parthiv-cc-tui";
          version = "0.1.0";
          src = ./tui;
          vendorHash = "sha256-uwBJAqN4sIepiiJf9lCDumLqfKJEowQO2tOiSWD3Fig=";
          subPackages = [
            "cmd/theme"
            "cmd/tui"
          ];
        };
      in
      {
        formatter = treefmtEval.config.build.wrapper;

        packages.tui = tuiPackage;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            treefmtEval.config.build.wrapper
            wrangler
            zola
          ];
        };

        # generate shared web and terminal theme artifacts
        apps.theme = {
          type = "app";
          program = lib.getExe (
            pkgs.writeShellApplication {
              name = "parthiv-cc-theme";
              runtimeInputs = [ tuiPackage ];
              text = ''
                theme -repo "$PWD" -theme ${themeSource}/theme/theme.toml
              '';
            }
          );
        };

        # local terminal UI for browsing the Zola content
        apps.tui = {
          type = "app";
          program = lib.getExe (
            pkgs.writeShellApplication {
              name = "parthiv-cc-tui";
              runtimeInputs = [ tuiPackage ];
              text = ''
                tui -content "$PWD/content"
              '';
            }
          );
        };

        # development server
        apps.default =
          let
            name = "parthiv-cc-site-dev";
          in
          {
            type = "app";
            program = lib.getExe (
              pkgs.writeShellApplication {
                inherit name;
                runtimeInputs = with pkgs; [
                  tuiPackage
                  zola
                ];
                text = ''
                  theme -repo "$PWD" -theme ${themeSource}/theme/theme.toml
                  zola serve
                '';
              }
            );
          };

        # production site
        packages.default = pkgs.stdenv.mkDerivation {
          name = "parthiv-cc-site";
          src = ./.;

          buildInputs = with pkgs; [
            tuiPackage
            zola
          ];

          buildPhase = ''
            theme -repo "$PWD" -theme ${themeSource}/theme/theme.toml
            zola build
          '';

          installPhase = ''
            cp -r public $out
          '';
        };

        # deploy production site
        apps.deploy =
          let
            name = "parthiv-cc-site-deploy";
            site = self.packages.${system}.default;
          in
          {
            type = "app";
            program = lib.getExe (
              pkgs.writeShellApplication {
                inherit name;
                runtimeInputs = with pkgs; [
                  wrangler
                ];
                text = ''
                  if [ -z "''${CLOUDFLARE_PROJECT_NAME:-}" ]; then
                    echo "Error: CLOUDFLARE_PROJECT_NAME environment variable is not set"
                    exit 1
                  fi
                  if [ -z "''${CLOUDFLARE_ACCOUNT_ID:-}" ]; then
                    echo "Error: CLOUDFLARE_ACCOUNT_ID environment variable is not set"
                    exit 1
                  fi
                  if [ -z "''${CLOUDFLARE_API_TOKEN:-}" ]; then
                    echo "Error: CLOUDFLARE_API_TOKEN environment variable is not set"
                    exit 1
                  fi

                  wrangler pages deploy ${site} \
                    --project-name="$CLOUDFLARE_PROJECT_NAME" \
                    --commit-dirty=true
                '';
              }
            );
          };
      }
    );
}
