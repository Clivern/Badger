# Copyright 2019 Silverbackhq
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Third Party Library
from django.db import models


class Subscriber(models.Model):

    TYPE_CHOICES = (
        ('email', 'EMAIL'),
        ('phone', 'PHONE'),
        ('endpoint', 'ENDPOINT')
    )

    STATUS_CHOICES = (
        ('pending', 'PENDING'),
        ('verified', 'VERIFIED'),
        ('unverified', 'UNVERIFIED')
    )

    email = models.CharField(max_length=60, verbose_name="Email Address")
    phone = models.CharField(max_length=60, verbose_name="Phone")
    endpoint = models.CharField(max_length=200, verbose_name="Endpoint")
    auth_token = models.CharField(max_length=200, verbose_name="Auth Token")
    type = models.CharField(max_length=50, choices=TYPE_CHOICES, default="email", verbose_name="Type")
    status = models.CharField(max_length=50, choices=STATUS_CHOICES, default="pending", verbose_name="Status")
    external_id = models.CharField(max_length=200, verbose_name="External ID")
    created_at = models.DateTimeField(auto_now_add=True, verbose_name="Created at")
    updated_at = models.DateTimeField(auto_now=True, verbose_name="Updated at")

    class Meta:
        db_table = "app_subscriber"
