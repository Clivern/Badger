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

from app.exceptions.error_codes import ErrorCodes


class ServerError(Exception):
    """ServerError Custom Exception"""

    def __init__(self, message, error_code=None, extra={}):
        Exception.__init__(self, message)
        self.error_code = ErrorCodes().SERVER_ERROR if error_code is None else error_code
        self.extra = extra

    def get_error_code(self):
        return self.error_code