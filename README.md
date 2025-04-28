# heatmap

This project generates heatmaps for testing estimates based on a specific folder
structure. The binary can generate a heatmap using fake data using the test flag
`-t` like shown below.

```
go run . -t
```

![Example Heatmap](./asset/heatmap.png)

### File Structure

The `heatmap` binary looks for files in the data dir, which defaults to
`./manual-testing/` and can be configured using the `HEATMAP_DATA_DIR` env var.
The file names within that data dir must be numerical, indicating the pull
request number of any given pull request.

```
% tree manual-testing
manual-testing
├── 0
├── 103
├── 12
├── 23
├── 277
├── 37
├── 41
├── 507
├── 52
├── 68
└── 7

1 directory, 11 files
```

The content of any given file must contain the estimated percentage of time
spent for manual testing. Optionally any relevant notes may be added to the file
below the percentage estimate. Providing more context may be useful to identify
the most important pain points.

```
% cat manual-testing/277
36%

- copying production data for testing took 3 hours
- verifying a single endpoint manually took 15 minutes
- we had to verify 8 different endpoints
```
