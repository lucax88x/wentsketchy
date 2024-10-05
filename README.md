# what is wentsketchy?

It provides a repo with building blocks for all the plugins to build complex setups.
It supports handling between complex programs, such as aerospace, and yabai.

demo:
https://github.com/user-attachments/assets/4a238f9f-dd4a-4e7b-9a6f-e458564acca7

today, it supports:

- aerospace (still couple of bugs)
- front-app
- sensors (temperatures to be refined), and need stats to be installed
- cpu (rewrite of the helper from FelixKratz) in go
- battery (one battery)
- calendar

## faq

- why go?
because I like and enjoy it, I find it a good compromise of stability and performance for low-level
- why not sbarlua?
I'm a fan of lua, but I don't like dynamic languages for long-term maintenance, and I needed parallelism
thus, I neded a "queue" for aerospace commands


## minimal setup to test it

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

"$HOME/bin/wentsketchy" start
```

and this in .aerospace.toml to test

```shell
exec-on-workspace-change = [
  '/bin/bash',
  '-c',
  'echo "aerospace_workspace_change { \"focused\": \"$AEROSPACE_FOCUSED_WORKSPACE\", \"prev\": \"$AEROSPACE_PREV_WORKSPACE\" } ¬" > /tmp/wentsketchy',
]
```

and put in ~/.config/sketchybar/config.yaml the wentsketchy configuration

```yaml
---
left:
  - main_icon
  - aerospace
right_notch:
  - front_app
right:
  - sensors
  - cpu
  - battery
  - calendar
```

### architecture

how it works?

wentsketchy is a simple go application to run by cli
it uses a fifo (named pipe) to handle communications between aerospace and sketchybar
it caches some aerospace data to keep it fast, and renews it every minute
sketchybar items will emit `update + sketchybar args`
aerospace items will emit specific events, such as `aerospace-workspace-change + sketchybar args`


### TODO
- how to get rid of echo commands not dieing
- get aerospace mode (layout, tabbed, etc)
- get aerospace fullscreen
- when workspace collapsed, show number of windows
- wifi item
- wifi https://github.com/FelixKratz/SketchyBar/discussions/12#discussioncomment-8908932
- vpn https://github.com/FelixKratz/SketchyBar/discussions/12#discussioncomment-1216869
- down & up speeds https://github.com/FelixKratz/SketchyBar/discussions/12#discussioncomment-8107907
- tests aerospace mocking it

# # Known limitations

- order of windows 
 there's no way to have a correct order of windows from aerospace
- highlight single window
  we have a front-app event from sketchybar, but no front-window events from anyone
  we can get the focused window, but no event to react to
- click on window
  we cannot select a window from aeropsace
