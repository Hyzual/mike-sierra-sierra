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
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const TerserPlugin = require("terser-webpack-plugin");

const common_configurations = require("./webpack.common.js");

const css_minimizer_plugin = new CssMinimizerPlugin({
    minimizerOptions: {
        preset: ["default", { discardComments: { removeAll: true } }],
    },
});

const prod_configurations = common_configurations.map((config) =>
    merge(config, {
        mode: "production",
        optimization: {
            minimize: true,
            minimizer: [css_minimizer_plugin, new TerserPlugin()],
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
