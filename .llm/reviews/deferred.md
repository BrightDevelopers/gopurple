# Deferred Findings

## Pre-existing issues not introduced by this change

### HIGH: URL parameter injection in internal/services/rdws.go
User-controlled parameters (serial, path, fileName, newName) are interpolated directly into URLs without URL encoding. Should use `url.PathEscape()` for path components and `url.QueryEscape()` for query parameters. Pre-dates this change entirely.

### LOW: Personal email in configs/lfn-control.json
Contains `gherlein@brightsign.biz`. Pre-existing developer config file. Consider adding to .gitignore.
