module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: [
    'plugin:vue/essential',
    '@vue/airbnb',
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'semi': 0,
    'max-len' : ['error', {code : 300}],
    'no-underscore-dangle':  ["off", "always"],
    'no-constant-condition': ["off", "always"],
    'no-plusplus': ["off", "always"],
    'func-names': ["off", "always"],
    "import/prefer-default-export": 0,
  },
  parserOptions: {
    parser: 'babel-eslint',
  },
};
