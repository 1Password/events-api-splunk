"use strict";

import * as Config from "./setup_configuration";
import * as Splunk from "./splunk_helpers";
import * as StoragePasswords from "./storage_passwords";

const INPUTS_CONF = "inputs";

const SIGNIN_INPUT =
  "script://$SPLUNK_HOME/etc/apps/onepassword_events_api/bin/signin_attempts";
const ITEMUSAAGE_INPUT =
  "script://$SPLUNK_HOME/etc/apps/onepassword_events_api/bin/item_usages";

const CUSTOM_CONF = "events_reporting";
const CUSTOM_CONF_STANZA = "config";

const SECRET_REALM = "events_reporting_realm";
const SECRET_NAME = "events_api_token";

export const VERSION = "1.7.0";
export const HOST = "1password.com";

export const onepassword_name_space = {
  owner: "nobody",
  app: "onepassword_events_api",
  sharing: "app",
};

export async function getIndexes(splunk_js_sdk) {
  // Create the Splunk JS SDK Service object
  const splunk_js_sdk_service = Config.create_splunk_js_sdk_service(
    splunk_js_sdk,
    onepassword_name_space
  );

  const indexes = await Config.get_indexes(splunk_js_sdk_service);

  return indexes;
}

export async function updateInputs(
  splunk_js_sdk_service,
  input,
  inputOptions,
  reload
) {
  await Splunk.update_configuration_file(
    splunk_js_sdk_service,
    INPUTS_CONF,
    input,
    inputOptions
  );

  if (reload) {
    await Config.reload_splunk_app(
      splunk_js_sdk_service,
      onepassword_name_space.app
    );
  }
}

export async function perform(
  splunk_js_sdk,
  authToken,
  setup_options,
  signin_options,
  itemusage_options
) {
  // Create the Splunk JS SDK Service object
  const splunk_js_sdk_service = Config.create_splunk_js_sdk_service(
    splunk_js_sdk,
    onepassword_name_space
  );

  // Get conf and do stuff to it
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

  // Disable and reload the app to update the config
  // Only reload when both inputs have been updated
  await updateInputs(
    splunk_js_sdk_service,
    SIGNIN_INPUT,
    { disabled: 1 },
    false
  );
  await updateInputs(
    splunk_js_sdk_service,
    ITEMUSAAGE_INPUT,
    { disabled: 1 },
    true
  );
  await updateInputs(
    splunk_js_sdk_service,
    SIGNIN_INPUT,
    signin_options,
    false
  );
  await updateInputs(
    splunk_js_sdk_service,
    ITEMUSAAGE_INPUT,
    itemusage_options,
    false
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
