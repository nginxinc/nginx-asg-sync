"""Describe project settings"""

import os

BASEDIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
PROJECT_ROOT = os.path.abspath(os.path.dirname(__file__))
TEST_DATA = f"{PROJECT_ROOT}/data"
DEFAULT_AWS_REGION = "us-east-2"
# Time in seconds to ensure reconfiguration changes in cluster
RECONFIGURATION_DELAY = 60
NGINX_API_VERSION = 4
