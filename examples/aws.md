# Configuration for AWS

## Setting up Access to AWS API

nginx-asg-sync uses the AWS API to get the list of IP addresses of the instances of an Auto Scaling group. To access the AWS API, nginx-asg-sync must have credentials. To provide credentials to nginx-asg-sync:

1. [Create an IAM role](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) and attach the predefined `AmazonEC2ReadOnlyAccess` policy to it. This policy allows read-only access to EC2 APIs.
2. When you launch the NGINX Plus instance, add this IAM role to the instance.

## nginx-asg-sync Configuration

nginx-asg-sync is configured in **/etc/nginx/config.yaml**.


```yaml
region: us-west-2
api_endpoint: http://127.0.0.1:8080/api
sync_interval_in_seconds: 5
cloud_provider: AWS
upstreams:
 - name: backend-one
   autoscaling_group: backend-one-group
   port: 80
   kind: http
 - name: backend-two
   autoscaling_group: backend-two-group
   port: 80
   kind: http
```

* The `api_endpoint` key defines the NGINX Plus API endpoint.
* The `sync_interval_in_seconds` key defines the synchronization interval: nginx-asg-sync checks for scaling updates every 5 seconds.
* The `cloud_provider` key defines a cloud provider that will be used. The default is `AWS`. This means the key can be empty if using AWS. Possible values are: `AWS`, `Azure`.
* The `region` key defines the AWS region where we deploy NGINX Plus and the Auto Scaling groups. Setting `region` to `self` will use the EC2 Metadata service to retreive the region of the current instance.
* The `upstreams` key defines the list of upstream groups. For each upstream group we specify:
  * `name` – The name we specified for the upstream block in the NGINX Plus configuration.
  * `autoscaling_group` – The name of the corresponding Auto Scaling group.
  * `port` – The port on which our backend applications are exposed.
  * `kind` – The protocol of the traffic NGINX Plus load balances to the backend application, here `http`. If the application uses TCP/UDP, specify `stream` instead.