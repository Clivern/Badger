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

# Standard Library
import json

# Third Party Library
from pyumetric import NewRelic_Provider

# Local Library
from app.modules.entity.metric_entity import MetricEntity
from app.modules.entity.option_entity import OptionEntity


class Metric():

    def __init__(self):
        self.__option_entity = OptionEntity()
        self.__metric_entity = MetricEntity()
        new_relic_api = self.__option_entity.get_one_by_key("newrelic_api_key")
        self.__newrelic = None
        if new_relic_api:
            self.__newrelic = NewRelic_Provider(new_relic_api.value)

    def get_one_by_id(self, id):
        metric = self.__metric_entity.get_one_by_id(id)

        if not metric:
            return False

        return {
            "id": metric.id,
            "title": metric.title,
            "description": metric.description,
            "source": metric.source,
            "data": metric.data,
            "x_axis": metric.x_axis,
            "y_axis": metric.y_axis
        }

    def get_one_by_title(self, title):
        metric = self.__metric_entity.get_one_by_title(title)

        if not metric:
            return False

        return {
            "id": metric.id,
            "title": metric.title,
            "description": metric.description,
            "source": metric.source,
            "data": metric.data,
            "x_axis": metric.x_axis,
            "y_axis": metric.y_axis
        }

    def insert_one(self, metric):
        return self.__metric_entity.insert_one(metric)

    def update_one_by_id(self, id, metric_data):
        return self.__metric_entity.update_one_by_id(id, metric_data)

    def count_all(self):
        return self.__metric_entity.count_all()

    def get_all(self, offset=None, limit=None):
        return self.__metric_entity.get_all(offset, limit)

    def delete_one_by_id(self, id):
        return self.__metric_entity.delete_one_by_id(id)

    def get_new_relic_apps(self):
        result = []
        try:
            response = self.__newrelic.get_apps()
        except Exception as e:
            raise Exception(e)
        apps = json.loads(response)
        for app in apps["applications"]:
            result.append({"key": app["id"], "value": app["name"]})
        return result
