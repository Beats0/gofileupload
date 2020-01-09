const backEnd = process.env.NODE_ENV === 'production'
  ? 'http://47.94.16.206:8081'
  : 'http://localhost:8081'

const checkSecretKey = process.env.NODE_ENV === 'production'
  ? 'this_is_some_release_CheckSecretKey'
  : 'this_is_some_CheckSecretKey'

export {
  backEnd,
  checkSecretKey,
}
