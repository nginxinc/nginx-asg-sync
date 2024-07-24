# Tests

The project includes automated tests for testing the ASG Sync tool in an AWS environment (Azure support will be added
later). The tests are written in Python3, use the pytest framework to run the tests and utilize boto3 to call the AWS
API.

Note: for now this is for internal use only, as AWS stack configuration is done outside of testing framework.

Below you will find the instructions on how to run the tests against a cloud provider.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Running Tests](#running-tests)
  - [Prerequisites](#prerequisites)
    - [Step 1 - Set up the environment](#step-1---set-up-the-environment)
    - [Step 2 - Run the Tests](#step-2---run-the-tests)
- [Configuring the Tests](#configuring-the-tests)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Running Tests

### Prerequisites

- AWS stack prepared.
- AWS access key and AWS secret access key.
- Python3 or Docker.

#### Step 1 - Set up the environment

- Either create|update ~/.aws/credentials file or set the AWS_SHARED_CREDENTIALS_FILE environment variable pointing to
  your own location. This file is an INI formatted file with section names corresponding to profiles. Tests use
  'default' profile. The file [credentials](data/credentials) is a minimal example of such a file.

#### Step 2 - Run the Tests

Run the tests:

- Use local Python3 installation:

    ```bash
    cd tests
    pip3 install -r requirements.txt
    python3 -m pytest --nginx-api=nginx_plus_api_url
    ```

- Use Docker:

    ```bash
    cd tests
    make run-tests AWS_CREDENTIALS=abs_path_to_creds_file NGINX_API=nginx_plus_api_url
    ```

## Configuring the Tests

The table below shows various configuration options for the tests. If you use Python3 to run the tests, use the
command-line arguments. If you use Docker, use the [Makefile](Makefile) variables.

| Command-line Argument | Makefile Variable | Description | Default |
| :----------------------- | :------------ | :------------ | :----------------------- |
| `--nginx-api` | `NGINX_API` | The NGINX Plus API url. | `N/A` |
| `--aws-region` | `AWS_REGION` | The AWS stack region. | `us-east-2` |
| `N/A` | `PYTEST_ARGS` | Any additional pytest command-line arguments (i.e `-k TestSmoke`) | `""` |

If you would like to use an IDE (such as PyCharm) to run the tests, use the [pytest.ini](pytest.ini) file to set the
command-line arguments.
