# Push Deploy

**Warning: This tool is still in early development and not ready for use yet!**

Deploy jekyll websites using custom plugins to github pages without having to build them locally first.

Push Deploy listen to webhookss from github & bitbucket to build your Jekyll site and deploy it to a git repo, s3 bucket or FTP server of your choice.


## Installing Dependencies

### Install git

```
apt-get update
apt-get install -y git
```

### Install ruby (RVM)

jekyll-deploy will work with just ruby, however using RVM is recommended.

```
gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
\curl -sSL https://get.rvm.io | bash -s stable --ruby

```

If you'd prefer not to use RVM you can install ruby with following command.

```
apt-get install ruby-full

```
