const Storage = {}

Storage.get = function (name, isObject = false) {
  return isObject ? JSON.parse(localStorage.getItem(name)) : localStorage.getItem(name)
}

Storage.set = function (name, val, isObject = false) {
  localStorage.setItem(name, isObject ? JSON.stringify(val) : val)
}

Storage.add = function (name, addVal) {
  const oldVal = Storage.get(name)
  const newVal = oldVal.concat(addVal)
  Storage.set(name, newVal)
}

Storage.delete = function (name) {
  localStorage.removeItem(name)
}

export default Storage
