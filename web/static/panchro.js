/// <reference lib="dom" />

/**  @type {typeof document.querySelector} */
const $ = document.querySelector.bind(document)
/**  @type {typeof document.querySelectorAll} */
const $$ = document.querySelectorAll.bind(document)

$("form#upload").addEventListener("submit", async (ev) => {
  ev.preventDefault()
  /**  @type {HTMLFormElement} */
  const form = ev.target
  const formdata = new FormData(form)
  const response = await fetch(form.action, {
    method: form.method,
    headers: {
      Authorization: "Bearer " + formdata.get("token"),
    },
    body: formdata,
  })
  const status = $("div#upload-status")
  if (response.ok) {
    const photo = await response.json()
    status.innerHTML = `Upload successful: <a href="/${photo.id}">${photo.id}</a>`
    status.classList.add("success")
    status.classList.remove("error")
    form.reset()
  } else {
    status.innerHTML = `Error: ${await response.text()}`
    status.classList.add("error")
    status.classList.remove("success")
  }
})

$("form#delete").addEventListener("submit", async (ev) => {
  ev.preventDefault()
  /**  @type {HTMLFormElement} */
  const form = ev.target
  const formdata = new FormData(form)
  const response = await fetch(`${form.action}/${formdata.get("id")}`, {
    method: "delete",
    headers: {
      Authorization: "Bearer " + formdata.get("token"),
    },
    body: formdata,
  })
  const status = $("div#delete-status")
  if (response.ok) {
    status.innerHTML = `Delete successful</a>`
    status.classList.add("success")
    status.classList.remove("error")
    form.reset()
  } else {
    status.innerHTML = `Error: ${await response.text()}`
    status.classList.add("error")
    status.classList.remove("success")
  }
})
