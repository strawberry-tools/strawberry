# Strawberry Development Information

*This is a temporary file. Once the Development section of Strawberry Docs is up, this information should be moved there.*

## Cloning

Cloning is straight forward.

1. First clone as you would normally: `git clone https://github.com/strawberry-tools/strawberry.git`
1. `cd strawberry`


## Dependencies

- g++ is needed. On Ubuntu, this can be installed via the "build-essential" package.
- rst2html is needed. On macOS that be can install via `brew install docutils`.
- asciidoctor is optional. Having it will allow additional tests to run. `brew install asciidoctor`
- pandoc is optional. Having it will allow additional tests to run. `brew install pandoc`


## Notes

- **RAM usage** - Running the full test suite (included RACE tests) can consume more than 4GB of RAM.
- Generating image files for Golden is not automatic. It needs to be done manually.
  - In `sb/resources/image_test.go`, the devMode variable needs to manually be set to `true` in the `TestImageOperationsGolden` function.
  - Clean test cache: `go clean -testcache`
  - go test ./...
  - In the output, you'll see a directory in `/tmp` where the files were generated. Find the image subdirectory, and copy those files to `sb/resources/testdata/golden/`.
  - revert the devMode change in the test file.
  - on macOS the temp dir would be in `$TMPDIR`.
  - https://discourse.gohugo.io/t/how-to-regenerate-files-in-testdata-golden/33944/2?u=felicianotech
