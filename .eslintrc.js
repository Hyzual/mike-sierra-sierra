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

module.exports = {
    plugins: ["@typescript-eslint", "no-unsanitized", "jest"],
    extends: [
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended",
        "prettier",
        "prettier/@typescript-eslint",
    ],
    parser: "@typescript-eslint/parser",
    env: {
        es6: true,
        browser: true,
    },
    reportUnusedDisableDirectives: true,
    rules: {
        // Possible Errors
        "no-template-curly-in-string": "error",
        // Best Practices
        "array-callback-return": "warn",
        "consistent-return": "warn",
        curly: "error",
        "default-case": "error",
        "dot-notation": "error",
        eqeqeq: "error",
        "no-alert": "error",
        "no-console": "error",
        "no-caller": "error",
        "no-div-regex": "error",
        "no-else-return": "warn",
        "no-eval": "error",
        "no-extend-native": "error",
        "no-extra-bind": "error",
        "no-implicit-coercion": "error",
        "no-implied-eval": "error",
        "no-iterator": "error",
        "no-labels": "error",
        "no-lone-blocks": "error",
        "no-loop-func": "warn",
        "no-multi-str": "error",
        "no-new": "warn",
        "no-new-func": "error",
        "no-new-wrappers": "error",
        "no-param-reassign": "warn",
        "no-proto": "error",
        "no-return-assign": "error",
        "no-return-await": "error",
        "no-self-compare": "error",
        "no-sequences": "error",
        "no-throw-literal": "error",
        "no-unmodified-loop-condition": "error",
        "no-useless-call": "error",
        "no-useless-concat": "error",
        "no-useless-return": "warn",
        "no-void": "error",
        "no-with": "error",
        radix: "error",
        "require-await": "error",
        // Typescript
        "@typescript-eslint/camelcase": "off",
        "@typescript-eslint/class-literal-property-style": "error",
        "@typescript-eslint/consistent-type-assertions": [
            "error",
            { assertionStyle: "never" },
        ],
        "@typescript-eslint/explicit-function-return-type": "error",
        "@typescript-eslint/no-explicit-any": "error",
        "@typescript-eslint/no-non-null-assertion": "error",
        "no-unused-vars": "off", // Typescript rule requires disabling base rule
        "@typescript-eslint/no-unused-vars": "error",
        "@typescript-eslint/no-use-before-define": [
            "error",
            { functions: false, typedefs: false },
        ],
        // No-unsanitized
        "no-unsanitized/property": [
            "error",
            {
                escape: {
                    methods: ["sanitize"],
                },
            },
        ],
        "no-unsanitized/method": [
            "error",
            {
                escape: {
                    methods: ["sanitize"],
                },
            },
        ],
    },
    overrides: [
        {
            files: ["*.test.ts"],
            extends: ["plugin:jest/recommended"],
            rules: {
                "jest/consistent-test-it": "error",
                "jest/valid-title": "error",
                "jest/no-restricted-matchers": [
                    "error",
                    { resolves: "Use `expect(await promise)` instead." },
                ],
                "jest/prefer-spy-on": "error",
                "jest/require-top-level-describe": "error",
                "jest/prefer-hooks-on-top": "error",
            },
        },
        {
            files: [".eslintrc.js", "jest.config.js", "webpack.*.js"],
            env: {
                node: true,
            },
            rules: {
                "no-console": "off",
                "@typescript-eslint/no-var-requires": "off",
            },
        },
    ],
};
