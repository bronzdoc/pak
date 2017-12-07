# Pak

> Build, inspect and promote artifacts with ease

Pak is a tool to generate build artifacts with custom metadata, that you can later promote.

## Usage

You define your artifacts using a `Pakfile.json` that looks like this:

```json
{
  "name": "ALPHA-ARTIFACT",
  "path": "./bin",
  "metadata": {
      "mantainer": "bronzdoc@mail.com",
      "BUILD_DATE": "${BUILD_DATE}"
    },

  "promote":{
    "rc":{
      "name": "RELEASE-ARTIFACT",
      "metadata": {
      "REL_DATE": "${BUILD_DATE}"
      }
    }
  }
}
```

> pak can read environment variables so you can use them inside your Pakfile.json as `${VAR_NAME}`

### Build

```shell
$ export BUILD_DATE=$(date)
$ pak build
```

this will generate a .tar file named `ALPHA-ARTIFACT.tar`

### Inspect

After you build an artifact using 'pak' you can inspect it's metadata with the inspect subcommand

```shell
$ pak inspect ALPHA-ARTIFACT.tar

{
  "build": {
    "metadata": {
      "BUILD_DATE": "lun dic  4 16:57:23 CST 2017",
      "mantainer": "bronzdoc@mail.com"
    },
    "name": "ALPHA-ARTIFACT"
  },
  "rc": {
    "metadata": {
      "REL_DATE": "${REL_DATE}"
    },
    "name": "RELEASE-ARTIFACT"
  }
}
```

You can pass the `--key-value` flag to the inspect subcommand to display the data in a key=value manner, i.e:

```shell
$ pak inspect --key-value ALPHA-ARTIFACT.tar

#build
  BUILD_DATE="lun dic  4 16:57:23 CST 2017"
  mantainer="bronzdoc@mail.com"
  name="ALPHA-ARTIFACT"
#rc
  REL_DATE="${REL_DATE}"
  name="RELEASE-ARTIFACT"
```

you can notice there is a `build` label that we didn't define, this is special and is the metadata `pak` created when you ran `pak build`.

You can access specific metadata if you pass to inspect the lable of the metadata you want to access, i.e:
```shell
$ pak inspect ALPHA-ARTIFACT.tar build

{
  "metadata": {
    "BUILD_DATE": "lun dic  4 16:57:23 CST 2017",
    "mantainer": "bronzdoc@mail.com"
  },
  "name": "ALPHA-ARTIFACT"
}
```

### Promote

```shell
$ export REL_DATE=$(date)
$ pak promote ALPHA-ARTIFACT.tar rc
```

This will generate a new artifact named `RELEASE-ARTIFACT.tar` with the old metadata and the newone.
Notice we passed the `promote label` we wanted the artifact to be promoted to, in this case this is `rc`.

If you inspect the new artifact you'll notice the artifact has the old and new metadata stored in it

```shell
$ pak inspect RELEASE-ARTIFACT.tar
{
  "build": {
    "metadata": {
      "BUILD_DATE": "lun dic  4 16:57:23 CST 2017",
      "mantainer": "bronzdoc@mail.com"
    },
    "name": "ALPHA-ARTIFACT"
  },
  "rc": {
    "metadata": {
      "REL_DATE": "mar dic  5 15:50:29 CST 2017",
    },
    "name": "RELEASE-ARTIFACT"
  }
```

## Install

### Binaries

- **linux** [386](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-linux-386) / [amd64](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-linux-amd64) / [arm](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-linux-arm) / [arm64](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-linux-arm64)
- **darwin** [386](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-darwin-386) / [amd64](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-darwin-amd64)
- **freebsd** [386](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-freebsd-386) / [amd64](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-freebsd-amd64)
- **windows** [386](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-windows-386) / [amd64](https://github.com/bronzdoc/pak/releases/download/v0.1.0/pak-windows-amd64)
