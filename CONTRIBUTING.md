# Contributing Guidelines

## Code styles

Make sure to check and fix code formatting issues before contributing:

- Chaincode and Go SDK

  Use [`gofmt`](https://go.dev/blog/gofmt) to format your code and
  [`golangci-lint`](https://golangci-lint.run/) to check for code style issues.

- Java SDK

  Run `mvn spotless:check` and `mvn checkstyle:check` to check for code style issues and
  `mvn spotless:apply` to reformat code.

- JavaScript SDK

  Run `npm run format` and `npm run lint` to format code and fix obvious issues.
  Run `npm run format:check` and `npm run lint:check` to check for them.

## Project versioning

We use [SemVer](http://semver.org/) for versioning.
Any changes to the code base should not be released in an existing version.

## Commit message format

The commit message should follow the
[Bluejava commit message format](https://github.com/bluejava/git-commit-guide).
The supported scopes are:

- **chaincode** for chaincode commits
- **go** for Go SDK-related commits
- **java** for Java SDK-related commits
- **javascript** for JavaScript SDK-related commits
- **build** for build scripts, CI, other development or deployment related commits
- use **\*** or leave empty to refer to commits that do not have a clear scope

## Submit changes

1. Please review and accept our [Code of Conduct](CODE_OF_CONDUCT.md).
2. Fork this repository, make changes to it, and run the tests to see if
   anything is broken along the way.
3. Submit a pull request to our repository.
   Please make sure you followed the instructions above.
4. Our team will review the pull request on regularly. We will try our best to
   inform you of our questions about your pull request and any changes that
   need to be made before we can merge your pull request.
5. We expect your response within two weeks, after which your pull request may
   be closed if no activity is shown.
