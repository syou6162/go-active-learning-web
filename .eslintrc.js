module.exports = {
  parserOptions: {
    parser: "@typescript-eslint/parser"
  },
  globals: {
    process: true
  },
  plugins: ["@typescript-eslint"],
  extends: [
    // add more generic rulesets here, such as:
    'eslint:recommended',
    'plugin:vue/recommended'
  ],
  rules: {
    // override/add rules settings here, such as:
    // 'vue/no-unused-vars': 'error'
    "no-unused-vars": ["error", { "args": "none" }],
    "@typescript-eslint/no-unused-vars": ["error", { "args": "none" }]
  }
}
