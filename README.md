# Fzf-Image

A simple implementation of Fzf that specializes in picking images.
So far only Kitty is supported (PR's welcome for Sixel, iTerm, Chafa etc...)
Fzf is not a dependency of this project, fzfimg uses a go-fzf which is an implementation
of Fzf not related to the original project.

- Does NOT require Fzf by June-gunn
- Is fast af (faster than ueberzug)
- Small and portable

- How to run:
  ![](cli.png)
  fzfimg accepts a list of images separated by newlines on standard in. It accepts all image formats
  that kitty's icat supports.

- Preview:
  ![](preview.png)

### this software is still in early development!

<b>TODO</b>

- fix exit IO to allow piping and process substitution
- add support for sixel
- add more control over images
- call terminal size on SIGWINCH
- add esacape functionality when theres no selection
