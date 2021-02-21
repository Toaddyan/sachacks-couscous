# couscous

Deploy [`code-server`](https://github.com/cdr/code-server) to [Google Cloud Shell](https://cloud.google.com/shell).

# Usage

``` bash
ansible-playbook -i inventory.yaml couscous.yaml
```

# `code-server`

## Cloud Program

1) Start `code-server` with the `--link [optional name]` flag

``` bash
code-server --link couscous
```

2) Authorize with Github
3) Click the link
## [Live Share Extension](https://marketplace.visualstudio.com/items?itemName=MS-vsliveshare.vsliveshare)

The Live Share Extension enables developers to collaborate realtime from within VS Code.

### Installation

Due to the closed-source nature of the VS Code Extension marketplace, [`code-server` is unable to connect to the _actual_ VS Code Extension Marketplace](https://github.com/cdr/code-server/issues/30#issuecomment-470146035). Only [open source extensions](https://github.com/cdr/code-server/pull/113) are available in the extensions explorer; however, closed source extensions can still be manually installed via a `.vsix` package.

### `Enable Proposed API`

The `--enable-proposed-api` option does not seem to work, success was found by [directly modifying the `product.json` file](https://github.com/MicrosoftDocs/live-share/issues/262).

### Authentication via User Code

Authentication for the Live Share extension can sometimes be unreliable. I have found it easier to use a "user code" for authentication which can be authorized at:
<https://insiders.liveshare.vsengsaas.visualstudio.com/auth/login>

# Cloud Shell

Cloud Shell is an online development and operations environment accessible from a browser. Cloud Shell conveniently creates a temporary development environment with 5GB of persistent disk storage.

## Limitations / Shortcomings

### Short-lived Sessions

Cloud Shell instances terminate after 20 minutes of inactivity. When using `gcloud` to connect to Cloud Shell, the `ServerAliveInterval` configuration can prolong the session.

``` bash
gcloud cloud-shell ssh --authorize-session --ssh-flag="-o ServerAliveInterval=30"
```

### Default password

[There is no default password](https://serverfault.com/a/813296), it must be manually configured.

```bash
sudo passwd $USER
```

### Persistent storage is minimal / utilizing `rclone`

Persistent storage is limited to 5GB located at `$HOME`. `rclone` can be utilized to overcome this limitation.

``` bash
sudo -s # switch to root
vim /etc/fuse.conf # change /etc/fuse.conf to allow_other
mkdir /mnt/cloudshell # create mounting point
chown $SUDO_USER:$SUDO_USER /mnt/cloudshell
exit # exit root
rclone --vfs-cache-mode full --cache-dir /tmp/rclone-cache -vv mount cloudshell: /mnt/cloudshell/ --allow-other
```

#### Problems

`rclone` mounts are not perfect. Problems may frequently occur; for instance, this message frequently appears in `vim`:

> WARNING: The file has been changed since reading it!!!
> Do you really want to write to it (y/n)?

#### A good example of how to fix these problems!

It may be possible to fix these problems by taking inspiration from [this project](https://github.com/animosity22/homescripts). Which was found courtesy of [this](https://forum.rclone.org/t/mounting-google-drive-potential-pitfalls/11600/9) forum post.
