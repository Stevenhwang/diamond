module.exports = {
  root: true,
  parserOptions: {
    parser: 'babel-eslint',
    sourceType: 'module'
  },
  env: {
    browser: true,
    node: true,
    es6: true
  },
  extends: ['plugin:vue/recommended', 'eslint:recommended', 'prettier'],

  rules: {
    'space-before-function-paren': 0,
    'vue/multi-word-component-names': 0,
    "vue/first-attribute-linebreak": 0,
    "vue/component-definition-name-casing": 0
  }
}
