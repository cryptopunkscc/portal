# Configuring KVM with Windows for testing

Experimental support for end-to-end testing on Windows OS.

## Requirements

* ssh client
* [kvm & virsh](/doc/kvm.md)

## Setup

1. Download `Win10_22H2_English_x64v1.iso` and put into [./iso/static](./iso/static).
2. Download `virtio-win-0.1.271.iso` and put into [./iso/static](./iso/static)
3. Run [vm_factory_test.go](./vm_factory_test.go) to create Window virtual machine for testing.
   * This process may take a while.
   * After installation only the admin account is activated by default.
   * Activating the user's account requires manually logging in to it.
   * End-to-end testing depends on the user's account.
   * Subsequent run of `vm_factory_test.go` will overwrite previously created virtual drive.

## Development

* To adjust virtual machine setup, modify [vm_factory.go](./vm_factory.go)
* Windows OS installation is automated via [autounattend.xml](./iso/config/autounattend/autounattend.xml).
* Use https://schneegans.de/windows/unattend-generator/ to generate or modify `autounattend.xml`.

