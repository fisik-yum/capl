# `capl`

WIP Implementation of the `.chnm` file format. This program is descended from a program called `skimmer` that I was initially working on, and still shares much of the code with it. However, I have discontinued my work on that project in favor of this, with the goal to increase performance with the same format used in that of its predecessor.

## Current State

This program plays Bad Apple (72x54) at full speed.

## Building and Testing

`./build.sh file`

options for `file`  
`loadscreen.chnm`

Unfortunately, keyboard controls have been removed and will be reimplemented in future commits.

## Adding new animations

Drop all `.chnm` files into `data/`. These files will be copies over to the build directory when you run `./build.sh` next time

## Creating new animations

Each frame is a single line of constant length, no escape sequences and indents. The first line of each file is a descriptor arranged in `<width> <height> <fps>`

## Tooling

See [https://github.com/fisik-yum/capl/tree/master/tools](tools)

## Troubleshooting

If passed a nonexistent file or invalid parameters, the program will not display anything, and after quitting your cursor will remain hidden. To fix this, start up a new terminal session, or run `capl` with valid parameters  
To quit, press the Escape key.

### Copyright Attribution

All code in this repository is licensed under the GPLv3 License.
