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

### Build site

```bash
nix build
```

### Deploy to Cloudflare Pages

Make sure to set `CLOUDFLARE_PROJECT_NAME`, `CLOUDFLARE_ACCOUNT_ID`, and `CLOUDFLARE_API_TOKEN` environment variables.

```bash
nix run .#deploy
```
