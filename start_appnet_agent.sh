#!/usr/bin/bash

# Quit on error.
set -e
# Treat undefined variables as errors.
set -u

export APPMESH_RESOURCE_ARN=mesh/test_mesh/virtualNode/test_vn

export AGENT_ADMIN_MODE=uds

export ENVOY_ADMIN_MODE=uds

export AGENT_ADMIN_UDS_PATH=/tmp/appnet_admin.sock

docker run --rm \
	-v /tmp:/tmp \
	-e AWS_REGION=us-west-2 \
	-e APPNET_ENVOY_RESTART_COUNT \
	-e APPMESH_RESOURCE_ARN \
	-e AWS_ACCESS_KEY_ID \
	-e AWS_SECRET_ACCESS_KEY \
	-e AWS_SESSION_TOKEN \
	-e ENVOY_LOG_LEVEL \
	-e AGENT_ADMIN_MODE \
	-e AGENT_ADMIN_UDS_PATH \
	-p 9901:9901 \
	-p 9902:9902 \
	   appnet-agent
