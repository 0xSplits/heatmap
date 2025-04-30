# painpoints

This project helps to expose certain costs of the engineering process. Teams can
define pain points that they want to assess in a data driven way and integrate
the `painpoints` binary together with its Github Action into any pull request
workflow.

The idea is for engineers to provide insights into their respective pain points
with every pull request they work on. The collected data can then be visualized
as heatmap by the `painpoints` binary. The binary can generate a heatmap using
fake data via the test flag `-t` like shown below.

```
go run . -t
```

![Heatmap Default](./asset/heatmap_default.png)

The buckets used to create the colour gradient default to `0,10,20,30,40,50`.
The brightest green is used for any number within the half open interval `[0,
10)`, which means that `10` is excluded, because it starts the next bucket.
Those buckets can be configured using the env var `PAINPOINTS_DATA_BUCKET`. Note
that the number of buckets must match the number of colours, which is 6.

```
PAINPOINTS_DATA_BUCKET="0,5,10,15,20,25" go run . -t
```

![Heatmap Bucket](./asset/heatmap_bucket.png)

### Data Dir

The `painpoints` binary looks for files in the data dir, which defaults to
`./painpoints/`. This data dir can be configured using the `PAINPOINTS_DATA_DIR`
env var. The file names within that data dir must be numerical, indicating the
pull request number of any given pull request for which data is being collected.
Note that an optional `README.md` can be placed in the data dir, which may be
useful to document any team specific objective. Other than that, no other file
types are allowed.

```
 % tree ./painpoints
./painpoints
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
├── 7
└── README.md

1 directory, 12 files
```

### Data Files

Every data file must start with the numerical data being collected. The
collected data is allowed to have an optional suffix, which defaults to `%` and
can be modified using the `PAINPOINTS_DATA_SUFFIX` env var. Further, data files
may contain any relevant notes after the first line. Providing more context
about the data being collected per pull request may be useful to identify the
most important pain points over time.

```
% cat ./painpoints/277
36%

- copying production data for testing took 3 hours
- verifying a single endpoint manually took 15 minutes
- we had to verify 8 different endpoints
```
