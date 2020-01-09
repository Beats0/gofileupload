const makeCancelable = (promise) => {
  let hasCanceled_ = false;
  const wrappedPromise = new Promise((resolve, reject) => {
    promise.then(val => (hasCanceled_ ? reject({ isCanceled: true }) : resolve(val)));
    promise.catch(error => (hasCanceled_ ? reject({ isCanceled: true }) : reject(error)));
  });
  return {
    promise: wrappedPromise,
    cancel() {
      hasCanceled_ = true;
    },
  };
};

// 保存到本地
// FIXME: 多个文件保存失败
const saveToDisk = (fileURL, fileName) => {
  // for non-IE
  if (!window.ActiveXObject) {
    const save = document.createElement('a')
    save.href = fileURL
    save.download = fileName || 'unknown'
    save.style = 'display:none;opacity:0;color:transparent;';
    (document.body || document.documentElement).appendChild(save)
    if (typeof save.click === 'function') {
      save.click()
    } else {
      save.target = '_blank'
      const event = document.createEvent('Event')
      event.initEvent('click', true, true)
      save.dispatchEvent(event)
    }
    (window.URL || window.webkitURL).revokeObjectURL(save.href)
  } else if (!!window.ActiveXObject && document.execCommand) {
    // for IE
    const _window = window.open(fileURL, '_blank')
    _window.document.close()
    _window.document.execCommand('SaveAs', true, fileName || fileURL)
    _window.close()
  }
}

export {
  makeCancelable,
  saveToDisk,
}
