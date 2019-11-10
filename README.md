# i

> i is an Operator for checking due date for policy check.

This is not for production usage. It is under PoC.

## Design

### Architecture

The controller runs syncer by the resource from input and watches the pod.

```
Controller --> Syncer -(watch)-> Pod
``` 

And then delete the pod and recreates.  
At the time, the controller will notify via slack which will be configured by Secret resource.

### Install CRD

First, you will need to install this CRD.  
Run this command.

```terminal
$ curl "" | kubectl apply -f 
```

Then, check crd status.

```terminal
$ kubectl get crd
```

### Create secret for slack notification

Create a secret.

```yaml
apiVersion:
Kind: Secret
metadata:
    app: slack-channel-1
data:
    SLACK_URL: "https://hook.xxx.xxx"
```

And then apply.

```terminal
$ kubectl apply -f slack-channel-1.yml
```

### Eye resource

Create a custom resource. You can create many rules.

```yaml
apiVersion: i.keisukeyamashita.com/alphav1
Kind: Eye
metadata:
    app: my-eye
    msg: "I see you"
spec:
    lifetime: "100m"
    secretRef:
        name: YOUR_SECRET_NAME
```

Then apply your resource.

```terminal
$ kubectl aaply -f my-eye.yml
```

Check you status.

```terminal
$ kubectl get eye
```

## Author

* [KeisukeYamashita](https://github.com/KeisukeYamashita)
