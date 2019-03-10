"""
History Web Controller
"""

# standard library
import os

# Django
from django.views import View
from django.http import HttpResponse

# local Django
from app.modules.core.context import Context
from app.modules.entity.option_entity import Option_Entity
from app.modules.core.decorators import redirect_if_not_installed

# Third party
from feedgen.feed import FeedGenerator


class AtomHistory(View):

    __context = None
    __option_entity = None
    __fg = None

    @redirect_if_not_installed
    def get(self, request):

        self.__fg = FeedGenerator()
        self.__context = Context()
        self.__option_entity = Option_Entity()

        self.__context.autoload_options()
        self.__context.push({
            "page_title": self.__context.get("app_name", os.getenv("APP_NAME", "Silverback")),
            "is_authenticated": request.user and request.user.is_authenticated
        })

        self.__fg.id('http://silverbackhq.org')
        self.__fg.title('Some Testfeed')
        self.__fg.author({'name': 'John Doe', 'email': 'john@silverbackhq.org'})
        self.__fg.link(href='http://example.com', rel='alternate')
        self.__fg.logo('http://ex.com/logo.jpg')
        self.__fg.subtitle('This is a cool feed!')
        self.__fg.link(href='http://silverbackhq.org/test.atom', rel='self')
        self.__fg.language('en')

        return HttpResponse(self.__fg.atom_str(), content_type='text/xml')


class RssHistory(View):

    __context = None
    __option_entity = None
    __fg = None

    @redirect_if_not_installed
    def get(self, request):

        self.__fg = FeedGenerator()
        self.__context = Context()
        self.__option_entity = Option_Entity()

        self.__context.autoload_options()
        self.__context.push({
            "page_title": self.__context.get("app_name", os.getenv("APP_NAME", "Silverback")),
            "is_authenticated": request.user and request.user.is_authenticated
        })

        self.__fg.id('http://silverbackhq.org')
        self.__fg.title('Some Testfeed')
        self.__fg.author({'name': 'John Doe', 'email': 'john@silverbackhq.org'})
        self.__fg.link(href='http://example.com', rel='alternate')
        self.__fg.logo('http://ex.com/logo.jpg')
        self.__fg.subtitle('This is a cool feed!')
        self.__fg.link(href='http://silverbackhq.org/test.atom', rel='self')
        self.__fg.language('en')

        return HttpResponse(self.__fg.atom_str(), content_type='text/xml')