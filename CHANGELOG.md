# Changelog

## v0.2.1

- Fix bug that crashed `hosted-models` example when `query` input was JSON literal on Windows.

## v0.2.0

Error handling improvements and updates to `hosted-models` example.

- Rewrite errors and error handling to provide more context. This may introduce backwards incompatible changes.
- Add `HostedModel.GetURL()` method.
- Rename `hosted-model` example and binary to `hosted-models`.
- Print everything but model output to stderr in `hosted-models` example.
- Replace `ErrInvlaidURL` (misspelled) with `ErrInvalidURL`.
- Add `--help` flag to `hosted-models` example.
- Fix error in `--url` description of `hosted-models` example.
- Use POSIX style flags for all CLI arguments in `examples/`.

## v0.1.2

- Fix 400 error on request retries.
- Add `CHANGELOG`.

## v0.1.1

- Add MIT License.

## v0.1.0

Initial public release.
