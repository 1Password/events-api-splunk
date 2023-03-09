/**
 * This is an example using pure react, with no JSX
 * If you would like to use JSX, you will need to use Babel to transpile your code
 * from JSK to JS. You will also need to use a task runner/module bundler to
 * help build your app before it can be used in the browser.
 * Some task runners/module bundlers are : gulp, grunt, webpack, and Parcel
 */

import React from "react";
import * as Setup from "./setup_page.js";
import { SetupWizard } from "../components/wizard.js";

const audienceDEPRECATED = "com.1password.streamingservice";
const e = React.createElement;
export default class SetupPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      indexes: [],
    };
  }

  async componentDidMount() {
    const indexes = await Setup.getIndexes(splunkjs);
    this.setState({
      indexes,
    });
  }

  handleValidate = (authToken) => {
    const errorMessage = this.validateJWT(authToken);
    if (typeof errorMessage !== "undefined") {
      return {
        error: errorMessage,
        success: false,
      };
    }

    return {
      error: "",
      success: true,
    };
  };

  handleSubmit = async (
    authToken,
    signInOption,
    itemUsageOption,
    auditEventsOption
  ) => {
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
      auditEventsCursorFile:
        '"/etc/apps/onepassword_events_api/local/auditevents_cursor_store"',
    };

    try {
      await Setup.perform(
        splunkjs,
        authToken,
        options,
        signInOption,
        itemUsageOption,
        auditEventsOption
      );
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
    const invalidJWTError = "This doesn't look like a valid JSON Web Token.";
    const deprecatedJWTErrror = "Please generate a new token.";

    const tokenComponents = token.split(".");
    if (tokenComponents.length !== 3) {
      return invalidJWTError;
    }
    let payload;
    try {
      payload = JSON.parse(atob(tokenComponents[1]));
    } catch (error) {
      return invalidJWTError;
    }
    if (!payload.aud || payload.aud.length !== 1) {
      return invalidJWTError;
    }
    if (payload.aud[0] === audienceDEPRECATED) {
      return deprecatedJWTErrror;
    }
  }

  render() {
    return e(
      SetupWizard,
      {
        handleSubmit: this.handleSubmit,
        handleValidate: this.handleValidate,
        indexes: this.state.indexes,
      },
      null
    );
  }
}
