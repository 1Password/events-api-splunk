const path = require("path");

module.exports = {
    entry: "./appserver/static/javascript/setup_page.js",
    output: {
        filename: "main.js",
        path: path.resolve(__dirname, "appserver/static/build"),
    },
};
