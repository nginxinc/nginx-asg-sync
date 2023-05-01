"""Describe overall framework configuration."""
import pytest
from boto3 import Session
from botocore.client import BaseClient

from tests.settings import DEFAULT_AWS_REGION


def pytest_addoption(parser) -> None:
    """Get cli-arguments.

    :param parser: pytest parser
    :return:
    """
    parser.addoption("--nginx-api", action="store", default="", help="The NGINX Plus API url.")
    parser.addoption(
        "--aws-region",
        action="store",
        default=DEFAULT_AWS_REGION,
        help="The AWS region name.",
    )


class CLIArguments:
    """
    Encapsulate CLI arguments.

    Attributes:
        nginx_api (str): NGINX Plus API url
        aws_region (str): AWS region name
    """

    def __init__(self, nginx_api: str, aws_region: str):
        self.nginx_api = nginx_api
        self.aws_region = aws_region


@pytest.fixture(scope="session", autouse=True)
def cli_arguments(request) -> CLIArguments:
    """
    Verify the CLI arguments.

    :param request: pytest fixture
    :return: CLIArguments
    """
    nginx_api = request.config.getoption("--nginx-api")
    assert nginx_api != "", "Empty NGINX Plus API url is not allowed"
    print(f"\nTests will use NGINX Plus API url: {nginx_api}")

    aws_region = request.config.getoption("--aws-region")
    print(f"\nTests will use AWS region: {aws_region}")

    return CLIArguments(nginx_api, aws_region)


@pytest.fixture(scope="session")
def autoscaling_client(cli_arguments) -> BaseClient:
    """
    Set up kubernets-client to operate in cluster.

    boto3 looks for AWS credentials file and uses a `default` profile from it.

    :param cli_arguments: a set of command-line arguments
    :return:
    """
    session = Session(profile_name="default", region_name=cli_arguments.aws_region)
    return session.client("autoscaling")
