Vagrant.configure("2") do |config|
    config.vm.box = "fr123k/ubuntu21-pulumi"
    config.ssh.extra_args = ["-t", "cd /vagrant/; rm -rf ./build; mkdir -p /tmp/vagrant/target/; ln -s /tmp/vagrant/target/ ./build; bash --login"]
end
