# Digital Ocean Janitor

A GitHub Action to cleanup Digital Ocean resources that have exceeded a TTL.

> By default the action will not perform the delete (i.e. it will be a dry-run). You need to explicitly set commit to `true`.

It supports cleaning up the following services:

- Droplets

## Inputs

| Name              | Required | Description                                                                            |
| ----------------- | -------- | -------------------------------------------------------------------------------------- |
| regions           | Y        | A comma separated list of regions to clean resources in. Use the slug (i.e. nyc1,lon1). You can use * for all regions. |
| allow-all-regions | N        | Set to true if use * from regions.                                                     |
| ttl               | Y        | The duration that a resource can live for. For example, use 24h for 1 day.             |
| commit            | N        | Whether to perform the delete. Defaults to `false` which is a dry run                  |
| token             | Y        | The Digital Ocean API token to use. It must have Write scope                           |
| ignore-tag        | N        | The name of the tag that indicates a resource should not be deleted.                   |

## Example Usage

For all regions:

```yaml
jobs:
  cleanup:
    runs-on: ubuntu-latest
    name: Cleanup resource groups
    steps:
      - name: Cleanup
        uses: rancher-sandbox/digitalocean-janitor@v0.1.0
        with:
            regions: *
            allow-all-regions: true
            ttl: 168h
            token: {{ secrets.DO_TOKEN }}
            ignore-tag: DO_NOT_DELETE
```

For specific regions:

```yaml
jobs:
  cleanup:
    runs-on: ubuntu-latest
    name: Cleanup resource groups
    steps:
      - name: Cleanup
        uses: rancher-sandbox/digitalocean-janitor@v0.1.0
        with:
            regions: nyc1,lon1
            ttl: 168h
            token: {{ secrets.DO_TOKEN }}
            ignore-tag: DO_NOT_DELETE
    

## Implementation Notes

It currently assumes that an instance of a service will have some form of creation date. This means that the implementation can be simpler as it doesn't need to adopt a "mark & sweep" pattern that requires saving state between runs of the action.

