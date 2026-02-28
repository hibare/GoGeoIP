export default [
  {
    ignores: ['dist/**', 'node_modules/**', 'pnpm-lock.yaml']
  },
  {
    rules: {
      'no-unused-vars': 'warn',
      'no-console': 'warn'
    }
  }
]
