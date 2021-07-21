/**
 * This is an example using pure react, with no JSX
 * If you would like to use JSX, you will need to use Babel to transpile your code
 * from JSK to JS. You will also need to use a task runner/module bundler to
 * help build your app before it can be used in the browser.
 * Some task runners/module bundlers are : gulp, grunt, webpack, and Parcel
 */

import * as Setup from "./setup_page_150.js";

define(["react", "splunkjs/splunk"], function (react, splunk_js_sdk) {
  const e = react.createElement;

  // Since we don't have TS, I'm using constants to help with key label enforcement
  const audienceDEPRECATED = "com.1password.streamingservice";

  class SetupPage extends react.Component {
    constructor(props) {
      super(props);

      this.state = {
        authToken: "",
        error: "",
        success: false,
        apps: [],
      };
    }

    async componentDidMount() {
      let appList = [];
      try {
        appList = await Setup.get_apps(splunk_js_sdk);
      } catch {
        this.setState({
          error: "Something went wrong - please refresh before continuing.",
        });
        return;
      }
      this.setState({
        apps: appList,
      });
    }

    handleSubmit = async (event) => {
      event.preventDefault();

      const errorMessage = this.validateJWT(this.state.authToken);
      if (typeof errorMessage !== "undefined") {
        this.setState({
          error: errorMessage,
          success: false,
        });
        return;
      }

      const options = {
        limit: 100,
        startAt: "2020-01-01T00:00:00Z",
        signInCursorFile:
          '"/etc/apps/onepassword_events_api/local/signin_cursor_store"',
        itemUsageCursorFile:
          '"/etc/apps/onepassword_events_api/local/itemusage_cursor_store"',
      };

      try {
        await Setup.perform(splunk_js_sdk, this.state.authToken, options);
      } catch (error) {
        console.log(error);
        this.setState({
          error:
            "Something went wrong while storing your token - please try again.",
          success: false,
        });
        return;
      }

      this.setState({
        error: "",
        success: true,
      });
    };

    // validateJWT verifies that the token has 3 parts -
    // the header, payload, and signature
    // validateJWT only attempts to parse the payload to catch potential issues
    validateJWT(token) {
      const tokenComponents = token.split(".");
      if (tokenComponents.length !== 3) {
        return "Invalid JSON Web Token - too short";
      }
      let payload;
      try {
        payload = JSON.parse(atob(tokenComponents[1]));
      } catch (error) {
        return "Invalid JSON Web Token - " + error.message;
      }
      if (!payload.aud || payload.aud.length !== 1) {
        return "Invalid JSON Web Token - missing aud";
      }
      if (payload.aud[0] === audienceDEPRECATED) {
        return "Please generate a new token";
      }
    }

    render() {
      return e("div", null, [
        e(
          "h1",
          null,
          "1Password Events API for Splunk Setup Page - Version 1.5.0"
        ),
        this.state.apps.length > 0 &&
          e("div", null, [
            e("div", { class: "warning" }, [
              e("h3", null, [
                "WARNING: Any installed app could gain access to your token. Before saving it below, make sure you know and trust the following applications:",
              ]),
              e("div", null, this.state.apps.join(", ")),
            ]),
          ]),
        e("div", null, [
          e("form", { onSubmit: this.handleSubmit }, [
            e("label", null, [
              e("h3", null, ["Events API Token"]),
              e("input", {
                type: "text",
                value: this.state.authToken,
                onChange: (e) => {
                  this.setState({
                    authToken: e.target.value,
                  });
                },
              }),
            ]),
            e("input", { type: "submit", value: "Submit" }),
          ]),
        ]),
        this.state.error && e("div", { class: "error" }, this.state.error),
        this.state.success &&
          e("div", { class: "success" }, [
            "Your token has been successfully updated. If this is the first time you're setting up 1Password Events API for Splunk, you'll have to enable the scripted inputs. If 1Password Events API for Splunk has already been setup, you'll have to disable and re-enable the scripted inputs for the changes to take effect.",
            e("br"),
            e("br"),
            "For more information, check out the support article ",
            e(
              "a",
              {
                href: "https://support.1password.com/events-reporting-splunk",
              },
              "here."
            ),
          ]),
      ]);
    }
  }

  return e(SetupPage);
});