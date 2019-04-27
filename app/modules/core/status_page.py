"""
Status Page Module
"""

import os
import json
from django.utils import timezone
from datetime import datetime
from datetime import timedelta
from app.modules.entity.option_entity import Option_Entity
from app.modules.entity.incident_entity import Incident_Entity
from app.modules.entity.incident_update_entity import Incident_Update_Entity
from app.modules.entity.incident_update_component_entity import Incident_Update_Component_Entity
from app.modules.entity.metric_entity import Metric_Entity
from django.utils.translation import gettext as _
from app.modules.entity.component_group_entity import Component_Group_Entity
from app.modules.entity.component_entity import Component_Entity
from dateutil.relativedelta import relativedelta
from django.forms.fields import DateTimeField


class Status_Page():

    __option_entity = None
    __incident_entity = None
    __incident_update_entity = None
    __incident_update_component_entity = None
    __component_group_entity = None
    __component_entity = None
    __metric_entity = None

    def __init__(self):
        self.__option_entity = Option_Entity()
        self.__incident_entity = Incident_Entity()
        self.__incident_update_entity = Incident_Update_Entity()
        self.__incident_update_component_entity = Incident_Update_Component_Entity()
        self.__component_group_entity = Component_Group_Entity()
        self.__component_entity = Component_Entity()
        self.__metric_entity = Metric_Entity()

    def get_system_status(self):
        return "operational"

    def get_about_site(self):
        option = self.__option_entity.get_one_by_key("builder_about")
        return option.value if option else ""

    def get_logo_url(self):
        option = self.__option_entity.get_one_by_key("builder_logo_url")
        return option.value if option else ""

    def get_favicon_url(self):
        option = self.__option_entity.get_one_by_key("builder_favicon_url")
        return option.value if option else ""

    def get_incident_by_uri(self, uri):
        incident = self.__incident_entity.get_one_by_uri(uri)
        app_name = self.__option_entity.get_one_by_key("app_name")

        if incident:
            incident_data = {
                "headline": incident.name,
                "headline_class": "text-danger",
                "status": incident.status,
                "sub_headline": _("Incident Report for %s") % (app_name.value),
                "affected_components": [],
                "updates": []
            }

            updates = self.__incident_update_entity.get_all(incident.id)

            for update in updates:
                incident_data["updates"].append({
                    "type": update.status.title(),
                    "body": update.message,
                    "date": "%(date)s %(tz)s" % {
                        "date": update.datetime.strftime("%B %d, %H:%M"),
                        "tz": os.getenv("APP_TIMEZONE", "UTC")
                    }
                })

                components = self.__incident_update_component_entity.get_all(update.id)

                for component in components:
                    if component.component.name not in incident_data["affected_components"]:
                        incident_data["affected_components"].append(component.component.name)

            incident_data["affected_components"] = ", ".join(incident_data["affected_components"])

            return incident_data

        return False

    def get_incidents_for_period(self, period):

        today = timezone.now()
        datem = datetime(today.year, today.month, 1)

        from_date = datem - relativedelta(months=+(period - 1) * 3)
        to_date = datem - relativedelta(months=+(period * 3))

        period = "%(from)s - %(to)s" % {
            "from": from_date.strftime("%B %Y"),
            "to": (to_date + relativedelta(months=+1)).strftime("%B %Y")
        }

        from_date = datetime(from_date.year, from_date.month, 1)
        to_date = datetime(to_date.year, to_date.month, 1)

        incidents = []
        while from_date > to_date:
            current_incidents = []

            incidents_list = self.__incident_entity.get_incident_on_month(DateTimeField().clean(from_date))

            for incident in incidents_list:
                current_incidents.append({
                    "uri": incident.uri,
                    "subject": incident.name,
                    "class": "text-danger",
                    "status": incident.status,
                    "final_update": _("This incident has been resolved.") if incident.status == "closed" else _("This incident is still open."),
                    "period": self.__get_incident_period(incident)
                })

            current_date = from_date.strftime("%B %Y")
            incidents.append({
                "date": current_date,
                "incidents": current_incidents
            })
            from_date -= relativedelta(months=+1)

        return {
            "period": period,
            "incidents": incidents
        }

    def __get_incident_period(self, incident):
        updates = self.__get_incident_updates(incident.id)
        if len(updates):
            return "%(from)s %(tz)s - %(to)s" % {
                "from": incident.datetime.strftime("%B %d, %H:%M"),
                "tz": os.getenv("APP_TIMEZONE", "UTC"),
                "to": updates[len(updates)-1]["date"]
            }
        return "%(from)s %(tz)s" % {
            "from": incident.datetime.strftime("%B %d, %H:%M"),
            "tz": os.getenv("APP_TIMEZONE", "UTC")
        }

    def get_past_incidents(self, days=7):
        i = 0
        past_incidents = []
        while days > i:
            date = (datetime.now() - timedelta(days=i))
            incidents_result = []
            incidents = self.__incident_entity.get_incident_from_days(i)
            for incident in incidents:
                incidents_result.append({
                    "uri": incident.uri,
                    "subject": incident.name,
                    "class": "text-danger",
                    "status": incident.status,
                    "updates": self.__get_incident_updates(incident.id)
                })

            past_incidents.append({
                "date": date.strftime("%B %d, %Y"),
                "incidents": incidents_result
            })
            i += 1
        return past_incidents

    def __get_incident_updates(self, incident_id):
        updates_result = []
        updates = self.__incident_update_entity.get_all(incident_id)
        for update in updates:
            updates_result.append({
                "type": update.status.title(),
                "date": "%(date)s %(tz)s" % {
                    "date": update.datetime.strftime("%B %d, %H:%M"),
                    "tz": os.getenv("APP_TIMEZONE", "UTC")
                },
                "body": update.message
            })
        return updates_result

    def get_system_metrics(self):
        metrics = []
        option = self.__option_entity.get_one_by_key("builder_metrics")
        if option:
            items = json.loads(option.value)
            for item in items:
                if "m-" in item:
                    item = int(item.replace("m-", ""))
                    if item:
                        metric = self.__metric_entity.get_one_by_id(item)
                        if metric:
                            metrics.append({
                                "id": "metric_container_%d" % (metric.id),
                                "title": metric.title,
                                "xtitle": metric.x_axis,
                                "ytitle": metric.y_axis,
                                "day_data": self.__get_metics(metric, "day"),
                                "week_data": self.__get_metics(metric, "week"),
                                "month_data": self.__get_metics(metric, "month")
                            })
        return metrics

    def __get_metics(self, metric, period):
        if metric.source == "newrelic":
            print(metric.data["application"])
            print(metric.data["response_time"])
        return []

    def get_services(self):
        services = []
        option = self.__option_entity.get_one_by_key("builder_components")
        if option:
            items = json.loads(option.value)
            for item in items:
                if "c-" in item:
                    component = self.__component_entity.get_one_by_id(item.replace("c-", ""))
                    if component:
                        services.append({
                            "name": component.name,
                            "description": component.description,
                            "current_status": self.get_status(component.id, "component"),
                            "current_status_class": "bg-green",
                            "uptime_chart": self.get_uptime_chart(component.id, "component"),
                            "sub_services": []
                        })
                elif "g-" in item:
                    group = self.__component_group_entity.get_one_by_id(item.replace("g-", ""))
                    services.append({
                        "name": group.name,
                        "description": group.description,
                        "current_status": self.get_status(group.id, "group"),
                        "current_status_class": "bg-green",
                        "uptime_chart": self.get_uptime_chart(group.id, "group"),
                        "sub_services": self.get_sub_services(group.id)
                    })

        return services

    def get_sub_services(self, group_id):
        services = []
        items = self.__component_entity.get_all_components_by_group(group_id)
        for item in items:
            services.append({
                "name": item.name,
                "description": item.description,
                "current_status": self.get_status(item.id, "component"),
                "current_status_class": "bg-green",
                "uptime_chart": self.get_uptime_chart(item.id, "component"),
                "sub_services": []
            })
        return services

    def get_status(self, id, type):
        if type == "component":
            return "Operational"
        elif type == "group":
            return "Operational"

    def get_uptime_chart(self, id, type):
        return []
