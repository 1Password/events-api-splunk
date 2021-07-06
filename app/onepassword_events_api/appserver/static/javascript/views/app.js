/**
 * This is an example using pure react, with no JSX
 * If you would like to use JSX, you will need to use Babel to transpile your code
 * from JSK to JS. You will also need to use a task runner/module bundler to
 * help build your app before it can be used in the browser.
 * Some task runners/module bundlers are : gulp, grunt, webpack, and Parcel
 */

import React from "react";
import * as Setup from "./setup_page.js";
import { version } from "../../../../package.json";

// Since we don't have TS, I'm using constants to help with key label enforcement
const authToken = "authToken";
const limit = "limit";
const startAt = "startAt";
const signInCursorFile = "signInCursorFile";
const itemUsageCursorFile = "itemUsageCursorFile";
const error = "error";
const aud = "aud";
const audienceDEPRECATED = "com.1password.streamingservice";

const e = React.createElement;

export default class SetupPage extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			[authToken]: "",
			[error]: "",
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

		try {
			this.validateJWT(this.state[authToken]);
		} catch (err) {
			return this.setState({
				...this.state,
				[error]: err,
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
	}

	// validateJWT verifies that the token has 3 parts -
	// the header, payload, and signature
	// validateJWT only attempts to parse the payload to catch potential issues
	validateJWT(token) {
		const tokenComponents = token.split(".");
		if (tokenComponents.length !== 3) {
			throw "Invalid JSON Web Token";
		}
		const payload = JSON.parse(atob(tokenComponents[1]))
		if (!payload[aud] || payload[aud].length !== 1) {
			throw "Invalid JSON Web Token";
		}
		if (payload[aud][0] === audienceDEPRECATED) {
			throw "Please generate a new token";
		}
	}

	render() {
		return e("div", null, [
			e("h2", null, `1Password Events API for Splunk Setup Page - Version ${version}`),
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
			this.state.error && e("div", null, this.state.error)
		]);
	}
}
