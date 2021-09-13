/**
 * This is an example using pure react, with no JSX
 * If you would like to use JSX, you will need to use Babel to transpile your code
 * from JSK to JS. You will also need to use a task runner/module bundler to
 * help build your app before it can be used in the browser.
 * Some task runners/module bundlers are : gulp, grunt, webpack, and Parcel
 */

import React from "react";
import * as Setup from "./setup_page.js";
import Wizard from "../components/wizard.js";

const audienceDEPRECATED = "com.1password.streamingservice";
const e = React.createElement;
export default class SetupPage extends React.Component {
  constructor(props) {
    super(props);
  }

  handleSubmit = async (authToken) => {
    const errorMessage = this.validateJWT(authToken);
    if (typeof errorMessage !== "undefined") {
      return {
        error: errorMessage,
        success: false,
      };
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
      await Setup.perform(splunkjs, authToken, options);
    } catch (error) {
      console.log(error);
      return {
        error:
          "Something went wrong while storing your token - please try again.",
        success: false,
      };
    }

    return {
      error: "",
      success: true,
    };
  };

  // validateJWT verifies that the token has 3 parts -
  // the header, payload, and signature
  // validateJWT only attempts to parse the payload to catch potential issues
  validateJWT(token) {
    const tokenComponents = token.split(".");
    if (tokenComponents.length !== 3) {
      return "This doesn't look like a valid JSON Web Token.";
    }
    let payload;
    try {
      payload = JSON.parse(atob(tokenComponents[1]));
    } catch (error) {
      return "This doesn't look like a valid JSON Web Token.";
    }
    if (!payload.aud || payload.aud.length !== 1) {
      return "This doesn't look like a valid JSON Web Token.";
    }
    if (payload.aud[0] === audienceDEPRECATED) {
      return "Please generate a new token";
    }
  }

  render() {
    return e(
      Wizard,
      {
        steps: [
          {
            description: e("span", null, [
              "To get started, you'll need to generate an Events API token.",
              e("br"),
              e("br"),
              'Click "Generate an Events API token", sign in to your account on',
              e("br"),
              "1Password.com, then follow the onscreen instructions.",
              e("br"),
              e("br"),
              "After you get your token, come back here to enter it.",
            ]),
            warning: true,
            redirect: true,
          },
          {
            description: "Enter the token you got from 1Password.com:",
          },
        ],
        handleSubmit: this.handleSubmit,
      },
      null
    );
  }
}
