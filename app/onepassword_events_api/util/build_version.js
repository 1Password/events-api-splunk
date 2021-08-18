const replace = require("replace-in-file");
const package = require("../package.json");
const manifest = require("../appserver/static/build/manifest.json");

const appConfVersion = {
  files: "./default/app.conf",
  from: /version \= .*/,
  to: `version = ${package.version}`,
};

const appConfBuild = {
  files: "./default/app.conf",
  from: /build \= .*/,
  to: `build = ${package.version.replace(/\./g, "")}`,
};

const setupXMLJSHash = {
  files: "./default/data/ui/views/setup_page_dashboard.xml",
  from: /script\=\"build\/.*\.js\"/,
  to: `script="build/${manifest["main.js"]}"`,
};

const setupXMLCSSHash = {
  files: "./default/data/ui/views/setup_page_dashboard.xml",
  from: /stylesheet\=\"build\/.*\.css\"/,
  to: `stylesheet="build/${manifest["main.css"]}"`,
};

const wizardVersion = {
  files: "./appserver/static/javascript/components/wizard.js",
  from: /const VERSION \= ".*";/,
  to: `const VERSION = "${package.version}";`,
};

replace(appConfVersion);
replace(appConfBuild);
replace(setupXMLJSHash);
replace(setupXMLCSSHash);
replace(wizardVersion);
