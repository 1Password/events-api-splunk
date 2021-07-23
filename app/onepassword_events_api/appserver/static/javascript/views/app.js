/**
 * This is an example using pure react, with no JSX
 * If you would like to use JSX, you will need to use Babel to transpile your code
 * from JSK to JS. You will also need to use a task runner/module bundler to
 * help build your app before it can be used in the browser.
 * Some task runners/module bundlers are : gulp, grunt, webpack, and Parcel
 */

import React from "react";
import * as Setup from "./setup_page.js";

// Since we don't have TS, I'm using constants to help with key label enforcement
const authToken = "authToken";
const limit = "limit";
const startAt = "startAt";
const signInCursorFile = "signInCursorFile";
const itemUsageCursorFile = "itemUsageCursorFile";
const error = "error";
const aud = "aud";
const audienceDEPRECATED = "com.1password.streamingservice";
const success = "success";

const e = React.createElement;

export default class SetupPage extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      [authToken]: "",
      [error]: "",
      [success]: false,
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({
      ...this.state,
      [event.target.name]: event.target.value,
    });
  }

  async handleSubmit(event) {
    event.preventDefault();

    const jwtError = this.validateJWT(this.state[authToken]);
    if (typeof jwtError !== "undefined") {
      return this.setState({
        ...this.state,
        [error]: jwtError,
        success: false,
      });
    }

    // Normalize inputs
    const options = {
      [authToken]: `"${this.state[authToken]}"`,
      [limit]: 100,
      [startAt]: "2020-01-01T00:00:00Z",
      [signInCursorFile]:
        '"/etc/apps/onepassword_events_api/local/signin_cursor_store"',
      [itemUsageCursorFile]:
        '"/etc/apps/onepassword_events_api/local/itemusage_cursor_store"',
    };

    await Setup.perform(splunkjs, options);

    return this.setState({
      ...this.state,
      [error]: "",
      [success]: true,
    });
  }

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
    if (!payload[aud] || payload[aud].length !== 1) {
      return "Invalid JSON Web Token - missing aud";
    }
    if (payload[aud][0] === audienceDEPRECATED) {
      return "Please generate a new token";
    }
  }

  render() {
    return e("div", null, [
      e(
        "h2",
        null,
        "1Password Events API for Splunk Setup Page - Version 1.4.2"
      ),
      e("div", null, [
        e("form", { onSubmit: this.handleSubmit }, [
          e("label", null, [
            e("div", null, ["Events API Token"]),
            e("input", {
              type: "text",
              name: authToken,
              value: this.state[authToken],
              onChange: this.handleChange,
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
