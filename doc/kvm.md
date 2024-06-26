# Installation for debian based distributions

source: https://www.linuxtechi.com/how-to-install-kvm-on-ubuntu-22-04/

Validate virtualization enabled i enabled, output > 0

```shell
egrep -c '(vmx|svm)' /proc/cpuinfo
```

Validate kvm ok

```shell
kvm-ok
```

Install kvm-ok if needed

```shell
sudo apt install -y cpu-checker
```

Install kvm dependencies

```shell
$ sudo apt install -y qemu-kvm virt-manager libvirt-daemon-system virtinst libvirt-clients bridge-utils
```

* qemu-kvm  – An opensource emulator and virtualization package that provides hardware emulation.
* virt-manager – A Qt-based graphical interface for managing virtual machines via the libvirt daemon.
* libvirt-daemon-system – A package that provides configuration files required to run the libvirt daemon.
* virtinst – A set of command-line utilities for provisioning and modifying virtual machines.
* libvirt-clients – A set of client-side libraries and APIs for managing and controlling virtual machines & hypervisors from the command line.
* bridge-utils – A set of tools for creating and managing bridge devices.

Enable and start the Libvirt daemon.

```shell
sudo systemctl enable --now libvirtd
```
```shell
sudo systemctl start libvirtd
```

Confirm that the virtualization daemon is running.

```shell
sudo systemctl status libvirtd
```

Add the currently logged-in user to the kvm and libvirt groups.

```shell
sudo usermod -aG kvm $USER
```
```shell
sudo usermod -aG libvirt $USER
```
(Optionally?) To create a network bridge, create the file `01-netcfg.yaml` with following content under the folder `/etc/netplan`.

```shell
sudo vim /etc/netplan/01-netcfg.yaml
```
but replace the IP address entries, interface name and mac address
```yaml
network:
  ethernets:
    enp0s3:
      dhcp4: false
      dhcp6: false
  # add configuration for bridge interface
  bridges:
    br0:
      interfaces: [enp0s3]
      dhcp4: false
      addresses: [192.168.1.162/24]
      macaddress: 08:00:27:4b:1d:45
      routes:
        - to: default
          via: 192.168.1.1
          metric: 100
      nameservers:
        addresses: [4.2.2.2]
      parameters:
        stp: false
      dhcp6: false
  version: 2
```

apply config

```shell
sudo netplan apply
```

Verify the network bridge `br0`

```shell
ip add show
```

Install guest OS through Virtual Machine Manager GUI