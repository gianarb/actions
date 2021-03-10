---
slug: syslinux
name: syslinux
tags: disk, boot loader
maintainers: Dan Finneran <daniel.finneran@gmail.com>
description: "This action can be used to install syslinux as a boot loader, in two parts
it will install the MBR at `DEST_DISK` and install the boot loader on the specified 
partition at `DEST_PARTITION`. Currently the only version of syslinux available is
`3.86`."
version: v1.0.0
createdAt: "2021-03-01T12:41:45.14Z"
---

The below example will use the install the [syslinux](https://wiki.archlinux.org/index.php/syslinux) boot loader to a
specified block device. The installation consists of installing an MBE to the
first few sectors of the disk, then installing the boot loader on the FAT32
filesystem (on the specified partition).

```yaml
actions:
    - name: "install syslinux to /dev/sda"
      image: quay.io/tinkerbell-actions/syslinux:v1.0.0
      timeout: 90
      environment:
          DEST_DISK: /dev/sdb
          DEST_PARTITION: /dev/sdb1
          SYSLINUX_VERSION: 3.86
```
