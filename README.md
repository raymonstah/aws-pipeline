# AWS Pipeline
* A tool to help deploy lambdas to AWS

## Getting started:


```
rho:aws-pipeline rho$ go run main.go 
NAME:
   main - a deployment tool for AWS

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   lambda          upload zipped lambdas from a path to s3
   cloudformation  deploy cloudformation stack
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
   ```


### Lambdas
The tool assumes that there are zip files in the specified directory and will upload them a given s3 bucket. If the s3 bucket does not exist, it will create a versioned S3 bucket. The bucket **must** be versioned to keep track of Lambdas.
```bash
go run main.go lambda help
NAME:
   main lambda - upload zipped lambdas from a path to s3

USAGE:
   main lambda [command options] [arguments...]

OPTIONS:
   --bucket value       location of where the zipped lambdas should be stored
   --target-path value  path to artifact resources (zip files for lambdas)
   --help, -h           show help (default: false)
```

Example:
```
go run main.go lambda --target-path /example-diretory --bucket nameOfBucket
```


### Cloudformation
The tool will parse a given cloudformation template, and pass any necessary parameters to Cloudformation.

It's responsible for creating or updating a given stack.
```
go run main.go cloudformation help
NAME:
   main cloudformation - deploy cloudformation stack

USAGE:
   main cloudformation [command options] [arguments...]

OPTIONS:
   --path value            path to cloudformation template
   --stack-name value      name of the cloudformation stack
   --lambdas-bucket value  optional -- location of the zipped lambdas
   --help, -h              show help (default: false)
```

Define a parameter `LambdasBucket` in your cloudformation template to specify the name of the S3 bucket where the zipped up lambda files can be found. This can be omitted if the stack doesn't use Lambda. Use the parameter `--lambdas-bucket` when running the tool to set the value.

Sample Cloudformation script:
```yaml
AWSTemplateFormatVersion: "2010-09-09"

Parameters:
  # The name of the bucket where all the lambdas are, if any.. this gets injected from the pipeline tool
  LambdasBucket:
    Type: String
  # Create a parameter for every lambda like so below. The pipeline tool will inject the version into this parameter so that it can be used when defining the `S3ObjectVersion`  
  SampleFunctionZip:
    Type: String

Resources:
  SampleFunction:
    Type: "AWS::Lambda::Function"
    Properties:
      Code:
        S3Bucket: !Ref LambdasBucket
        S3Key: "sample-function.zip"
        # The version of the lambda that gets injected by the pipeline tool
        S3ObjectVersion: !Ref SampleFunctionZip

```
Example:
```
go run main.go cloudformation --path cloudformation.yml --stack-name sample-stack --lambdas-bucket lamda-bucket-here
```