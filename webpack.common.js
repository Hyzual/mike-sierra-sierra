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
const RemoveEmptyScriptsPlugin = require("webpack-remove-empty-scripts");
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

const clean_plugin = new CleanWebpackPlugin({
    cleanAfterEveryBuildPatterns: ["!css-assets/", "!css-assets/**"],
});

const manifest_plugin = new WebpackAssetsManifest({
    output: "manifest.json",
    writeToDisk: true,
});

const typescript_type_check_plugin = new ForkTsCheckerWebpackPlugin();

const remove_empty_style_js_file_plugin = new RemoveEmptyScriptsPlugin();

const typescript_rule = {
    test: /\.ts$/,
    loader: "esbuild-loader",
    options: {
        loader: "ts",
        target: "es2020",
    },
};

const css_rule = {
    test: /\.css$/,
    use: [MiniCssExtractPlugin.loader, "css-loader"],
};

const svg_regexp = /\.svg$/;

const svg_separate_file_rule = {
    test: svg_regexp,
    type: "asset/resource",
    generator: {
        filename: "css-assets/[name]-[hash][ext][query]",
    },
};

const svg_data_uri_rule = {
    test: svg_regexp,
    resourceQuery: /data/,
    type: "asset/inline",
};

const svg_raw_rule = {
    test: svg_regexp,
    resourceQuery: /raw/,
    type: "asset/source",
};

const configuration = {
    entry: {
        index: "./scripts/index.ts",
        style: "./styles/style.css",
    },
    context,
    target: ["web"],
    output,
    module: {
        rules: [
            typescript_rule,
            css_rule,
            svg_separate_file_rule,
            svg_data_uri_rule,
            svg_raw_rule,
        ],
    },
    plugins: [
        clean_plugin,
        manifest_plugin,
        typescript_type_check_plugin,
        css_extract_plugin,
        remove_empty_style_js_file_plugin,
    ],
    resolve: {
        extensions: [".ts", ".js"],
        fallback: { util: false },
    },
};

module.exports = [configuration];
