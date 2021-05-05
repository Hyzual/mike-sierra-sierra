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

const { merge } = require("webpack-merge");
const { ESBuildMinifyPlugin } = require("esbuild-loader");

const common_configurations = require("./webpack.common.js");

const prod_configurations = common_configurations.map((config) =>
    merge(config, {
        mode: "production",
        optimization: {
            minimize: true,
            minimizer: [
                new ESBuildMinifyPlugin({
                    target: "es2020",
                    css: true,
                }),
            ],
        },
        stats: {
            all: false,
            assets: true,
            errors: true,
            errorDetails: true,
            performance: true,
            timings: true,
        },
    })
);

module.exports = prod_configurations;
