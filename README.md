# API microservice
A REST API microservice orchestrated by Amazon ECS on AWS Fargate that writes tweets and parses emojis to an Amazon Aurora PostgreSQL database.

## How to create this service?
1. Install the AWS Copilot CLI [https://aws.github.io/copilot-cli/](https://aws.github.io/copilot-cli/)
2. Run
   ```bash
   $ copilot init
   ```
3. Enter "in" for the name of your application.
4. Select "Backend Service" for the service type.
5. Enter "backendpu" for the name of the service.
6. Say "Y" to deploying to a "mpk" environment 🚀

Once deployed, your service will be accessible at `http://backendpu.in.local:8080` within your VPC.  

## What does it do?
AWS Copilot uses AWS CloudFormation under the hood to provision your infrastructure resources.
You should see two different stacks created for you:
1. `in-mpk-backendpu`: Holds your ECS Service.
2. `in-mpk-backendpu-AddonsStack-<RandomString>`: Holds the Aurora database.

Take a look at the resources in the stacks to see all that's generated for you.

## How does it work?
Copilot stores the infrastructure-as-code for your service under the `copilot/` directory.
```
copilot
└── api
    ├── addons
    │   └── db.template.yaml
    └── manifest.yml
```
The `manifest.yml` file under `api/` holds the common configuration for a "backend service" pattern.
For example, in this manifest we demo how you can set up autoscaling for your service as well as container healthchecks.

The Aurora database is defined under the `api/addons/` directory which can hold any arbitrary CloudFormation template.  
The "addons" features allows you to integrate with any AWS services that are not provided by default with the Copilot manifest.

You can find out in more details how Copilot works from our [documentation](https://aws.github.io/copilot-cli/).

## Deleting the service
If you'd like to delete only the service from the "in" application.
```bash
$ copilot svc delete
```
If you'd like to delete the entire application including other services and deployment environments:
```bash
$ copilot app delete
```

## How to test locally
```
docker-compose build
docker-compose up

python manual_test.py
```
