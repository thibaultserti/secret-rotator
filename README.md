# Secret rotator

## Badges

[![Build Status](https://github.com/thibaultserti/secret-rotator/actions/workflows/release.yaml/badge.svg)](https://github.com/thibaultserti/secret-rotator/actions/workflows/release.yaml)
[![License](https://img.shields.io/github/license/thibaultserti/secret-rotator)](/LICENSE)
[![Release](https://img.shields.io/github/release/thibaultserti/secret-rotator.svg)](https://github.com/thibaultserti/secret-rotator/releases/latest)
[![GitHub Releases Stats of secret-rotator](https://img.shields.io/github/downloads/thibaultserti/secret-rotator/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=thibaultserti&repository=secret-rotator)

[![Maintainability](https://api.codeclimate.com/v1/badges/4133d7da3d73fa0c0884/maintainability)](https://codeclimate.com/github/thibaultserti/secret-rotator/maintainability)
[![codecov](https://codecov.io/gh/thibaultserti/secret-rotator/branch/main/graph/badge.svg?token=5BO47LR632)](https://codecov.io/gh/thibaultserti/secret-rotator)
[![Go Report Card](https://goreportcard.com/badge/github.com/thibaultserti/test-saas-ci)](https://goreportcard.com/report/github.com/thibaultserti/secret-rotator)


## Run locally

To run the code locally, you can use the makefile:

```bash
# load env vars from file
export $(cat .env | xargs) 2>&1 > /dev/null
# build code
$ make build
# build & run
$ make run
```

Configuration is retrieved from `secretrotator.yaml` by default. (See structure below)
Path can be configured from the env var `SECRETROTATOR_CONFIG_PATH`.
You can take example from the `secretrotator.example.yaml` file or from the configuration file in `test/`

`.env` contains the secrets used to connect to the multiple backends. (See `.env.example` and the configuration section below)

## Run with docker

To run the code with docker, you can use the makefile:

```bash
# build & run with docker-compose
$ make docker-compose
```

As for the local mode, you'll need to fill the `secretrotator.yaml` configuration file and the `.env` file.

## Configuration

| **input**                                | **required** | **default** | **supported value**                   | **description**                                                                                                      |
| ---------------------------------------- | ------------ | ----------- | ------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| env                                      | False        | prod        | any string                            | Environment, used to configure logger format.                                                                        |
| logLevel                                 | False        | info        | `trace`/`debug`/`info`/`warn`/`error` | Log level for the logger.                                                                                            |
| providers.\<id\>                         | True         | N/A         | any string                            | Arbitrary ID, used to match configuration in the secrets section.                                                    |
| providers.\<id\>.type                    | True         | N/A         | `gitlab`/`gcp`                        | Backend use either as a secret source or a secret destination.                                                       |
| providers.\<id\>.repoUrl                 | False        | ""          | any string                            | URL used to communicate with the backend API. Not used for GCP.                                                      |
| providers.\<id\>.apiToken                | True         | ""          | any string                            | Token used to communicate with the backend API. Needs to be set to "" to be overrides by env vars. Not used for GCP. |
| secrets.name                             | True         | N/A         |                                       | Arbitrary ame of the secret rotation operation.                                                                      |
| secrets[].source.id                      | True         | N/A         | any string                            | Must match the ID of a previously declared provider.                                                                 |
| secrets[].source.type                    | True         | N/A         | `gcp_sa`                              | Source type to read the secret.                                                                                      |
| secrets[].source.path                    | True         | N/A         | any string                            | Path where to read the secret: GCP projectId.                                                                        |
| secrets[].source.secretName              | True         | N/A         | any string                            | Secret name: GCP SA email.                                                                                           |
| secrets[].source.options.nbMaxConcurrent | False        | 2           | any int                               | Max number of versions of the secret to keep.                                                                        |
| secrets[].dest[].id                      | True         | N/A         | any int                               | Must match the ID of a previously declared provider.                                                                 |
| secrets[].dest[].type                    | True         | N/A         | `gitlab_project`                      | Source type to write the secret.                                                                                     |
| secrets[].dest[].path                    | True         | N/A         | any int                               | Path where to write the secret: GitLab repo path.                                                                    |
| secrets[].dest[].secretName              | True         | N/A         | any int                               | Secret name: GitLab CI variable key.                                                                                 |
| secrets[].dest[].options.base64          | False        | false       | any int                               | Enables or disables base64 encoding of secrets.                                                                      |

Every non array option has a corresponding environment variable.
For example `providers.gitlab.type` corresponds to `SECRETROTATOR_PROVIDERS_GITLAB_TYPE`.
When using env vars, corresponding configuration must be set to `""` in the config file so that the variable is evaluated.

## Development


### Run the tests

Fill `.env` secret in `test/`

```bash
$ make test
```

### Run quality

```bash
$ make quality
```
