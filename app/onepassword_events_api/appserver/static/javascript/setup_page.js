"use strict";

import React from "react";
import ReactDOM from "react-dom";
import SetupPage from "./views/app";
import "../styles/setup_page.css";
import "../styles/switch.css";

ReactDOM.render(
  React.createElement(SetupPage),
  document.getElementById("main_container")
);
