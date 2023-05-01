import json
import time

import pytest
import requests
from botocore.waiter import WaiterModel, create_waiter_with_client

from tests.settings import NGINX_API_VERSION, RECONFIGURATION_DELAY


def wait_for_changes_in_api(req_url, desired_capacity) -> None:
    resp = requests.get(req_url)
    nginx_upstream = json.loads(resp.text)
    counter = 0
    while len(nginx_upstream["peers"]) != desired_capacity and counter < 10:
        time.sleep(RECONFIGURATION_DELAY)
        counter = counter + 1
        resp = requests.get(req_url)
        nginx_upstream = json.loads(resp.text)


def wait_for_changes_in_aws(autoscaling_client, group_name, desired_capacity) -> None:
    waiter_name = "autoscaling_completed"
    argument = (
        f"contains(AutoScalingGroups[?(starts_with(AutoScalingGroupName, `{group_name}`) == `true`)]."
        f"[length(Instances[?LifecycleState=='InService']) == `{desired_capacity}`][], `false`)"
    )
    waiter_config = {
        "version": 2,
        "waiters": {
            "autoscaling_completed": {
                "acceptors": [
                    {
                        "argument": argument,
                        "expected": True,
                        "matcher": "path",
                        "state": "success",
                    },
                    {
                        "argument": argument,
                        "expected": False,
                        "matcher": "path",
                        "state": "retry",
                    },
                ],
                "delay": 5,
                "maxAttempts": 20,
                "operation": "DescribeAutoScalingGroups",
            }
        },
    }
    waiter_model = WaiterModel(waiter_config)
    custom_waiter = create_waiter_with_client(waiter_name, waiter_model, autoscaling_client)
    custom_waiter.wait()


def scale_aws_group(autoscaling_client, group_name, desired_capacity) -> dict:
    counter = 0
    while counter < 10:
        try:
            response = autoscaling_client.set_desired_capacity(
                AutoScalingGroupName=group_name,
                DesiredCapacity=desired_capacity,
                HonorCooldown=True,
            )
            print(f"Scaling activity started: {response}")
            return response
        except autoscaling_client.exceptions.ScalingActivityInProgressFault:
            print("Scaling activity is in progress, wait for 60 seconds then retry.")
            counter = counter + 1
            time.sleep(RECONFIGURATION_DELAY)
    pytest.fail(f"Failed to scale the group {group_name}")


def get_aws_group_name(autoscaling_client, group_name) -> str:
    """
    Get AWS unique group name.

    :param autoscaling_client: AWS API
    :param group_name:
    :return: str
    """
    groups = autoscaling_client.describe_auto_scaling_groups()["AutoScalingGroups"]
    return list(filter(lambda group: group_name in group["AutoScalingGroupName"], groups))[0]["AutoScalingGroupName"]


class TestSmoke:
    @pytest.mark.parametrize(
        "test_data",
        [
            pytest.param(
                {
                    "group_name": "WebserverGroup1",
                    "api_url": "/http/upstreams/backend1",
                },
                id="backend1",
            ),
            pytest.param(
                {
                    "group_name": "WebserverGroup2",
                    "api_url": "/http/upstreams/backend2",
                },
                id="backend2",
            ),
            pytest.param(
                {
                    "group_name": "WebserverGroup3",
                    "api_url": "/stream/upstreams/tcp-backend",
                },
                id="tcp-backend",
            ),
        ],
    )
    def test_aws_scale_up(self, cli_arguments, autoscaling_client, test_data):
        desired_capacity = 5
        group_name = get_aws_group_name(autoscaling_client, test_data["group_name"])
        scale_aws_group(autoscaling_client, group_name, desired_capacity)
        wait_for_changes_in_aws(autoscaling_client, group_name, desired_capacity)
        wait_for_changes_in_api(
            f"{cli_arguments.nginx_api}/{NGINX_API_VERSION}{test_data['api_url']}",
            desired_capacity,
        )
        resp = requests.get(f"{cli_arguments.nginx_api}/{NGINX_API_VERSION}{test_data['api_url']}")
        nginx_upstream = json.loads(resp.text)
        assert (
            len(nginx_upstream["peers"]) == desired_capacity
        ), f"Expected {desired_capacity} servers, found: {nginx_upstream['peers']}"

    @pytest.mark.parametrize(
        "test_data",
        [
            pytest.param(
                {
                    "group_name": "WebserverGroup1",
                    "api_url": "/http/upstreams/backend1",
                },
                id="backend1",
            ),
            pytest.param(
                {
                    "group_name": "WebserverGroup2",
                    "api_url": "/http/upstreams/backend2",
                },
                id="backend2",
            ),
            pytest.param(
                {
                    "group_name": "WebserverGroup3",
                    "api_url": "/stream/upstreams/tcp-backend",
                },
                id="tcp-backend",
            ),
        ],
    )
    def test_aws_scale_down(self, cli_arguments, autoscaling_client, test_data):
        desired_capacity = 1
        group_name = get_aws_group_name(autoscaling_client, test_data["group_name"])
        scale_aws_group(autoscaling_client, group_name, desired_capacity)
        wait_for_changes_in_aws(autoscaling_client, group_name, desired_capacity)
        wait_for_changes_in_api(
            f"{cli_arguments.nginx_api}/{NGINX_API_VERSION}{test_data['api_url']}",
            desired_capacity,
        )
        resp = requests.get(f"{cli_arguments.nginx_api}/{NGINX_API_VERSION}{test_data['api_url']}")
        nginx_upstream = json.loads(resp.text)
        assert (
            len(nginx_upstream["peers"]) == desired_capacity
        ), f"Expected {desired_capacity} servers, found: {nginx_upstream['peers']}"
