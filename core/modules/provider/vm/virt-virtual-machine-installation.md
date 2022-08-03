## virt virtual machine installation process

1. make sure that virtualization is enabled

```shell
egrep 'vmx|svm' /proc/cpuinfo --color=auto
```

2. install relevant virtualization packages

```shell
## centos
yum install -y qemu-kvm libvirt virt-install libvirt-devel
systemctl start libvirtd && systemctl enable libvirtd

## ubuntu 
sudo apt install qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils virtinst virt-manager libvirt-dev
sudo systemctl is-active libvirtd
sudo usermod -aG libvirt $USER
sudo usermod -aG kvm $USER

```

3. make a virtual machine image

3.1 Download the original image from the official website 

example: `CentOS-7-x86_64-Minimal-2009.iso`,`ubuntu-20.04.3-live-server-amd64.iso`

3.2 virsh start the virtual environment

```shell
## 1 create hard disk space
qemu-img create -f raw /data/CentOS7.raw 50G
virt-install --virt-type kvm --name centos7 --ram 1024 --cdrom=./CentOS-7-x86_64-Minimal-2009.iso --disk path=/data/centos7.raw --network network=default --graphics vnc,listen=0.0.0.0 --noautoconsole

## ubuntu
sudo virt-install --virt-type kvm --name test --vcpus 1 --ram 2048 --disk path=/home/gr/app/ubuntu20.qcow2 --network network=default --graphics vnc,listen=0.0.0.0 --noautoconsole --boot hd

```


3.3 Use vnc view to connect the visualization page to complete the installation

3.4 After completing the installation, update the system and install `qemu-ga`

```shell
# rhel/centos
yum install qemu-guest-agent

# windows，latest virtio-win iso
https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/latest-virtio/
# windows，the latest qemu ga installation package
https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/latest-qemu-ga/

## centos 
# start the qemu ga daemon
systemctl start qemu-guest-agent

# join boot
systemctl enable qemu-guest-agent
```

3.5 shutdown
```shell
virsh destory centos7
```

4.Convert raw hard disk format to qcow2 (smaller size)

```shell
yum -y install libguestfs-tools
qemu-img convert -f raw -O qcow2 /data/CentOS7.raw /data/CentOS7.qcow2
```

5. export run template
```shell
virsh dumpxml centos > centos.xml
```

modify startup configdriver
```shell
<disk type='file' device='disk'>
<driver name='qemu' type='raw' cache='none'/>
<source file='/var/lib/libvirt/images/rhel62-2.img'/>
## >>> 
<disk type='file' device='disk'>
<driver name='qemu' type='qcow2' cache='none'/>
<source file='/var/lib/libvirt/images/rhel62-2.qcow2'/>


### device increase in segment
<channel type='unix'>
  <source mode='bind' path='/tmp/channel.sock'/>
  <target type='virtio' name='org.qemu.guest_agent.0'/>
</channel>

```

6. start the virtual machine

6.1 change password
```shell
virt-customize -a /data/CentOS7.qcow2 --root-password password:StrongRootPassword
```

6.2 injection public key
```shell
virt-customize -a /data/CentOS7.qcow2   --ssh-inject root:file:/root/.ssh/id_rsa.pub
```

6.3 register the virtual machine
```shell
virsh define centos.xml
```

7. virtual machine management

```shell
## define a virtual machine
virsh define centos.xml

## start the virtual machine
virsh start centos

## restart the virtual machine
virsh reboot centos

## Shut down the virtual machine (soft shutdown)
virsh shutdown centos

## Shut down the virtual machine (hard shutdown)
virsh destroy centos

## undefine a virtual machine
virsh undefine centos
```