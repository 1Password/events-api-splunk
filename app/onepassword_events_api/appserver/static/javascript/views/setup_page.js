"use strict";

import * as Config from "./setup_configuration";
import * as Splunk from "./splunk_helpers";
import * as StoragePasswords from "./storage_passwords";

const CUSTOM_CONF = "events_reporting";
const CUSTOM_CONF_STANZA = "config";

const SECRET_REALM = "events_reporting_realm";
const SECRET_NAME = "events_api_token";

export const onepassword_name_space = {
  owner: "nobody",
  app: "onepassword_events_api",
  sharing: "app",
};

export async function perform(splunk_js_sdk, authToken, setup_options) {
  // Create the Splunk JS SDK Service object
  const splunk_js_sdk_service = Config.create_splunk_js_sdk_service(
    splunk_js_sdk,
    onepassword_name_space
  );

  // // Get conf and do stuff to it
  await Splunk.update_configuration_file(
    splunk_js_sdk_service,
    CUSTOM_CONF,
    CUSTOM_CONF_STANZA,
    setup_options
  );

  await StoragePasswords.write_secret(
    splunk_js_sdk_service,
    SECRET_REALM,
    SECRET_NAME,
    authToken
  );

  // Completes the setup, by access the app.conf's [install]
  // stanza and then setting the `is_configured` to true
  await Config.complete_setup(splunk_js_sdk_service);

  // Reloads the splunk app so that splunk is aware of the
  // updates made to the file system
  await Config.reload_splunk_app(
    splunk_js_sdk_service,
    onepassword_name_space.app
  );
}

export async function get_apps(splunk_js_sdk) {
  // Create the Splunk JS SDK Service object
  const splunk_js_sdk_service = Config.create_splunk_js_sdk_service(
    splunk_js_sdk,
    onepassword_name_space
  );

  const apps = await Splunk.get_apps(splunk_js_sdk_service);
  return apps.filter((appName) => appName !== onepassword_name_space.app);
}
