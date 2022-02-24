# Contributing to `appctl` #

Thanks for your help improving the appctl!

## Getting Help ##

If you have questions around `appctl`, please contact the platform9 support or if you encounter any problem using it do raise an [Issue](https://github.com/platform9/appctl/issues), you can reach us via email at [support@platform9.com](support@platform9.com).


## Workflow ##

Contribute to `appctl` by following the guidelines below:
- [Clone](https://github.com/platform9/appctl.git)
- [Setup local environment](https://github.com/platform9/appctl#building-appctl-locally) 
- Work on the cloned repository
- Open a pull request to [merge into appctl repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request)


## Building ## 

Follow these steps to build from source and deploy:

1. Hookup appctl with a [local or a hosted fast-path](https://github.com/platform9/fast-path#readme)
2. Run the appropriate `make` target
- Linux

  ```sh
  make build-linux64
  ```
- Windows
  ```sh
  make build-win64
  ```
- MacOS
  ```sh
  make build-mac
  ```
3. To generate all the binaries run the default target. The executables are placed in `bin` directory.


## Running the unit tests and manual test ##
1. To run unit tests:
```sh
make test
```
2. When submitting the Pull request, perform manual tests of all commands and attach snapshots or output files.
```sh
# Commands to test
1. appctl login
2. appctl deploy
3. appctl list
4. appctl describe 
5. appctl delete
```

## Committing ###

Please follow the Pull request template, before raising a PR, so reviewers will gain a deeper understanding as they review. If an outstanding issue has been fixed, please include the Fixes Issue # in your commit message.
