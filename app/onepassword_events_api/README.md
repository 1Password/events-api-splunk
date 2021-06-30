# Overview

This application includes a scripted input to ingest data into Splunk from the Events Reporting API. After completing setup you will be able to monitor and alert on important 1Password event data.

# Binary File Declaration

bin/item_usages - binary source code has been included in the lib folder
bin/signin_attempts - binary source code has also been included in the lib folder

The source code is actually the same for these two binaries, but has been provided twice to meet the Splunk naming specification.

## Program Flow

This program starts in the `onepassword_events_api/default/app.conf`, where the `[install]` stanza's `is_configured` property is set to `false`. This causes Splunk to redirect to it's setup page that is specified so that an admin/user can configure it for use.

In the `onepassword_events_api/default/app.conf`'s, `[ui]` stanza there is a `setup_page` property that points to which resource should be used for the setup page. In this case it's pointing to `default/data/ui/views/setup_page_dashboard.xml`.

Once setup is finished, Splunk will need to be restarted in order to be aware of the new configuration variables. On startup,
the scripted inputs (included in `onepassword_events_api/bin/`) will be triggered and Splunk will index the retrieved 1Password events.

## Setup

Click on the 1Password Application in the Apps navigation pane and follow the setup instructions, making sure the input is valid (there is currently no validation on this page). Some possible mistakes might be:

1. Incorrect token
2. Incorrect Events API URL

Restart Splunk by choosing Settings -> Server Controls -> Restart Splunk

Click on the 1Password App again, this time you will be navigated to Search. Start seeing what data has already been ingested by filtering by the sign in attempts source type, ie `sourcetype="1password:insights:signin_attempts"`. If you don't see any events, try increasing the length of time to "All time".

## Debugging

If you've gone through the installation steps and do not see any ingested events, take a look at the logs at `$SPLUNK_HOME/var/log/splunk/splunkd.log` to see if there are any actionable steps.

Common errors:

```
ERROR ExecProcessor - message from "/opt/splunk/etc/apps/onepassword_events_api/bin/signin_attempts" panic: introspect request failed: could not unmarshal response: 404 page not found
```

- There is something wrong with your JWT token. A common issue is not copying the entire token during the setup flow.
