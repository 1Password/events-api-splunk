import React, { useState } from "react";

import { onepassword_name_space } from "../views/setup_page";

const VERSION = "1.6.0";

const e = React.createElement;
export default function Wizard(props) {
  const totalSteps = props.steps.length;
  const [currentStep, setCurrentStep] = useState(1);
  const [result, setResult] = useState({ success: false, error: "" });
  const [authToken, setAuthToken] = useState("");

  const handleNext = async () => {
    setCurrentStep(currentStep + 1);
  };

  const handlePrev = async () => {
    setCurrentStep(currentStep - 1);
    setResult({
      success: false,
      error: "",
    });
  };

  const handleSkip = async () => {
    setCurrentStep(totalSteps);
  };

  const handleSubmit = async () => {
    const result = await props.handleSubmit(authToken);
    setResult(result);
  };

  // currentStep is using one-based indexing
  const stepDetails = props.steps[currentStep - 1];

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
            e("div", { className: "step-number" }, [currentStep]),
            e("div", { className: "step-desc" }, [stepDetails.step]),
          ]),
        stepDetails.img &&
          e("img", {
            src: `/static/app/${onepassword_name_space.app}/img/${stepDetails.img}`,
          }),
        currentStep === totalSteps &&
          e("div", { className: "token block" }, [
            e("input", {
              type: "text",
              placeholder: "Enter your token here",
              value: authToken,
              onChange: (e) => setAuthToken(e.target.value),
            }),
          ]),
        stepDetails.warning &&
          (!result.success || result.error) &&
          e("div", { className: "warning block" }, [
            "Your other Splunk apps or add-ons may be able to access your Events API token. Make sure you trust them before you add your token.",
          ]),
        result.error && e("div", { className: "error block" }, [result.error]),
        result.success &&
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
        e("a", { onClick: handleSkip }, [
          currentStep === 1 && "I already have my Events API token",
        ]),
        e("div", { className: "controls-buttons" }, [
          currentStep !== 1 &&
            e("button", { className: "prev btn", onClick: handlePrev }, [
              "Previous",
            ]),
          currentStep !== totalSteps
            ? e("button", { className: "next btn", onClick: handleNext }, [
                "Next",
              ])
            : !result.success
            ? e("button", { className: "next btn", onClick: handleSubmit }, [
                "Add Token",
              ])
            : e(
                "a",
                {
                  href: `/app/${onepassword_name_space.app}`,
                },
                [e("button", { className: "next btn" }, "Finish")]
              ),
        ]),
      ]),
    ]),
    e("div", { className: "progress-container" }, [
      e("div", { className: "step-short block" }, [stepDetails.stepShort]),
      e("div", { className: "step-count block" }, [
        `Step ${currentStep} of ${totalSteps}`,
      ]),
      e("div", { className: "progress block" }, [
        e(
          "div",
          {
            className: "progress-bar",
            role: "progressbar",
            style: {
              width: `${(currentStep * 100) / totalSteps}%`,
            },
          },
          null
        ),
      ]),
    ]),
  ]);
}
