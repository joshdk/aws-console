[![License][license-badge]][license-link]
[![Actions][github-actions-badge]][github-actions-link]
[![Releases][github-release-badge]][github-release-link]

# AWS Console

🔗 Generate a temporary login URL for the AWS Console

## Installation

Prebuilt binaries for several architectures can be found attached to any of the available [releases][github-release-link].

For Linux:
```shell
wget https://github.com/joshdk/aws-console/releases/download/v0.1.0/aws-console-linux-amd64.tar.gz
tar -xf aws-console-linux-amd64.tar.gz
sudo install aws-console /usr/bin/aws-console
```

For Mac:
```shell
brew tab joshdk/tap
brew install aws-console
brew upgrade aws-console
```

A development version can also be built directly from this repository.
Requires that you already have a functional Go toolchain.
```shell
go install github.com/joshdk/aws-console@master
```

## Usage

### Configs and Credentials

This tool generates temporary login URLs for the AWS Console using the credentials from a named AWS cli profile.

The configuration files for these named profiles are located at `~/.aws/credentials` and `~/.aws/config`.
For more information on these two file and configuring profiles, please take a look at:

- https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
- https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html

### User Federation

In the likely event that a named profile provides credentials for an IAM user (opposed to an IAM role), that user must first be federated to obtain temporary credentials.
AWS does not permit generating a Console login URL using IAM user credentials, which is why federating users is necessary.
For more information on federating credentials, please take a look at:

- https://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html

This tool will detect and automatically federate IAM users transparently.

### Examples

Generate an AWS Console login URL for the default profile:
```shell
aws-console
```

Or for the named "production" profile:
```shell
aws-console production
```

Or from the output of the aws cli itself:
```shell
aws sts assume-role … | aws-console -
```

---

Open the generated URL using the default browser:
```shell
aws-console --browser
```

Or copy the URL to the system clipboard:
```shell
aws-console --clipboard
```

---

Display the generated URL in the terminal as a QR code:
```shell
aws-console --qr
```

Or save it as an image to a file:
```shell
aws-console --qr > qr.png
```

---

Limit session duration to half an hour:
```shell
aws-console --duration 30m
```

Redirect to the IAM service after logging in:
```shell
aws-console --redirect iam
```

---

Federate the user and use the name "audit":
```shell
aws-console --name audit
```

Attach a readonly policy to the federated user:
```shell
aws-console --policy readonly
```

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[github-actions-badge]:  https://github.com/joshdk/aws-console/workflows/Build/badge.svg
[github-actions-link]:   https://github.com/joshdk/aws-console/actions
[github-release-badge]:  https://img.shields.io/github/release/joshdk/aws-console/all.svg
[github-release-link]:   https://github.com/joshdk/aws-console/releases
[license-badge]:         https://img.shields.io/badge/license-MIT-green.svg
[license-file]:          https://github.com/joshdk/aws-console/blob/master/LICENSE.txt
[license-link]:          https://opensource.org/licenses/MIT
