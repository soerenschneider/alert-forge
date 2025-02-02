import random
import json
import requests
from datetime import datetime, timedelta

alertmanager_webhook_url = 'http://localhost:8888/webhook'
#alertmanager_webhook_url = 'https://am.svc.ez.soeren.cloud/webhook'


# Sample severities, alertnames, instances, and jobs to randomly assign to alerts
severities = ["Critical", "High", "Major", "Warning", "Medium", "Low", "Info"]
alertnames = ["CPU_Usage", "Memory_Usage", "Disk_Space", "Network_Latency", "Service_Down", "Database_Error", "App_Error", "API_Failure", "Load_Avg"]
instances = ["server1:9100", "server2:9100", "server3:9100", "database1:9100", "webapp1:9100", "api_gateway:443", "load_balancer:80", "cache_server"]
jobs = ["prometheus", "api_service", "web_service", "database_service", "load_balancer_service", "auth_service"]

n = 0

def generate_alert():
    # Randomly select severity, alertname, instance, job
    severity = random.choice(severities)
    alertname = random.choice(alertnames)
    instance = random.choice(instances)
    job = random.choice(jobs)

    # Create the current time in ISO format
    global n
    start = datetime.utcnow() + timedelta(hours=n)
    n+=1
    start_time = start.isoformat() + "Z"  # Use UTC time for the alert start time
    if random.random() >= 0.75:
        end_time = (start + timedelta(hours=1)).isoformat() + "Z"  # Prometheus' 'inactive' end time format
    else:
        end_time = "0001-01-01T00:00:00Z"  # Prometheus' 'inactive' end time format

# Construct the alert in the required format
    alert = {
        "status": "firing",  # The alert status (could also be 'resolved')
        "labels": {
            "alertname": alertname,
            "dc": "eu-west-1",  # Example datacenter label
            "instance": instance,
            "job": job,
            "severity": severity
        },
        "annotations": {
            "description": f"Description for {alertname} alert with severity {severity} on {instance}."
        },
        "startsAt": start_time,
        "endsAt": end_time,
        "generatorURL": f"http://example.com:9090/graph?g0.expr={alertname}+%3E+0\u0026g0.tab=1"
    }

    return alert

def send_alert_to_alertmanager(alerts):
    headers = {
        'Content-Type': 'application/json',
    }

    # Construct the payload for the webhook request
    payload = {
        "receiver": "webhook",
        "status": "firing",
        "alerts": alerts,
        "groupLabels": {
            "alertname": alerts[0]["labels"]["alertname"],
            "job": alerts[0]["labels"]["job"]
        },
        "commonLabels": alerts[0]["labels"],
        "commonAnnotations": {
            "description": "This is a generated alert"
        },
        "externalURL": "http://example.com:9093",
        "version": "4",
        "groupKey": f"{{}}:{{alertname=\"{alerts[0]['labels']['alertname']}\", job=\"{alerts[0]['labels']['job']}\"}}"
    }

    try:
        response = requests.post(alertmanager_webhook_url, data=json.dumps(payload), headers=headers)

        if response.status_code == 200:
            print(f"Alerts successfully posted to Alertmanager")
        else:
            print(f"Failed to post alerts: {response.status_code}, {response.text}")

    except Exception as e:
        print(f"Error posting alerts: {e}")

def generate_and_send_alerts(num_alerts=50):
    alerts = []
    for _ in range(num_alerts):
        alert = generate_alert()
        alerts.append(alert)

    send_alert_to_alertmanager(alerts)

if __name__ == "__main__":
    generate_and_send_alerts(20)
