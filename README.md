# Pak - WIP

> Metadata ready packages

Pak is a tool to generate packaes with custom metadata, you define your packages using a `pakfile.json` that looks something like this:
```json
{
  "artifact_name": "artifactName",
  "path": "./bin",
  "metadata": {
      "mantainer": "bronzdoc@mail.com",
      "build_id": "${BUILD_ID}"
    }
}
```

> pak can read environment variables so you can use them inside your pakfile.json as `${VAR_NAME}`

# Usage

```shell
$ pak build
```

pak will look for a pakfile.json and build a packae from that

