/*
 *   Copyright (c) 2020 Joris MASSON

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

const mergeWebpackConfig = require("webpack-merge");
const OptimizeCssAssetsPlugin = require("optimize-css-assets-webpack-plugin");
const cssNano = require("cssnano");

const common_configurations = require("./webpack.common.js");

const cssOptimizerPlugin = new OptimizeCssAssetsPlugin({
    cssProcessor: cssNano,
    cssProcessorPluginOptions: {
        preset: ["default", { discardComments: { removeAll: true } }],
    },
});

const prod_configurations = common_configurations.map((config) =>
    mergeWebpackConfig(config, {
        mode: "production",
        plugins: [cssOptimizerPlugin],
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