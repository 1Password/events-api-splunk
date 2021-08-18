import React from "react";

const VERSION = "1.6.0";

const e = React.createElement;
export default class Wizard extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      currentStep: 1,
      totalSteps: props.steps.length,
      success: false,
      error: "",
      authToken: "",
    };
  }

  handleNext = async (event) => {
    event.preventDefault();
    this.setState({
      currentStep: this.state.currentStep + 1,
    });
  };

  handlePrev = async (event) => {
    event.preventDefault();
    this.setState({
      currentStep: this.state.currentStep - 1,
      success: false,
      error: "",
    });
  };

  handleSkip = async (event) => {
    event.preventDefault();
    this.setState({
      currentStep: this.state.totalSteps,
    });
  };

  handleSubmit = async (event) => {
    event.preventDefault();
    const result = await this.props.handleSubmit(this.state.authToken);
    this.setState(result);
  };

  render() {
    const stepDetails = this.props.steps[this.state.currentStep - 1];
    return e("div", { className: "container" }, [
      e("div", { className: "main-contents" }, [
        e("div", { className: "version" }, [`Version ${VERSION}`]),
        e("div", { className: "title" }, [
          e(
            "h1",
            { className: "block" },
            "1Password Events Reporting for Splunk"
          ),
        ]),
        e("div", { className: "content" }, [
          e(
            "div",
            { className: "description" },
            e("div", { className: "block" }, [stepDetails.description])
          ),
          stepDetails.step !== "" &&
            e("div", { className: "step block" }, [
              e("div", { className: "step-number" }, [this.state.currentStep]),
              e("div", { className: "step-desc" }, [stepDetails.step]),
            ]),
          stepDetails.img &&
            e("img", {
              src: `/static/app/onepassword_events_api/img/${stepDetails.img}`,
            }),
          this.state.currentStep === this.state.totalSteps &&
            e("div", { className: "token block" }, [
              e("input", {
                type: "text",
                placeholder: "Enter your token here",
                value: this.state.authToken,
                onChange: (e) => {
                  this.setState({
                    authToken: e.target.value,
                  });
                },
              }),
            ]),
          stepDetails.warning &&
            (!this.state.success || this.state.error) &&
            e("div", { className: "warning block" }, [
              "Your other Splunk apps or add-ons may be able to access your Events API token. Make sure you trust them before you add your token.",
            ]),
          this.state.error &&
            e("div", { className: "error block" }, [this.state.error]),
          this.state.success &&
            e("div", { className: "success block" }, [
              "Your token has been successfully updated. If this is ",
              e(
                "strong",
                null,
                "the first time you're setting up 1Password Events Reporting for Splunk, you'll have to enable the scripted inputs."
              ),
              e("br"),
              e("br"),
              "If 1Password Events Reporting for Splunk ",
              e(
                "strong",
                null,
                "has already been setup, you'll have to disable and re-enable the scripted inputs for the changes to take effect."
              ),
              e("br"),
              e("br"),
              "For more information, check out the ",
              e(
                "a",
                {
                  target: "_blank",
                  href: "https://support.1password.com/events-reporting-splunk",
                },
                "support article here."
              ),
            ]),
        ]),
        e("div", { className: "controls block" }, [
          e("a", { onClick: this.handleSkip }, [
            this.state.currentStep === 1
              ? "I already have my Events API token"
              : "",
          ]),
          e("div", { className: "controls-buttons" }, [
            this.state.currentStep !== 1 &&
              e("button", { className: "prev btn", onClick: this.handlePrev }, [
                "Previous",
              ]),
            this.state.currentStep !== this.state.totalSteps
              ? e(
                  "button",
                  { className: "next btn", onClick: this.handleNext },
                  ["Next"]
                )
              : !this.state.success
              ? e(
                  "button",
                  { className: "next btn", onClick: this.handleSubmit },
                  ["Add Token"]
                )
              : e(
                  "a",
                  {
                    href: "/app/onepassword_events_api",
                  },
                  [e("button", { className: "next btn" }, "Finish")]
                ),
          ]),
        ]),
      ]),
      e("div", { className: "progress-container" }, [
        e("div", { className: "step-short block" }, [stepDetails.stepShort]),
        e("div", { className: "step-count block" }, [
          `Step ${this.state.currentStep} of ${this.state.totalSteps}`,
        ]),
        e("div", { className: "progress block" }, [
          e(
            "div",
            {
              className: "progress-bar",
              role: "progressbar",
              style: {
                width: `${(this.state.currentStep * 100) /
                  this.state.totalSteps}%`,
              },
            },
            null
          ),
        ]),
      ]),
    ]);
  }
}
