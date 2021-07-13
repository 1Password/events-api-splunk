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
    to: `build = ${package.version.replace(/\./g, "")}`
}

const setupXMLHash = {
    files: "./default/data/ui/views/setup_page_dashboard.xml",
    from: /script\=\"build\/.*\.js\"/,
    to: `script="build/${manifest["main.js"]}"`,
};

const appjsVersion = {
    files: "./appserver/static/javascript/views/app.js",
    from: /\"1Password Events API for Splunk Setup Page - Version .*\"/,
    to: `"1Password Events API for Splunk Setup Page - Version ${package.version}"`
};

const tarballVersion = {
    files: "../../Makefile",
    from: /VERSION=.*/,
    to: `VERSION=${package.version}`
}

const userAgentVersion = {
    files: "../../src/api/events.go",
    from: /const DefaultUserAgent \= \"1Password Events API for Splunk \/ .*\"/,
    to: `const DefaultUserAgent = "1Password Events API for Splunk / ${package.version}"`
}

replace(appConfVersion);
replace(appConfBuild);
replace(setupXMLHash);
replace(appjsVersion);
replace(tarballVersion);
replace(userAgentVersion);
