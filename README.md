

# Grade assistant (Graas)

## Some goals

  * Try to make future-proof (select stable, well adapted, popular dependencies)
  * Task generation should be based on strong cryptography
  * Follow current best practices of Golang
  * Proper testing is important and use formatters to improve code quality
    * https://bitfieldconsulting.com/golang/tdd
    * https://staticcheck.io/
    * If unsure about IDE, GoLand and Fleet will do it


## Current work

Initial remake with cobra and viper has been made for long term support.
Example https://github.com/carolynvs/stingoftheviper




## Potential libraries for future development

* GitHub API access https://github.com/cli/go-gh (Might not needed much)
* Building Containers https://github.com/containers/buildah
 * No root required, we want that
* Possibly also some other build tool should be considered for general binaries. For every task, there should be some common ground to define building of the binary and extract the output
  * Maybe Task https://github.com/go-task/task
  * Most likely for Go, C and Rust binaries at first
  * Also PDF, .docx generation is possible in the future, very versatile is required
