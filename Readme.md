# champ

[![Codewake](https://www.codewake.com/badges/ask_question_flat_square.svg)](https://www.codewake.com/p/champ)

Experiment into 2nd screen player for Plex protocol on single board
computers like RPi. Build instructions https://github.com/dz0ny/champ/wiki/Building-for-RPi

## Install

App is packaged with snap and published in app store and under https://github.com/dz0ny/champ/releases

## Usage

The main thing
```
➜ champ [<flags>]

Minimalistic Plex 2nd screen client

Flags:
  --help                  Show context-sensitive help (also try --help-long and
                          --help-man).
  --debug                 Verbose mode.
  --title="Champ Player"  Name of this player
  --port="32016"          HTTP server port
  --version               Show application version.
```

Dev helper for MPV integration, this will later be integration into champ.
```
➜ spinwheel [<flags>]

Shuffle player which also plays from YouTube(tm)

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
      --debug              Verbose mode.
  -p, --playlist=PLAYLIST  Path to playlist file (.yaml)
      --version            Show application version.

```
