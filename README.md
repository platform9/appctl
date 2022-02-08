# appctl
**Run apps, Not Clusters.** Deploy your app on kubernetes in seconds,with no clusters required. Check for more information at https://platform9.com/appctl/

* Read the docs: [getting started with appctl](https://platform9.com/docs/appctl/getting-started)

### Purpose
CLI to deploy & manage apps in Platform9 environment.

### A better way to run apps on K8s
*  **appctl** exposes the high value app orchestration capabilities available from Kubernetes and k-native, while hiding infrastructure complexity. 
* As a result, it is much faster to run apps while also running them more cost effectively in the cloud

### How appctl works
![flow-diagram](images/graphic_how-appctl-works.png)

### Usage
- Downloading the CLI can be done from [appctl website](https://platform9.com/appctl/) or using curl.

**For Linux**
```sh
curl -O https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/linux/appctl

chmod +x appctl
```

**For Mac**
```sh
curl -O https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/macos/appctl

chmod +x appctl
```

**For Windows**
```sh
curl -O https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/windows/appctl

After successfull download give the executable permission to appctl.
```
