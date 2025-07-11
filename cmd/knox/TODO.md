# CLI init sub-command

- GIVEN knox is run with the 'init' sub-command
- WHEN there is no existing knox project for this dir
- THEN create a 'vault' file for this project

## init breakdown
- How do we determine if a project already exists? 
Skip for now?
- Where do we create these vault files? Respecting XDG Dirs
would be nice, but maybe a feature for later?
