score-processor (sp) is a tool for processing sheet music into images, inspired by edwardx999's [ScoreProcessor](https://github.com/edwardx999/ScoreProcessor).

# Install
## Go
Make sure you have OpenCV (and optionally mupdf) installed
```bash
go install github.com/htsee/score-processor
```

## Nix
Add `github:htsee/score-processor` as an input to your `flake.nix`

Then add the package `score-processor.packages."x86_64-linux".default` to the package list of your configuration

# Usage
Type `sp` to get a list of subcommands

All commands require inputs and a destination folder. They can accept multiple inputs. A new folder will be created if the folder specified does not exist.

All subcommands accept `-h` flag to access a list of flags.

## Subcommands

`convert`: convert PDFs into images (PNG). Requires mupdf.

`cut`: cut the score into staves.

`denoise`: remove noise from the score.

`deskew`: deskew the score. 

`fit`: fit the image into specified aspect ratio.

`pad`: add a white border around the score.

`rotate`: rotate the score by the specified angle.

`splice`: combine multiple staves together.

`trim`: trim the score. Can be used to remove marginal elements.

`vsplice`: combine two pages. Useful for orchestral scores.

# Build
Clone the repo:
```bash
git clone github.com/htsee/score-processor
cd score-processor
```

## Go
Make sure you have OpenCV installed
```bash
go build
```

## Nix
```bash
nix build
```

Feel free to open pull requests
