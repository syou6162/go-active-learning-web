module.exports = {
  parserOptions: {
    parser: "@typescript-eslint/parser"
  },
  globals: {
    process: true
  },
  env: {
    "jest/globals": true
  },  
  plugins: [
    "@typescript-eslint",
    "jest"
  ],
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
