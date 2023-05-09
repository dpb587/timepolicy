# timepolicy

Implement retention policies on the command line across a variety of object and storage types.

## Command Line Usage

Pipe lines that start with a time and configure `--policy` specifications to control which line entries match.

```shell
( echo '2023-01-01 backup-20230101.tar.gz' ; echo '2023-01-02 backup-20230102.tar.gz' ) \
  | go run ./cmd/timepolicy \
    --policy='1y;by=month' \
    --time=YYYY-MM-DD \
    --write='$2'
#> backup-20230102.tar.gz
```

See [`timepolicy --help`](https://github.com/dpb587/timepolicy/blob/resources/release/latest/timepolicy-help.txt) for full documentation and features.

### Installation

Binaries for Linux, macOS, and Windows can be downloaded from the [Releases](https://github.com/dpb587/timepolicy/releases/latest) page. A [Homebrew](https://brew.sh/) recipe is also available for Linux and macOS.

```
brew install dpb587/tap/timepolicy
```

### Examples

Prune Google Cloud snapshots based on creation time...

```shell
gcloud compute snapshots list --format='value(creationTimestamp, name)' \
  | timepolicy \
      --policy='1y;by=month // within 1 year, keep newest per month' \
      --policy='28d;by=day  // within 28 days, keep newest per day' \
      --policy='7d;by=hour  // within 7 days, keep newest per hour' \
      --write='$2' \
      --invert \
  | xargs -- \
      echo gcloud compute snapshots delete
```

Prune local backup files based on date in file name (e.g. `Hubitat_2023-05-04~2.3.5.131.lzf`)...

```shell
find . -name '*.lzf' \
  | sed -E 's#./(Hubitat_(..........).+)#\2 \1#' \
  | timepolicy \
      --policy='1y;by=month // within 1 year, keep newest by month' \
      --policy='14d;by=day  // within 14 days, keep newest by day' \
      --ts='YYYY-MM-DD' \
      --write='$2' \
      --invert \
  | xargs -n1 -- \
      echo rm
```

List selected files from Amazon S3 and total size based on modification date...

```shell
aws s3api list-objects --bucket acme-backup-us-west-1 \
  --output=text \
  --query='Contents[*].[LastModified, Size, Key]' \
  | timepolicy \
      --policy='14d;by=day;oldest // within 14 days, keep oldest by day' \
  | tee >( cut -f3 ) \
  | cut -f2 | paste -sd+ - | bc | numfmt --to=iec-i
```

## Futures

* support reading policies from files
* expand unit tests

## License

[MIT License](LICENSE)
