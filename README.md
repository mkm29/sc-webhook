# SecurityContext Kubernetes Admission Webhook

```yaml
Author: Mitch Murphy
Date: 2022-06-26
```

The following contains a simple Go application that both validates and mutates incoming Pod creation requests to ensure that no Pod is run as root.  

## Setup

There are a few things that must be done in order to configure your Kubernetes cluser to use these webhooks. I have breifly defined these below.  

1. Build the webhook Docker image and push to a registry. You must update the `Deployment` to reference this image in `./dev/manifests/webhook/webhook.deploy.yaml`  

2. Since Admission Webhooks can only communicate via TLS, certs must be generated. Here we name the service `security-webhook` and deploy it in the `default` namcespace, so that must be listed as the CN in our certs. Furthermore, we need to specify what the DNS will be in the SAN for our cert (in this case `security-webhook.default.svc`). Once the certs have been generated, the corresponding entries must be updated in the mutating/validating webhook configuration manifests, as well creating a TLS secret in Kubernetes (that will be mounted to the webhook container). All of this has been codified in `./dev/gen-certs.sh`  

3. The order that you deploy resources to Kubernetes does matter:  

  > `kubectl apply -f dev/manifests/webhook/webhook.tls.secret.yaml`  
  > `kubectl apply -f dev/manifests/webhook/webhook.svc.yaml`  
  > `kubectl apply -f dev/manifests/webhook/webhook.deploy.yaml`  
  > `kubectl apply -f dev/manifests/webhook/cluster-config/mutating.config.yaml`  
  > `kubectl apply -f dev/manifests/webhook/cluster-config/validating.config.yaml` 

4. Now that the webhook has been deployed and configured, you are ready to test it. Notice that in both of our webhooks we added a `namespaceSelector` so that our webhook is only applied to namespaces that have the label  `admission-webhook: enabled`. 

  > `kubectl apply -f ./dev/manifests/cluser-config/apps.ns.yaml`  

5. Try running an NGINX Pod and observe what happens. Since NGINX listens on ports 80/443 (protected) it must run as root, therefore after deploying you will notice that it adds the security context to our Pod, however, it will not be able to run.  

  > `kubectl run root-nginx --image nginx --restart Never --namespace apps`  
  > `kubectl describe pod/root-nginx -n apps`

6. I have modified NGINX to run as the built in nginx user (see the [Dockerfile](dev/nginx/Dockerfile) as well as the corresponding configuration files). Therefore, the webhook adds the security context and the container is able to run fine. (Build and push the image prior to running.)

  > `kubectl run nginx --image localhost:5000/my-nginx:latest --restart Never --namespace apps`  
  > `kubectl get po nginx -n apps -o jsonpath='{.spec.securityContext}'`  

7. I have also included a little Flask app (and [Dockerfile](dev/flask/Dockerfile)). One thing to note that in all Docker files you need to change the permissions on certain directories, as well as changing the running user with the `USER` statement (this must be the numeric value for the user). Note that in the NGINX case the creators already created an nginx user, you just need to configure it to use that.

  > `kubectl run flask --image localhost:5000/my-flask:latest --restart Never --namespace apps`  
  > `kubectl get po flask -n apps -o jsonpath='{.spec.securityContext}'` 