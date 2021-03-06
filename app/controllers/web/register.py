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
import os

# Third Party Library
from django.views import View
from django.shortcuts import render
from django.contrib import messages
from django.shortcuts import redirect
from django.utils.translation import gettext as _

# Local Library
from app.controllers.controller import Controller
from app.modules.core.user import User as UserModule
from app.modules.core.decorators import redirect_if_authenticated
from app.modules.core.decorators import redirect_if_not_installed


class Register(View, Controller):
    """Register Page Controller"""

    template_name = 'templates/register.html'

    @redirect_if_not_installed
    @redirect_if_authenticated
    def get(self, request, token):

        self.__user = UserModule()

        if not self.__user.check_register_request(token):
            messages.error(request, _("Register token is expired or invalid, Please contact system administrator!"))
            return redirect("app.web.login")

        self.autoload_options()
        self.context_push({
            "page_title": _("Register · %s") % self.context_get("app_name", os.getenv("APP_NAME", "Silverback")),
            "register_request": self.__user.get_register_request_by_token(token)
        })

        return render(request, self.template_name, self.context_get())
