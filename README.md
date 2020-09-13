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
  live_rotate: false
  service_account_remotes:
    '/opt/rclone/service_accounts/crop':
      - tv
      - movies
      - music
      - 4k_movies
      - source_4k_movies
      - staging
  global_params:
    default:
      move:
        - '--order-by=modtime,ascending'
        - '--transfers=8'
        - '--delete-empty-src-dirs'
      sync:
        - '--fast-list'
        - '--tpslimit-burst=50'
        - '--max-backlog=2000000'
        - '--track-renames'
        - '--use-mmap'
        - '--no-update-modtime'
        - '--drive-chunk-size=128M'
      dedupe:
        - '--dedupe-mode=newest'
        - '--tpslimit=5'
uploader:
  - name: cloudbox_unionfs
    enabled: true
    check:
      limit: 360
      type: age
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
    rclone_params:
      global_move: default
      move_server_side:
        - '--delete-empty-src-dirs'
      global_dedupe: default
  - name: tv
    enabled: true
    check:
      limit: 1440
      type: age
    local_folder: /mnt/local/Media/TV
    remotes:
      move: 'tv:/Media/TV'
    rclone_params:
      global_move: default
  - name: movies
    enabled: true
    check:
      limit: 720
      type: age
    local_folder: /mnt/local/Media/Movies
    remotes:
      move: 'movies:/Media/Movies'
    rclone_params:
      global_move: default
syncer:
  - name: 4k_movies
    enabled: true
    source_remote: 'source_4k_movies:/'
    remotes:
      sync:
        - '4k_movies:/'
      dedupe:
        - '4k_movies:/'
    rclone_params:
      global_sync: default
      global_dedupe: default
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

`crop sync -p 2`

- Manual - Perform manual sync/copy job(s)

`crop manual --copy --src remote1:/Backups --dst remote2:/Backups --sa /opt/service_accounts -- --dry-run`

`crop manual --sync --src remote1:/Backups --dst remote2:/Backups --sa /opt/service_accounts --dedupe --`

***

## Notes

- Make use of `--dry-run` and `-vv` to ensure your configuration is correct and yielding expected results.

- `live_rotate` will enable on-demand live-rotation of service accounts for a customized build of rclone / gclone.


## Credits

- [rclone](https://github.com/rclone/rclone) - Without this awesome tool, this project would not exist!
- [sasync](https://github.com/88lex/sasync) - Sync ideas and service account technique originated from here.

# Donate

If you find this project helpful, feel free to make a small donation to the developer:

  - [Monzo](https://monzo.me/today): Credit Cards, Apple Pay, Google Pay

  - [Paypal: l3uddz@gmail.com](https://www.paypal.me/l3uddz)
  
  - [GitHub Sponsor](https://github.com/sponsors/l3uddz): GitHub matches contributions for first 12 months.

  - BTC: 3CiHME1HZQsNNcDL6BArG7PbZLa8zUUgjL
