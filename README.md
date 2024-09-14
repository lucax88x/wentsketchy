# minimal setup to test it

minimal setup

```shell
bun install
```

build the project and copy the bin to  ~/bin/ (remember to have ~/bin in the $PATH)

```
bun run build
bun run cp
```


then use this .sketchybarrc to test

```shell
#!/bin/bash

"$HOME/bin/wentsketchy" init

# Forcing all item scripts to run (never do this outside of sketchybarrc)
sketchybar --update

echo "sketchybar configuation loaded.."
```


## Architecture

wip


## TODO

- find out why wentsketchy hangs
- workspace focused event
