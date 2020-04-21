[![made-with-golang](https://img.shields.io/badge/Made%20with-Golang-blue.svg?style=flat-square)](https://golang.org/)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%203-blue.svg?style=flat-square)](https://github.com/l3uddz/crop/blob/master/LICENSE.md)
[![last commit (develop)](https://img.shields.io/github/last-commit/l3uddz/crop/develop.svg?colorB=177DC1&label=Last%20Commit&style=flat-square)](https://github.com/l3uddz/crop/commits/develop)
[![Discord](https://img.shields.io/discord/381077432285003776.svg?colorB=177DC1&label=Discord&style=flat-square)](https://discord.io/cloudbox)
[![Contributing](https://img.shields.io/badge/Contributing-gray.svg?style=flat-square)](CONTRIBUTING.md)
[![Donate](https://img.shields.io/badge/Donate-gray.svg?style=flat-square)](#donate)

# crop

CLI tool to run upload/sync jobs with rclone.

## Example Configuration

```yaml
rclone:
  config: /home/seed/.config/rclone/rclone.conf
  path: /usr/bin/rclone
  stats: 30s
uploader:
  google:
    check:
      limit: 360
      type: age
    enabled: true
    hidden:
      cleanup: true
      enabled: true
      folder: /mnt/local/.unionfs-fuse
      type: unionfs
    local_folder: /mnt/local/Media
    remotes:
      clean:
        - 'gdrive:'
        - 'staging:'
      move: 'staging:/Media'
      move_server_side:
        - from: 'staging:/Media'
          to: 'gdrive:/Media'
    sa_folder: /opt/service_account_maker/service_accounts/crop
    rclone_params:
      move:
        - '--transfers=8'
        - '--delete-empty-src-dirs'
      move_server_side:
        - '--delete-empty-src-dirs'
      dedupe:
        - '--tpslimit=50'
syncer:
  4k_movies:
    enabled: true
    sa_folder: /opt/service_account_maker/service_accounts/crop
    source_remote: 'source_4k_movies:/'
    remotes:
      sync:
        - '4k_movies:/'
      dedupe:
        - '4k_movies:/'
    rclone_params:
      sync:
        - '--fast-list'
        - '--tpslimit-burst=50'
        - '--max-backlog=2000000'
        - '--track-renames'
        - '--use-mmap'
      dedupe:
        - '--tpslimit=5'
```

## Example Commands

- Clean - Perform clean for associated uploader job(s).

`crop clean --dry-run`

`crop clean -u google`

`crop clean`

- Upload - Perform uploader job(s)

`crop upload --dry-run`

`crop upload -u google`

`crop upload -u google --no-check`

`crop upload`

- Sync - Perform syncer job(s)

`crop sync --dry-run`

`crop sync -s google`

`crop sync`

***

## Notes

Make use of `--dry-run` and `-vv` to ensure your configuration is correct and yielding expected results.

## Credits

- [rclone](https://github.com/rclone/rclone) - Without this awesome tool, this project would not exist!

# Donate

If you find this project helpful, feel free to make a small donation to the developer:

  - [Monzo](https://monzo.me/today): Credit Cards, Apple Pay, Google Pay

  - [Paypal: l3uddz@gmail.com](https://www.paypal.me/l3uddz)
  
  - [GitHub Sponsor](https://github.com/sponsors/l3uddz): GitHub matches contributions for first 12 months.

  - BTC: 3CiHME1HZQsNNcDL6BArG7PbZLa8zUUgjL