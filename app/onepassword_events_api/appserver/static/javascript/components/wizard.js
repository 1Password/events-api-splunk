import React, { useState } from "react";

import { VERSION, HOST, onepassword_name_space } from "../views/setup_page";

const e = React.createElement;

export const SetupWizard = (props) => {
  const [currentStep, setCurrentStep] = useState(1);
  const [result, setResult] = useState({ success: false, error: "" });
  const [authToken, setAuthToken] = useState("");
  const [index, setIndex] = useState({
    signIn: "main",
    itemUsage: "main",
    auditEvents: "main",
  });
  const [enabled, setEnabled] = useState({
    signIn: true,
    itemUsage: true,
    auditEvents: true,
  });
  const [loading, setLoading] = useState(false);

  const handleNext = () => {
    setCurrentStep(currentStep + 1);
  };

  const handleValidate = () => {
    const r = props.handleValidate(authToken);
    setResult({ ...result, error: r.error });
    if (r.success) {
      handleNext();
    }
  };

  const handleSubmit = async () => {
    setLoading(true);
    const signInOptions = {
      index: index.signIn,
      disabled: enabled.signIn ? 0 : 1,
    };
    const itemUsageOptions = {
      index: index.itemUsage,
      disabled: enabled.itemUsage ? 0 : 1,
    };
    const auditEventsOptions = {
      index: index.auditEvents,
      disabled: enabled.auditEvents ? 0 : 1,
    };
    const result = await props.handleSubmit(
      authToken,
      signInOptions,
      itemUsageOptions,
      auditEventsOptions
    );
    setLoading(false);
    setResult(result);
  };

  const steps = [
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
      children: e(React.Fragment, null, [
        (!result.success || result.error) &&
          e(
            "div",
            { className: "warning block" },
            "Your other Splunk apps or add-ons may be able to access your Events API token. Make sure you trust them before you add your token."
          ),
        e(
          "a",
          {
            target: "_blank",
            href: `https://start.${HOST}/signin?landing-page=%2Fintegrations%2Fevents_reporting%2Fcreate%3Ftype%3Dsplunk%26name%3D${location.hostname}`,
          },
          e("button", { className: "generate", onClick: handleNext }, [
            e("img", {
              className: "plus",
              src: "/static/app/onepassword_events_api/img/plus.svg",
            }),
            "Generate an Events API token",
          ])
        ),
      ]),
      controls: e(React.Fragment, null, [
        e(
          "a",
          { onClick: handleNext },
          currentStep === 1 && "I already have my Events API token"
        ),
        e("div"),
      ]),
    },
    {
      description: "Enter the token you got from 1Password.com:",
      children: e("div", { className: "token block" }, [
        e("input", {
          type: "text",
          placeholder: "Events API Token",
          value: authToken,
          onChange: (e) => setAuthToken(e.target.value),
        }),
      ]),
      controls: e(React.Fragment, null, [
        e("div"),
        e(
          "button",
          {
            className: "next btn",
            onClick: handleValidate,
          },
          "Next"
        ),
      ]),
    },
    {
      description: "Configure your Events Reporting for Splunk Scripted Inputs",
      children: e("div", { className: "indexes block" }, [
        e("div", { className: "index" }, [
          e("label", { className: "switch" }, [
            e("input", {
              type: "checkbox",
              checked: enabled.signIn,
              onChange: (e) =>
                setEnabled({ ...enabled, signIn: e.target.checked }),
            }),
            e("span"),
          ]),
          e("div", { className: "title" }, "Sign-in Attempts"),
          e(
            "select",
            {
              id: "index",
              value: index.signIn,
              disabled: !enabled.signIn,
              onChange: (e) => setIndex({ ...index, signIn: e.target.value }),
            },
            [
              props.indexes.map((i) => {
                return e("option", { value: i.name }, i.name);
              }),
            ]
          ),
        ]),
        e("div", { className: "index" }, [
          e("label", { className: "switch" }, [
            e("input", {
              type: "checkbox",
              checked: enabled.itemUsage,
              onChange: (e) =>
                setEnabled({ ...enabled, itemUsage: e.target.checked }),
            }),
            e("span"),
          ]),
          e("div", { className: "title" }, "Item Usage Events"),
          e(
            "select",
            {
              id: "index",
              value: index.itemUsage,
              disabled: !enabled.itemUsage,
              onChange: (e) =>
                setIndex({ ...index, itemUsage: e.target.value }),
            },
            [
              props.indexes.map((i) => {
                return e("option", { value: i.name }, i.name);
              }),
            ]
          ),
        ]),
        e("div", { className: "index" }, [
          e("label", { className: "switch" }, [
            e("input", {
              type: "checkbox",
              checked: enabled.auditEvents,
              onChange: (e) =>
                setEnabled({ ...enabled, auditEvents: e.target.checked }),
            }),
            e("span"),
          ]),
          e("div", { className: "title" }, "Audit Events"),
          e(
            "select",
            {
              id: "index",
              value: index.auditEvents,
              disabled: !enabled.auditEvents,
              onChange: (e) =>
                setIndex({ ...index, auditEvents: e.target.value }),
            },
            [
              props.indexes.map((i) => {
                return e("option", { value: i.name }, i.name);
              }),
            ]
          ),
        ]),
      ]),
      controls: e(React.Fragment, null, [
        e("div"),
        result.success
          ? e(
              "a",
              {
                href: `/app/${onepassword_name_space.app}`,
              },
              e(
                "button",
                {
                  className: result.success ? "finish btn" : "btn",
                  disabled: !result.success,
                },
                "Finish"
              )
            )
          : loading
          ? e("button", { className: "btn", disabled: true }, "Saving ...")
          : e(
              "button",
              { className: "next btn", onClick: handleSubmit },
              "Submit"
            ),
      ]),
    },
  ];

  // currentStep is using one-based indexing
  const stepDetails = steps[currentStep - 1];

  return e("div", { className: "container" }, [
    e("div", { className: "main-contents" }, [
      e("div", { className: "version block" }, [`Version ${VERSION}`]),
      e(
        "div",
        { className: "title" },
        e("h1", { className: "block" }, "1Password Events Reporting for Splunk")
      ),
      e("div", { className: "content" }, [
        e(
          "div",
          { className: "description" },
          e("div", { className: "block" }, stepDetails.description)
        ),
        stepDetails.children,
        result.error && e("div", { className: "error block" }, result.error),
        result.success &&
          e("div", { className: "success block" }, [
            "Events API Reporting for Splunk has been successfully setup.",
          ]),
      ]),
      e("div", { className: "controls block" }, [stepDetails.controls]),
    ]),
  ]);
};
