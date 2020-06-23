#!/usr/bin/env bash
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -eEuo pipefail

export PROJ=mco-fyi
export APP=$PROJ
export TAG="gcr.io/$PROJ/$APP"

docker build --tag "$TAG" .
docker push "$TAG"
gcloud beta run deploy "$APP" \
  --image "$TAG"              \
  --platform "managed"        \
  --region "us-central1"      \
  --project "$PROJ"           \
  --concurrency 80            \
  --memory=1Gi                \
  --allow-unauthenticated
