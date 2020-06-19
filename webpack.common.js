/*
 *   Copyright (C) 2020  Joris MASSON
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

const path = require("path");
const WebpackAssetsManifest = require("webpack-assets-manifest");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");
const FixStyleOnlyEntriesPlugin = require("webpack-fix-style-only-entries");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const ForkTsCheckerWebpackPlugin = require("fork-ts-checker-webpack-plugin");

const context = __dirname;
const output = {
    path: path.resolve(__dirname, "./assets"),
    filename: "[name]-[chunkhash].js",
    publicPath: "/assets/",
};

const css_extract_plugin = new MiniCssExtractPlugin({
    filename: "[name]-[chunkhash].css",
});

const css_style_only_plugin = new FixStyleOnlyEntriesPlugin({
    extensions: ["css"],
    silent: true,
});

const manifest_plugin = new WebpackAssetsManifest({
    output: "manifest.json",
    writeToDisk: true,
});

const typescript_type_check_plugin = new ForkTsCheckerWebpackPlugin();

const typescript_rule = {
    test: /\.ts(x?)$/,
    exclude: /node_modules/,
    use: [
        {
            loader: "ts-loader",
            options: { transpileOnly: true },
        },
    ],
};

const css_rule = {
    test: /\.css$/,
    use: [MiniCssExtractPlugin.loader, "css-loader"],
};

const configuration = {
    entry: {
        index: "./scripts/index.ts",
        style: "./styles/style.css",
    },
    context,
    output,
    module: {
        rules: [typescript_rule, css_rule],
    },
    plugins: [
        new CleanWebpackPlugin(),
        manifest_plugin,
        typescript_type_check_plugin,
        css_style_only_plugin,
        css_extract_plugin,
    ],
    resolve: {
        extensions: [".ts", ".js"],
    },
};

module.exports = [configuration];
