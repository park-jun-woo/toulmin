## tangl:Subject
- this document is `comail`

## tangl:See
- see `check` from `mail/check`

## tangl:Cases
- in case of `can accept`
  - `spf ok` is a general rule using `check`.`SpfPass` grounded on `source-body`
  - `source lacks email` is a counter rule using `check`.`SourceLacksEmail` grounded on `source-body` and `mx`, level fail
