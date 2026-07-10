# ptgen

[![Test-Build-Push](https://github.com/FatmanUK/ptgen/actions/workflows/deploy.yml/badge.svg)](https://github.com/FatmanUK/ptgen/actions/workflows/deploy.yml)

I wanted to "compile" my `.mod` file from plain text, so it could be almost entirely versioned by git. I have a plan to implement samples and instruments.

`.mod` trackers (eg. MilkyTracker) are great, but they tend to use old-school DOS interfaces, and as such lack copy-n-paste. Entering notes (patterns) by hand in Milky is fun, but after a while it wears a bit. I'd rather enter the data in a modern text editor using a specific sparse format.

To be clear, by `.mod` file I mean an old-school 90s musical composition format. They were mainly used by Amiga games. A `.mod` file had to be incredibly small and sparse --- several of them had to fit comfortably alongside a full game in just 880 *kilobytes* of disk space. Producing catchy tunes under these conditions required all the skills of the compositional geniuses of the day.

See works in progress:
  - https://github.com/FatmanUK/mod_cold_boot
  - https://github.com/FatmanUK/mod_signal_lost

## Lint/Test
   
    ❯ make lint
    ❯ make test

## Build
   
    ❯ make build
    ❯ make run

## Podman
   
    ❯ make podman-build
    ❯ make podman-run
    ❯ podman pull ghcr.io/fatmanuk/ptgen

## Run

Then generate a test `.mod` with:
   
    ❯ <test_data/sample2/metadata.json ./ptgen \
          -p test_data/sample2 \
          -m /dev/stdin \
          >/tmp/sample2.mod

Note the `.mod` data is output directly to stdout, so you must redirect it. I did that because it simplifies the command syntax and feels more "Unix-y".

## Samples

I'm currently using samples from ST-01. They're tiny *and* sound good. Besides, they're (likely?) unencumbered by copyright, unlike the rest of the ST-XX sample archive.

To get the ST-01 archive on Debian-based Linux, issue these commands:
   
    ❯ sudo apt update
    ❯ sudo apt install lhasa wget
    ❯ mkdir samples
    ❯ cd samples
    ❯ ST01='https://aminet.net/mods/inst/st-01.lha'
    ❯ wget -O st-01.lha "${ST01}"
    ❯ lha x st-01.lha

These files are in IFF format and AmigaOS doesn't use file extensions, so to make MilkyTracker see them you have to rename the ones you want with '.iff' extensions. Like this:
   
    ❯ mv Stabs Stabs.iff
