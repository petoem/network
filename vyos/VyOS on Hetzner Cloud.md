# Hetzner Cloud with VyOS

## Install latest rolling VyOS

- Create Cloud VPS with any OS you want.
- Go to VPS settings and mount old VyOS ISO from Library.
- Boot VPS with mounted ISO.
- Login to VyOS and run "install image".
- Take default values.
- Install ...
- Unmount and Reboot.
- Fix Network connection with:

```sh
sudo ip link set dev eth0 up
sudo ip route add 172.31.1.1/32 eth0
sudo ip route add default via 172.31.1.1
```

- Enable DHCP and DHCPv6 on eth0.
- Internet should work now! ;)
- Add latest VyOS image to system.

```sh
add system image https://downloads.vyos.io/rolling/current/amd64/vyos-rolling-latest.iso
```

- Reboot server. (You should now be in the latest VyOS system. If not use `set system image default-boot`.)
- Change password of vyos user or existing config!

## Use Apt repos on VyOS

Add the following to `/config/scripts/vyos-postconfig-bootup.script`:

```sh
# Append to sources.list the debian buster mirrors
tee -a /etc/apt/sources.list << END
deb http://deb.debian.org/debian buster main contrib non-free
deb-src http://deb.debian.org/debian buster main contrib non-free

deb http://deb.debian.org/debian-security/ buster/updates main contrib non-free
deb-src http://deb.debian.org/debian-security/ buster/updates main contrib non-free

deb http://deb.debian.org/debian buster-updates main contrib non-free
deb-src http://deb.debian.org/debian buster-updates main contrib non-free

deb http://deb.debian.org/debian buster-backports main contrib non-free
deb-src http://deb.debian.org/debian buster-backports main contrib non-free
END
```
