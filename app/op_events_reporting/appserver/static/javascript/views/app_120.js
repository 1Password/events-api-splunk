/**
 * This is an example using pure react, with no JSX
 * If you would like to use JSX, you will need to use Babel to transpile your code
 * from JSK to JS. You will also need to use a task runner/module bundler to
 * help build your app before it can be used in the browser.
 * Some task runners/module bundlers are : gulp, grunt, webpack, and Parcel
 */

import * as Setup from "./setup_page.js";

define(["react", "splunkjs/splunk"], function (react, splunk_js_sdk) {
	const e = react.createElement;

	// Since we don't have TS, I'm using constants to help with key label enforcement
	const authToken = "authToken";
	const limit = "limit";
	const startAt = "startAt";
	const url = "url";
	const signInCursorFile = "signInCursorFile";
	const itemUsageCursorFile = "itemUsageCursorFile";

	class SetupPage extends react.Component {
		constructor(props) {
			super(props);

			this.state = {
				[authToken]: "",
				[url]: "https://events.1password.com",
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

			// Normalize inputs
			const options = {
				[authToken]: `"${this.state[authToken]}"`,
				[url]: `"${this.state[url]}"`,
				[limit]: 100,
				[startAt]: "2020-01-01T00:00:00Z",
				[signInCursorFile]:
					'"/etc/apps/op_events_reporting/local/signin_cursor_store"',
				[itemUsageCursorFile]:
					'"/etc/apps/op_events_reporting/local/itemusage_cursor_store"',
			};

			await Setup.perform(splunk_js_sdk, options);
		}

		render() {
			return e("div", null, [
				e("h2", null, "Events Reporting Setup Page - Version 1.2.0"),
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
						e("label", null, [
							e("div", null, ["Events API URL"]),
							e("input", {
								type: "text",
								name: url,
								value: this.state[url],
								onChange: this.handleChange,
							}),
						]),
						e("input", { type: "submit", value: "Submit" }),
					]),
				]),
			]);
		}
	}

	return e(SetupPage);
});
