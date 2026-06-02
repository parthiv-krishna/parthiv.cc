# parthiv.cc

Personal website for Parthiv Krishna, built using [Zola](https://www.getzola.org/) and [Nix](https://nixos.org/), and based on the [Papaya theme](https://github.com/justint/papaya).

## Nix Apps / Packages

### Format

```bash
nix fmt
```

### Run development server

```bash
nix run
```

This generates the shared theme artifacts and starts `zola serve`.

### Generate theme artifacts

```bash
nix run .#theme
```

### Run terminal UI

```bash
nix run .#tui
```

### Build site

```bash
nix build
```

The build generates `static/generated/theme.css` before running `zola build`.

### Deploy to Cloudflare Pages

Make sure to set `CLOUDFLARE_PROJECT_NAME`, `CLOUDFLARE_ACCOUNT_ID`, and `CLOUDFLARE_API_TOKEN` environment variables.

```bash
nix run .#deploy
```
