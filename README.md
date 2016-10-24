# Concourse Calendar Resource

This is an [implementation of a Concourse resource](https://concourse.ci/implementing-resources.html) which allows you to trigger Concourse jobs and create calendar events using a Google Calendar. It uses a generic calendar client interface so can be extended to work with other calendar providers.

## Requirements

* A Google account. This account will be used by the resource and is not necessarily the account containing the calendar it will manage. Inside the developer console of this account you will need to generate the following:
    * A [service account](https://developers.google.com/identity/protocols/OAuth2ServiceAccount). This is what the calendar resource acts on the behalf of. It will have its own email address, which you will delegate access to later.
    * Credentials. This comes in the form of JSON you download. Supply the JSON to the Concourse resource in your pipeline (see [Examples: Resource configuration](#resource-configuration) below).
* A Google Calendar. This can be in a separate account. You will need to delegate read and/or write access to the service account, which will have its own email address. Read access is required for triggering jobs, whereas write access is required for creating events in a calendar.

## Examples

### Resource configuration

Here is an example resource configuration for a pipeline:

```
resource_types:
- name: calendar
  type: docker-image
  source:
    repository: henryknott/calendar-resource
    tag: release-v1.1

resources:
- name: team-calendar
  type: calendar
  source:
    provider: google
    calendar_id: calendar-containing-events@gmail.com
    event_name: Test
    credentials: {"type": "service_account","project_id": "some-random-145222","private_key_id": "REDACTED","private_key": "-----BEGIN PRIVATE KEY-----REDACTED-----END PRIVATE KEY-----\n","client_email": "calendar-dev@some-random-145222.iam.gserviceaccount.com","client_id": "REDACTED","auth_uri": "https://accounts.google.com/o/oauth2/auth","token_uri": "https://accounts.google.com/o/oauth2/token","auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/calendar-dev%40some-random-145222.iam.gserviceaccount.com"}
```

`provider`: The only provider supported currently is Google Calendars.

`calendar_id`: This uniquely identifies the calendar within your account. For Google calendars this will be the email address of the account.

`event_name`: Events with this name will trigger a Concourse job when they are happening.

`credentials`: For Google calendars this will be your service account credentials as JSON. The JSON should be minified and supplied as a single-line value.

### Pipelines

#### Trigger a job

Based on the example resource configuration above you can use a Google calendar event to trigger a job like this:

```
jobs:
- name: calendar-resource-example
  plan:
  - get: team-calendar
    trigger: true
  - task: print-calendar-event-details
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: busybox
      inputs:
        - name: team-calendar
      run:
        path: sh
        args:
        - -c
        - |
          cat team-calendar/input
```

Concourse will poll for new versions of the resource. Any event in the calendar called `Test` will be considered a new version of the resource, and will therefore trigger this job.

#### Add a calendar event

The resource can also be used to add an event to a Google calendar:

```
jobs:
- name: calendar-resource-example-two
  plan:
  - task: add-calendar-event
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: busybox
      run:
        path: sh
        args:
        - -c
        - |
          echo "This job does nothing in particular, other than add an event to a Google calendar"
    on_success:
      put: team-calendar
      params:
        summary: This is the event name/title
        description: Put useful descriptive information here.
        time_zone: Europe/London
        start_time: 2016-10-24T13:00:00+01:00
        end_time: 2016-10-24T14:00:00+01:00
```

The start and end time values are strings formatted to RFC3339.

