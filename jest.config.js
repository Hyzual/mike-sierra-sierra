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

// For a detailed explanation regarding each configuration property, visit:
// https://jestjs.io/docs/en/configuration.html
module.exports = {
    preset: "ts-jest/presets/js-with-ts",

    // The root directory that Jest should scan for tests and modules within
    rootDir: __dirname,

    // The glob patterns Jest uses to detect test files
    testMatch: ["**/?(*.)+(test).ts"],

    // Automatically clear mock calls and instances between every test
    clearMocks: true,

    // Reset the module registry before running each individual test
    resetModules: true,

    // Make calling deprecated APIs throw helpful error messages
    errorOnDeprecated: true,

    // An array of glob patterns indicating a set of files for which coverage information should be collected
    collectCoverageFrom: ["**/*.ts"],

    // The directory where Jest should output its coverage files
    coverageDirectory: "coverage",

    coverageReporters: ["text"],

    // An array of regexp pattern strings used to skip coverage collection
    coveragePathIgnorePatterns: ["/node_modules/", "/assets/", "/tests/"],

    // The test environment that will be used for testing
    testEnvironment: "jest-environment-jsdom",

    // Transpile lit-html for tests because it uses ES modules
    transformIgnorePatterns: [
        "/node_modules/(?!(lit-html|lit-element)).+\\.js",
    ],

    moduleNameMapper: {
        // Ignore SVG imports in code under test
        "^.+\\.svg": "identity-obj-proxy",
    },
};
